//go:build integration

package integration

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	"github.com/edgelesssys/constellation/coordinator/atls"
	"github.com/edgelesssys/constellation/coordinator/core"
	"github.com/edgelesssys/constellation/coordinator/kms"
	"github.com/edgelesssys/constellation/coordinator/pubapi/pubproto"
	"github.com/edgelesssys/constellation/coordinator/role"
	"github.com/edgelesssys/constellation/coordinator/store"
	"github.com/edgelesssys/constellation/coordinator/storewrapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/*
Notes regarding the integration test implementation:

Scale:
'numberPeers' should be < 30, otherwise activation might stuck, because something with the docker network
doesn't scale well (maybe > 50 wireguard kernel interfaces are the reason).
With over 150 nodes, the node activation will fail due to Docker internal network naming issues.
This could be further extended, but currently the number of possible nodes is enough for this test.

Usage of docker library:
Sometimes the API calls are slower than using the 'sh docker ...' commands. This is specifically the case
for the termination. However, to keep the code clean, we accept this tradeoff and use the library functions.
*/

const (
	publicgRPCPort            = "9000"
	constellationImageName    = "constellation:latest"
	etcdImageName             = "bitnami/etcd:3.5.2"
	etcdOverlayNetwork        = "constellationIntegrationTest"
	masterSecret              = "ConstellationIntegrationTest"
	localLogDirectory         = "/tmp/coordinator/logs"
	numberFirstActivation     = 3
	numberSecondaryActivation = 3
	numberThirdActivation     = 3
)

var (
	hostconfigConstellationPeer = &container.HostConfig{
		Binds:      []string{"/dev/net/tun:/dev/net/tun"}, // necessary for wireguard interface creation
		CapAdd:     strslice.StrSlice{"NET_ADMIN"},        // necessary for wireguard interface creation
		AutoRemove: true,
	}
	configConstellationPeer = &container.Config{
		Image:        constellationImageName,
		AttachStdout: true, // necessary to attach to the container log
		AttachStderr: true, // necessary to attach to the container log
		Tty:          true, // necessary to attach to the container log
	}

	hostconfigEtcd = &container.HostConfig{
		AutoRemove: true,
	}
	configEtcd = &container.Config{
		Image: etcdImageName,
		Env: []string{
			"ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379",
			"ETCD_ADVERTISE_CLIENT_URLS=http://127.0.0.1:2379",
			"ETCD_LOG_LEVEL=debug",
			"ETCD_DATA_DIR=/bitnami/etcd/data",
		},
		Entrypoint:   []string{"/opt/bitnami/etcd/bin/etcd"},
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	constellationDockerImageBuildOptions = types.ImageBuildOptions{
		Dockerfile:     "test/Dockerfile",
		Tags:           []string{constellationImageName},
		Remove:         true,
		ForceRemove:    true,
		SuppressOutput: false,
		PullParent:     true,
	}
	containerLogConfig = types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	}

	wgExecConfig = types.ExecConfig{
		Cmd:          []string{"wg"},
		AttachStdout: true,
		AttachStderr: true,
	}
	pingExecConfig = types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
	}
	activeCoordinators []string
	coordinatorCounter int
	nodeCounter        int
)

type peerInfo struct {
	dockerData    container.ContainerCreateCreatedBody
	isCoordinator bool
	vpnIP         string
}

func TestMain(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	activePeers := make(map[string]peerInfo)

	defer goleak.VerifyNone(t,
		// https://github.com/census-instrumentation/opencensus-go/issues/1262
		goleak.IgnoreTopFunction("go.opencensus.io/stats/view.(*worker).start"),
		// https://github.com/kubernetes/klog/issues/282, https://github.com/kubernetes/klog/issues/188
		goleak.IgnoreTopFunction("k8s.io/klog/v2.(*loggingT).flushDaemon"),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cwd, err := os.Getwd()
	require.NoError(err)
	require.NoError(os.Chdir(filepath.Join(cwd, "..")))
	require.NoError(createTempDir())
	// setup Docker containers
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	require.NoError(err)
	defer cli.Close()

	versionInfo, err := cli.Info(ctx)
	require.NoError(err)
	t.Logf("start integration test, local docker version %v", versionInfo.ServerVersion)

	require.NoError(imageBuild(ctx, cli))
	defer cli.ImageRemove(ctx, constellationImageName, types.ImageRemoveOptions{Force: true, PruneChildren: true})

	reader, err := cli.ImagePull(ctx, etcdImageName, types.ImagePullOptions{})
	require.NoError(err)
	_, err = io.Copy(os.Stdout, reader)
	require.NoError(err)
	require.NoError(reader.Close())

	// Add another docker network to be able to resolve etcd-storage from the coordinator.
	// This is not possible in the default "bridge" network.
	dockerNetwork, err := cli.NetworkCreate(ctx, etcdOverlayNetwork, types.NetworkCreate{Driver: "bridge", Internal: true})
	require.NoError(err)
	defer cli.NetworkRemove(ctx, etcdOverlayNetwork)

	// setup etcd
	t.Log("create etcd container...")
	respEtcd, err := cli.ContainerCreate(ctx, configEtcd, hostconfigEtcd, nil, nil, "etcd-storage")
	require.NoError(err)
	require.NoError(cli.ContainerStart(ctx, respEtcd.ID, types.ContainerStartOptions{}))
	defer killDockerContainer(ctx, cli, respEtcd)
	require.NoError(cli.NetworkConnect(ctx, dockerNetwork.ID, respEtcd.ID, nil))
	etcdData, err := cli.ContainerInspect(ctx, respEtcd.ID)
	require.NoError(err)
	etcdIPAddr := etcdData.NetworkSettings.DefaultNetworkSettings.IPAddress
	etcdstore, err := store.NewEtcdStore(net.JoinHostPort(etcdIPAddr, "2379"), false, zap.NewNop())
	require.NoError(err)
	defer etcdstore.Close()

	defer killDockerContainers(ctx, cli, activePeers)
	// setup coordinator containers
	t.Log("create 1st coordinator container...")
	require.NoError(createCoordinatorContainer(ctx, cli, "master-1", dockerNetwork.ID, activePeers))
	t.Log("create 2nd coordinator container...")
	require.NoError(createCoordinatorContainer(ctx, cli, "master-2", dockerNetwork.ID, activePeers))
	// 1st activation phase
	ips, err := spawnContainers(ctx, cli, numberFirstActivation, activePeers)
	require.NoError(err)

	t.Logf("node ips: %v", ips)
	t.Log("activate coordinator...")
	start := time.Now()
	assert.NoError(startCoordinator(ctx, activeCoordinators[0], ips))
	elapsed := time.Since(start)
	t.Logf("activation took %v", elapsed)
	// activate additional coordinator
	require.NoError(addNewCoordinatorToCoordinator(ctx, activeCoordinators[1], activeCoordinators[0]))
	require.NoError(updateVPNIPs(activePeers, etcdstore))

	t.Log("count peers in instances")
	countPeersTest(ctx, t, cli, wgExecConfig, activePeers)
	t.Log("start ping test")
	pingTest(ctx, t, cli, pingExecConfig, activePeers, etcdstore)

	// 2nd activation phase
	ips, err = spawnContainers(ctx, cli, numberSecondaryActivation, activePeers)
	require.NoError(err)
	t.Logf("node ips: %v", ips)
	t.Log("add additional nodes")
	start = time.Now()
	assert.NoError(addNewNodesToCoordinator(ctx, activeCoordinators[1], ips))
	elapsed = time.Since(start)
	t.Logf("adding took %v", elapsed)
	require.NoError(updateVPNIPs(activePeers, etcdstore))

	t.Log("count peers in instances")
	countPeersTest(ctx, t, cli, wgExecConfig, activePeers)
	t.Log("start ping test")
	pingTest(ctx, t, cli, pingExecConfig, activePeers, etcdstore)

	// ----------------------------------------------------------------
	t.Log("create 3nd coordinator container...")
	require.NoError(createCoordinatorContainer(ctx, cli, "master-3", dockerNetwork.ID, activePeers))
	// activate additional coordinator
	require.NoError(addNewCoordinatorToCoordinator(ctx, activeCoordinators[2], activeCoordinators[1]))
	require.NoError(updateVPNIPs(activePeers, etcdstore))

	// 3rd activation phase
	ips, err = spawnContainers(ctx, cli, numberThirdActivation, activePeers)
	require.NoError(err)
	t.Logf("node ips: %v", ips)
	t.Log("add additional nodes")
	start = time.Now()
	assert.NoError(addNewNodesToCoordinator(ctx, activeCoordinators[2], ips))
	elapsed = time.Since(start)
	t.Logf("adding took %v", elapsed)
	require.NoError(updateVPNIPs(activePeers, etcdstore))

	t.Log("count peers in instances")
	countPeersTest(ctx, t, cli, wgExecConfig, activePeers)
	t.Log("start ping test")
	pingTest(ctx, t, cli, pingExecConfig, activePeers, etcdstore)
}

// helper methods
func startCoordinator(ctx context.Context, coordinatorAddr string, ips []string) error {
	tlsConfig, err := atls.CreateAttestationClientTLSConfig([]atls.Validator{&core.MockValidator{}})
	if err != nil {
		return err
	}

	conn, err := grpc.DialContext(ctx, net.JoinHostPort(coordinatorAddr, publicgRPCPort), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pubproto.NewAPIClient(conn)
	adminKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return err
	}
	adminKey = adminKey.PublicKey()

	stream, err := client.ActivateAsCoordinator(ctx, &pubproto.ActivateAsCoordinatorRequest{
		AdminVpnPubKey:     adminKey[:],
		NodePublicIps:      ips,
		MasterSecret:       []byte(masterSecret),
		KmsUri:             kms.ClusterKMSURI,
		StorageUri:         kms.NoStoreURI,
		KeyEncryptionKeyId: "",
		UseExistingKek:     false,
	})
	if err != nil {
		return err
	}

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func createTempDir() error {
	if err := os.RemoveAll(localLogDirectory); err != nil && !os.IsNotExist(err) {
		return err
	} 		
	return os.MkdirAll(localLogDirectory, 0o755)
}

func addNewCoordinatorToCoordinator(ctx context.Context, newCoordinatorAddr, oldCoordinatorAddr string) error {
	tlsConfig, err := atls.CreateAttestationClientTLSConfig([]atls.Validator{&core.MockValidator{}})
	if err != nil {
		return err
	}

	conn, err := grpc.DialContext(ctx, net.JoinHostPort(oldCoordinatorAddr, publicgRPCPort), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pubproto.NewAPIClient(conn)

	_, err = client.ActivateAdditionalCoordinator(ctx, &pubproto.ActivateAdditionalCoordinatorRequest{
		CoordinatorPublicIp: newCoordinatorAddr,
	})
	if err != nil {
		return err
	}
	return nil
}

func addNewNodesToCoordinator(ctx context.Context, coordinatorAddr string, ips []string) error {
	tlsConfig, err := atls.CreateAttestationClientTLSConfig([]atls.Validator{&core.MockValidator{}})
	if err != nil {
		return err
	}

	conn, err := grpc.DialContext(ctx, net.JoinHostPort(coordinatorAddr, publicgRPCPort), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pubproto.NewAPIClient(conn)

	stream, err := client.ActivateAdditionalNodes(ctx, &pubproto.ActivateAdditionalNodesRequest{NodePublicIps: ips})
	if err != nil {
		return err
	}
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func spawnContainers(ctx context.Context, cli *client.Client, count int, activeContainers map[string]peerInfo) ([]string, error) {
	tmpPeerIPs := make([]string, 0, count)
	// spawn client container(s) and obtain their docker network ip address
	for i := 0; i < count; i++ {
		resp, err := createNewNode(ctx, cli)
		if err != nil {
			return nil, err
		}
		attachDockerContainerStdoutStderrToFile(ctx, cli, resp.containerResponse.ID, role.Node)
		tmpPeerIPs = append(tmpPeerIPs, resp.dockerIPAddr)
		containerData, err := cli.ContainerInspect(ctx, resp.containerResponse.ID)
		if err != nil {
			return nil, err
		}
		activeContainers[containerData.NetworkSettings.DefaultNetworkSettings.IPAddress] = peerInfo{dockerData: resp.containerResponse, isCoordinator: false}
	}
	return tmpPeerIPs, blockUntilUp(ctx, tmpPeerIPs)
}

func createCoordinatorContainer(ctx context.Context, cli *client.Client, name, dockerNetworkID string, activePeers map[string]peerInfo) error {
	resp, err := cli.ContainerCreate(ctx, configConstellationPeer, hostconfigConstellationPeer, nil, nil, name)
	if err != nil {
		return err
	}
	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	attachDockerContainerStdoutStderrToFile(ctx, cli, resp.ID, role.Coordinator)
	coordinatorData, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return err
	}
	activePeers[coordinatorData.NetworkSettings.DefaultNetworkSettings.IPAddress] = peerInfo{dockerData: resp, isCoordinator: true}
	activeCoordinators = append(activeCoordinators, coordinatorData.NetworkSettings.DefaultNetworkSettings.IPAddress)
	return cli.NetworkConnect(ctx, dockerNetworkID, resp.ID, nil)
}

// Make the port forward binding, so we can access the coordinator from the host
func makeBinding(ip, internalPort string, externalPort string) nat.PortMap {
	binding := nat.PortBinding{
		HostIP:   ip,
		HostPort: externalPort,
	}
	bindingMap := map[nat.Port][]nat.PortBinding{nat.Port(fmt.Sprintf("%s/tcp", internalPort)): {binding}}
	return nat.PortMap(bindingMap)
}

func killDockerContainers(ctx context.Context, cli *client.Client, activeContainers map[string]peerInfo) {
	for _, v := range activeContainers {
		killDockerContainer(ctx, cli, v.dockerData)
	}
}

func killDockerContainer(ctx context.Context, cli *client.Client, container container.ContainerCreateCreatedBody) {
	fmt.Print("Kill container ", container.ID[:10], "... ")
	if err := cli.ContainerKill(ctx, container.ID, "9"); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Success")
}

func attachDockerContainerStdoutStderrToFile(ctx context.Context, cli *client.Client, id string, peerRole role.Role) {
	resp, err := cli.ContainerLogs(ctx, id, containerLogConfig)
	if err != nil {
		panic(err)
	}
	var file *os.File
	switch peerRole {
	case role.Node:
		file, err = os.Create(fmt.Sprintf("%s/node-%d", localLogDirectory, nodeCounter))
		nodeCounter += 1
	case role.Coordinator:
		file, err = os.Create(fmt.Sprintf("%s/coordinator-%d", localLogDirectory, coordinatorCounter))
		coordinatorCounter += 1
	default:
		panic("invalid role")
	}
	if err != nil {
		panic(err)
	}
	go io.Copy(file, resp) // TODO: this goroutine leaks
}

func imageBuild(ctx context.Context, dockerClient *client.Client) error {
	// Docker need a BuildContext, generate it...
	tar, err := archive.TarWithOptions(".", &archive.TarOptions{})
	if err != nil {
		return err
	}

	resp, err := dockerClient.ImageBuild(ctx, tar, constellationDockerImageBuildOptions)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		return err
	}
	// Block until EOF, so the build has finished if we continue
	_, err = io.Copy(io.Discard, resp.Body)

	return err
}

// count number of wireguard peers within all active docker containers
func countPeersTest(ctx context.Context, t *testing.T, cli *client.Client, execConfig types.ExecConfig, activeContainers map[string]peerInfo) {
	t.Run("countPeerTest", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)
		for ip, id := range activeContainers {
			respIDExecCreate, err := cli.ContainerExecCreate(ctx, id.dockerData.ID, execConfig)
			require.NoError(err)
			respID, err := cli.ContainerExecAttach(ctx, respIDExecCreate.ID, types.ExecStartCheck{})
			require.NoError(err)
			output, err := io.ReadAll(respID.Reader)
			require.NoError(err)
			respID.Close()
			countedPeers := strings.Count(string(output), "peer")
			fmt.Printf("% 3d peers in container %s [%s] out of % 3d total nodes \n", countedPeers, id.dockerData.ID, ip, len(activeContainers))

			assert.Equal(len(activeContainers), countedPeers)
		}
	})
}

func pingTest(ctx context.Context, t *testing.T, cli *client.Client, execConfig types.ExecConfig, activeContainers map[string]peerInfo, etcdstore store.Store) {
	t.Run("pingTest", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)
		peerVPNIPsWithoutAdmins, err := getPeerVPNIPsFromEtcd(etcdstore)
		require.NoError(err)
		// all nodes + coordinator are peers
		require.Equal(len(peerVPNIPsWithoutAdmins), len(activeContainers))

		for i := 0; i < len(peerVPNIPsWithoutAdmins); i++ {
			execConfig.Cmd = []string{"ping", "-q", "-c", "1", "-W", "1", peerVPNIPsWithoutAdmins[i]}
			for _, id := range activeContainers {
				fmt.Printf("Ping from container %v | % 19s to container % 19s", id.dockerData.ID, id.vpnIP, peerVPNIPsWithoutAdmins[i])

				respIDExecCreate, err := cli.ContainerExecCreate(ctx, id.dockerData.ID, execConfig)
				require.NoError(err)

				err = cli.ContainerExecStart(ctx, respIDExecCreate.ID, types.ExecStartCheck{})
				require.NoError(err)

				resp, err := cli.ContainerExecInspect(ctx, respIDExecCreate.ID)
				require.NoError(err)
				assert.Equal(0, resp.ExitCode)
				if resp.ExitCode == 0 {
					fmt.Printf(" ...Success\n")
				} else {
					fmt.Printf(" ...Failure\n")
				}
			}
		}
	})
}

type newNodeData struct {
	containerResponse container.ContainerCreateCreatedBody
	dockerIPAddr      string
}

// pass error one level up
func createNewNode(ctx context.Context, cli *client.Client) (*newNodeData, error) {
	resp, err := cli.ContainerCreate(ctx, configConstellationPeer, hostconfigConstellationPeer, nil, nil, "")
	if err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	containerData, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return nil, err
	}
	fmt.Printf("created Node %v\n", containerData.ID)
	return &newNodeData{resp, containerData.NetworkSettings.IPAddress}, nil
}

func awaitPeerResponse(ctx context.Context, ip string, tlsConfig *tls.Config) error {
	// Block, so the connection gets established/fails immediately
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, net.JoinHostPort(ip, publicgRPCPort), grpc.WithBlock(), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		return err
	}
	return conn.Close()
}

func blockUntilUp(ctx context.Context, peerIPs []string) error {
	tlsConfig, err := atls.CreateAttestationClientTLSConfig([]atls.Validator{&core.MockValidator{}})
	if err != nil {
		return err
	}
	for _, ip := range peerIPs {
		// Block, so the connection gets established/fails immediately
		if err := awaitPeerResponse(ctx, ip, tlsConfig); err != nil {
			return err
		}
	}
	return nil
}

func getPeerVPNIPsFromEtcd(etcdstore store.Store) ([]string, error) {
	peers, err := storewrapper.StoreWrapper{Store: etcdstore}.GetPeers()
	if err != nil {
		return nil, err
	}

	vpnIPS := make([]string, 0, len(peers))

	for _, peer := range peers {
		if peer.Role != role.Admin {
			vpnIPS = append(vpnIPS, peer.VPNIP)
		}
	}
	return vpnIPS, nil
}

func translatePublicToVPNIP(publicIP string, etcdstore store.Store) (string, error) {
	peers, err := storewrapper.StoreWrapper{Store: etcdstore}.GetPeers()
	if err != nil {
		return "", err
	}
	for _, peer := range peers {
		if peer.PublicIP == publicIP {
			return peer.VPNIP, nil
		}
	}
	return "", errors.New("Did not found VPN IP")
}

func updateVPNIPs(activeContainers map[string]peerInfo, etcdstore store.Store) error {
	for publicIP, v := range activeContainers {
		vpnIP, err := translatePublicToVPNIP(publicIP, etcdstore)
		if err != nil {
			return err
		}
		v.vpnIP = vpnIP
		activeContainers[publicIP] = v
	}
	return nil
}

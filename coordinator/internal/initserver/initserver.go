package initserver

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/edgelesssys/constellation/coordinator/config"
	"github.com/edgelesssys/constellation/coordinator/initproto"
	"github.com/edgelesssys/constellation/coordinator/internal/diskencryption"
	"github.com/edgelesssys/constellation/coordinator/internal/kubernetes"
	"github.com/edgelesssys/constellation/coordinator/nodestate"
	"github.com/edgelesssys/constellation/coordinator/role"
	"github.com/edgelesssys/constellation/coordinator/util"
	attestationtypes "github.com/edgelesssys/constellation/internal/attestation/types"
	"github.com/edgelesssys/constellation/internal/file"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	nodeLock *sync.Mutex

	kube        ClusterInitializer
	disk        EncryptedDisk
	fileHandler file.Handler
	grpcServer  *grpc.Server

	logger *zap.Logger

	initproto.UnimplementedAPIServer
}

func New(nodeLock *sync.Mutex, kube ClusterInitializer, logger *zap.Logger) *Server {
	logger = logger.Named("initServer")
	server := &Server{
		nodeLock: nodeLock,
		disk:     diskencryption.New(),
		kube:     kube,
		logger:   logger,
	}

	grpcLogger := logger.Named("gRPC")
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(grpcLogger),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(grpcLogger),
		)),
	)
	initproto.RegisterAPIServer(grpcServer, server)

	server.grpcServer = grpcServer

	return server
}

func (s *Server) Serve(ip, port string) error {
	lis, err := net.Listen("tcp", net.JoinHostPort(ip, port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	return s.grpcServer.Serve(lis)
}

func (s *Server) Init(ctx context.Context, req *initproto.InitRequest) (*initproto.InitResponse, error) {
	if ok := s.nodeLock.TryLock(); !ok {
		go s.grpcServer.GracefulStop()
		return nil, status.Error(codes.FailedPrecondition, "node is already being activated")
	}

	id, err := s.deriveAttestationID(req.MasterSecret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	if err := s.setupDisk(req.MasterSecret); err != nil {
		return nil, status.Errorf(codes.Internal, "setting up disk: %s", err)
	}

	state := nodestate.NodeState{
		Role:      role.Coordinator,
		OwnerID:   id.Owner,
		ClusterID: id.Cluster,
	}
	if err := state.ToFile(s.fileHandler); err != nil {
		return nil, status.Errorf(codes.Internal, "persisting node state: %s", err)
	}

	kubeconfig, err := s.kube.InitCluster(ctx,
		req.AutoscalingNodeGroups,
		req.CloudServiceAccountUri,
		req.KubernetesVersion,
		id,
		kubernetes.KMSConfig{
			MasterSecret:       req.MasterSecret,
			KMSURI:             req.KmsUri,
			StorageURI:         req.StorageUri,
			KeyEncryptionKeyID: req.KeyEncryptionKeyId,
			UseExistingKEK:     req.UseExistingKek,
		},
		sshProtoKeysToMap(req.SshUserKeys),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "initializing cluster: %s", err)
	}

	return &initproto.InitResponse{
		Kubeconfig: kubeconfig,
		OwnerId:    id.Owner,
		ClusterId:  id.Cluster,
	}, nil
}

func (s *Server) setupDisk(masterSecret []byte) error {
	if err := s.disk.Open(); err != nil {
		return fmt.Errorf("opening encrypted disk: %w", err)
	}
	defer s.disk.Close()

	uuid, err := s.disk.UUID()
	if err != nil {
		return fmt.Errorf("retrieving uuid of disk: %w", err)
	}
	uuid = strings.ToLower(uuid)

	// TODO: Choose a way to salt the key derivation
	diskKey, err := util.DeriveKey(masterSecret, []byte("Constellation"), []byte("key"+uuid), 32)
	if err != nil {
		return err
	}

	return s.disk.UpdatePassphrase(string(diskKey))
}

func (s *Server) deriveAttestationID(masterSecret []byte) (attestationtypes.ID, error) {
	clusterID, err := util.GenerateRandomBytes(config.RNGLengthDefault)
	if err != nil {
		return attestationtypes.ID{}, err
	}

	// TODO: Choose a way to salt the key derivation
	ownerID, err := util.DeriveKey(masterSecret, []byte("Constellation"), []byte("id"), config.RNGLengthDefault)
	if err != nil {
		return attestationtypes.ID{}, err
	}

	return attestationtypes.ID{Owner: ownerID, Cluster: clusterID}, nil
}

func sshProtoKeysToMap(keys []*initproto.SSHUserKey) map[string]string {
	keyMap := make(map[string]string)
	for _, key := range keys {
		keyMap[key.Username] = key.PublicKey
	}
	return keyMap
}

type ClusterInitializer interface {
	InitCluster(
		ctx context.Context,
		autoscalingNodeGroups []string,
		cloudServiceAccountURI string,
		kubernetesVersion string,
		id attestationtypes.ID,
		config kubernetes.KMSConfig,
		sshUserKeys map[string]string,
	) ([]byte, error)
}

// EncryptedDisk manages the encrypted state disk.
type EncryptedDisk interface {
	// Open prepares the underlying device for disk operations.
	Open() error
	// Close closes the underlying device.
	Close() error
	// UUID gets the device's UUID.
	UUID() (string, error)
	// UpdatePassphrase switches the initial random passphrase of the encrypted disk to a permanent passphrase.
	UpdatePassphrase(passphrase string) error
}

package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"

	"github.com/edgelesssys/constellation/cli/internal/cloudcmd"
	"github.com/edgelesssys/constellation/coordinator/util"
	"github.com/edgelesssys/constellation/internal/atls"
	"github.com/edgelesssys/constellation/internal/cloud/cloudprovider"
	"github.com/edgelesssys/constellation/internal/constants"
	"github.com/edgelesssys/constellation/internal/file"
	"github.com/edgelesssys/constellation/internal/grpc/dialer"
	"github.com/edgelesssys/constellation/verify/verifyproto"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// NewVerifyCmd returns a new cobra.Command for the verify command.
func NewVerifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify {aws|azure|gcp}",
		Short: "Verify the confidential properties of a Constellation cluster",
		Long: `Verify the confidential properties of a Constellation cluster.

If arguments aren't specified, values are read from ` + "`" + constants.ClusterIDsFileName + "`.",
		Args: cobra.MatchAll(
			cobra.ExactArgs(1),
			isCloudProvider(0),
			warnAWS(0),
		),
		RunE: runVerify,
	}
	cmd.Flags().String("owner-id", "", "verify using the owner identity derived from the master secret")
	cmd.Flags().String("unique-id", "", "verify using the unique cluster identity")
	cmd.Flags().StringP("node-endpoint", "e", "", "endpoint of the node to verify, passed as HOST[:PORT]")
	return cmd
}

func runVerify(cmd *cobra.Command, args []string) error {
	provider := cloudprovider.FromString(args[0])
	fileHandler := file.NewHandler(afero.NewOsFs())
	verifyClient := &constellationVerifier{dialer: dialer.New(nil, nil, &net.Dialer{})}
	return verify(cmd, provider, fileHandler, verifyClient)
}

func verify(
	cmd *cobra.Command, provider cloudprovider.Provider, fileHandler file.Handler, verifyClient verifyClient,
) error {
	flags, err := parseVerifyFlags(cmd, fileHandler)
	if err != nil {
		return err
	}

	config, err := readConfig(cmd.OutOrStdout(), fileHandler, flags.configPath, provider)
	if err != nil {
		return fmt.Errorf("reading and validating config: %w", err)
	}

	validators, err := cloudcmd.NewValidators(provider, config)
	if err != nil {
		return err
	}

	if err := validators.UpdateInitPCRs(flags.ownerID, flags.clusterID); err != nil {
		return err
	}
	if validators.Warnings() != "" {
		cmd.Print(validators.Warnings())
	}

	nonce, err := util.GenerateRandomBytes(32)
	if err != nil {
		return err
	}
	userData, err := util.GenerateRandomBytes(32)
	if err != nil {
		return err
	}

	if err := verifyClient.Verify(
		cmd.Context(),
		flags.endpoint,
		&verifyproto.GetAttestationRequest{
			Nonce:    nonce,
			UserData: userData,
		},
		validators.V()[0],
	); err != nil {
		return err
	}

	cmd.Println("OK")
	return nil
}

func parseVerifyFlags(cmd *cobra.Command, fileHandler file.Handler) (verifyFlags, error) {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return verifyFlags{}, fmt.Errorf("parsing config path argument: %w", err)
	}
	ownerID, err := cmd.Flags().GetString("owner-id")
	if err != nil {
		return verifyFlags{}, fmt.Errorf("parsing owner-id argument: %w", err)
	}
	clusterID, err := cmd.Flags().GetString("unique-id")
	if err != nil {
		return verifyFlags{}, fmt.Errorf("parsing unique-id argument: %w", err)
	}
	endpoint, err := cmd.Flags().GetString("node-endpoint")
	if err != nil {
		return verifyFlags{}, fmt.Errorf("parsing node-endpoint argument: %w", err)
	}

	// Get empty values from ID file
	emptyEndpoint := endpoint == ""
	emptyIDs := ownerID == "" && clusterID == ""
	if emptyEndpoint || emptyIDs {
		if details, err := readIds(fileHandler); err == nil {
			if emptyEndpoint {
				cmd.Printf("Using endpoint from %q. Specify --node-endpoint to override this.\n", constants.ClusterIDsFileName)
				endpoint = details.Endpoint
			}
			if emptyIDs {
				cmd.Printf("Using IDs from %q. Specify --owner-id  and/or --unique-id to override this.\n", constants.ClusterIDsFileName)
				ownerID = details.OwnerID
				clusterID = details.ClusterID
			}
		} else if !errors.Is(err, fs.ErrNotExist) {
			return verifyFlags{}, err
		}
	}

	// Validate
	if ownerID == "" && clusterID == "" {
		return verifyFlags{}, errors.New("neither owner-id nor unique-id provided to verify the cluster")
	}
	endpoint, err = validateEndpoint(endpoint, constants.CoordinatorPort)
	if err != nil {
		return verifyFlags{}, fmt.Errorf("validating endpoint argument: %w", err)
	}

	return verifyFlags{
		endpoint:   endpoint,
		configPath: configPath,
		ownerID:    ownerID,
		clusterID:  clusterID,
	}, nil
}

type verifyFlags struct {
	endpoint   string
	ownerID    string
	clusterID  string
	configPath string
}

func readIds(fileHandler file.Handler) (clusterIDsFile, error) {
	det := clusterIDsFile{}
	if err := fileHandler.ReadJSON(constants.ClusterIDsFileName, &det); err != nil {
		return clusterIDsFile{}, fmt.Errorf("reading cluster ids: %w", err)
	}
	return det, nil
}

// verifyCompletion handles the completion of CLI arguments. It is frequently called
// while the user types arguments of the command to suggest completion.
func verifyCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		return []string{"gcp", "azure"}, cobra.ShellCompDirectiveNoFileComp
	default:
		return []string{}, cobra.ShellCompDirectiveError
	}
}

type constellationVerifier struct {
	dialer grpcInsecureDialer
}

// Verify retrieves an attestation statement from the Constellation and verifies it using the validator.
func (v *constellationVerifier) Verify(
	ctx context.Context, endpoint string, req *verifyproto.GetAttestationRequest, validator atls.Validator,
) error {
	conn, err := v.dialer.DialInsecure(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("dialing init server: %w", err)
	}
	defer conn.Close()

	client := verifyproto.NewAPIClient(conn)

	resp, err := client.GetAttestation(ctx, req)
	if err != nil {
		return fmt.Errorf("getting attestation: %w", err)
	}

	signedData, err := validator.Validate(resp.Attestation, req.Nonce)
	if err != nil {
		return fmt.Errorf("validating attestation: %w", err)
	}

	if !bytes.Equal(signedData, req.UserData) {
		return errors.New("signed data in attestation does not match provided user data")
	}
	return nil
}

type verifyClient interface {
	Verify(ctx context.Context, endpoint string, req *verifyproto.GetAttestationRequest, validator atls.Validator) error
}

type grpcInsecureDialer interface {
	DialInsecure(ctx context.Context, endpoint string) (conn *grpc.ClientConn, err error)
}

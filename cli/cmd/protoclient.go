package cmd

import (
	"context"

	"github.com/edgelesssys/constellation/cli/proto"
	"github.com/edgelesssys/constellation/coordinator/atls"
	"github.com/edgelesssys/constellation/coordinator/state"
)

type protoClient interface {
	Connect(ip, port string, validators []atls.Validator) error
	Close() error
	GetState(ctx context.Context) (state.State, error)
	Activate(ctx context.Context, userPublicKey, masterSecret []byte, nodeIPs, coordinatorIPs, autoscalingNodeGroups []string, cloudServiceAccountURI string) (proto.ActivationResponseClient, error)
}

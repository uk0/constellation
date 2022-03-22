package aws

import "github.com/edgelesssys/constellation/coordinator/core"

// CloudControllerManager holds the AWS cloud-controller-manager configuration.
type CloudControllerManager struct{}

// Image returns the container image used to provide cloud-controller-manager for the cloud-provider.
func (c CloudControllerManager) Image() string {
	return "us.gcr.io/k8s-artifacts-prod/provider-aws/cloud-controller-manager:v1.22.0-alpha.0"
}

// Path returns the path used by cloud-controller-manager executable within the container image.
func (c CloudControllerManager) Path() string {
	return "/aws-cloud-controller-manager"
}

// Name returns the cloud-provider name as used by k8s cloud-controller-manager (k8s.gcr.io/cloud-controller-manager).
func (c CloudControllerManager) Name() string {
	return "aws"
}

// PrepareInstance is called on every instance before deploying the cloud-controller-manager.
// Allows for cloud-provider specific hooks.
func (c CloudControllerManager) PrepareInstance(instance core.Instance, vpnIP string) error {
	// no specific hook required.
	return nil
}

// Supported is used to determine if cloud controller manager is implemented for this cloud provider.
func (c CloudControllerManager) Supported() bool {
	return false
}

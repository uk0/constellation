package gcp

// Autoscaler holds the GCP cluster-autoscaler configuration.
type Autoscaler struct{}

// Name returns the cloud-provider name as used by k8s cluster-autoscaler.
func (a *Autoscaler) Name() string {
	return "gce"
}

// Supported is used to determine if we support autoscaling for the cloud provider.
func (a *Autoscaler) Supported() bool {
	return true
}

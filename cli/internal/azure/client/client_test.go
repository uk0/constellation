package client

import (
	"testing"

	"github.com/edgelesssys/constellation/internal/cloud/cloudprovider"
	"github.com/edgelesssys/constellation/internal/cloud/cloudtypes"
	"github.com/edgelesssys/constellation/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestSetGetState(t *testing.T) {
	testCases := map[string]struct {
		state   state.ConstellationState
		wantErr bool
	}{
		"valid state": {
			state: state.ConstellationState{
				CloudProvider: cloudprovider.Azure.String(),
				AzureWorkers: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip1",
						PrivateIP: "ip2",
					},
				},
				AzureControlPlane: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip3",
						PrivateIP: "ip4",
					},
				},
				Name:                       "name",
				UID:                        "uid",
				BootstrapperHost:           "bootstrapper-host",
				AzureResourceGroup:         "resource-group",
				AzureLocation:              "location",
				AzureSubscription:          "subscription",
				AzureTenant:                "tenant",
				AzureSubnet:                "azure-subnet",
				AzureNetworkSecurityGroup:  "network-security-group",
				AzureWorkersScaleSet:       "worker-scale-set",
				AzureControlPlanesScaleSet: "controlplane-scale-set",
			},
		},
		"missing workers": {
			state: state.ConstellationState{
				CloudProvider: cloudprovider.Azure.String(),
				AzureControlPlane: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip3",
						PrivateIP: "ip4",
					},
				},
				Name:                       "name",
				UID:                        "uid",
				BootstrapperHost:           "bootstrapper-host",
				AzureResourceGroup:         "resource-group",
				AzureLocation:              "location",
				AzureSubscription:          "subscription",
				AzureTenant:                "tenant",
				AzureSubnet:                "azure-subnet",
				AzureNetworkSecurityGroup:  "network-security-group",
				AzureWorkersScaleSet:       "worker-scale-set",
				AzureControlPlanesScaleSet: "controlplane-scale-set",
			},
			wantErr: true,
		},
		"missing controlplane": {
			state: state.ConstellationState{
				CloudProvider: cloudprovider.Azure.String(),
				AzureWorkers: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip1",
						PrivateIP: "ip2",
					},
				},
				Name:                       "name",
				UID:                        "uid",
				BootstrapperHost:           "bootstrapper-host",
				AzureResourceGroup:         "resource-group",
				AzureLocation:              "location",
				AzureSubscription:          "subscription",
				AzureTenant:                "tenant",
				AzureSubnet:                "azure-subnet",
				AzureNetworkSecurityGroup:  "network-security-group",
				AzureWorkersScaleSet:       "worker-scale-set",
				AzureControlPlanesScaleSet: "controlplane-scale-set",
			},
			wantErr: true,
		},
		"missing name": {
			state: state.ConstellationState{
				CloudProvider: cloudprovider.Azure.String(),
				AzureWorkers: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip1",
						PrivateIP: "ip2",
					},
				},
				AzureControlPlane: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip3",
						PrivateIP: "ip4",
					},
				},
				UID:                        "uid",
				BootstrapperHost:           "bootstrapper-host",
				AzureResourceGroup:         "resource-group",
				AzureLocation:              "location",
				AzureSubscription:          "subscription",
				AzureTenant:                "tenant",
				AzureSubnet:                "azure-subnet",
				AzureNetworkSecurityGroup:  "network-security-group",
				AzureWorkersScaleSet:       "worker-scale-set",
				AzureControlPlanesScaleSet: "controlplane-scale-set",
			},
			wantErr: true,
		},
		"missing uid": {
			state: state.ConstellationState{
				CloudProvider: cloudprovider.Azure.String(),
				AzureWorkers: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip1",
						PrivateIP: "ip2",
					},
				},
				AzureControlPlane: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip3",
						PrivateIP: "ip4",
					},
				},
				Name:                       "name",
				BootstrapperHost:           "bootstrapper-host",
				AzureResourceGroup:         "resource-group",
				AzureLocation:              "location",
				AzureSubscription:          "subscription",
				AzureTenant:                "tenant",
				AzureSubnet:                "azure-subnet",
				AzureNetworkSecurityGroup:  "network-security-group",
				AzureWorkersScaleSet:       "worker-scale-set",
				AzureControlPlanesScaleSet: "controlplane-scale-set",
			},
			wantErr: true,
		},
		"missing bootstrapper host": {
			state: state.ConstellationState{
				CloudProvider: cloudprovider.Azure.String(),
				AzureWorkers: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip1",
						PrivateIP: "ip2",
					},
				},
				AzureControlPlane: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip3",
						PrivateIP: "ip4",
					},
				},
				Name:                       "name",
				UID:                        "uid",
				AzureResourceGroup:         "resource-group",
				AzureLocation:              "location",
				AzureSubscription:          "subscription",
				AzureTenant:                "tenant",
				AzureSubnet:                "azure-subnet",
				AzureNetworkSecurityGroup:  "network-security-group",
				AzureWorkersScaleSet:       "worker-scale-set",
				AzureControlPlanesScaleSet: "controlplane-scale-set",
			},
			wantErr: true,
		},
		"missing resource group": {
			state: state.ConstellationState{
				CloudProvider: cloudprovider.Azure.String(),
				AzureWorkers: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip1",
						PrivateIP: "ip2",
					},
				},
				AzureControlPlane: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip3",
						PrivateIP: "ip4",
					},
				},
				Name:                       "name",
				UID:                        "uid",
				BootstrapperHost:           "bootstrapper-host",
				AzureLocation:              "location",
				AzureSubscription:          "subscription",
				AzureTenant:                "tenant",
				AzureSubnet:                "azure-subnet",
				AzureNetworkSecurityGroup:  "network-security-group",
				AzureWorkersScaleSet:       "worker-scale-set",
				AzureControlPlanesScaleSet: "controlplane-scale-set",
			},
			wantErr: true,
		},
		"missing location": {
			state: state.ConstellationState{
				CloudProvider: cloudprovider.Azure.String(),
				AzureWorkers: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip1",
						PrivateIP: "ip2",
					},
				},
				AzureControlPlane: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip3",
						PrivateIP: "ip4",
					},
				},
				Name:                       "name",
				UID:                        "uid",
				BootstrapperHost:           "bootstrapper-host",
				AzureResourceGroup:         "resource-group",
				AzureSubscription:          "subscription",
				AzureTenant:                "tenant",
				AzureSubnet:                "azure-subnet",
				AzureNetworkSecurityGroup:  "network-security-group",
				AzureWorkersScaleSet:       "worker-scale-set",
				AzureControlPlanesScaleSet: "controlplane-scale-set",
			},
			wantErr: true,
		},
		"missing subscription": {
			state: state.ConstellationState{
				CloudProvider: cloudprovider.Azure.String(),
				AzureWorkers: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip1",
						PrivateIP: "ip2",
					},
				},
				AzureControlPlane: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip3",
						PrivateIP: "ip4",
					},
				},
				Name:                       "name",
				UID:                        "uid",
				BootstrapperHost:           "bootstrapper-host",
				AzureResourceGroup:         "resource-group",
				AzureTenant:                "tenant",
				AzureLocation:              "location",
				AzureSubnet:                "azure-subnet",
				AzureNetworkSecurityGroup:  "network-security-group",
				AzureWorkersScaleSet:       "worker-scale-set",
				AzureControlPlanesScaleSet: "controlplane-scale-set",
			},
			wantErr: true,
		},
		"missing tenant": {
			state: state.ConstellationState{
				CloudProvider: cloudprovider.Azure.String(),
				AzureWorkers: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip1",
						PrivateIP: "ip2",
					},
				},
				AzureControlPlane: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip3",
						PrivateIP: "ip4",
					},
				},
				Name:                       "name",
				UID:                        "uid",
				BootstrapperHost:           "bootstrapper-host",
				AzureResourceGroup:         "resource-group",
				AzureSubscription:          "subscription",
				AzureLocation:              "location",
				AzureSubnet:                "azure-subnet",
				AzureNetworkSecurityGroup:  "network-security-group",
				AzureWorkersScaleSet:       "worker-scale-set",
				AzureControlPlanesScaleSet: "controlplane-scale-set",
			},
			wantErr: true,
		},
		"missing subnet": {
			state: state.ConstellationState{
				CloudProvider: cloudprovider.Azure.String(),
				AzureWorkers: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip1",
						PrivateIP: "ip2",
					},
				},
				AzureControlPlane: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip3",
						PrivateIP: "ip4",
					},
				},
				Name:                       "name",
				UID:                        "uid",
				BootstrapperHost:           "bootstrapper-host",
				AzureResourceGroup:         "resource-group",
				AzureLocation:              "location",
				AzureSubscription:          "subscription",
				AzureTenant:                "tenant",
				AzureNetworkSecurityGroup:  "network-security-group",
				AzureWorkersScaleSet:       "worker-scale-set",
				AzureControlPlanesScaleSet: "controlplane-scale-set",
			},
			wantErr: true,
		},
		"missing network security group": {
			state: state.ConstellationState{
				CloudProvider: cloudprovider.Azure.String(),
				AzureWorkers: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip1",
						PrivateIP: "ip2",
					},
				},
				AzureControlPlane: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip3",
						PrivateIP: "ip4",
					},
				},
				Name:                       "name",
				UID:                        "uid",
				BootstrapperHost:           "bootstrapper-host",
				AzureResourceGroup:         "resource-group",
				AzureLocation:              "location",
				AzureSubscription:          "subscription",
				AzureTenant:                "tenant",
				AzureSubnet:                "azure-subnet",
				AzureWorkersScaleSet:       "worker-scale-set",
				AzureControlPlanesScaleSet: "controlplane-scale-set",
			},
			wantErr: true,
		},
		"missing worker scale set": {
			state: state.ConstellationState{
				CloudProvider: cloudprovider.Azure.String(),
				AzureWorkers: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip1",
						PrivateIP: "ip2",
					},
				},
				AzureControlPlane: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip3",
						PrivateIP: "ip4",
					},
				},
				Name:                       "name",
				UID:                        "uid",
				BootstrapperHost:           "bootstrapper-host",
				AzureResourceGroup:         "resource-group",
				AzureLocation:              "location",
				AzureSubscription:          "subscription",
				AzureTenant:                "tenant",
				AzureSubnet:                "azure-subnet",
				AzureNetworkSecurityGroup:  "network-security-group",
				AzureControlPlanesScaleSet: "controlplane-scale-set",
			},
			wantErr: true,
		},
		"missing controlplane scale set": {
			state: state.ConstellationState{
				CloudProvider: cloudprovider.Azure.String(),
				AzureWorkers: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip1",
						PrivateIP: "ip2",
					},
				},
				AzureControlPlane: cloudtypes.Instances{
					"0": {
						PublicIP:  "ip3",
						PrivateIP: "ip4",
					},
				},
				Name:                      "name",
				UID:                       "uid",
				BootstrapperHost:          "bootstrapper-host",
				AzureResourceGroup:        "resource-group",
				AzureLocation:             "location",
				AzureSubscription:         "subscription",
				AzureTenant:               "tenant",
				AzureSubnet:               "azure-subnet",
				AzureNetworkSecurityGroup: "network-security-group",
				AzureWorkersScaleSet:      "worker-scale-set",
			},
			wantErr: true,
		},
	}

	t.Run("SetState", func(t *testing.T) {
		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				assert := assert.New(t)

				client := Client{}
				if tc.wantErr {
					assert.Error(client.SetState(tc.state))
				} else {
					assert.NoError(client.SetState(tc.state))
					assert.Equal(tc.state.AzureWorkers, client.workers)
					assert.Equal(tc.state.AzureControlPlane, client.controlPlanes)
					assert.Equal(tc.state.Name, client.name)
					assert.Equal(tc.state.UID, client.uid)
					assert.Equal(tc.state.AzureResourceGroup, client.resourceGroup)
					assert.Equal(tc.state.AzureLocation, client.location)
					assert.Equal(tc.state.AzureSubscription, client.subscriptionID)
					assert.Equal(tc.state.AzureTenant, client.tenantID)
					assert.Equal(tc.state.AzureSubnet, client.subnetID)
					assert.Equal(tc.state.AzureNetworkSecurityGroup, client.networkSecurityGroup)
					assert.Equal(tc.state.AzureWorkersScaleSet, client.workerScaleSet)
					assert.Equal(tc.state.AzureControlPlanesScaleSet, client.controlPlaneScaleSet)
				}
			})
		}
	})

	t.Run("GetState", func(t *testing.T) {
		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				assert := assert.New(t)

				client := Client{
					workers:              tc.state.AzureWorkers,
					controlPlanes:        tc.state.AzureControlPlane,
					name:                 tc.state.Name,
					uid:                  tc.state.UID,
					loadBalancerPubIP:    tc.state.BootstrapperHost,
					resourceGroup:        tc.state.AzureResourceGroup,
					location:             tc.state.AzureLocation,
					subscriptionID:       tc.state.AzureSubscription,
					tenantID:             tc.state.AzureTenant,
					subnetID:             tc.state.AzureSubnet,
					networkSecurityGroup: tc.state.AzureNetworkSecurityGroup,
					workerScaleSet:       tc.state.AzureWorkersScaleSet,
					controlPlaneScaleSet: tc.state.AzureControlPlanesScaleSet,
				}
				if tc.wantErr {
					_, err := client.GetState()
					assert.Error(err)
				} else {
					state, err := client.GetState()
					assert.NoError(err)
					assert.Equal(tc.state, state)
				}
			})
		}
	})
}

func TestSetStateCloudProvider(t *testing.T) {
	assert := assert.New(t)

	client := Client{}
	stateMissingCloudProvider := state.ConstellationState{
		AzureWorkers: cloudtypes.Instances{
			"0": {
				PublicIP:  "ip1",
				PrivateIP: "ip2",
			},
		},
		AzureControlPlane: cloudtypes.Instances{
			"0": {
				PublicIP:  "ip3",
				PrivateIP: "ip4",
			},
		},
		Name:                       "name",
		UID:                        "uid",
		AzureResourceGroup:         "resource-group",
		AzureLocation:              "location",
		AzureSubscription:          "subscription",
		AzureSubnet:                "azure-subnet",
		AzureNetworkSecurityGroup:  "network-security-group",
		AzureWorkersScaleSet:       "worker-scale-set",
		AzureControlPlanesScaleSet: "controlplane-scale-set",
	}
	assert.Error(client.SetState(stateMissingCloudProvider))
	stateIncorrectCloudProvider := state.ConstellationState{
		CloudProvider: "incorrect",
		AzureWorkers: cloudtypes.Instances{
			"0": {
				PublicIP:  "ip1",
				PrivateIP: "ip2",
			},
		},
		AzureControlPlane: cloudtypes.Instances{
			"0": {
				PublicIP:  "ip3",
				PrivateIP: "ip4",
			},
		},
		Name:                       "name",
		UID:                        "uid",
		AzureResourceGroup:         "resource-group",
		AzureLocation:              "location",
		AzureSubscription:          "subscription",
		AzureSubnet:                "azure-subnet",
		AzureNetworkSecurityGroup:  "network-security-group",
		AzureWorkersScaleSet:       "worker-scale-set",
		AzureControlPlanesScaleSet: "controlplane-scale-set",
	}
	assert.Error(client.SetState(stateIncorrectCloudProvider))
}

func TestInit(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	client := Client{}
	require.NoError(client.init("location", "name"))
	assert.Equal("location", client.location)
	assert.Equal("name", client.name)
	assert.NotEmpty(client.uid)
}

/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package state

import (
	"testing"

	"github.com/edgelesssys/constellation/v2/internal/constants"
	"github.com/edgelesssys/constellation/v2/internal/file"
	"github.com/siderolabs/talos/pkg/machinery/config/encoder"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

var defaultState = &State{
	Version: "v1",
	Infrastructure: Infrastructure{
		UID:             "123",
		ClusterEndpoint: "test-cluster-endpoint",
		InitSecret:      []byte{0x41},
		APIServerCertSANs: []string{
			"api-server-cert-san-test",
			"api-server-cert-san-test-2",
		},
		Azure: &Azure{
			ResourceGroup:            "test-rg",
			SubscriptionID:           "test-sub",
			NetworkSecurityGroupName: "test-nsg",
			LoadBalancerName:         "test-lb",
			UserAssignedIdentity:     "test-uami",
			AttestationURL:           "test-maaUrl",
		},
		GCP: &GCP{
			ProjectID:  "test-project",
			IPCidrNode: "test-cidr-node",
			IPCidrPod:  "test-cidr-pod",
		},
	},
	ClusterValues: ClusterValues{
		ClusterID:       "test-cluster-id",
		OwnerID:         "test-owner-id",
		MeasurementSalt: []byte{0x41},
	},
}

func TestWriteToFile(t *testing.T) {
	prepareFs := func(existingFiles ...string) file.Handler {
		fs := afero.NewMemMapFs()
		fh := file.NewHandler(fs)
		for _, name := range existingFiles {
			if err := fh.Write(name, []byte{0x41}); err != nil {
				t.Fatalf("failed to create file %s: %v", name, err)
			}
		}
		return fh
	}

	testCases := map[string]struct {
		state   *State
		fh      file.Handler
		wantErr bool
	}{
		"success": {
			state: defaultState,
			fh:    prepareFs(),
		},
		"overwrite": {
			state: defaultState,
			fh:    prepareFs(constants.StateFilename),
		},
		"empty state": {
			state: &State{},
			fh:    prepareFs(),
		},
		"rofs": {
			state:   defaultState,
			fh:      file.NewHandler(afero.NewReadOnlyFs(afero.NewMemMapFs())),
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			err := tc.state.WriteToFile(tc.fh, constants.StateFilename)

			if tc.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(mustMarshalYaml(t, tc.state), mustReadFromFile(t, tc.fh))
			}
		})
	}
}

func TestReadFromFile(t *testing.T) {
	prepareFs := func(existingFiles map[string][]byte) file.Handler {
		fs := afero.NewMemMapFs()
		fh := file.NewHandler(fs)
		for name, content := range existingFiles {
			if err := fh.Write(name, content); err != nil {
				t.Fatalf("failed to create file %s: %v", name, err)
			}
		}
		return fh
	}

	testCases := map[string]struct {
		existingFiles map[string][]byte
		wantErr       bool
	}{
		"success": {
			existingFiles: map[string][]byte{
				constants.StateFilename: mustMarshalYaml(t, defaultState),
			},
		},
		"no state file present": {
			existingFiles: map[string][]byte{},
			wantErr:       true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			fh := prepareFs(tc.existingFiles)

			state, err := ReadFromFile(fh, constants.StateFilename)

			if tc.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(tc.existingFiles[constants.StateFilename], mustMarshalYaml(t, state))
			}
		})
	}
}

func mustMarshalYaml(t *testing.T, v any) []byte {
	t.Helper()
	b, err := encoder.NewEncoder(v).Encode()
	if err != nil {
		t.Fatalf("failed to marshal yaml: %v", err)
	}
	return b
}

func mustReadFromFile(t *testing.T, fh file.Handler) []byte {
	t.Helper()
	b, err := fh.Read(constants.StateFilename)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	return b
}

func TestMerge(t *testing.T) {
	testCases := map[string]struct {
		state    *State
		other    *State
		expected *State
		wantErr  bool
	}{
		"success": {
			state: &State{
				Infrastructure: Infrastructure{
					ClusterEndpoint: "test-cluster-endpoint",
					UID:             "123",
				},
			},
			other: &State{
				Version: "v1",
				Infrastructure: Infrastructure{
					UID: "456",
				},
				ClusterValues: ClusterValues{
					ClusterID: "test-cluster-id",
				},
			},
			expected: &State{
				Version: "v1",
				Infrastructure: Infrastructure{
					ClusterEndpoint: "test-cluster-endpoint",
					UID:             "456",
				},
				ClusterValues: ClusterValues{
					ClusterID: "test-cluster-id",
				},
			},
		},
		"empty state": {
			state: &State{},
			other: &State{
				Version: "v1",
				Infrastructure: Infrastructure{
					UID: "456",
				},
				ClusterValues: ClusterValues{
					ClusterID: "test-cluster-id",
				},
			},
			expected: &State{
				Version: "v1",
				Infrastructure: Infrastructure{
					UID: "456",
				},
				ClusterValues: ClusterValues{
					ClusterID: "test-cluster-id",
				},
			},
		},
		"empty other": {
			state: &State{
				Version: "v1",
				Infrastructure: Infrastructure{
					UID: "456",
				},
				ClusterValues: ClusterValues{
					ClusterID: "test-cluster-id",
				},
			},
			other: &State{},
			expected: &State{
				Version: "v1",
				Infrastructure: Infrastructure{
					UID: "456",
				},
				ClusterValues: ClusterValues{
					ClusterID: "test-cluster-id",
				},
			},
		},
		"empty state and other": {
			state:    &State{},
			other:    &State{},
			expected: &State{},
		},
		"identical": {
			state: &State{
				Version: "v1",
				Infrastructure: Infrastructure{
					UID: "456",
				},
				ClusterValues: ClusterValues{
					ClusterID: "test-cluster-id",
				},
			},
			other: &State{
				Version: "v1",
				Infrastructure: Infrastructure{
					UID: "456",
				},
				ClusterValues: ClusterValues{
					ClusterID: "test-cluster-id",
				},
			},
			expected: &State{
				Version: "v1",
				Infrastructure: Infrastructure{
					UID: "456",
				},
				ClusterValues: ClusterValues{
					ClusterID: "test-cluster-id",
				},
			},
		},
		"nested pointer": {
			state: &State{
				Version: "v1",
				Infrastructure: Infrastructure{
					UID: "123",
					Azure: &Azure{
						AttestationURL: "test-maaUrl",
					},
				},
				ClusterValues: ClusterValues{
					ClusterID: "test-cluster-id",
				},
			},
			other: &State{
				Version: "v1",
				Infrastructure: Infrastructure{
					UID: "456",
					Azure: &Azure{
						AttestationURL: "test-maaUrl-2",
					},
				},
				ClusterValues: ClusterValues{
					ClusterID: "test-cluster-id",
				},
			},
			expected: &State{
				Version: "v1",
				Infrastructure: Infrastructure{
					UID: "456",
					Azure: &Azure{
						AttestationURL: "test-maaUrl-2",
					},
				},
				ClusterValues: ClusterValues{
					ClusterID: "test-cluster-id",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			_, err := tc.state.Merge(tc.other)

			if tc.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(tc.expected, tc.state)
			}
		})
	}
}

func TestMarshalHexBytes(t *testing.T) {
	testCases := map[string]struct {
		in       HexBytes
		expected string
		wantErr  bool
	}{
		"success": {
			in:       []byte{0xab, 0xcd, 0xef},
			expected: "abcdef\n",
		},
		"empty": {
			in:       []byte{},
			expected: "\"\"\n",
		},
		"nil": {
			in:       nil,
			expected: "\"\"\n",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			actual, err := yaml.Marshal(tc.in)

			if tc.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(tc.expected, string(actual))
			}
		})
	}
}

func TestUnmarshalHexBytes(t *testing.T) {
	testCases := map[string]struct {
		in       string
		expected HexBytes
		wantErr  bool
	}{
		"success": {
			in:       "abcdef",
			expected: []byte{0xab, 0xcd, 0xef},
		},
		"empty": {
			in:       "",
			expected: nil,
		},
		"byte slice compat": {
			in:       "[0xab, 0xcd, 0xef]",
			expected: []byte{0xab, 0xcd, 0xef},
		},
		"byte slice compat 2": {
			in:       "[00, 12, 34]",
			expected: []byte{0x00, 0x0c, 0x22},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			var actual HexBytes
			err := yaml.Unmarshal([]byte(tc.in), &actual)

			if tc.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(tc.expected, actual)
			}
		})
	}
}

func TestMarshalUnmarshalHexBytes(t *testing.T) {
	in := HexBytes{0xab, 0xcd, 0xef}
	expected := "abcdef\n"

	actual, err := yaml.Marshal(in)
	require.NoError(t, err)
	assert.Equal(t, expected, string(actual))

	var actual2 HexBytes
	err = yaml.Unmarshal(actual, &actual2)
	require.NoError(t, err)
	assert.Equal(t, in, actual2)
}
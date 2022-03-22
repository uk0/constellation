package cmd

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGCPNonCVMValidator(t *testing.T) {
	testCases := map[string]struct {
		ownerID     string
		clusterID   string
		errExpected bool
	}{
		"no input": {
			ownerID:     "",
			clusterID:   "",
			errExpected: true,
		},
		"unencoded secret ID": {
			ownerID:     "owner-id",
			clusterID:   base64.StdEncoding.EncodeToString([]byte("unique-id")),
			errExpected: true,
		},
		"unencoded cluster ID": {
			ownerID:     base64.StdEncoding.EncodeToString([]byte("owner-id")),
			clusterID:   "unique-id",
			errExpected: true,
		},
		"correct input": {
			ownerID:     base64.StdEncoding.EncodeToString([]byte("owner-id")),
			clusterID:   base64.StdEncoding.EncodeToString([]byte("unique-id")),
			errExpected: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			cmd := newVerifyGCPNonCVMCmd()
			cmd.Flags().String("owner-id", "", "")
			cmd.Flags().String("unique-id", "", "")
			require.NoError(cmd.Flags().Set("owner-id", tc.ownerID))
			require.NoError(cmd.Flags().Set("unique-id", tc.clusterID))
			var out bytes.Buffer
			cmd.SetOut(&out)
			var errOut bytes.Buffer
			cmd.SetErr(&errOut)

			_, err := getGCPNonCVMValidator(cmd, map[uint32][]byte{})
			if tc.errExpected {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

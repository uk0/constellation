package kms

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStore(t *testing.T) {
	testCases := map[string]struct {
		uri         string
		errExpected bool
	}{
		"no store": {
			uri:         NoStoreURI,
			errExpected: false,
		},
		"aws s3": {
			uri:         fmt.Sprintf(AWSS3URI, ""),
			errExpected: true,
		},
		"azure blob": {
			uri:         fmt.Sprintf(AzureBlobURI, "", ""),
			errExpected: true,
		},
		"gcp storage": {
			uri:         fmt.Sprintf(GCPStorageURI, "", ""),
			errExpected: true,
		},
		"unknown store": {
			uri:         "storage://unknown",
			errExpected: true,
		},
		"invalid scheme": {
			uri:         ClusterKMSURI,
			errExpected: true,
		},
		"not a url": {
			uri:         ":/123",
			errExpected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			_, err := getStore(context.Background(), tc.uri)
			if tc.errExpected {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestGetKMS(t *testing.T) {
	testCases := map[string]struct {
		uri         string
		errExpected bool
	}{
		"cluster kms": {
			uri:         ClusterKMSURI,
			errExpected: false,
		},
		"aws kms": {
			uri:         fmt.Sprintf(AWSKMSURI, ""),
			errExpected: true,
		},
		"azure kms": {
			uri:         fmt.Sprintf(AzureKMSURI, "", ""),
			errExpected: true,
		},
		"azure hsm": {
			uri:         fmt.Sprintf(AzureHSMURI, ""),
			errExpected: true,
		},
		"gcp kms": {
			uri:         fmt.Sprintf(GCPKMSURI, "", "", "", ""),
			errExpected: true,
		},
		"unknown kms": {
			uri:         "kms://unknown",
			errExpected: true,
		},
		"invalid scheme": {
			uri:         NoStoreURI,
			errExpected: true,
		},
		"not a url": {
			uri:         ":/123",
			errExpected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			kms, err := getKMS(context.Background(), tc.uri, nil)
			if tc.errExpected {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.NotNil(kms)
			}
		})
	}
}

func TestSetUpKMS(t *testing.T) {
	assert := assert.New(t)

	kms, err := SetUpKMS(context.TODO(), "storage://unknown", "kms://unknown")
	assert.Error(err)
	assert.Nil(kms)

	kms, err = SetUpKMS(context.Background(), "storage://no-store", "kms://cluster-kms")
	assert.NoError(err)
	assert.NotNil(kms)
}

func TestGetAWSKMSConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	policy := "{keyPolicy: keyPolicy}"
	escapedPolicy := url.QueryEscape(policy)
	uri, err := url.Parse(fmt.Sprintf(AWSKMSURI, escapedPolicy))
	require.NoError(err)
	policyProducer, err := getAWSKMSConfig(uri)
	require.NoError(err)
	keyPolicy, err := policyProducer.CreateKeyPolicy("")
	require.NoError(err)
	assert.Equal(policy, keyPolicy)
}

func TestGetAzureBlobConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	connStr := "DefaultEndpointsProtocol=https;AccountName=test;AccountKey=Q29uc3RlbGxhdGlvbg==;EndpointSuffix=core.windows.net"
	escapedConnStr := url.QueryEscape(connStr)
	container := "test"
	uri, err := url.Parse(fmt.Sprintf(AzureBlobURI, container, escapedConnStr))
	require.NoError(err)
	rContainer, rConnStr, err := getAzureBlobConfig(uri)
	require.NoError(err)
	assert.Equal(container, rContainer)
	assert.Equal(connStr, rConnStr)
}

func TestGetGCPKMSConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	project := "test-project"
	location := "global"
	keyRing := "test-ring"
	protectionLvl := "2"
	uri, err := url.Parse(fmt.Sprintf(GCPKMSURI, project, location, keyRing, protectionLvl))
	require.NoError(err)
	rProject, rLocation, rKeyRing, rProtectionLvl, err := getGCPKMSConfig(uri)
	require.NoError(err)
	assert.Equal(project, rProject)
	assert.Equal(location, rLocation)
	assert.Equal(keyRing, rKeyRing)
	assert.Equal(2, rProtectionLvl)

	uri, err = url.Parse(fmt.Sprintf(GCPKMSURI, project, location, keyRing, "invalid"))
	require.NoError(err)
	_, _, _, _, err = getGCPKMSConfig(uri)
	assert.Error(err)
}

func TestGetConfig(t *testing.T) {
	const testUri = "test://config?name=test-name&data=test-data&value=test-value"

	testCases := map[string]struct {
		uri         string
		keys        []string
		errExpected bool
	}{
		"success": {
			uri:         testUri,
			keys:        []string{"name", "data", "value"},
			errExpected: false,
		},
		"less keys than capture groups": {
			uri:         testUri,
			keys:        []string{"name", "data"},
			errExpected: false,
		},
		"invalid regex": {
			uri:         testUri,
			keys:        []string{"name", "data", "test-value"},
			errExpected: true,
		},
		"missing value": {
			uri:         "test://config?name=test-name&data=test-data&value",
			keys:        []string{"name", "data", "value"},
			errExpected: true,
		},
		"more keys than expected": {
			uri:         testUri,
			keys:        []string{"name", "data", "value", "anotherValue"},
			errExpected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			uri, err := url.Parse(tc.uri)
			require.NoError(err)

			res, err := getConfig(uri.Query(), tc.keys)
			if tc.errExpected {
				assert.Error(err)
				assert.Len(res, len(tc.keys))
			} else {
				assert.NoError(err)
				require.Len(res, len(tc.keys))
				for i := range tc.keys {
					assert.NotEmpty(res[i])
				}
			}
		})
	}
}

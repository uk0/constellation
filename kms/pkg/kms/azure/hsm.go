package azure

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/crypto"
	"github.com/edgelesssys/constellation/kms/pkg/config"
	"github.com/edgelesssys/constellation/kms/pkg/kms"
	"github.com/edgelesssys/constellation/kms/pkg/kms/util"
	"github.com/edgelesssys/constellation/kms/pkg/storage"
)

// Suffix for HSM Vaults.
const HSMDefaultCloud VaultSuffix = ".managedhsm.azure.net/"

// HSMClient implements the CloudKMS interface for Azure managed HSM.
type HSMClient struct {
	credentials azcore.TokenCredential
	client      *azkeys.Client
	storage     kms.Storage
	vaultURL    string
	opts        *Opts
}

// NewHSM initializes a KMS client for Azure manged HSM Key Vault.
func NewHSM(ctx context.Context, vaultName string, store kms.Storage, opts *Opts) (*HSMClient, error) {
	if opts == nil {
		opts = &Opts{}
	}
	cred, err := azidentity.NewDefaultAzureCredential(opts.credentials)
	if err != nil {
		return nil, fmt.Errorf("loading credentials: %w", err)
	}

	vaultURL := vaultPrefix + vaultName + string(HSMDefaultCloud)
	client, err := azkeys.NewClient(vaultURL, cred, (*azkeys.ClientOptions)(opts.client))
	if err != nil {
		return nil, fmt.Errorf("creating HSM client: %w", err)
	}

	// `azkeys.NewClient()` does not error if the vault is non existant
	// Test here if we can reach the vault, and error otherwise
	pager := client.ListKeys(&azkeys.ListKeysOptions{MaxResults: to.Int32Ptr(2)})
	pager.NextPage(ctx)
	if pager.Err() != nil {
		return nil, fmt.Errorf("HSM not reachable: %w", err)
	}

	if store == nil {
		store = storage.NewMemMapStorage()
	}

	return &HSMClient{vaultURL: vaultURL, client: client, credentials: cred, storage: store, opts: opts}, nil
}

// CreateKEK creates a new Key Encryption Key using Azure managed HSM.
//
// If no key material is provided, a new key is generated by the HSM, otherwise the key material is used to import the key.
func (c *HSMClient) CreateKEK(ctx context.Context, keyID string, key []byte) error {
	if len(key) == 0 {
		if _, err := c.client.CreateOCTKey(ctx, keyID, &azkeys.CreateOCTKeyOptions{
			HardwareProtected: true,
			KeySize:           to.Int32Ptr(config.SymmetricKeyLength * 8),
			Tags:              config.KmsTags,
		}); err != nil {
			return fmt.Errorf("creating new KEK: %w", err)
		}
		return nil
	}

	jwk := azkeys.JSONWebKey{
		K: key,
		KeyOps: []*string{
			to.StringPtr("wrapKey"),
			to.StringPtr("unwrapKey"),
		},
		KeyType: (*azkeys.KeyType)(to.StringPtr(string(azkeys.OctHSM))),
	}
	importOpts := &azkeys.ImportKeyOptions{
		Hsm: to.BoolPtr(true),
		KeyAttributes: &azkeys.KeyAttributes{
			Attributes: azkeys.Attributes{
				Enabled: to.BoolPtr(true),
			},
		},
		Tags: config.KmsTags,
	}

	if _, err := c.client.ImportKey(ctx, keyID, jwk, importOpts); err != nil {
		return fmt.Errorf("importing KEK to Azure HSM: %w", err)
	}
	return nil
}

// GetDEK loads an encrypted DEK from storage and unwraps it using an HSM-backed key.
func (c *HSMClient) GetDEK(ctx context.Context, kekID string, keyID string, dekSize int) ([]byte, error) {
	encryptedDEK, err := c.storage.Get(ctx, keyID)
	if err != nil {
		if !errors.Is(err, storage.ErrDEKUnset) {
			return nil, fmt.Errorf("loading encrypted DEK from storage: %w", err)
		}

		// If the DEK does not exist we generate a new random DEK and save it to storage
		newDEK, err := util.GetRandomKey(dekSize)
		if err != nil {
			return nil, fmt.Errorf("could not generate key: %w", err)
		}
		if err := c.putDEK(ctx, kekID, keyID, newDEK); err != nil {
			return nil, fmt.Errorf("creating new DEK: %w", err)
		}

		return newDEK, nil
	}

	version, err := c.getKeyVersion(ctx, kekID)
	if err != nil {
		return nil, fmt.Errorf("unable to detect key version: %w", err)
	}

	cryptoClient, err := crypto.NewClient(fmt.Sprintf("%skeys/%s/%s", c.vaultURL, kekID, version), c.credentials, (*crypto.ClientOptions)(c.opts.client))
	if err != nil {
		return nil, fmt.Errorf("creating crypto client for KEK: %s: %w", kekID, err)
	}

	res, err := cryptoClient.UnwrapKey(ctx, crypto.AES256, encryptedDEK, nil)
	if err != nil {
		return nil, fmt.Errorf("unwrapping key: %w", err)
	}

	return res.Result, nil
}

// putDEK wraps a key using an HSM-backed key and saves it to storage.
func (c *HSMClient) putDEK(ctx context.Context, kekID, keyID string, plainDEK []byte) error {
	version, err := c.getKeyVersion(ctx, kekID)
	if err != nil {
		return fmt.Errorf("unable to detect key version: %w", err)
	}
	cryptoClient, err := crypto.NewClient(fmt.Sprintf("%skeys/%s/%s", c.vaultURL, kekID, version), c.credentials, nil)
	if err != nil {
		return fmt.Errorf("creating crypto client for KEK: %s: %w", kekID, err)
	}

	res, err := cryptoClient.WrapKey(ctx, crypto.AES256, plainDEK, &crypto.WrapKeyOptions{})
	if err != nil {
		return fmt.Errorf("wrapping key: %w", err)
	}

	return c.storage.Put(ctx, keyID, res.Result)
}

// getKeyVersion detects the latests version number of a given key.
func (c *HSMClient) getKeyVersion(ctx context.Context, kekID string) (string, error) {
	kek, err := c.client.GetKey(ctx, kekID, &azkeys.GetKeyOptions{})
	if err != nil {
		return "", err
	}

	parsed, err := url.Parse(*kek.Key.ID)
	if err != nil {
		return "", err
	}

	path := strings.Split(strings.TrimPrefix(parsed.Path, "/keys/"), "/")
	if len(path) != 2 {
		return "", fmt.Errorf("invalid key ID URL: %s", *kek.Key.ID)
	}

	return path[1], nil
}

/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package azure

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
	"github.com/edgelesssys/constellation/v2/kms/internal/config"
	"github.com/edgelesssys/constellation/v2/kms/internal/storage"
	"github.com/edgelesssys/constellation/v2/kms/kms"
	"github.com/edgelesssys/constellation/v2/kms/kms/util"
)

type hsmClientAPI interface {
	CreateKey(ctx context.Context, name string, parameters azkeys.CreateKeyParameters, options *azkeys.CreateKeyOptions) (azkeys.CreateKeyResponse, error)
	ImportKey(ctx context.Context, name string, parameters azkeys.ImportKeyParameters, options *azkeys.ImportKeyOptions) (azkeys.ImportKeyResponse, error)
	GetKey(ctx context.Context, name string, version string, options *azkeys.GetKeyOptions) (azkeys.GetKeyResponse, error)
	UnwrapKey(ctx context.Context, name string, version string, parameters azkeys.KeyOperationsParameters, options *azkeys.UnwrapKeyOptions) (azkeys.UnwrapKeyResponse, error)
	WrapKey(ctx context.Context, name string, version string, parameters azkeys.KeyOperationsParameters, options *azkeys.WrapKeyOptions) (azkeys.WrapKeyResponse, error)
}

// HSMDefaultCloud is the suffix for HSM Vaults.
const HSMDefaultCloud VaultSuffix = ".managedhsm.azure.net/"

// HSMClient implements the CloudKMS interface for Azure managed HSM.
type HSMClient struct {
	credentials azcore.TokenCredential
	client      hsmClientAPI
	storage     kms.Storage
	vaultURL    string
}

// NewHSM initializes a KMS client for Azure manged HSM Key Vault.
func NewHSM(ctx context.Context, vaultName string, store kms.Storage, opts *Opts) (*HSMClient, error) {
	if opts == nil {
		opts = &Opts{}
	}
	cred, err := azidentity.NewDefaultAzureCredential(opts.Credentials)
	if err != nil {
		return nil, fmt.Errorf("loading credentials: %w", err)
	}

	vaultURL := vaultPrefix + vaultName + string(HSMDefaultCloud)
	client, err := azkeys.NewClient(vaultURL, cred, opts.Keys)
	if err != nil {
		return nil, fmt.Errorf("creating azure key vault client: %w", err)
	}

	// `azkeys.NewClient()` does not error if the vault is non existent
	// Test here if we can reach the vault, and error otherwise
	pager := client.NewListKeysPager(&azkeys.ListKeysOptions{MaxResults: to.Ptr[int32](2)})
	if _, err := pager.NextPage(ctx); err != nil {
		return nil, fmt.Errorf("HSM not reachable: %w", err)
	}

	if store == nil {
		store = storage.NewMemMapStorage()
	}

	return &HSMClient{
		vaultURL:    vaultURL,
		client:      client,
		credentials: cred,
		storage:     store,
	}, nil
}

// CreateKEK creates a new Key Encryption Key using Azure managed HSM.
//
// If no key material is provided, a new key is generated by the HSM, otherwise the key material is used to import the key.
func (c *HSMClient) CreateKEK(ctx context.Context, keyID string, key []byte) error {
	if len(key) == 0 {
		if _, err := c.client.CreateKey(ctx, keyID, azkeys.CreateKeyParameters{
			Kty:     to.Ptr(azkeys.JSONWebKeyTypeOctHSM),
			KeySize: to.Ptr[int32](config.SymmetricKeyLength * 8),
			Tags:    toAzureTags(config.KmsTags),
		}, &azkeys.CreateKeyOptions{}); err != nil {
			return fmt.Errorf("creating new KEK: %w", err)
		}
		return nil
	}

	jwk := azkeys.JSONWebKey{
		K: key,
		KeyOps: []*string{
			to.Ptr("wrapKey"),
			to.Ptr("unwrapKey"),
		},
		Kty: to.Ptr(azkeys.JSONWebKeyTypeOctHSM),
	}
	importParams := azkeys.ImportKeyParameters{
		HSM: to.Ptr(true),
		KeyAttributes: &azkeys.KeyAttributes{
			Enabled: to.Ptr(true),
		},
		Tags: toAzureTags(config.KmsTags),
		Key:  &jwk,
	}

	if _, err := c.client.ImportKey(ctx, keyID, importParams, &azkeys.ImportKeyOptions{}); err != nil {
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
			return nil, fmt.Errorf("key generation: %w", err)
		}
		if err := c.putDEK(ctx, kekID, keyID, newDEK); err != nil {
			return nil, fmt.Errorf("creating new DEK: %w", err)
		}

		return newDEK, nil
	}

	params := azkeys.KeyOperationsParameters{
		Algorithm: to.Ptr(azkeys.JSONWebKeyEncryptionAlgorithmA256KW),
		Value:     encryptedDEK,
	}
	res, err := c.client.UnwrapKey(ctx, kekID, "", params, &azkeys.UnwrapKeyOptions{})
	if err != nil {
		return nil, fmt.Errorf("unwrapping key: %w", err)
	}

	return res.Result, nil
}

// putDEK wraps a key using an HSM-backed key and saves it to storage.
func (c *HSMClient) putDEK(ctx context.Context, kekID, keyID string, plainDEK []byte) error {
	params := azkeys.KeyOperationsParameters{
		Algorithm: to.Ptr(azkeys.JSONWebKeyEncryptionAlgorithmA256KW),
		Value:     plainDEK,
	}
	res, err := c.client.WrapKey(ctx, kekID, "", params, &azkeys.WrapKeyOptions{})
	if err != nil {
		return fmt.Errorf("wrapping key: %w", err)
	}

	return c.storage.Put(ctx, keyID, res.Result)
}

// toAzureTags converts a map of tags to map of tag pointers.
func toAzureTags(tags map[string]string) map[string]*string {
	tagsOut := make(map[string]*string)
	for k, v := range tags {
		tagsOut[k] = to.Ptr(v)
	}
	return tagsOut
}

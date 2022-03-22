package aws

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/edgelesssys/constellation/kms/pkg/config"
	kmsInterface "github.com/edgelesssys/constellation/kms/pkg/kms"
	"github.com/edgelesssys/constellation/kms/pkg/kms/util"
	"github.com/edgelesssys/constellation/kms/pkg/storage"
)

const (
	// DEKContext is used as the encryption context in AWS KMS.
	DEKContext = "aws:ebs:id"
)

// ClientAPI satisfies the Amazons KMS client's methods we need.
// This allows us to mock the actual client, see https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/
type ClientAPI interface {
	CreateAlias(ctx context.Context, params *kms.CreateAliasInput, optFns ...func(*kms.Options)) (*kms.CreateAliasOutput, error)
	CreateKey(ctx context.Context, params *kms.CreateKeyInput, optFns ...func(*kms.Options)) (*kms.CreateKeyOutput, error)
	Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error)
	DeleteAlias(ctx context.Context, params *kms.DeleteAliasInput, optFns ...func(*kms.Options)) (*kms.DeleteAliasOutput, error)
	DescribeKey(ctx context.Context, params *kms.DescribeKeyInput, optFns ...func(*kms.Options)) (*kms.DescribeKeyOutput, error)
	Encrypt(ctx context.Context, params *kms.EncryptInput, optFns ...func(*kms.Options)) (*kms.EncryptOutput, error)
	GenerateDataKey(ctx context.Context, params *kms.GenerateDataKeyInput, optFns ...func(*kms.Options)) (*kms.GenerateDataKeyOutput, error)
	GenerateDataKeyWithoutPlaintext(ctx context.Context, params *kms.GenerateDataKeyWithoutPlaintextInput, optFns ...func(*kms.Options)) (*kms.GenerateDataKeyWithoutPlaintextOutput, error)
	GetParametersForImport(ctx context.Context, params *kms.GetParametersForImportInput, optFns ...func(*kms.Options)) (*kms.GetParametersForImportOutput, error)
	ImportKeyMaterial(ctx context.Context, params *kms.ImportKeyMaterialInput, optFns ...func(*kms.Options)) (*kms.ImportKeyMaterialOutput, error)
	PutKeyPolicy(ctx context.Context, params *kms.PutKeyPolicyInput, optFns ...func(*kms.Options)) (*kms.PutKeyPolicyOutput, error)
	ScheduleKeyDeletion(ctx context.Context, params *kms.ScheduleKeyDeletionInput, optFns ...func(*kms.Options)) (*kms.ScheduleKeyDeletionOutput, error)
}

// KeyPolicyProducer allows to have callbacks for generating key policies at runtime.
type KeyPolicyProducer interface {
	// CreateKeyPolicy returns a key policy for a given key ID.
	CreateKeyPolicy(keyID string) (string, error)
}

// KMSClient implements the CloudKMS interface for AWS.
type KMSClient struct {
	awsClient      ClientAPI
	policyProducer KeyPolicyProducer
	storage        kmsInterface.Storage
}

// New creates and initializes a new KMSClient for AWS.
//
// The parameter client needs to be initialized with valid AWS credentials (https://aws.github.io/aws-sdk-go-v2/docs/getting-started).
// If storage is nil, the default MemMapStorage is used.
func New(ctx context.Context, policyProducer KeyPolicyProducer, store kmsInterface.Storage, optFns ...func(*awsconfig.LoadOptions) error) (*KMSClient, error) {
	if store == nil {
		store = storage.NewMemMapStorage()
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return nil, err
	}
	client := kms.NewFromConfig(cfg)

	return &KMSClient{
		awsClient:      client,
		policyProducer: policyProducer,
		storage:        store,
	}, nil
}

// CreateKEK creates a new KEK with the given key material and policy. If successful, the key can be referenced by keyID in the KMS in accordance to the policy.
// https://docs.aws.amazon.com/kms/latest/developerguide/importing-keys.html
func (c *KMSClient) CreateKEK(ctx context.Context, keyID string, key []byte) error {
	alias := "alias/" + keyID

	// Check whether key with keyID already exists
	describeKeyInput := &kms.DescribeKeyInput{
		KeyId: aws.String(alias),
	}
	// If the keyID parameter is used for a key in the AWS KMS, the response includes the keyID generated by AWS at creation
	var awsGeneratedKeyID string
	newKeyCreationNeeded := true
	var nfe *types.NotFoundException
	describeKeyOutput, err := c.awsClient.DescribeKey(ctx, describeKeyInput)
	if err == nil {
		// The request is valid and a key with the keyID exists
		awsGeneratedKeyID = *describeKeyOutput.KeyMetadata.KeyId
		newKeyCreationNeeded = false
	} else if !errors.As(err, &nfe) {
		return err
	}

	// If it is not needed to create a new key, the steps to create the key and the alias can be skipped
	if newKeyCreationNeeded {
		// specifies that the key should be empty at creation and the material will be imported afterward
		origin := types.OriginTypeExternal
		if len(key) == 0 {
			origin = types.OriginTypeAwsKms
		}
		// Creates new AWS KMS key with empty key material
		var tags []types.Tag
		for tagKey, tagValue := range config.KmsTags {
			tags = append(tags, types.Tag{
				TagKey:   aws.String(tagKey),
				TagValue: aws.String(tagValue),
			})
		}
		createKeyInput := &kms.CreateKeyInput{
			Description: aws.String("Constellation Key Encryption Key"),
			Origin:      origin,
			Tags:        tags,
		}
		kekMetadata, err := c.awsClient.CreateKey(ctx, createKeyInput)
		if err != nil {
			return err
		}
		// Use the keyId of the created key in the following
		awsGeneratedKeyID = *kekMetadata.KeyMetadata.KeyId

		// Creates Alias for the KEK, so the key can be accessed by specifying the keyID
		createAliasInput := &kms.CreateAliasInput{
			AliasName:   aws.String(alias),
			TargetKeyId: &awsGeneratedKeyID,
		}
		if _, err = c.awsClient.CreateAlias(ctx, createAliasInput); err != nil {
			c.tryCleanUpResources(ctx, newKeyCreationNeeded, awsGeneratedKeyID, alias)
			return err
		}
	}

	// Only import key if the key is not empty
	if len(key) != 0 {
		// Retrieves token and public AWS key to encrypt the KEK for transmitting to AWS KMS
		getImportParameterInput := &kms.GetParametersForImportInput{
			KeyId: &awsGeneratedKeyID,
			// if supported, it is recommended to use 'RSAES_OAEP_SHA_256': https://docs.aws.amazon.com/kms/latest/developerguide/importing-keys-get-public-key-and-token.html
			WrappingAlgorithm: types.AlgorithmSpecRsaesOaepSha256,
			WrappingKeySpec:   types.WrappingKeySpecRsa2048,
		}
		getParametersForImportOutput, err := c.awsClient.GetParametersForImport(ctx, getImportParameterInput)
		if err != nil {
			c.tryCleanUpResources(ctx, newKeyCreationNeeded, awsGeneratedKeyID, alias)
			return err
		}

		// Encrypt the private key with the public key provided by AWS KMS
		// From the AWS KMS get-public-key documentation:
		// The value is a DER-encoded X.509 public key, also known as SubjectPublicKeyInfo (SPKI), as defined in RFC 5280.
		// When you use the HTTP API or the Amazon Web Services CLI, the value is Base64-encoded. Otherwise, it is not Base64-encoded.
		publicKey, err := util.ParseDERtoPublicKeyRSA(getParametersForImportOutput.PublicKey)
		if err != nil {
			c.tryCleanUpResources(ctx, newKeyCreationNeeded, awsGeneratedKeyID, alias)
			return err
		}
		encryptedKEK, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, key, nil)
		if err != nil {
			c.tryCleanUpResources(ctx, newKeyCreationNeeded, awsGeneratedKeyID, alias)
			return err
		}

		// Pushes the key material for the created KEK to AWS KMS
		// In case the key has already key material, the importKeyMaterial operation only succeeds if the newly imported key material matches the previous imported one
		// Otherwise it responds with a IncorrectKeyMaterialException
		importKeyMaterialInput := &kms.ImportKeyMaterialInput{
			EncryptedKeyMaterial: encryptedKEK,
			ImportToken:          getParametersForImportOutput.ImportToken,
			KeyId:                &awsGeneratedKeyID,
			ExpirationModel:      types.ExpirationModelTypeKeyMaterialDoesNotExpire,
		}
		if _, err = c.awsClient.ImportKeyMaterial(ctx, importKeyMaterialInput); err != nil {
			c.tryCleanUpResources(ctx, newKeyCreationNeeded, awsGeneratedKeyID, alias)
			return err
		}
	}

	// Pushes key policy for KEK
	// Since the default policy of the KEK does not allow decryption for an IAM role, one has to include that in the key policy when importing the KEK.
	// Decryption is needed for retrieving the DEKs from the storage.
	policy, err := c.policyProducer.CreateKeyPolicy(awsGeneratedKeyID)
	if err != nil {
		c.tryCleanUpResources(ctx, newKeyCreationNeeded, awsGeneratedKeyID, alias)
		return err
	}
	putKeyPolicyInput := &kms.PutKeyPolicyInput{
		KeyId:      &awsGeneratedKeyID,
		Policy:     &policy,
		PolicyName: aws.String("default"),
	}
	if _, err = c.awsClient.PutKeyPolicy(ctx, putKeyPolicyInput); err != nil {
		c.tryCleanUpResources(ctx, newKeyCreationNeeded, awsGeneratedKeyID, alias)
		return err
	}
	return nil
}

// GetDEK returns the DEK for dekID and kekID from the KMS.
func (c *KMSClient) GetDEK(ctx context.Context, kekID, keyID string, dekSize int) ([]byte, error) {
	// The KEK should be identified by its alias. The alias always has the same scheme: 'alias/<kekId>'
	kekID = "alias/" + kekID

	// If a key for keyID exists in the storage, decrypt the key using the KEK.
	dek, err := c.decryptDEKFromStorage(ctx, kekID, keyID)
	if err == nil {
		return dek, nil
	}
	if !errors.Is(err, storage.ErrDEKUnset) {
		return nil, err
	}
	return c.putNewDEKToStorage(ctx, kekID, keyID, dekSize)
}

func (c *KMSClient) decryptDEKFromStorage(ctx context.Context, kekID, keyID string) ([]byte, error) {
	encryptedKey, err := c.storage.Get(ctx, keyID)
	if err != nil {
		return nil, err
	}
	decryptInput := &kms.DecryptInput{
		CiphertextBlob:    encryptedKey,
		EncryptionContext: map[string]string{DEKContext: keyID},
		KeyId:             &kekID,
	}
	decryptOutput, err := c.awsClient.Decrypt(ctx, decryptInput)
	if err != nil {
		return nil, err
	}
	return decryptOutput.Plaintext, nil
}

func (c *KMSClient) putNewDEKToStorage(ctx context.Context, kekID, keyID string, dekSize int) ([]byte, error) {
	// GenerateDataKey always generates a new unique key, even if the input stays the same.
	input := &kms.GenerateDataKeyInput{
		KeyId: &kekID,
		// The encryption context is used for encryption. It must be the same when decrypting the ciphertext output
		EncryptionContext: map[string]string{DEKContext: keyID}, // https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html#encrypt_context
		NumberOfBytes:     aws.Int32(int32(dekSize)),
	}
	output, err := c.awsClient.GenerateDataKey(ctx, input)
	if err != nil {
		return nil, err
	}
	// store encrypted key in storage
	if err := c.storage.Put(ctx, keyID, output.CiphertextBlob); err != nil {
		return nil, err
	}
	return output.Plaintext, nil
}

func (c *KMSClient) tryCleanUpResources(ctx context.Context, generatedNewKey bool, awsGeneratedKeyID, alias string) {
	if !generatedNewKey {
		return
	}
	// Delete Alias
	deleteAliasInput := &kms.DeleteAliasInput{
		AliasName: &alias,
	}
	_, _ = c.awsClient.DeleteAlias(ctx, deleteAliasInput) // Might fail, ignoring the error.

	// Delete Key
	scheduleKeyDeletionInput := &kms.ScheduleKeyDeletionInput{
		KeyId:               &awsGeneratedKeyID,
		PendingWindowInDays: aws.Int32(7),
	}
	_, _ = c.awsClient.ScheduleKeyDeletion(ctx, scheduleKeyDeletionInput) // Might fail, ignoring the error.
}

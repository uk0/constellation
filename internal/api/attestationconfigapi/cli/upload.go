/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: BUSL-1.1
*/
package main

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/edgelesssys/constellation/v2/internal/api/attestationconfigapi"
	"github.com/edgelesssys/constellation/v2/internal/api/attestationconfigapi/cli/client"
	"github.com/edgelesssys/constellation/v2/internal/api/fetcher"
	"github.com/edgelesssys/constellation/v2/internal/attestation/variant"
	"github.com/edgelesssys/constellation/v2/internal/file"
	"github.com/edgelesssys/constellation/v2/internal/logger"
	"github.com/edgelesssys/constellation/v2/internal/staticupload"
	"github.com/edgelesssys/constellation/v2/internal/verify"
	"github.com/google/go-tdx-guest/proto/tdx"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newUploadCmd() *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:   "upload VARIANT KIND FILE",
		Short: "Upload an object to the attestationconfig API",

		Long: fmt.Sprintf("Upload a new object to the attestationconfig API. For snp-reports the new object is added to a cache folder first.\n"+
			"The CLI then determines the lowest version within the cache-window present in the cache and writes that value to the config api if necessary. "+
			"For guest-firmware objects the object is added to the API directly.\n"+
			"Please authenticate with AWS through your preferred method (e.g. environment variables, CLI) "+
			"to be able to upload to S3. Set the %s and %s environment variables to authenticate with cosign.",
			envCosignPrivateKey, envCosignPwd,
		),
		Example: "COSIGN_PASSWORD=$CPW COSIGN_PRIVATE_KEY=$CKEY cli upload azure-sev-snp attestation-report /some/path/report.json",

		Args:    cobra.MatchAll(cobra.ExactArgs(3), arg0isAttestationVariant(), isValidKind(1)),
		PreRunE: envCheck,
		RunE:    runUpload,
	}
	uploadCmd.Flags().StringP("upload-date", "d", "", "upload a version with this date as version name.")
	uploadCmd.Flags().BoolP("force", "f", false, "Use force to manually push a new latest version."+
		" The version gets saved to the cache but the version selection logic is skipped.")
	uploadCmd.Flags().IntP("cache-window-size", "s", versionWindowSize, "Number of versions to be considered for the latest version.")

	return uploadCmd
}

func envCheck(_ *cobra.Command, _ []string) error {
	if os.Getenv(envCosignPrivateKey) == "" || os.Getenv(envCosignPwd) == "" {
		return fmt.Errorf("please set both %s and %s environment variables", envCosignPrivateKey, envCosignPwd)
	}
	cosignPwd = os.Getenv(envCosignPwd)
	privateKey = os.Getenv(envCosignPrivateKey)
	return nil
}

func runUpload(cmd *cobra.Command, args []string) (retErr error) {
	ctx := cmd.Context()
	log := logger.NewTextLogger(slog.LevelDebug).WithGroup("attestationconfigapi")

	uploadCfg, err := newConfig(cmd, ([3]string)(args[:3]))
	if err != nil {
		return fmt.Errorf("parsing cli flags: %w", err)
	}

	client, clientClose, err := client.New(ctx,
		staticupload.Config{
			Bucket:         uploadCfg.bucket,
			Region:         uploadCfg.region,
			DistributionID: uploadCfg.distribution,
		},
		[]byte(cosignPwd), []byte(privateKey),
		false, uploadCfg.cacheWindowSize, log,
	)

	defer func() {
		err := clientClose(cmd.Context())
		if err != nil {
			retErr = errors.Join(retErr, fmt.Errorf("failed to invalidate cache: %w", err))
		}
	}()

	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}

	return uploadReport(ctx, client, uploadCfg, file.NewHandler(afero.NewOsFs()), log)
}

func uploadReport(
	ctx context.Context, apiClient *client.Client,
	cfg uploadConfig, fs file.Handler, log *slog.Logger,
) error {
	if cfg.kind != attestationReport {
		return fmt.Errorf("kind %s not supported", cfg.kind)
	}

	apiFetcher := attestationconfigapi.NewFetcherWithCustomCDNAndCosignKey(cfg.url, cfg.cosignPublicKey)
	latestVersionInAPI, err := apiFetcher.FetchLatestVersion(ctx, cfg.variant)
	if err != nil {
		var notFoundErr *fetcher.NotFoundError
		if errors.As(err, &notFoundErr) {
			log.Info("No versions found in API, but assuming that we are uploading the first version.")
		} else {
			return fmt.Errorf("fetching latest version: %w", err)
		}
	}

	var newVersion, latestVersion any
	switch cfg.variant {
	case variant.AWSSEVSNP{}, variant.AzureSEVSNP{}, variant.GCPSEVSNP{}:
		latestVersion = latestVersionInAPI.SEVSNPVersion

		log.Info(fmt.Sprintf("Reading SNP report from file: %s", cfg.path))
		newVersion, err = readSNPReport(cfg.path, fs)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("Input SNP report: %+v", newVersion))

	case variant.AzureTDX{}:
		latestVersion = latestVersionInAPI.TDXVersion

		log.Info(fmt.Sprintf("Reading TDX report from file: %s", cfg.path))
		newVersion, err = readTDXReport(cfg.path, fs)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("Input TDX report: %+v", newVersion))

	default:
		return fmt.Errorf("variant %s not supported", cfg.variant)
	}

	if err := apiClient.UploadLatestVersion(
		ctx, cfg.variant, newVersion, latestVersion, cfg.uploadDate, cfg.force,
	); err != nil && !errors.Is(err, client.ErrNoNewerVersion) {
		return fmt.Errorf("updating latest version: %w", err)
	}

	return nil
}

func convertTCBVersionToSNPVersion(tcb verify.TCBVersion) attestationconfigapi.SEVSNPVersion {
	return attestationconfigapi.SEVSNPVersion{
		Bootloader: tcb.Bootloader,
		TEE:        tcb.TEE,
		SNP:        tcb.SNP,
		Microcode:  tcb.Microcode,
	}
}

func convertQuoteToTDXVersion(quote *tdx.QuoteV4) attestationconfigapi.TDXVersion {
	return attestationconfigapi.TDXVersion{
		QESVN:      binary.LittleEndian.Uint16(quote.Header.QeSvn),
		PCESVN:     binary.LittleEndian.Uint16(quote.Header.PceSvn),
		QEVendorID: [16]byte(quote.Header.QeVendorId),
		XFAM:       [8]byte(quote.TdQuoteBody.Xfam),
		TEETCBSVN:  [16]byte(quote.TdQuoteBody.TeeTcbSvn),
	}
}

type uploadConfig struct {
	variant         variant.Variant
	kind            objectKind
	path            string
	uploadDate      time.Time
	cosignPublicKey string
	region          string
	bucket          string
	distribution    string
	url             string
	force           bool
	cacheWindowSize int
}

func newConfig(cmd *cobra.Command, args [3]string) (uploadConfig, error) {
	dateStr, err := cmd.Flags().GetString("upload-date")
	if err != nil {
		return uploadConfig{}, fmt.Errorf("getting upload date: %w", err)
	}
	uploadDate := time.Now()
	if dateStr != "" {
		uploadDate, err = time.Parse(client.VersionFormat, dateStr)
		if err != nil {
			return uploadConfig{}, fmt.Errorf("parsing date: %w", err)
		}
	}

	region, err := cmd.Flags().GetString("region")
	if err != nil {
		return uploadConfig{}, fmt.Errorf("getting region: %w", err)
	}

	bucket, err := cmd.Flags().GetString("bucket")
	if err != nil {
		return uploadConfig{}, fmt.Errorf("getting bucket: %w", err)
	}

	testing, err := cmd.Flags().GetBool("testing")
	if err != nil {
		return uploadConfig{}, fmt.Errorf("getting testing flag: %w", err)
	}
	apiCfg := getAPIEnvironment(testing)

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return uploadConfig{}, fmt.Errorf("getting force: %w", err)
	}

	cacheWindowSize, err := cmd.Flags().GetInt("cache-window-size")
	if err != nil {
		return uploadConfig{}, fmt.Errorf("getting cache window size: %w", err)
	}

	variant, err := variant.FromString(args[0])
	if err != nil {
		return uploadConfig{}, fmt.Errorf("invalid attestation variant: %q: %w", args[0], err)
	}

	kind := kindFromString(args[1])
	path := args[2]

	return uploadConfig{
		variant:         variant,
		kind:            kind,
		path:            path,
		uploadDate:      uploadDate,
		cosignPublicKey: apiCfg.cosignPublicKey,
		region:          region,
		bucket:          bucket,
		url:             apiCfg.url,
		distribution:    apiCfg.distribution,
		force:           force,
		cacheWindowSize: cacheWindowSize,
	}, nil
}

/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/edgelesssys/constellation/v2/cli/internal/cloudcmd"
	"github.com/edgelesssys/constellation/v2/cli/internal/helm"
	"github.com/edgelesssys/constellation/v2/cli/internal/kubecmd"
	"github.com/edgelesssys/constellation/v2/cli/internal/state"
	"github.com/edgelesssys/constellation/v2/cli/internal/terraform"
	"github.com/edgelesssys/constellation/v2/internal/api/attestationconfigapi"
	"github.com/edgelesssys/constellation/v2/internal/attestation/variant"
	"github.com/edgelesssys/constellation/v2/internal/cloud/cloudprovider"
	"github.com/edgelesssys/constellation/v2/internal/compatibility"
	"github.com/edgelesssys/constellation/v2/internal/config"
	"github.com/edgelesssys/constellation/v2/internal/constants"
	"github.com/edgelesssys/constellation/v2/internal/file"
	"github.com/edgelesssys/constellation/v2/internal/kms/uri"
	"github.com/edgelesssys/constellation/v2/internal/versions"
	"github.com/rogpeppe/go-internal/diff"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

const (
	// skipInfrastructurePhase skips the terraform apply of the upgrade process.
	skipInfrastructurePhase skipPhase = "infrastructure"
	// skipHelmPhase skips the helm upgrade of the upgrade process.
	skipHelmPhase skipPhase = "helm"
	// skipImagePhase skips the image upgrade of the upgrade process.
	skipImagePhase skipPhase = "image"
	// skipK8sPhase skips the k8s upgrade of the upgrade process.
	skipK8sPhase skipPhase = "k8s"
)

// skipPhase is a phase of the upgrade process that can be skipped.
type skipPhase string

func newUpgradeApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply an upgrade to a Constellation cluster",
		Long:  "Apply an upgrade to a Constellation cluster by applying the chosen configuration.",
		Args:  cobra.NoArgs,
		RunE:  runUpgradeApply,
	}

	cmd.Flags().BoolP("yes", "y", false, "run upgrades without further confirmation\n"+
		"WARNING: might delete your resources in case you are using cert-manager in your cluster. Please read the docs.\n"+
		"WARNING: might unintentionally overwrite measurements in the running cluster.")
	cmd.Flags().Duration("timeout", 5*time.Minute, "change helm upgrade timeout\n"+
		"Might be useful for slow connections or big clusters.")
	cmd.Flags().Bool("conformance", false, "enable conformance mode")
	cmd.Flags().Bool("skip-helm-wait", false, "install helm charts without waiting for deployments to be ready")
	cmd.Flags().StringSlice("skip-phases", nil, "comma-separated list of upgrade phases to skip\n"+
		"one or multiple of { infrastructure | helm | image | k8s }")
	if err := cmd.Flags().MarkHidden("timeout"); err != nil {
		panic(err)
	}

	return cmd
}

type upgradeApplyFlags struct {
	rootFlags
	yes            bool
	upgradeTimeout time.Duration
	conformance    bool
	helmWaitMode   helm.WaitMode
	skipPhases     skipPhases
}

func (f *upgradeApplyFlags) parse(flags *pflag.FlagSet) error {
	if err := f.rootFlags.parse(flags); err != nil {
		return err
	}

	rawSkipPhases, err := flags.GetStringSlice("skip-phases")
	if err != nil {
		return fmt.Errorf("parsing skip-phases flag: %w", err)
	}
	var skipPhases []skipPhase
	for _, phase := range rawSkipPhases {
		switch skipPhase(phase) {
		case skipInfrastructurePhase, skipHelmPhase, skipImagePhase, skipK8sPhase:
			skipPhases = append(skipPhases, skipPhase(phase))
		default:
			return fmt.Errorf("invalid phase %s", phase)
		}
	}
	f.skipPhases = skipPhases

	f.yes, err = flags.GetBool("yes")
	if err != nil {
		return fmt.Errorf("getting 'yes' flag: %w", err)
	}

	f.upgradeTimeout, err = flags.GetDuration("timeout")
	if err != nil {
		return fmt.Errorf("getting 'timeout' flag: %w", err)
	}

	f.conformance, err = flags.GetBool("conformance")
	if err != nil {
		return fmt.Errorf("getting 'conformance' flag: %w", err)
	}
	skipHelmWait, err := flags.GetBool("skip-helm-wait")
	if err != nil {
		return fmt.Errorf("getting 'skip-helm-wait' flag: %w", err)
	}
	f.helmWaitMode = helm.WaitModeAtomic
	if skipHelmWait {
		f.helmWaitMode = helm.WaitModeNone
	}

	return nil
}

func runUpgradeApply(cmd *cobra.Command, _ []string) error {
	log, err := newCLILogger(cmd)
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}
	defer log.Sync()

	fileHandler := file.NewHandler(afero.NewOsFs())
	upgradeID := generateUpgradeID(upgradeCmdKindApply)

	kubeUpgrader, err := kubecmd.New(cmd.OutOrStdout(), constants.AdminConfFilename, fileHandler, log)
	if err != nil {
		return err
	}

	configFetcher := attestationconfigapi.NewFetcher()

	var flags upgradeApplyFlags
	if err := flags.parse(cmd.Flags()); err != nil {
		return err
	}

	// Set up terraform upgrader
	upgradeDir := filepath.Join(constants.UpgradeDir, upgradeID)
	clusterUpgrader, err := cloudcmd.NewClusterUpgrader(
		cmd.Context(),
		constants.TerraformWorkingDir,
		upgradeDir,
		flags.tfLogLevel,
		fileHandler,
	)
	if err != nil {
		return fmt.Errorf("setting up cluster upgrader: %w", err)
	}

	helmClient, err := helm.NewClient(constants.AdminConfFilename, log)
	if err != nil {
		return fmt.Errorf("creating Helm client: %w", err)
	}

	applyCmd := upgradeApplyCmd{
		kubeUpgrader:    kubeUpgrader,
		helmApplier:     helmClient,
		clusterUpgrader: clusterUpgrader,
		configFetcher:   configFetcher,
		fileHandler:     fileHandler,
		flags:           flags,
		log:             log,
	}
	return applyCmd.upgradeApply(cmd, upgradeDir)
}

type upgradeApplyCmd struct {
	helmApplier     helmApplier
	kubeUpgrader    kubernetesUpgrader
	clusterUpgrader clusterUpgrader
	configFetcher   attestationconfigapi.Fetcher
	fileHandler     file.Handler
	flags           upgradeApplyFlags
	log             debugLog
}

func (u *upgradeApplyCmd) upgradeApply(cmd *cobra.Command, upgradeDir string) error {
	conf, err := config.New(u.fileHandler, constants.ConfigFilename, u.configFetcher, u.flags.force)
	var configValidationErr *config.ValidationError
	if errors.As(err, &configValidationErr) {
		cmd.PrintErrln(configValidationErr.LongMessage())
	}
	if err != nil {
		return err
	}
	if cloudcmd.UpgradeRequiresIAMMigration(conf.GetProvider()) {
		cmd.Println("WARNING: This upgrade requires an IAM migration. Please make sure you have applied the IAM migration using `iam upgrade apply` before continuing.")
		if !u.flags.yes {
			yes, err := askToConfirm(cmd, "Did you upgrade the IAM resources?")
			if err != nil {
				return fmt.Errorf("asking for confirmation: %w", err)
			}
			if !yes {
				cmd.Println("Skipping upgrade.")
				return nil
			}
		}
	}
	conf.KubernetesVersion, err = validK8sVersion(cmd, string(conf.KubernetesVersion), u.flags.yes)
	if err != nil {
		return err
	}

	stateFile, err := state.ReadFromFile(u.fileHandler, constants.StateFilename)
	if err != nil {
		return fmt.Errorf("reading state file: %w", err)
	}

	if err := u.confirmAndUpgradeAttestationConfig(cmd, conf.GetAttestationConfig(), stateFile.ClusterValues.MeasurementSalt); err != nil {
		return fmt.Errorf("upgrading measurements: %w", err)
	}

	// If infrastructure phase is skipped, we expect the new infrastructure
	// to be in the Terraform configuration already. Otherwise, perform
	// the Terraform migrations.
	if !u.flags.skipPhases.contains(skipInfrastructurePhase) {
		migrationRequired, err := u.planTerraformMigration(cmd, conf)
		if err != nil {
			return fmt.Errorf("planning Terraform migrations: %w", err)
		}

		if migrationRequired {
			postMigrationInfraState, err := u.migrateTerraform(cmd, conf, upgradeDir)
			if err != nil {
				return fmt.Errorf("performing Terraform migrations: %w", err)
			}

			// Merge the pre-upgrade state with the post-migration infrastructure values
			if _, err := stateFile.Merge(
				// temporary state with post-migration infrastructure values
				state.New().SetInfrastructure(postMigrationInfraState),
			); err != nil {
				return fmt.Errorf("merging pre-upgrade state with post-migration infrastructure values: %w", err)
			}

			// Write the post-migration state to disk
			if err := stateFile.WriteToFile(u.fileHandler, constants.StateFilename); err != nil {
				return fmt.Errorf("writing state file: %w", err)
			}
		}
	}

	// extend the clusterConfig cert SANs with any of the supported endpoints:
	// - (legacy) public IP
	// - fallback endpoint
	// - custom (user-provided) endpoint
	sans := append([]string{stateFile.Infrastructure.ClusterEndpoint, conf.CustomEndpoint}, stateFile.Infrastructure.APIServerCertSANs...)
	if err := u.kubeUpgrader.ExtendClusterConfigCertSANs(cmd.Context(), sans); err != nil {
		return fmt.Errorf("extending cert SANs: %w", err)
	}

	if conf.GetProvider() != cloudprovider.Azure && conf.GetProvider() != cloudprovider.GCP && conf.GetProvider() != cloudprovider.AWS {
		cmd.PrintErrln("WARNING: Skipping service and image upgrades, which are currently only supported for AWS, Azure, and GCP.")
		return nil
	}

	var upgradeErr *compatibility.InvalidUpgradeError
	if !u.flags.skipPhases.contains(skipHelmPhase) {
		err = u.handleServiceUpgrade(cmd, conf, stateFile, upgradeDir)
		switch {
		case errors.As(err, &upgradeErr):
			cmd.PrintErrln(err)
		case err == nil:
			cmd.Println("Successfully upgraded Constellation services.")
		case err != nil:
			return fmt.Errorf("upgrading services: %w", err)
		}
	}
	skipImageUpgrade := u.flags.skipPhases.contains(skipImagePhase)
	skipK8sUpgrade := u.flags.skipPhases.contains(skipK8sPhase)
	if !(skipImageUpgrade && skipK8sUpgrade) {
		err = u.kubeUpgrader.UpgradeNodeVersion(cmd.Context(), conf, u.flags.force, skipImageUpgrade, skipK8sUpgrade)
		switch {
		case errors.Is(err, kubecmd.ErrInProgress):
			cmd.PrintErrln("Skipping image and Kubernetes upgrades. Another upgrade is in progress.")
		case errors.As(err, &upgradeErr):
			cmd.PrintErrln(err)
		case err != nil:
			return fmt.Errorf("upgrading NodeVersion: %w", err)
		}
	}
	return nil
}

func diffAttestationCfg(currentAttestationCfg config.AttestationCfg, newAttestationCfg config.AttestationCfg) (string, error) {
	// cannot compare structs directly with go-cmp because of unexported fields in the attestation config
	currentYml, err := yaml.Marshal(currentAttestationCfg)
	if err != nil {
		return "", fmt.Errorf("marshalling remote attestation config: %w", err)
	}
	newYml, err := yaml.Marshal(newAttestationCfg)
	if err != nil {
		return "", fmt.Errorf("marshalling local attestation config: %w", err)
	}
	diff := string(diff.Diff("current", currentYml, "new", newYml))
	return diff, nil
}

// planTerraformMigration checks if the Constellation version the cluster is being upgraded to requires a migration.
func (u *upgradeApplyCmd) planTerraformMigration(cmd *cobra.Command, conf *config.Config) (bool, error) {
	u.log.Debugf("Planning Terraform migrations")

	vars, err := cloudcmd.TerraformUpgradeVars(conf)
	if err != nil {
		return false, fmt.Errorf("parsing upgrade variables: %w", err)
	}
	u.log.Debugf("Using Terraform variables:\n%v", vars)

	// Check if there are any Terraform migrations to apply

	// Add manual migrations here if required
	//
	// var manualMigrations []terraform.StateMigration
	// for _, migration := range manualMigrations {
	// 	  u.log.Debugf("Adding manual Terraform migration: %s", migration.DisplayName)
	// 	  u.upgrader.AddManualStateMigration(migration)
	// }

	return u.clusterUpgrader.PlanClusterUpgrade(cmd.Context(), cmd.OutOrStdout(), vars, conf.GetProvider())
}

// migrateTerraform checks if the Constellation version the cluster is being upgraded to requires a migration
// of cloud resources with Terraform. If so, the migration is performed and the post-migration infrastructure state is returned.
// If no migration is required, the current (pre-upgrade) infrastructure state is returned.
func (u *upgradeApplyCmd) migrateTerraform(cmd *cobra.Command, conf *config.Config, upgradeDir string,
) (state.Infrastructure, error) {
	// If there are any Terraform migrations to apply, ask for confirmation
	fmt.Fprintln(cmd.OutOrStdout(), "The upgrade requires a migration of Constellation cloud resources by applying an updated Terraform template. Please manually review the suggested changes below.")
	if !u.flags.yes {
		ok, err := askToConfirm(cmd, "Do you want to apply the Terraform migrations?")
		if err != nil {
			return state.Infrastructure{}, fmt.Errorf("asking for confirmation: %w", err)
		}
		if !ok {
			cmd.Println("Aborting upgrade.")
			// User doesn't expect to see any changes in his workspace after aborting an "upgrade apply",
			// therefore, roll back to the backed up state.
			if err := u.clusterUpgrader.RestoreClusterWorkspace(); err != nil {
				return state.Infrastructure{}, fmt.Errorf(
					"restoring Terraform workspace: %w, restore the Terraform workspace manually from %s ",
					err,
					filepath.Join(upgradeDir, constants.TerraformUpgradeBackupDir),
				)
			}
			return state.Infrastructure{}, fmt.Errorf("cluster upgrade aborted by user")
		}
	}
	u.log.Debugf("Applying Terraform migrations")

	infraState, err := u.clusterUpgrader.ApplyClusterUpgrade(cmd.Context(), conf.GetProvider())
	if err != nil {
		return state.Infrastructure{}, fmt.Errorf("applying terraform migrations: %w", err)
	}

	cmd.Printf("Infrastructure migrations applied successfully and output written to: %s\n"+
		"A backup of the pre-upgrade state has been written to: %s\n",
		u.flags.pathPrefixer.PrefixPrintablePath(constants.StateFilename),
		u.flags.pathPrefixer.PrefixPrintablePath(filepath.Join(upgradeDir, constants.TerraformUpgradeBackupDir)),
	)
	return infraState, nil
}

// validK8sVersion checks if the Kubernetes patch version is supported and asks for confirmation if not.
func validK8sVersion(cmd *cobra.Command, version string, yes bool) (validVersion versions.ValidK8sVersion, err error) {
	validVersion, err = versions.NewValidK8sVersion(version, true)
	if versions.IsPreviewK8sVersion(validVersion) {
		cmd.PrintErrf("Warning: Constellation with Kubernetes %v is still in preview. Use only for evaluation purposes.\n", validVersion)
	}
	valid := err == nil

	if !valid && !yes {
		confirmed, err := askToConfirm(cmd, fmt.Sprintf("WARNING: The Kubernetes patch version %s is not supported. If you continue, Kubernetes upgrades will be skipped. Do you want to continue anyway?", version))
		if err != nil {
			return validVersion, fmt.Errorf("asking for confirmation: %w", err)
		}
		if !confirmed {
			return validVersion, fmt.Errorf("aborted by user")
		}
	}

	return validVersion, nil
}

// confirmAndUpgradeAttestationConfig checks if the locally configured measurements are different from the cluster's measurements.
// If so the function will ask the user to confirm (if --yes is not set) and upgrade the cluster's config.
func (u *upgradeApplyCmd) confirmAndUpgradeAttestationConfig(
	cmd *cobra.Command, newConfig config.AttestationCfg, measurementSalt []byte,
) error {
	clusterAttestationConfig, err := u.kubeUpgrader.GetClusterAttestationConfig(cmd.Context(), newConfig.GetVariant())
	if err != nil {
		return fmt.Errorf("getting cluster attestation config: %w", err)
	}

	// If the current config is equal, or there is an error when comparing the configs, we skip the upgrade.
	equal, err := newConfig.EqualTo(clusterAttestationConfig)
	if err != nil {
		return fmt.Errorf("comparing attestation configs: %w", err)
	}
	if equal {
		return nil
	}
	cmd.Println("The configured attestation config is different from the attestation config in the cluster.")
	diffStr, err := diffAttestationCfg(clusterAttestationConfig, newConfig)
	if err != nil {
		return fmt.Errorf("diffing attestation configs: %w", err)
	}
	cmd.Println("The following changes will be applied to the attestation config:")
	cmd.Println(diffStr)
	if !u.flags.yes {
		ok, err := askToConfirm(cmd, "Are you sure you want to change your cluster's attestation config?")
		if err != nil {
			return fmt.Errorf("asking for confirmation: %w", err)
		}
		if !ok {
			return errors.New("aborting upgrade since attestation config is different")
		}
	}

	if err := u.kubeUpgrader.ApplyJoinConfig(cmd.Context(), newConfig, measurementSalt); err != nil {
		return fmt.Errorf("updating attestation config: %w", err)
	}
	cmd.Println("Successfully updated the cluster's attestation config")
	return nil
}

func (u *upgradeApplyCmd) handleServiceUpgrade(
	cmd *cobra.Command, conf *config.Config, stateFile *state.State, upgradeDir string,
) error {
	var secret uri.MasterSecret
	if err := u.fileHandler.ReadJSON(constants.MasterSecretFilename, &secret); err != nil {
		return fmt.Errorf("reading master secret: %w", err)
	}
	serviceAccURI, err := cloudcmd.GetMarshaledServiceAccountURI(conf.GetProvider(), conf, u.flags.pathPrefixer, u.log, u.fileHandler)
	if err != nil {
		return fmt.Errorf("getting service account URI: %w", err)
	}
	options := helm.Options{
		Force:        u.flags.force,
		Conformance:  u.flags.conformance,
		HelmWaitMode: u.flags.helmWaitMode,
	}

	prepareApply := func(allowDestructive bool) (helm.Applier, bool, error) {
		options.AllowDestructive = allowDestructive
		executor, includesUpgrades, err := u.helmApplier.PrepareApply(conf, stateFile, options, serviceAccURI, secret)
		var upgradeErr *compatibility.InvalidUpgradeError
		switch {
		case errors.As(err, &upgradeErr):
			cmd.PrintErrln(err)
		case err != nil:
			return nil, false, fmt.Errorf("getting chart executor: %w", err)
		}
		return executor, includesUpgrades, nil
	}

	executor, includesUpgrades, err := prepareApply(helm.DenyDestructive)
	if err != nil {
		if !errors.Is(err, helm.ErrConfirmationMissing) {
			return fmt.Errorf("upgrading charts with deny destructive mode: %w", err)
		}
		if !u.flags.yes {
			cmd.PrintErrln("WARNING: Upgrading cert-manager will destroy all custom resources you have manually created that are based on the current version of cert-manager.")
			ok, askErr := askToConfirm(cmd, "Do you want to upgrade cert-manager anyway?")
			if askErr != nil {
				return fmt.Errorf("asking for confirmation: %w", err)
			}
			if !ok {
				cmd.Println("Skipping upgrade.")
				return nil
			}
		}
		executor, includesUpgrades, err = prepareApply(helm.AllowDestructive)
		if err != nil {
			return fmt.Errorf("upgrading charts with allow destructive mode: %w", err)
		}
	}

	// Save the Helm charts for the upgrade to disk
	chartDir := filepath.Join(upgradeDir, "helm-charts")
	if err := executor.SaveCharts(chartDir, u.fileHandler); err != nil {
		return fmt.Errorf("saving Helm charts to disk: %w", err)
	}
	u.log.Debugf("Helm charts saved to %s", chartDir)

	if includesUpgrades {
		u.log.Debugf("Creating backup of CRDs and CRs")
		crds, err := u.kubeUpgrader.BackupCRDs(cmd.Context(), upgradeDir)
		if err != nil {
			return fmt.Errorf("creating CRD backup: %w", err)
		}
		if err := u.kubeUpgrader.BackupCRs(cmd.Context(), crds, upgradeDir); err != nil {
			return fmt.Errorf("creating CR backup: %w", err)
		}
	}
	if err := executor.Apply(cmd.Context()); err != nil {
		return fmt.Errorf("applying Helm charts: %w", err)
	}

	return nil
}

// skipPhases is a list of phases that can be skipped during the upgrade process.
type skipPhases []skipPhase

// contains returns true if the list of phases contains the given phase.
func (s skipPhases) contains(phase skipPhase) bool {
	for _, p := range s {
		if strings.EqualFold(string(p), string(phase)) {
			return true
		}
	}
	return false
}

type kubernetesUpgrader interface {
	UpgradeNodeVersion(ctx context.Context, conf *config.Config, force, skipImage, skipK8s bool) error
	ExtendClusterConfigCertSANs(ctx context.Context, alternativeNames []string) error
	GetClusterAttestationConfig(ctx context.Context, variant variant.Variant) (config.AttestationCfg, error)
	ApplyJoinConfig(ctx context.Context, newAttestConfig config.AttestationCfg, measurementSalt []byte) error
	BackupCRs(ctx context.Context, crds []apiextensionsv1.CustomResourceDefinition, upgradeDir string) error
	BackupCRDs(ctx context.Context, upgradeDir string) ([]apiextensionsv1.CustomResourceDefinition, error)
}

type clusterUpgrader interface {
	PlanClusterUpgrade(ctx context.Context, outWriter io.Writer, vars terraform.Variables, csp cloudprovider.Provider) (bool, error)
	ApplyClusterUpgrade(ctx context.Context, csp cloudprovider.Provider) (state.Infrastructure, error)
	RestoreClusterWorkspace() error
}

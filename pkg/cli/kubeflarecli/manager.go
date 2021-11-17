package kubeflarecli

import (
	"os"

	"github.com/replicatedhq/kubeflare/pkg/apis"
	accessapplicationcontroller "github.com/replicatedhq/kubeflare/pkg/controller/accessapplication"
	dnsrecordcontroller "github.com/replicatedhq/kubeflare/pkg/controller/dnsrecord"
	pagerulecontroller "github.com/replicatedhq/kubeflare/pkg/controller/pagerule"
	wafrulecontroller "github.com/replicatedhq/kubeflare/pkg/controller/webapplicationfirewallrule"
	workerroutecontroller "github.com/replicatedhq/kubeflare/pkg/controller/workerroute"
	zonecontroller "github.com/replicatedhq/kubeflare/pkg/controller/zone"
	"github.com/replicatedhq/kubeflare/pkg/logger"
	"github.com/replicatedhq/kubeflare/pkg/version"
	"github.com/replicatedhq/kubeflare/pkg/webhook"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func ManagerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "manager",
		Short:         "runs the kubeflare manager (in cluster controller)",
		Long:          `...`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Infof("Starting kubeflare manager version %+v", version.GetBuild())

			v := viper.GetViper()

			if v.GetString("log-level") == "debug" {
				logger.Info("setting log level to debug")
				logger.SetDebug()
			}

			// Get a config to talk to the apiserver
			cfg, err := config.GetConfig()
			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			// Create a new Cmd to provide shared dependencies and start components
			options := manager.Options{
				MetricsBindAddress: v.GetString("metrics-addr"),
				LeaderElection:     v.GetBool("leader-elect"),
				LeaderElectionID:   "leaderelection.kubeflare.io",
			}

			mgr, err := manager.New(cfg, options)
			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			// Setup Scheme for all resources
			if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			if err := zonecontroller.Add(mgr); err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			if err := dnsrecordcontroller.Add(mgr); err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			if err := pagerulecontroller.Add(mgr); err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			if err := accessapplicationcontroller.Add(mgr); err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			if err := wafrulecontroller.Add(mgr); err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			if err := workerroutecontroller.Add(mgr); err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			if err := webhook.AddToManager(mgr); err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			// Start the Cmd
			if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().String("metrics-addr", ":8088", "The address the metric endpoint binds to.")
	cmd.Flags().Bool("leader-elect", false, "Enable leader election for controller manager. "+
		"Enabling this will ensure there is only one active controller manager.")

	return cmd
}

package managercli

import (
	"os"
	"strings"

	"github.com/replicatedhq/kubeflare/pkg/apis"
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

func RunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "run",
		Short:         "runs the kubeflare manager",
		Long:          `...`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Infof("Starting kubeflare version %+v", version.GetBuild())

			v := viper.GetViper()

			if v.GetString("log-level") == "debug" {
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

	return cmd
}

func defaultManagerTag() string {
	tag := version.Version()
	if strings.HasPrefix(tag, "v") {
		tag = strings.TrimPrefix(tag, "v")
	}

	return tag
}

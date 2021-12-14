package integrationcli

import (
	"context"
	"io/ioutil"

	"github.com/cloudflare/cloudflare-go"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	kubeflarescheme "github.com/replicatedhq/kubeflare/pkg/client/kubeflareclientset/scheme"
	"github.com/replicatedhq/kubeflare/pkg/controller/zone"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func RunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "run",
		Short:         "runs a single integration test from a definition",
		Long:          `...`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()

			b, err := ioutil.ReadFile(v.GetString("spec"))
			if err != nil {
				return err
			}

			decode := kubeflarescheme.Codecs.UniversalDeserializer().Decode
			obj, _, err := decode(b, nil, nil)
			if err != nil {
				return err
			}
			instance := obj.(*crdsv1alpha1.Zone)
			instance.Name = v.GetString("zone-name")

			ctx := context.TODO()

			cf, err := cloudflare.New(v.GetString("key"), v.GetString("email"))
			if err != nil {
				return err
			}

			err = zone.ReconcileSettings(ctx, instance, cf)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String("spec", "", "the spec file to apply")
	cmd.MarkFlagRequired("spec")
	cmd.Flags().String("email", "", "email to use with cloudflare")
	cmd.MarkFlagRequired("email")
	cmd.Flags().String("key", "", "cloudflare api key")
	cmd.MarkFlagRequired("key")
	cmd.Flags().String("zone-name", "", "zone name (domain) to use")
	cmd.MarkFlagRequired("zone-name")

	return cmd
}

package kubeflarecli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	kubeflarescheme "github.com/replicatedhq/kubeflare/pkg/client/kubeflareclientset/scheme"
	"github.com/replicatedhq/kubeflare/pkg/cloudflare/dns"
	"github.com/replicatedhq/kubeflare/pkg/cloudflare/pagerules"
	"github.com/replicatedhq/kubeflare/pkg/cloudflare/workerroute"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func ImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "import",
		Short:         "import existing settings from cloudflare into custom resources",
		Long:          `...`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()

			_, err := os.Stat(v.GetString("output-dir"))
			if os.IsNotExist(err) {
				if err := os.MkdirAll(v.GetString("output-dir"), 0755); err != nil {
					return errors.Wrap(err, "mkdir")
				}
			} else if err != nil {
				return errors.Wrap(err, "stat")
			}

			cf, err := cloudflare.NewWithAPIToken(v.GetString("api-token"))
			if err != nil {
				return errors.Wrap(err, "create clouflare client")
			}

			zoneID, err := cf.ZoneIDByName(v.GetString("zone"))
			if err != nil {
				return errors.Wrap(err, "get zone id")
			}

			kubeflarescheme.AddToScheme(scheme.Scheme)
			s := serializer.NewYAMLSerializer(serializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

			if v.GetBool("dns-records") {
				dnsRecords, err := dns.FetchDNSRecordsForZone(v.GetString("api-token"), v.GetString("zone"), zoneID)
				if err != nil {
					return errors.Wrap(err, "fetch dns records")
				}

				for _, dnsRecord := range dnsRecords {
					buf := bytes.NewBuffer(nil)
					err := s.Encode(dnsRecord, buf)
					if err != nil {
						return errors.Wrap(err, "encode")
					}
					outputFile := filepath.Join(v.GetString("output-dir"), fmt.Sprintf("%s.yaml", dnsRecord.Name))
					if err := ioutil.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
						return errors.Wrap(err, "write file")
					}
				}
			}

			if v.GetBool("page-rules") {
				pageRules, err := pagerules.FetchPageRulesForZone(v.GetString("api-token"), v.GetString("zone"), zoneID)
				if err != nil {
					return errors.Wrap(err, "fetch page rules")
				}

				for _, pageRule := range pageRules {
					buf := bytes.NewBuffer(nil)
					err := s.Encode(pageRule, buf)
					if err != nil {
						return errors.Wrap(err, "encode")
					}
					outputFile := filepath.Join(v.GetString("output-dir"), fmt.Sprintf("%s.yaml", pageRule.Name))
					if err := ioutil.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
						return errors.Wrap(err, "write file")
					}
				}
			}

			if v.GetBool("worker-routes") {
				workerRoutes, err := workerroute.FetchWorkerRoutesForZone(v.GetString("api-token"), v.GetString("zone"), zoneID)
				if err != nil {
					return errors.Wrap(err, "fetch worker routes")
				}

				for _, workerRoute := range workerRoutes {
					buf := bytes.NewBuffer(nil)
					err := s.Encode(workerRoute, buf)
					if err != nil {
						return errors.Wrap(err, "encode")
					}
					outputFile := filepath.Join(v.GetString("output-dir"), fmt.Sprintf("%s.yaml", workerRoute.Name))
					if err := ioutil.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
						return errors.Wrap(err, "write file")
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().String("api-token", "", "cloudflare api token")
	cmd.MarkFlagRequired("api-token")

	cmd.Flags().String("zone", "", "dns zone to import")
	cmd.MarkFlagRequired("zone")

	cmd.Flags().String("output-dir", filepath.Join(".", "imported"), "output dir to write files to")

	cmd.Flags().Bool("dns-records", true, "when set, import existing dns records from the zone")
	cmd.Flags().Bool("page-rules", true, "when set, import existing page rules from the zone")
	cmd.Flags().Bool("worker-routes", true, "when set, import existing worker routes from the zone")

	return cmd
}

package kubeflarecli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/replicatedhq/kubeflare/pkg/version"
)

func Version() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "version",
		Short:         "kubeflare version information",
		Long:          `...`,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Kubeflare %s\n", version.Version())
			return nil
		},
	}

	return cmd
}

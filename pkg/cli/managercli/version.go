package managercli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/replicatedhq/kubeflare/pkg/version"
)

func Version() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "version",
		Short:         "manager version information",
		Long:          `...`,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Kubeflare Manager %s\n", version.Version())
			return nil
		},
	}

	return cmd
}

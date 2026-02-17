package cli

import "github.com/spf13/cobra"

func newClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Generate cluster manifests",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunParentCommand(cmd, "cluster")
		},
	}

	cmd.AddCommand(newClusterCNPGCmd())

	return cmd
}

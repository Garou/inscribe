package cli

import "github.com/spf13/cobra"

func newBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Generate backup manifests",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunParentCommand(cmd, "backup")
		},
	}

	cmd.AddCommand(newBackupCNPGCmd())

	return cmd
}

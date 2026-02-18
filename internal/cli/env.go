package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newEnvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env [path]",
		Short: "Output shell configuration for INSCRIBE_TEMPLATE_DIR",
		Long: `Output shell export statements to configure the template directory.

Add to your shell profile:
  eval "$(inscribe env /path/to/templates)"

Or with the current --template-dir value:
  eval "$(inscribe env)"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := templateDir
			if len(args) == 1 {
				dir = args[0]
			}

			absDir, err := filepath.Abs(dir)
			if err != nil {
				return fmt.Errorf("resolving path: %w", err)
			}

			cmd.Printf("export INSCRIBE_TEMPLATE_DIR=%q\n", absDir)
			return nil
		},
	}

	return cmd
}

package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	templateDir string
	outputDir   string
)

// NewRootCmd creates the root inscribe command.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inscribe",
		Short: "Generate Kubernetes manifests from templates",
		Long:  "Inscribe is an interactive CLI tool for generating Kubernetes manifest files via templating.",
	}

	cmd.PersistentFlags().StringVar(&templateDir, "template-dir", getEnvOrDefault("INSCRIBE_TEMPLATE_DIR", "templates"), "Path to template directory")
	cmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "o", ".", "Output directory for generated manifests")

	cmd.AddCommand(newClusterCmd())
	cmd.AddCommand(newBackupCmd())
	cmd.AddCommand(newEnvCmd())

	return cmd
}

func getEnvOrDefault(env, defaultVal string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return defaultVal
}

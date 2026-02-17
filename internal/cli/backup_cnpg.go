package cli

import "github.com/spf13/cobra"

func newBackupCNPGCmd() *cobra.Command {
	var (
		name        string
		namespace   string
		schedule    string
		clusterName string
		method      string
		context     string
		filename    string
	)

	cmd := &cobra.Command{
		Use:   "cnpg",
		Short: "Generate a CNPG scheduled backup manifest",
		Long:  "Generate a CloudNativePG (CNPG) scheduled backup manifest from a template.",
		RunE: func(cmd *cobra.Command, args []string) error {
			flagValues := make(map[string]string)
			if cmd.Flags().Changed("name") {
				flagValues["name"] = name
			}
			if cmd.Flags().Changed("namespace") {
				flagValues["namespace"] = namespace
			}
			if cmd.Flags().Changed("schedule") {
				flagValues["schedule"] = schedule
			}
			if cmd.Flags().Changed("cluster-name") {
				flagValues["cnpg-clusters"] = clusterName
			}
			if cmd.Flags().Changed("method") {
				flagValues["backup-methods"] = method
			}

			return RunBridge(BridgeConfig{
				TemplateName: "cnpg-backup",
				TemplateDir:  templateDir,
				OutputDir:    outputDir,
				FlagValues:   flagValues,
				Filename:     filename,
				Context:      context,
			})
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Backup name (dns-name)")
	cmd.Flags().StringVar(&namespace, "namespace", "", "Kubernetes namespace")
	cmd.Flags().StringVar(&schedule, "schedule", "", "Cron schedule expression")
	cmd.Flags().StringVar(&clusterName, "cluster-name", "", "CNPG cluster name to back up")
	cmd.Flags().StringVar(&method, "method", "", "Backup method (barmanObjectStore, volumeSnapshot)")
	cmd.Flags().StringVar(&context, "context", "", "Kubernetes context")
	cmd.Flags().StringVar(&filename, "filename", "", "Output filename")

	return cmd
}

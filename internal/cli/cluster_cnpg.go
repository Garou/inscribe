package cli

import "github.com/spf13/cobra"

func newClusterCNPGCmd() *cobra.Command {
	var (
		name      string
		namespace string
		instances string
		resources string
		context   string
		filename  string
	)

	cmd := &cobra.Command{
		Use:   "cnpg",
		Short: "Generate a CNPG PostgreSQL cluster manifest",
		Long:  "Generate a CloudNativePG (CNPG) PostgreSQL cluster manifest from a template.",
		RunE: func(cmd *cobra.Command, args []string) error {
			flagValues := make(map[string]string)
			if cmd.Flags().Changed("name") {
				flagValues["name"] = name
			}
			if cmd.Flags().Changed("namespace") {
				flagValues["namespace"] = namespace
			}
			if cmd.Flags().Changed("instances") {
				flagValues["instances"] = instances
			}
			if cmd.Flags().Changed("resources") {
				flagValues["cnpg-resource-templates"] = resources
			}

			return RunBridge(BridgeConfig{
				TemplateName: "cnpg-cluster",
				TemplateDir:  templateDir,
				OutputDir:    outputDir,
				FlagValues:   flagValues,
				Filename:     filename,
				Context:      context,
			})
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Cluster name (dns-name)")
	cmd.Flags().StringVar(&namespace, "namespace", "", "Kubernetes namespace")
	cmd.Flags().StringVar(&instances, "instances", "", "Number of instances")
	cmd.Flags().StringVar(&resources, "resources", "", "Resource template (e.g., \"Production - 4Gi/2CPU\")")
	cmd.Flags().StringVar(&context, "context", "", "Kubernetes context")
	cmd.Flags().StringVar(&filename, "filename", "", "Output filename")

	return cmd
}

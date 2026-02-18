package cli

import (
	"fmt"
	"sort"
	"strings"

	"inscribe/internal/domain"
	"inscribe/internal/engine"

	"github.com/spf13/cobra"
)

// BuildDynamicCommands loads the template registry from dir and builds
// cobra commands dynamically from the registered templates.
// Returns nil gracefully if dir is invalid or contains no templates.
func BuildDynamicCommands(dir string) []*cobra.Command {
	reg, err := engine.NewRegistry(dir)
	if err != nil {
		return nil
	}

	templates := reg.ListTemplates()
	if len(templates) == 0 {
		return nil
	}

	parents := make(map[string]*cobra.Command)

	for _, tmpl := range templates {
		segments := strings.Fields(tmpl.Command)
		if len(segments) < 2 {
			continue
		}

		parentName := segments[0]
		parent, ok := parents[parentName]
		if !ok {
			parent = buildParentCommand(parentName)
			parents[parentName] = parent
		}

		leaf := buildLeafCommand(reg, tmpl)
		parent.AddCommand(leaf)
	}

	var cmds []*cobra.Command
	for _, cmd := range parents {
		cmds = append(cmds, cmd)
	}
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].Use < cmds[j].Use
	})
	return cmds
}

// buildParentCommand creates a grouping command that delegates to RunParentCommand.
func buildParentCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Generate %s manifests", name),
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunParentCommand(cmd, name, templateDir)
		},
	}
}

// buildLeafCommand creates a leaf command with dynamic flags from the template's fields.
func buildLeafCommand(reg domain.TemplateRegistry, tmpl domain.TemplateMeta) *cobra.Command {
	segments := strings.Fields(tmpl.Command)
	leafName := segments[len(segments)-1]

	// Extract fields to register flags
	parser := engine.NewParser(reg)
	fields, err := parser.ExtractFields(tmpl.Name)
	if err != nil {
		fields = nil
	}

	// Storage for flag values — one per field plus context, kubeconfig, and filename
	flagVars := make(map[string]*string)
	for _, f := range fields {
		val := ""
		flagVars[f.Name] = &val
	}
	var context, kubeconfig, filename string

	cmd := &cobra.Command{
		Use:   leafName,
		Short: tmpl.Description,
		Long:  tmpl.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			flagValues := make(map[string]string)
			for name, ptr := range flagVars {
				if cmd.Flags().Changed(name) {
					flagValues[name] = *ptr
				}
			}

			return RunBridge(BridgeConfig{
				TemplateName: tmpl.Name,
				TemplateDir:  templateDir,
				OutputDir:    outputDir,
				FlagValues:   flagValues,
				Filename:     filename,
				Context:      context,
				Kubeconfig:   kubeconfig,
			})
		},
	}

	// Register a flag per extracted field, deduplicating by name since templates
	// may reference the same field multiple times (e.g. {{ input "name" "dns-name" }}
	// used in both metadata.name and service names).
	registered := make(map[string]bool)
	for _, f := range fields {
		if registered[f.Name] {
			continue
		}
		registered[f.Name] = true
		cmd.Flags().StringVar(flagVars[f.Name], f.Name, "", flagDescription(reg, f))
	}

	// Standard flags
	cmd.Flags().StringVar(&context, "context", "", "Kubernetes context")
	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file")
	cmd.Flags().StringVar(&filename, "filename", "", "Output filename")

	return cmd
}

// flagDescription generates help text for a dynamic flag based on its field type.
func flagDescription(reg domain.TemplateRegistry, f domain.FieldDefinition) string {
	switch f.Type {
	case domain.FieldInput:
		return fmt.Sprintf("Value for %s (validated as %s)", f.Name, f.ValidationType)
	case domain.FieldAutoList:
		return fmt.Sprintf("Value for %s (auto-listed from cluster if omitted)", f.Source)
	case domain.FieldTemplateGroup:
		subs, err := reg.GetSubTemplates(f.Source)
		if err != nil {
			return fmt.Sprintf("Template group: %s", f.Source)
		}
		var descs []string
		for _, s := range subs {
			descs = append(descs, fmt.Sprintf("%q", s.Description))
		}
		return fmt.Sprintf("One of: %s", strings.Join(descs, ", "))
	case domain.FieldStaticList:
		list, err := reg.GetStaticList(f.Source)
		if err != nil {
			return fmt.Sprintf("Static list: %s", f.Source)
		}
		return fmt.Sprintf("One of: %s", strings.Join(list.Items, ", "))
	default:
		return f.Name
	}
}

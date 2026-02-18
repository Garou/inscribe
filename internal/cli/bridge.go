package cli

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"inscribe/internal/domain"
	"inscribe/internal/engine"
	"inscribe/internal/kubernetes"
	"inscribe/internal/output"
	"inscribe/internal/tui"
	"inscribe/internal/tui/components/atoms"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// BridgeConfig holds the configuration for running the bridge.
type BridgeConfig struct {
	TemplateName string
	TemplateDir  string
	OutputDir    string
	FlagValues   map[string]string // CLI flag name → value (only set flags)
	Filename     string
	Context      string
	Kubeconfig   string
}

// RunBridge orchestrates the template→TUI→render→write flow.
func RunBridge(cfg BridgeConfig) error {
	// 1. Load template registry
	reg, err := engine.NewRegistry(cfg.TemplateDir)
	if err != nil {
		return fmt.Errorf("loading templates from %q: %w", cfg.TemplateDir, err)
	}

	// 2. Parse template (pass 1) to extract fields
	parser := engine.NewParser(reg)
	fields, err := parser.ExtractFields(cfg.TemplateName)
	if err != nil {
		return fmt.Errorf("extracting fields: %w", err)
	}

	// 3. Check which fields are satisfied by flags
	allProvided := true
	values := make(map[string]string)

	// Copy provided flag values
	for k, v := range cfg.FlagValues {
		values[k] = v
	}

	// Add context if provided
	if cfg.Context != "" {
		values["context"] = cfg.Context
	}

	// Validate provided values and check completeness
	for _, f := range fields {
		v, ok := values[f.Name]
		if !ok || v == "" {
			allProvided = false
			continue
		}

		// Validate manual fields
		if f.Type == domain.FieldManual {
			validator, err := domain.GetValidator(f.ValidationType)
			if err == nil {
				if err := validator(v); err != nil {
					return fmt.Errorf("invalid value for %q: %w", f.Name, err)
				}
			}
		}

		// Resolve templateGroup values by matching description to content
		if f.Type == domain.FieldTemplateGroup {
			subs, err := reg.GetSubTemplates(f.Source)
			if err != nil {
				return fmt.Errorf("loading sub-templates for %q: %w", f.Source, err)
			}
			resolved := false
			for _, sub := range subs {
				if matchesSubTemplate(v, sub) {
					values[f.Name] = sub.Content
					resolved = true
					break
				}
			}
			if !resolved {
				return fmt.Errorf("no matching sub-template %q for group %q (available: %s)", v, f.Source, listSubTemplateOptions(subs))
			}
		}

		// Resolve list values
		if f.Type == domain.FieldList {
			list, err := reg.GetStaticList(f.Source)
			if err != nil {
				return fmt.Errorf("loading static list %q: %w", f.Source, err)
			}
			found := false
			for _, item := range list.Items {
				if item == v {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("invalid value %q for list %q (available: %v)", v, f.Source, list.Items)
			}
		}
	}

	// 4. Decision: all provided → render directly, otherwise TUI
	if !allProvided || cfg.Filename == "" {
		client := kubernetes.NewClient(cfg.Kubeconfig)
		result, err := tui.RunWizard(fields, values, reg, client, cfg.Filename)
		if err != nil {
			return fmt.Errorf("wizard: %w", err)
		}
		values = result.Values
		cfg.Filename = result.Filename
	}

	// 5. Render template (pass 2)
	rendered, err := parser.Render(cfg.TemplateName, values)
	if err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	// 6. Write output
	writer := output.NewWriter()
	path, err := writer.Write(rendered, cfg.OutputDir, cfg.Filename)
	if err != nil {
		return fmt.Errorf("writing manifest: %w", err)
	}

	fmt.Printf("Manifest written to: %s\n", path)
	fmt.Println()
	printColoredYAML(rendered)
	fmt.Println()
	return nil
}

// printColoredYAML writes syntax-highlighted YAML to stdout.
func printColoredYAML(yaml string) {
	var buf bytes.Buffer
	err := quick.Highlight(&buf, yaml, "yaml", "terminal256", "monokai")
	if err != nil {
		fmt.Print(yaml)
		return
	}
	buf.WriteTo(os.Stdout)
}

// RunParentCommand handles parent commands (e.g. "inscribe cluster") by scanning
// the registry for templates matching the command prefix, then either auto-selecting
// (one match) or showing an interactive picker before delegating to the leaf subcommand.
func RunParentCommand(cmd *cobra.Command, commandPrefix string) error {
	reg, err := engine.NewRegistry(templateDir)
	if err != nil {
		return fmt.Errorf("loading templates from %q: %w", templateDir, err)
	}

	matches := reg.ListTemplatesByCommandPrefix(commandPrefix)
	if len(matches) == 0 {
		return fmt.Errorf("no templates found for %q in %q", commandPrefix, templateDir)
	}

	var selected domain.TemplateMeta
	if len(matches) == 1 {
		selected = matches[0]
	} else {
		options := make([]huh.Option[string], len(matches))
		for i, m := range matches {
			options[i] = huh.NewOption(m.Description, m.Name)
		}

		var choice string
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select a template").
					Options(options...).
					Value(&choice),
			),
		).WithTheme(atoms.Theme()).Run()
		if err != nil {
			return fmt.Errorf("template selection: %w", err)
		}

		for _, m := range matches {
			if m.Name == choice {
				selected = m
				break
			}
		}
	}

	// Extract the subcommand name: "cluster cnpg" with prefix "cluster" → "cnpg"
	subName := strings.TrimPrefix(selected.Command, commandPrefix+" ")
	sub, _, err := cmd.Find([]string{subName})
	if err != nil {
		return fmt.Errorf("subcommand %q not found: %w", subName, err)
	}

	return sub.RunE(sub, nil)
}

func matchesSubTemplate(value string, sub domain.SubTemplateMeta) bool {
	// Match by description (case-insensitive friendly name like "prod", "qa", "test")
	return value == sub.Description || value == sub.FilePath
}

func listSubTemplateOptions(subs []domain.SubTemplateMeta) string {
	var options string
	for i, sub := range subs {
		if i > 0 {
			options += ", "
		}
		options += fmt.Sprintf("%q", sub.Description)
	}
	return options
}

package tui

import (
	"fmt"

	"inscribe/internal/domain"
	"inscribe/internal/tui/components/atoms"
	"inscribe/internal/tui/components/molecules"
	"inscribe/internal/tui/components/organisms"

	"github.com/charmbracelet/huh"
)

// WizardResult holds all collected values after the wizard completes.
type WizardResult struct {
	Values   map[string]string
	Filename string
}

// RunWizard orchestrates the TUI wizard flow:
// 1. Context/namespace selection (if autoDetect fields exist)
// 2. Field collection for remaining fields
// 3. Filename input
func RunWizard(
	fields []domain.FieldDefinition,
	prefilledValues map[string]string,
	registry domain.TemplateRegistry,
	client domain.KubeClient,
	defaultFilename string,
) (*WizardResult, error) {
	// Initialize value pointers map with pre-filled values
	valuePtrs := make(map[string]*string)
	for _, f := range fields {
		s := ""
		if v, ok := prefilledValues[f.Name]; ok {
			s = v
		}
		valuePtrs[f.Name] = &s
	}

	// Determine if we need k8s context selection
	needsK8s := false
	needsCNPGSelect := false
	for _, f := range fields {
		if f.Type == domain.FieldAutoDetect {
			needsK8s = true
			if f.Source == "cnpg-clusters" {
				needsCNPGSelect = true
			}
		}
	}

	var contextValue, namespaceValue string
	if v, ok := prefilledValues["context"]; ok {
		contextValue = v
	}
	if v, ok := prefilledValues["namespace"]; ok {
		namespaceValue = v
	}

	// Phase 1: Context selection (if needed and not pre-filled)
	if needsK8s && contextValue == "" {
		contextForm := huh.NewForm(
			organisms.ContextSelectGroup(client, &contextValue, &namespaceValue),
		).WithTheme(atoms.Theme())

		if err := contextForm.Run(); err != nil {
			return nil, fmt.Errorf("context selection: %w", err)
		}
	}

	// Phase 2: Namespace selection (if needed and not pre-filled)
	if needsK8s && namespaceValue == "" {
		nsForm := huh.NewForm(
			organisms.NamespaceSelectGroup(client, contextValue, &namespaceValue),
		).WithTheme(atoms.Theme())

		if err := nsForm.Run(); err != nil {
			return nil, fmt.Errorf("namespace selection: %w", err)
		}
	}

	// Set autoDetect values
	if needsK8s {
		if p, ok := valuePtrs["namespace"]; ok {
			*p = namespaceValue
		}
	}

	// Phase 2.5: CNPG cluster selection (if needed)
	if needsCNPGSelect {
		if p, ok := valuePtrs["cnpg-clusters"]; ok && *p == "" {
			var cnpgCluster string
			cnpgForm := huh.NewForm(
				huh.NewGroup(
					molecules.K8sCNPGClusterSelect(client, contextValue, namespaceValue, &cnpgCluster),
				).Title("CNPG Cluster Selection"),
			).WithTheme(atoms.Theme())

			if err := cnpgForm.Run(); err != nil {
				return nil, fmt.Errorf("CNPG cluster selection: %w", err)
			}
			*p = cnpgCluster
		}
	}

	// Phase 3: Collect remaining fields (manual, templateGroup, list)
	remainingFields := filterNonAutoDetectFields(fields, prefilledValues)
	if len(remainingFields) > 0 {
		fieldGroup := organisms.FieldGroup(remainingFields, valuePtrs, registry)
		fieldForm := huh.NewForm(fieldGroup).WithTheme(atoms.Theme())

		if err := fieldForm.Run(); err != nil {
			return nil, fmt.Errorf("field collection: %w", err)
		}
	}

	// Phase 4: Filename
	filename := defaultFilename
	if filename == "" {
		filenameForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Output filename").
					Placeholder("manifest.yaml").
					Value(&filename).
					Validate(func(s string) error {
						_, err := domain.NewFilename(s)
						return err
					}),
			).Title("Output"),
		).WithTheme(atoms.Theme())

		if err := filenameForm.Run(); err != nil {
			return nil, fmt.Errorf("filename input: %w", err)
		}
	}

	// Collect final values
	result := &WizardResult{
		Values:   make(map[string]string),
		Filename: filename,
	}
	for name, ptr := range valuePtrs {
		result.Values[name] = *ptr
	}

	return result, nil
}

// filterNonAutoDetectFields returns fields that are not autoDetect and not already pre-filled.
func filterNonAutoDetectFields(fields []domain.FieldDefinition, prefilled map[string]string) []domain.FieldDefinition {
	var result []domain.FieldDefinition
	for _, f := range fields {
		if f.Type == domain.FieldAutoDetect {
			continue
		}
		if _, ok := prefilled[f.Name]; ok {
			continue
		}
		result = append(result, f)
	}
	return result
}

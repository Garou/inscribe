package tui

import (
	"testing"

	"inscribe/internal/domain"
)

func TestFilterNonAutoListFields(t *testing.T) {
	fields := []domain.FieldDefinition{
		{Name: "name", Type: domain.FieldInput, ValidationType: "dns-name"},
		{Name: "namespace", Type: domain.FieldAutoList, Source: "namespace"},
		{Name: "resources", Type: domain.FieldTemplateGroup, Source: "res-group"},
		{Name: "method", Type: domain.FieldStaticList, Source: "methods"},
		{Name: "cnpg-clusters", Type: domain.FieldAutoList, Source: "cnpg-clusters"},
	}

	t.Run("no prefilled values", func(t *testing.T) {
		result := filterNonAutoListFields(fields, map[string]string{})
		if len(result) != 3 {
			t.Fatalf("expected 3 fields, got %d", len(result))
		}
		if result[0].Name != "name" {
			t.Errorf("result[0].Name = %q, want %q", result[0].Name, "name")
		}
		if result[1].Name != "resources" {
			t.Errorf("result[1].Name = %q, want %q", result[1].Name, "resources")
		}
		if result[2].Name != "method" {
			t.Errorf("result[2].Name = %q, want %q", result[2].Name, "method")
		}
	})

	t.Run("with prefilled value", func(t *testing.T) {
		prefilled := map[string]string{"name": "mydb"}
		result := filterNonAutoListFields(fields, prefilled)
		if len(result) != 2 {
			t.Fatalf("expected 2 fields, got %d", len(result))
		}
		if result[0].Name != "resources" {
			t.Errorf("result[0].Name = %q, want %q", result[0].Name, "resources")
		}
	})

	t.Run("all prefilled", func(t *testing.T) {
		prefilled := map[string]string{
			"name":      "mydb",
			"resources": "prod",
			"method":    "barman",
		}
		result := filterNonAutoListFields(fields, prefilled)
		if len(result) != 0 {
			t.Errorf("expected 0 fields, got %d", len(result))
		}
	})

	t.Run("empty fields", func(t *testing.T) {
		result := filterNonAutoListFields(nil, map[string]string{})
		if len(result) != 0 {
			t.Errorf("expected 0 fields, got %d", len(result))
		}
	})
}

package output

import (
	"fmt"
	"os"
	"path/filepath"

	"inscribe/internal/domain"
)

// Writer implements domain.ManifestWriter.
type Writer struct{}

var _ domain.ManifestWriter = (*Writer)(nil)

// NewWriter creates a new manifest file writer.
func NewWriter() *Writer {
	return &Writer{}
}

// Write writes the rendered content to a file at outputDir/filename.
// Returns the full path of the written file.
func (w *Writer) Write(content string, outputDir string, filename string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("creating output directory %q: %w", outputDir, err)
	}

	fullPath := filepath.Join(outputDir, filename)

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("writing manifest to %q: %w", fullPath, err)
	}

	return fullPath, nil
}

package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fabricatorsltd/go-tools-cli/internal/catalog"
	"github.com/fabricatorsltd/go-tools-cli/internal/wizard"
)

func generateGoMod(cfg wizard.Config) ([]string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("module %s\n\ngo 1.21\n\nrequire (\n", cfg.ModulePath))

	for _, name := range cfg.SelectedModules {
		mod := catalog.FindModule(name)
		if mod == nil {
			continue
		}
		sb.WriteString(fmt.Sprintf("\t%s v0.0.0-00010101000000-000000000000 // indirect — replace with real version\n", mod.ImportPath))
	}
	sb.WriteString(")\n")

	outPath := filepath.Join(cfg.OutputDir, "go.mod")
	if err := os.WriteFile(outPath, []byte(sb.String()), 0644); err != nil {
		return nil, fmt.Errorf("writing go.mod: %w", err)
	}
	return []string{outPath}, nil
}

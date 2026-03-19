package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fabricatorsltd/go-tools-cli/internal/catalog"
	"github.com/fabricatorsltd/go-tools-cli/internal/wizard"
)

func generateMain(cfg wizard.Config) ([]string, error) {
	var sb strings.Builder
	sb.WriteString("package main\n\nimport (\n")

	for _, name := range cfg.SelectedModules {
		mod := catalog.FindModule(name)
		if mod == nil {
			continue
		}
		sb.WriteString(fmt.Sprintf("\t_ \"%s\"\n", mod.ImportPath))
	}
	sb.WriteString(")\n\nfunc main() {\n\t// TODO: implement your application\n}\n")

	outPath := filepath.Join(cfg.OutputDir, "main.go")
	if err := os.WriteFile(outPath, []byte(sb.String()), 0644); err != nil {
		return nil, fmt.Errorf("writing main.go: %w", err)
	}
	return []string{outPath}, nil
}

func generateReadme(cfg wizard.Config) ([]string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\nScaffolded with [go-tools-cli](https://github.com/fabricatorsltd/go-tools-cli).\n\n## Modules\n\n", cfg.ProjectName))

	for _, name := range cfg.SelectedModules {
		mod := catalog.FindModule(name)
		if mod == nil {
			sb.WriteString(fmt.Sprintf("- `%s`\n", name))
			continue
		}
		sb.WriteString(fmt.Sprintf("- [`%s`](https://%s) — %s\n", mod.Name, mod.ImportPath, mod.Description))
	}

	sb.WriteString("\n## Getting started\n\n```bash\ngo mod tidy\ngo run .\n```\n")

	outPath := filepath.Join(cfg.OutputDir, "README.md")
	if err := os.WriteFile(outPath, []byte(sb.String()), 0644); err != nil {
		return nil, fmt.Errorf("writing README.md: %w", err)
	}
	return []string{outPath}, nil
}

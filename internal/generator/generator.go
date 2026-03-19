package generator

import (
	"github.com/fabricatorsltd/go-tools-cli/internal/wizard"
)

// Generator orchestrates file creation
type Generator struct {
	cfg wizard.Config
}

func New(cfg wizard.Config) *Generator {
	return &Generator{cfg: cfg}
}

// Generate creates all project files and returns the list of created file paths
func (g *Generator) Generate() ([]string, error) {
	var files []string

	// Always generate go.mod and main.go
	gomodFiles, err := generateGoMod(g.cfg)
	if err != nil {
		return nil, err
	}
	files = append(files, gomodFiles...)

	mainFiles, err := generateMain(g.cfg)
	if err != nil {
		return nil, err
	}
	files = append(files, mainFiles...)

	// README
	readmeFiles, err := generateReadme(g.cfg)
	if err != nil {
		return nil, err
	}
	files = append(files, readmeFiles...)

	// Boilerplate or Full
	if g.cfg.OutputDepth >= wizard.DepthBoilerplate {
		bpFiles, err := generateBoilerplate(g.cfg)
		if err != nil {
			return nil, err
		}
		files = append(files, bpFiles...)
	}

	if g.cfg.OutputDepth >= wizard.DepthFull {
		fullFiles, err := generateFull(g.cfg)
		if err != nil {
			return nil, err
		}
		files = append(files, fullFiles...)
	}

	return files, nil
}

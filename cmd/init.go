package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fabricatorsltd/go-tools-cli/internal/generator"
	"github.com/fabricatorsltd/go-tools-cli/internal/wizard"
	cli "github.com/mirkobrombin/go-cli-builder/v2/pkg/cli"
)

type RootCmd struct {
	cli.Base
	Init InitCmd `cmd:"init" help:"Initialize a new go-tools project"`
}

type InitCmd struct {
	cli.Base
	Output string `cli:"output,o" help:"Output directory for the new project" default:"."`
	Name   string `cli:"name,n"   help:"Project name (skips interactive prompt)"`
	Module string `cli:"module,m" help:"Go module path, e.g. github.com/user/app (skips interactive prompt)"`
	Preset string `cli:"preset,p" help:"Preset: api, worker, fullstack, minimal (skips preset selection)"`
	Depth  string `cli:"depth,d"  help:"Output depth: minimal, boilerplate, full (skips depth selection)"`
}

func (c *InitCmd) Run() error {
	absOut, err := filepath.Abs(c.Output)
	if err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}
	if err := os.MkdirAll(absOut, 0755); err != nil {
		return fmt.Errorf("cannot create output dir: %w", err)
	}

	pre := wizard.Prefill{
		ProjectName: c.Name,
		ModulePath:  c.Module,
		PresetName:  c.Preset,
		Depth:       c.Depth,
	}

	m := wizard.NewModel(absOut, pre)
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return err
	}

	finalModel, ok := result.(wizard.Model)
	if !ok || !finalModel.Done {
		fmt.Println("Aborted.")
		return nil
	}

	gen := generator.New(finalModel.Config)
	files, err := gen.Generate()
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	fmt.Println("\n✓ Project scaffolded successfully!\n")
	fmt.Println("Created files:")
	for _, f := range files {
		fmt.Printf("  • %s\n", f)
	}
	fmt.Printf("\nGet started:\n  cd %s && go mod tidy && go run .\n", absOut)
	return nil
}

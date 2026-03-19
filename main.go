package main

import (
	"fmt"
	"os"

	"github.com/fabricatorsltd/go-tools-cli/cmd"
	cli "github.com/mirkobrombin/go-cli-builder/v2/pkg/cli"
)

func main() {
	app, err := cli.New(&cmd.RootCmd{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	app.SetName("go-tools-cli")
	app.RootNode.Description = "Bootstrap new Go projects using the go-tools ecosystem.\nRun a subcommand to get started, or pass --help to any command."
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

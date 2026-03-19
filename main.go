package main

import (
	"fmt"
	"os"

	"github.com/fabricatorsltd/go-tools-cli/cmd"
	cli "github.com/mirkobrombin/go-cli-builder/v2/pkg/cli"
)

func main() {
	app := &cmd.RootCmd{}
	if err := cli.Run(app); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

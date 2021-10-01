package main

import (
	cli "github.com/mengdaming/tcr-cli/cmd"
	gui "github.com/mengdaming/tcr-gui/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"log"
)

const docDirectory = "../doc"

func main() {
	genDocFor(cli.GetRootCmd())
	genDocFor(gui.GetRootCmd())
}

func genDocFor(cmd *cobra.Command) {
	if doc.GenMarkdownTree(cmd, docDirectory) != nil {
		log.Fatal(doc.GenMarkdownTree(cmd, docDirectory))
	}
}

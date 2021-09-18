package main

import (
	"github.com/mengdaming/tcr-cli/cmd"

	"github.com/spf13/cobra/doc"

	"log"
)

func main() {
	tcr := cmd.GetRootCmd()
	err := doc.GenMarkdownTree(tcr, "../doc")
	if err != nil {
		log.Fatal(err)
	}
}

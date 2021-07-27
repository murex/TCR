package main

import (
	"github.com/mengdaming/tcr/cmd"

	"github.com/spf13/cobra/doc"

	"log"
)

func main() {
	tcr := cmd.GetRootCmd()
	err := doc.GenMarkdownTree(tcr, "./")
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"testing"
)

func Test_generateMermaidDiagram(t *testing.T) {
	mermaid, _ := generateMermaidDiagram(".")
	fmt.Println(mermaid)
}

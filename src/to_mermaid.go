package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func parseImports(file string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	var imports []string
	for _, s := range f.Imports {
		imports = append(imports, strings.Trim(s.Path.Value, `"`))
	}
	return imports, nil
}

func generateMermaidDiagram(dir string) (string, error) {
	diagram := "graph TD\n"

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == "_test_results" {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(info.Name()) != ".go" {
			return nil
		}

		moduleName := strings.TrimSuffix(info.Name(), ".go")
		imports, err := parseImports(path)
		if err != nil {
			return err
		}

		for _, imp := range imports {
			diagram += fmt.Sprintf("    %s --> %s\n", moduleName, imp)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return diagram, nil
}

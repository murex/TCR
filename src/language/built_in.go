/*
Copyright (c) 2023 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package language

import (
	"embed"
	"path"

	"github.com/murex/tcr/helpers"
)

const builtInDir = "built-in"

// builtInFS is the placeholder for the embedded language configuration filesystem
//
//go:embed "built-in"
var builtInFS embed.FS

func init() {
	loadBuiltInLanguages()
}

func loadBuiltInLanguages() {
	// helpers.SetSimpleTrace(os.Stdout)
	helpers.Trace("Loading built-in languages")
	entries, err := builtInFS.ReadDir(builtInDir)
	if err != nil {
		helpers.Trace("Error loading built-in languages: ", err)
	}
	// Loop on all YAML files in built-in language directory
	for _, entry := range entries {
		err := addBuiltIn(asLanguage(*loadBuiltInLanguage(entry.Name())))
		if err != nil {
			helpers.Trace("Error in ", entry.Name(), ": ", err)
		}
	}
}

func loadBuiltInLanguage(yamlFilename string) *configYAML {
	var languageCfg configYAML

	err := helpers.LoadFromYAMLFile(builtInFS, path.Join(builtInDir, yamlFilename), &languageCfg)
	if err != nil {
		helpers.Trace("Error in ", yamlFilename, ": ", err)
		return nil
	}
	languageCfg.Name = helpers.ExtractNameFromYAMLFilename(yamlFilename)
	return &languageCfg
}

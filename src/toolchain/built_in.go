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

package toolchain

import (
	"embed"
	"github.com/murex/tcr/utils"
	"path"
)

const builtInDir = "built-in"

// builtInFS is the placeholder for the embedded toolchain configuration filesystem
//
//go:embed "built-in"
var builtInFS embed.FS

func init() {
	loadBuiltInToolchains()
}

func loadBuiltInToolchains() {
	// utils.SetSimpleTrace(os.Stdout)
	utils.Trace("Loading built-in toolchains")
	entries, err := builtInFS.ReadDir(builtInDir)
	if err != nil {
		utils.Trace("Error loading built-in toolchains: ", err)
	}
	// Loop on all YAML files in built-in toolchain directory
	for _, entry := range entries {
		err := addBuiltIn(asToolchain(*loadBuiltInToolchain(entry.Name())))
		if err != nil {
			utils.Trace("Error in ", entry.Name(), ": ", err)
		}
	}
}

func loadBuiltInToolchain(yamlFilename string) *configYAML {
	var toolchainCfg configYAML

	err := utils.LoadFromYAMLFile(builtInFS, path.Join(builtInDir, yamlFilename), &toolchainCfg)
	if err != nil {
		utils.Trace("Error in ", yamlFilename, ": ", err)
		return nil
	}
	toolchainCfg.Name = utils.ExtractNameFromYAMLFilename(yamlFilename)
	return &toolchainCfg
}

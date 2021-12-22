/*
Copyright (c) 2021 Murex

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

package config

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

const (
	yamlIndent    = 2
	yamlExtension = "yml"
)

// loadFromYaml loads a structure configuration from a YAML file
func loadFromYaml(filename string, out interface{}) {
	// In case we need to use variables in yaml configuration files:
	// Cf. https://anil.io/blog/symfony/yaml/using-variables-in-yaml-files/
	// Cf. https://pkg.go.dev/os#Expand

	data, err := os.ReadFile(filename)
	if err != nil {
		trace("Error while reading configuration file: ", err)
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		trace("Error while unmarshalling configuration data: ", err)
	}
}

// saveToYaml saves a structure configuration into a YAML file
func saveToYaml(in interface{}, filename string) {
	// First we marshall the data
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(yamlIndent)
	err := yamlEncoder.Encode(&in)
	if err != nil {
		trace("Error while marshalling configuration data: ", err)
	}
	// Then we save it
	err = os.WriteFile(filename, b.Bytes(), 0644) //nolint:gosec // We want people to be able to share this
	if err != nil {
		trace("Error while saving configuration: ", err)
	}
}

func createConfigSubDir(dirPath string, description string) {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		trace("Creating ", description, ": ", dirPath)
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			trace("Error creating ", description, ": ", err)
		}
	}
}

func buildYamlFilePath(dirPath string, name string) string {
	return filepath.Join(dirPath, buildYamlFilename(name))
}

func buildYamlFilename(name string) string {
	return strings.ToLower(name + "." + yamlExtension)
}

func extractNameFromYamlFilename(filename string) string {
	return strings.TrimSuffix(strings.ToLower(filename), "."+yamlExtension)
}

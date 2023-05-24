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

package utils

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

// LoadFromYAML loads a structure configuration from a YAML file
func LoadFromYAML(filename string, out any) {
	// In case we need to use variables in yaml configuration files:
	// Cf. https://anil.io/blog/symfony/yaml/using-variables-in-yaml-files/
	// Cf. https://pkg.go.dev/os#Expand

	data, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		Trace("Error while reading configuration file: ", err)
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		Trace("Error while unmarshalling configuration data: ", err)
	}
}

// SaveToYAML saves a structure configuration into a YAML file
func SaveToYAML(in any, filename string) {
	// First we marshall the data
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(yamlIndent)
	err := yamlEncoder.Encode(&in)
	if err != nil {
		Trace("Error while marshalling configuration data: ", err)
	}
	// Then we save it
	err = os.WriteFile(filename, b.Bytes(), 0644) //nolint:gosec,revive // We want people to be able to share this
	if err != nil {
		Trace("Error while saving configuration: ", err)
	}
}

// CreateConfigSubDir creates a sub-directory in the configuration directory
func CreateConfigSubDir(dirPath string, description string) {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		Trace("Creating ", description, ": ", dirPath)
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			Trace("Error creating ", description, ": ", err)
		}
	}
}

// BuildYAMLFilePath creates a YAML file path
func BuildYAMLFilePath(dirPath string, name string) string {
	return filepath.Join(dirPath, buildYAMLFilename(name))
}

func buildYAMLFilename(name string) string {
	return strings.ToLower(name + "." + yamlExtension)
}

// ExtractNameFromYAMLFilename extracts the name from a YAML file (removing extension)
func ExtractNameFromYAMLFilename(filename string) string {
	return strings.TrimSuffix(strings.ToLower(filename), "."+yamlExtension)
}

// ListYAMLFilesIn lists all YAML files in the provided directory
func ListYAMLFilesIn(dirPath string) (list []string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil || len(entries) == 0 {
		// If we cannot open the directory or if it's empty, we don't go any further
		return nil
	}
	// Loop on all YAML files in the directory
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == "."+yamlExtension {
			list = append(list, entry.Name())
		}
	}
	return list
}

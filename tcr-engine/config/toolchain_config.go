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
	"fmt"
	"github.com/murex/tcr/tcr-engine/toolchain"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

const (
	toolchainDir = "toolchain"
	yamlIndent   = 2
)

var (
	toolchainDirPath string
)

type (
	// ToolchainCommandConfig defines the structure of a toolchain configuration.
	ToolchainCommandConfig struct {
		Os        []string `yaml:"os,flow"`
		Arch      []string `yaml:"arch,flow"`
		Command   string   `yaml:"command"`
		Arguments []string `yaml:"arguments,flow"`
	}

	// ToolchainConfig defines the structure of a toolchain configuration.
	ToolchainConfig struct {
		Name         string
		BuildCommand []ToolchainCommandConfig `yaml:"build"`
		TestCommand  []ToolchainCommandConfig `yaml:"test"`
	}
)

func initToolchainConfig() {
	initToolchainConfigDirPath()
}

func saveToolchainConfigs() {
	createToolchainConfigDir()
	trace("Saving toolchains configuration")
	// Loop on all existing toolchains
	for _, name := range toolchain.Names() {
		trace("- ", name)
		tchn, _ := toolchain.Get(name)
		saveToYaml(asToolchainConfig(tchn), buildYamlFilePath(name))
	}
}

func resetToolchainConfigs() {
	trace("Resetting toolchains configuration")
	// Loop on all existing toolchains
	for _, name := range toolchain.Names() {
		trace("- ", name)
		toolchain.Reset(name)
	}
}

func buildYamlFilePath(name string) string {
	filename := name + "." + configFileType
	return filepath.Join(toolchainDirPath, filename)
}

func asToolchainConfig(tchn *toolchain.Toolchain) ToolchainConfig {
	return ToolchainConfig{
		Name:         tchn.GetName(),
		BuildCommand: asToolchainCommandConfigTable(tchn.BuildCommands),
		TestCommand:  asToolchainCommandConfigTable(tchn.TestCommands),
	}
}

func asToolchainCommandConfigTable(commands []toolchain.Command) []ToolchainCommandConfig {
	var res []ToolchainCommandConfig
	for _, command := range commands {
		res = append(res, asToolchainCommandConfig(command))
	}
	return res
}

func asToolchainCommandConfig(command toolchain.Command) ToolchainCommandConfig {
	return ToolchainCommandConfig{
		Os:        asOsTableConfig(command.Os),
		Arch:      asArchTableConfig(command.Arch),
		Command:   command.Path,
		Arguments: command.Arguments,
	}
}

func asOsTableConfig(osNames []toolchain.OsName) []string {
	var res []string
	for _, osName := range osNames {
		res = append(res, string(osName))
	}
	return res
}

func asArchTableConfig(archNames []toolchain.ArchName) []string {
	var res []string
	for _, archName := range archNames {
		res = append(res, string(archName))
	}
	return res
}

// saveToYaml saves a structure configuration into a YAML file
func saveToYaml(tchn interface{}, filename string) {
	// First we marshall the data
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(yamlIndent)
	err := yamlEncoder.Encode(&tchn)
	if err != nil {
		trace("Error while marshalling configuration data: ", err)
	}
	// Then we save it
	err = os.WriteFile(filename, b.Bytes(), 0644) //nolint:gosec // We want people to be able to share this
	if err != nil {
		trace("Error while saving configuration: ", err)
	}
}

// loadFromYaml loads a structure configuration from a YAML file
func loadFromYaml() ToolchainConfig {
	// TODO see if we need to use variables in yaml configuration files
	// Cf. https://anil.io/blog/symfony/yaml/using-variables-in-yaml-files/
	// Cf. https://pkg.go.dev/os#Expand
	// TODO load data from file
	var data = `
name: dummy
build:
  command: dummy_build_command
  arguments:
  - arg1
  - arg2
test:
  command: dummy_test_command
  arguments:
  - arg3
  - arg4
`

	var tchn ToolchainConfig

	err := yaml.Unmarshal([]byte(data), &tchn)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m:\n%v\n\n", tchn)
	return tchn
}

func initToolchainConfigDirPath() {
	toolchainDirPath = filepath.Join(configDirPath, toolchainDir)
}

func createToolchainConfigDir() {
	_, err := os.Stat(toolchainDirPath)
	if os.IsNotExist(err) {
		trace("Creating TCR toolchain configuration directory: ", toolchainDirPath)
		err := os.MkdirAll(toolchainDirPath, os.ModePerm)
		if err != nil {
			trace("Error creating TCR toolchain configuration directory: ", err)
		}
	}
}

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
	"github.com/murex/tcr/tcr-engine/toolchain"
	"gopkg.in/yaml.v3"
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
	loadToolchainConfigs()
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

func loadToolchainConfigs() {

	entries, err := os.ReadDir(toolchainDirPath)
	if err != nil || len(entries) == 0 {
		// If we cannot open toolchain directory or if it's empty, we don't go any further
		return
	}

	// Loop on all files in toolchain directory
	trace("Loading toolchains configuration")
	for _, entry := range entries {
		if entry.IsDir() {
			break
		}
		trace("- ", entry.Name())
		toolchainCfg := loadFromYaml(filepath.Join(toolchainDirPath, entry.Name()))
		err := toolchain.Register(asToolchain(toolchainCfg))
		if err != nil {
			trace("Error in ", entry.Name(), ": ", err)
		}
	}
}

func asToolchain(toolchainCfg ToolchainConfig) toolchain.Toolchain {
	return toolchain.Toolchain{
		Name:          toolchainCfg.Name,
		BuildCommands: asToolchainCommandTable(toolchainCfg.BuildCommand),
		TestCommands:  asToolchainCommandTable(toolchainCfg.TestCommand),
	}
}

func asToolchainCommandTable(commandsCfg []ToolchainCommandConfig) []toolchain.Command {
	var res []toolchain.Command
	for _, commandCfg := range commandsCfg {
		res = append(res, asToolchainCommand(commandCfg))
	}
	return res
}

func asToolchainCommand(commandCfg ToolchainCommandConfig) toolchain.Command {
	return toolchain.Command{
		Os:        asOsTable(commandCfg.Os),
		Arch:      asArchTable(commandCfg.Arch),
		Path:      commandCfg.Command,
		Arguments: commandCfg.Arguments,
	}
}

func asOsTable(names []string) []toolchain.OsName {
	var res []toolchain.OsName
	for _, name := range names {
		res = append(res, toolchain.OsName(name))
	}
	return res
}

func asArchTable(names []string) []toolchain.ArchName {
	var res []toolchain.ArchName
	for _, name := range names {
		res = append(res, toolchain.ArchName(name))
	}
	return res
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
func loadFromYaml(filename string) ToolchainConfig {
	// In case we need to use variables in yaml configuration files:
	// Cf. https://anil.io/blog/symfony/yaml/using-variables-in-yaml-files/
	// Cf. https://pkg.go.dev/os#Expand

	var tchn ToolchainConfig

	data, err := os.ReadFile(filename)
	if err != nil {
		trace("Error while reading configuration file: ", err)
	}
	if err := yaml.Unmarshal(data, &tchn); err != nil {
		trace("Error while unmarshalling configuration data: ", err)
	}
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

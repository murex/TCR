/*
Copyright (c) 2022 Murex

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
	"github.com/murex/tcr/toolchain"
	"github.com/murex/tcr/utils"
	"path/filepath"
)

const (
	toolchainDir = "toolchain"
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
		Name          string                   `yaml:"-"`
		BuildCommand  []ToolchainCommandConfig `yaml:"build"`
		TestCommand   []ToolchainCommandConfig `yaml:"test"`
		TestResultDir string                   `yaml:"test-result-dir"`
	}
)

func initToolchainConfig() {
	initToolchainConfigDirPath()
	loadToolchainConfigs()
}

func saveToolchainConfigs() {
	createToolchainConfigDir()
	utils.Trace("Saving toolchains configuration")
	// Loop on all existing toolchains
	for _, name := range toolchain.Names() {
		utils.Trace("- ", name)
		saveToolchainConfig(name)
	}
}

func saveToolchainConfig(name string) {
	tchn, _ := toolchain.GetToolchain(name)
	utils.SaveToYAML(asToolchainConfig(tchn), utils.BuildYAMLFilePath(toolchainDirPath, name))
}

// GetToolchainConfigFileList returns the list of toolchain configuration files found in toolchain directory
func GetToolchainConfigFileList() (list []string) {
	return utils.ListYAMLFilesIn(toolchainDirPath)
}

func loadToolchainConfigs() {
	utils.Trace("Loading toolchains configuration")
	// Loop on all YAML files in toolchain directory
	for _, entry := range GetToolchainConfigFileList() {
		err := toolchain.Register(asToolchain(*loadToolchainConfig(entry)))
		if err != nil {
			utils.Trace("Error in ", entry, ": ", err)
		}
	}
}

func loadToolchainConfig(yamlFilename string) *ToolchainConfig {
	var toolchainCfg ToolchainConfig
	utils.LoadFromYAML(filepath.Join(toolchainDirPath, yamlFilename), &toolchainCfg)
	toolchainCfg.Name = utils.ExtractNameFromYAMLFilename(yamlFilename)
	return &toolchainCfg
}

func asToolchain(toolchainCfg ToolchainConfig) *toolchain.Toolchain {
	return toolchain.New(
		toolchainCfg.Name,
		asToolchainCommandTable(toolchainCfg.BuildCommand),
		asToolchainCommandTable(toolchainCfg.TestCommand),
		toolchainCfg.TestResultDir,
	)
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
	utils.Trace("Resetting toolchains configuration")
	// Loop on all existing toolchains
	for _, name := range toolchain.Names() {
		utils.Trace("- ", name)
		toolchain.Reset(name)
	}
}

func asToolchainConfig(tchn toolchain.TchnInterface) ToolchainConfig {
	return ToolchainConfig{
		Name:          tchn.GetName(),
		BuildCommand:  asToolchainCommandConfigTable(tchn.GetBuildCommands()),
		TestCommand:   asToolchainCommandConfigTable(tchn.GetTestCommands()),
		TestResultDir: tchn.GetTestResultDir(),
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

func initToolchainConfigDirPath() {
	toolchainDirPath = filepath.Join(configDirPath, toolchainDir)
}

// GetToolchainConfigDirPath returns the path to the toolchain configuration directory
func GetToolchainConfigDirPath() string {
	return toolchainDirPath
}

func createToolchainConfigDir() {
	utils.CreateConfigSubDir(toolchainDirPath, "TCR toolchain configuration directory")
}

func showToolchainConfigs() {
	utils.Trace("Configured toolchains:")
	entries := GetToolchainConfigFileList()
	if len(entries) == 0 {
		utils.Trace("- none (will use built-in toolchains)")
	}
	for _, entry := range entries {
		loadToolchainConfig(entry).show()
	}
}

func (t ToolchainConfig) show() {
	prefix := "toolchain." + t.Name
	for _, cmd := range t.BuildCommand {
		cmd.show(prefix + ".build")
	}
	for _, cmd := range t.TestCommand {
		cmd.show(prefix + ".test")
	}
	utils.TraceConfigValue(prefix+".test-result-dir", t.TestResultDir)
}

func (c ToolchainCommandConfig) show(prefix string) {
	utils.TraceConfigValue(prefix+".os", c.Os)
	utils.TraceConfigValue(prefix+".arch", c.Arch)
	utils.TraceConfigValue(prefix+".command", c.Command)
	utils.TraceConfigValue(prefix+".args", c.Arguments)
}

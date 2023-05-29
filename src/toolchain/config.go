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
	"github.com/murex/tcr/utils"
	"os"
	"path/filepath"
)

const (
	toolchainDir = "toolchain"
)

var (
	toolchainDirPath string
)

type (
	// commandConfigYAML defines the structure of a toolchain configuration.
	commandConfigYAML struct {
		Os        []string `yaml:"os,flow"`
		Arch      []string `yaml:"arch,flow"`
		Command   string   `yaml:"command"`
		Arguments []string `yaml:"arguments,flow"`
	}

	// configYAML defines the structure of a toolchain configuration.
	configYAML struct {
		Name          string              `yaml:"-"`
		BuildCommand  []commandConfigYAML `yaml:"build"`
		TestCommand   []commandConfigYAML `yaml:"test"`
		TestResultDir string              `yaml:"test-result-dir"`
	}
)

// InitConfig initializes the toolchain configuration
func InitConfig(configDirPath string) {
	initConfigDirPath(configDirPath)
	loadConfigs()
}

// SaveConfigs saves the toolchain configurations
func SaveConfigs() {
	createConfigDir()
	utils.Trace("Saving toolchains configuration")
	// Loop on all existing toolchains
	for _, name := range Names() {
		utils.Trace("- ", name)
		saveConfig(name)
	}
}

func saveConfig(name string) {
	tchn, _ := Get(name)
	utils.SaveToYAMLFile(appFS, asConfig(tchn), utils.BuildYAMLFilePath(toolchainDirPath, name))
}

// GetConfigFileList returns the list of toolchain configuration files found in toolchain directory
func GetConfigFileList() (list []string) {
	return utils.ListYAMLFilesIn(appFS, toolchainDirPath)
}

func loadConfigs() {
	utils.Trace("Loading toolchains configuration")
	// Loop on all YAML files in toolchain directory
	for _, entry := range GetConfigFileList() {
		err := Register(asToolchain(*loadConfig(entry)))
		if err != nil {
			utils.Trace("Error in ", entry, ": ", err)
		}
	}
}

func loadConfig(yamlFilename string) *configYAML {
	var toolchainCfg configYAML
	err := utils.LoadFromYAMLFile(os.DirFS(toolchainDirPath), yamlFilename, &toolchainCfg)
	if err != nil {
		utils.Trace("Error in ", yamlFilename, ": ", err)
		return nil
	}
	toolchainCfg.Name = utils.ExtractNameFromYAMLFilename(yamlFilename)
	return &toolchainCfg
}

func asToolchain(toolchainCfg configYAML) *Toolchain {
	return New(
		toolchainCfg.Name,
		asCommandTable(toolchainCfg.BuildCommand),
		asCommandTable(toolchainCfg.TestCommand),
		toolchainCfg.TestResultDir,
	)
}

func asCommandTable(commandsCfg []commandConfigYAML) []Command {
	var res []Command
	for _, commandCfg := range commandsCfg {
		res = append(res, asCommand(commandCfg))
	}
	return res
}

func asCommand(commandCfg commandConfigYAML) Command {
	return Command{
		Os:        asOsTable(commandCfg.Os),
		Arch:      asArchTable(commandCfg.Arch),
		Path:      commandCfg.Command,
		Arguments: commandCfg.Arguments,
	}
}

func asOsTable(names []string) []OsName {
	var res []OsName
	for _, name := range names {
		res = append(res, OsName(name))
	}
	return res
}

func asArchTable(names []string) []ArchName {
	var res []ArchName
	for _, name := range names {
		res = append(res, ArchName(name))
	}
	return res
}

// ResetConfigs resets the toolchains configuration
func ResetConfigs() {
	utils.Trace("Resetting toolchains configuration")
	// Loop on all existing toolchains
	for _, name := range Names() {
		utils.Trace("- ", name)
		Reset(name)
	}
}

func asConfig(tchn TchnInterface) configYAML {
	return configYAML{
		Name:          tchn.GetName(),
		BuildCommand:  asCommandConfigTable(tchn.GetBuildCommands()),
		TestCommand:   asCommandConfigTable(tchn.GetTestCommands()),
		TestResultDir: tchn.GetTestResultDir(),
	}
}

func asCommandConfigTable(commands []Command) []commandConfigYAML {
	var res []commandConfigYAML
	for _, command := range commands {
		res = append(res, asCommandConfig(command))
	}
	return res
}

func asCommandConfig(command Command) commandConfigYAML {
	return commandConfigYAML{
		Os:        asOsTableConfig(command.Os),
		Arch:      asArchTableConfig(command.Arch),
		Command:   command.Path,
		Arguments: command.Arguments,
	}
}

func asOsTableConfig(osNames []OsName) []string {
	var res []string
	for _, osName := range osNames {
		res = append(res, string(osName))
	}
	return res
}

func asArchTableConfig(archNames []ArchName) []string {
	var res []string
	for _, archName := range archNames {
		res = append(res, string(archName))
	}
	return res
}

func initConfigDirPath(configDirPath string) {
	toolchainDirPath = filepath.Join(configDirPath, toolchainDir)
}

// GetConfigDirPath returns the path to the toolchain configuration directory
func GetConfigDirPath() string {
	return toolchainDirPath
}

func createConfigDir() {
	utils.CreateSubDir(appFS, toolchainDirPath, "TCR toolchain configuration directory")
}

// ShowConfigs shows the toolchains configuration
func ShowConfigs() {
	utils.Trace("Configured toolchains:")
	entries := GetConfigFileList()
	if len(entries) == 0 {
		utils.Trace("- none (will use built-in toolchains)")
	}
	for _, entry := range entries {
		loadConfig(entry).show()
	}
}

func (t configYAML) show() {
	prefix := "toolchain." + t.Name
	for _, cmd := range t.BuildCommand {
		cmd.show(prefix + ".build")
	}
	for _, cmd := range t.TestCommand {
		cmd.show(prefix + ".test")
	}
	utils.TraceKeyValue(prefix+".test-result-dir", t.TestResultDir)
}

func (c commandConfigYAML) show(prefix string) {
	utils.TraceKeyValue(prefix+".os", c.Os)
	utils.TraceKeyValue(prefix+".arch", c.Arch)
	utils.TraceKeyValue(prefix+".command", c.Command)
	utils.TraceKeyValue(prefix+".args", c.Arguments)
}

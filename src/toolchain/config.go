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
	"github.com/murex/tcr/helpers"
	"github.com/murex/tcr/toolchain/command"
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
	helpers.Trace("Saving toolchains configuration")
	// Loop on all existing toolchains
	for _, name := range Names() {
		helpers.Trace("- ", name)
		saveConfig(name)
	}
}

func saveConfig(name string) {
	tchn, _ := Get(name)
	helpers.SaveToYAMLFile(appFS, asConfig(tchn), helpers.BuildYAMLFilePath(toolchainDirPath, name))
}

// GetConfigFileList returns the list of toolchain configuration files found in toolchain directory
func GetConfigFileList() (list []string) {
	return helpers.ListYAMLFilesIn(appFS, toolchainDirPath)
}

func loadConfigs() {
	helpers.Trace("Loading toolchains configuration")
	// Loop on all YAML files in toolchain directory
	for _, entry := range GetConfigFileList() {
		err := Register(asToolchain(*loadConfig(entry)))
		if err != nil {
			helpers.Trace("Error in ", entry, ": ", err)
		}
	}
}

func loadConfig(yamlFilename string) *configYAML {
	var toolchainCfg configYAML
	err := helpers.LoadFromYAMLFile(os.DirFS(toolchainDirPath), yamlFilename, &toolchainCfg)
	if err != nil {
		helpers.Trace("Error in ", yamlFilename, ": ", err)
		return nil
	}
	toolchainCfg.Name = helpers.ExtractNameFromYAMLFilename(yamlFilename)
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

func asCommandTable(commandsCfg []commandConfigYAML) []command.Command {
	var res []command.Command
	for _, commandCfg := range commandsCfg {
		res = append(res, asCommand(commandCfg))
	}
	return res
}

func asCommand(commandCfg commandConfigYAML) command.Command {
	return command.Command{
		Os:        asOsTable(commandCfg.Os),
		Arch:      asArchTable(commandCfg.Arch),
		Path:      commandCfg.Command,
		Arguments: commandCfg.Arguments,
	}
}

func asOsTable(names []string) []command.OsName {
	var res []command.OsName
	for _, name := range names {
		res = append(res, command.OsName(name))
	}
	return res
}

func asArchTable(names []string) []command.ArchName {
	var res []command.ArchName
	for _, name := range names {
		res = append(res, command.ArchName(name))
	}
	return res
}

// ResetConfigs resets the toolchains configuration
func ResetConfigs() {
	helpers.Trace("Resetting toolchains configuration")
	// Loop on all existing toolchains
	for _, name := range Names() {
		helpers.Trace("- ", name)
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

func asCommandConfigTable(commands []command.Command) []commandConfigYAML {
	var res []commandConfigYAML
	for _, c := range commands {
		res = append(res, asCommandConfig(c))
	}
	return res
}

func asCommandConfig(cmd command.Command) commandConfigYAML {
	return commandConfigYAML{
		Os:        asOsTableConfig(cmd.Os),
		Arch:      asArchTableConfig(cmd.Arch),
		Command:   cmd.Path,
		Arguments: cmd.Arguments,
	}
}

func asOsTableConfig(osNames []command.OsName) []string {
	var res []string
	for _, osName := range osNames {
		res = append(res, string(osName))
	}
	return res
}

func asArchTableConfig(archNames []command.ArchName) []string {
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
	helpers.CreateSubDir(appFS, toolchainDirPath, "TCR toolchain configuration directory")
}

// ShowConfigs shows the toolchains configuration
func ShowConfigs() {
	helpers.Trace("Configured toolchains:")
	entries := GetConfigFileList()
	if len(entries) == 0 {
		helpers.Trace("- none (will use built-in toolchains)")
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
	helpers.TraceKeyValue(prefix+".test-result-dir", t.TestResultDir)
}

func (c commandConfigYAML) show(prefix string) {
	helpers.TraceKeyValue(prefix+".os", c.Os)
	helpers.TraceKeyValue(prefix+".arch", c.Arch)
	helpers.TraceKeyValue(prefix+".command", c.Command)
	helpers.TraceKeyValue(prefix+".args", c.Arguments)
}

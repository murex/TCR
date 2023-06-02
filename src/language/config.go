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
	"github.com/murex/tcr/utils"
	"os"
	"path/filepath"
)

const (
	languageDir = "language"
)

var (
	languageDirPath string
)

type (
	// toolchainConfigYAML defines the structure for toolchains configuration related to a language
	toolchainConfigYAML struct {
		Default    string   `yaml:"default"`
		Compatible []string `yaml:"compatible-with,flow"`
	}

	// fileTreeFilterConfigYAML defines the structure for file tree filtering configuration related to a language
	fileTreeFilterConfigYAML struct {
		Directories  []string `yaml:"directories,flow"`
		FilePatterns []string `yaml:"patterns,flow"`
	}

	// configYAML defines the structure of a language configuration.
	configYAML struct {
		Name        string                   `yaml:"-"`
		Toolchains  toolchainConfigYAML      `yaml:"toolchains"`
		SourceFiles fileTreeFilterConfigYAML `yaml:"source-files"`
		TestFiles   fileTreeFilterConfigYAML `yaml:"test-files"`
	}
)

// InitConfig initializes the language configuration
func InitConfig(configDirPath string) {
	initConfigDirPath(configDirPath)
	loadConfigs()
}

// SaveConfigs saves the language configurations
func SaveConfigs() {
	createConfigDir()
	utils.Trace("Saving languages configuration")
	// Loop on all existing languages
	for _, name := range Names() {
		utils.Trace("- ", name)
		saveConfig(name)
	}
}

func saveConfig(name string) {
	lang, _ := Get(name)
	utils.SaveToYAMLFile(appFS, asConfig(lang), utils.BuildYAMLFilePath(languageDirPath, name))
}

// GetConfigFileList returns the list of language configuration files found in language directory
func GetConfigFileList() (list []string) {
	return utils.ListYAMLFilesIn(appFS, languageDirPath)
}

func loadConfigs() {
	utils.Trace("Loading languages configuration")
	// Loop on all YAML files in language directory
	for _, entry := range GetConfigFileList() {
		err := Register(asLanguage(*loadConfig(entry)))
		if err != nil {
			utils.Trace("Error in ", entry, ": ", err)
		}
	}
}

func loadConfig(yamlFilename string) *configYAML {
	var languageCfg configYAML
	err := utils.LoadFromYAMLFile(os.DirFS(languageDirPath), yamlFilename, &languageCfg)
	if err != nil {
		utils.Trace("Error in ", yamlFilename, ": ", err)
		return nil
	}
	languageCfg.Name = utils.ExtractNameFromYAMLFilename(yamlFilename)
	return &languageCfg
}

func asLanguage(languageCfg configYAML) *Language {
	return New(
		languageCfg.Name,
		asToolchains(languageCfg.Toolchains),
		asFileTreeFilter(languageCfg.SourceFiles),
		asFileTreeFilter(languageCfg.TestFiles),
	)
}

func asFileTreeFilter(filesCfg fileTreeFilterConfigYAML) FileTreeFilter {
	return FileTreeFilter{
		Directories:  asDirectoryTable(filesCfg.Directories),
		FilePatterns: asFilePatternTable(filesCfg.FilePatterns),
	}
}

func asToolchains(toolchainsCfg toolchainConfigYAML) Toolchains {
	return Toolchains{
		Default:    toolchainsCfg.Default,
		Compatible: asToolchainTable(toolchainsCfg.Compatible),
	}
}

func asToolchainTable(toolchainsCfg []string) []string {
	return append([]string{}, toolchainsCfg...)
}

func asDirectoryTable(directoryTableCfg []string) []string {
	return append([]string{}, directoryTableCfg...)
}

func asFilePatternTable(filePatternTableCfg []string) []string {
	return append([]string{}, filePatternTableCfg...)
}

// ResetConfigs resets the languages configuration
func ResetConfigs() {
	utils.Trace("Resetting languages configuration")
	// Loop on all existing languages
	for _, name := range Names() {
		utils.Trace("- ", name)
		Reset(name)
	}
}

func asConfig(lang LangInterface) configYAML {
	return configYAML{
		Name:        lang.GetName(),
		Toolchains:  asToolchainsConfig(lang.GetToolchains()),
		SourceFiles: asFileTreeFilterConfig(lang.GetSrcFileFilter()),
		TestFiles:   asFileTreeFilterConfig(lang.GetTestFileFilter()),
	}
}

func asFileTreeFilterConfig(files FileTreeFilter) fileTreeFilterConfigYAML {
	return fileTreeFilterConfigYAML{
		Directories:  asDirectoryTableConfig(files.Directories),
		FilePatterns: asFilePatternTableConfig(files.FilePatterns),
	}
}

func asToolchainsTableConfig(toolchains []string) []string {
	return append([]string(nil), toolchains...)
}

func asToolchainsConfig(toolchains Toolchains) toolchainConfigYAML {
	return toolchainConfigYAML{
		Default:    toolchains.Default,
		Compatible: asToolchainsTableConfig(toolchains.Compatible),
	}
}

func asDirectoryTableConfig(directoryTable []string) []string {
	return append([]string(nil), directoryTable...)
}

func asFilePatternTableConfig(filePatternTable []string) []string {
	return append([]string(nil), filePatternTable...)
}

func initConfigDirPath(configDirPath string) {
	languageDirPath = filepath.Join(configDirPath, languageDir)
}

// GetConfigDirPath returns the path to the language configuration directory
func GetConfigDirPath() string {
	return languageDirPath
}

func createConfigDir() {
	utils.CreateSubDir(appFS, languageDirPath, "TCR language configuration directory")
}

// ShowConfigs shows the languages configuration
func ShowConfigs() {
	utils.Trace("Configured languages:")
	entries := GetConfigFileList()
	if len(entries) == 0 {
		utils.Trace("- none (will use built-in languages)")
	}
	for _, entry := range entries {
		loadConfig(entry).show()
	}
}

func (l configYAML) show() {
	prefix := "language." + l.Name
	l.Toolchains.show(prefix + ".toolchains")
	l.SourceFiles.show(prefix + ".source-files")
	l.TestFiles.show(prefix + ".test-files")
}

func (lt toolchainConfigYAML) show(prefix string) {
	utils.TraceKeyValue(prefix+".default", lt.Default)
	utils.TraceKeyValue(prefix+".compatible-with", lt.Compatible)
}

func (ftf fileTreeFilterConfigYAML) show(prefix string) {
	utils.TraceKeyValue(prefix+".directories", ftf.Directories)
	utils.TraceKeyValue(prefix+".patterns", ftf.FilePatterns)
}

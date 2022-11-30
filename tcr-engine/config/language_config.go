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
	"github.com/murex/tcr/language"
	"path/filepath"
)

const (
	languageDir = "language"
)

var (
	languageDirPath string
)

type (
	// LanguageToolchainConfig defines the structure for toolchains configuration related to a language
	LanguageToolchainConfig struct {
		Default    string   `yaml:"default"`
		Compatible []string `yaml:"compatible-with,flow"`
	}

	// LanguageFileTreeFilterConfig defines the structure for file tree filtering configuration related to a language
	LanguageFileTreeFilterConfig struct {
		Directories  []string `yaml:"directories,flow"`
		FilePatterns []string `yaml:"patterns,flow"`
	}

	// LanguageConfig defines the structure of a language configuration.
	LanguageConfig struct {
		Name        string                       `yaml:"-"`
		Toolchains  LanguageToolchainConfig      `yaml:"toolchains"`
		SourceFiles LanguageFileTreeFilterConfig `yaml:"source-files"`
		TestFiles   LanguageFileTreeFilterConfig `yaml:"test-files"`
	}
)

func initLanguageConfig() {
	initLanguageConfigDirPath()
	loadLanguageConfigs()
}

func saveLanguageConfigs() {
	createLanguageConfigDir()
	trace("Saving languages configuration")
	// Loop on all existing languages
	for _, name := range language.Names() {
		trace("- ", name)
		saveLanguageConfig(name)
	}
}

func saveLanguageConfig(name string) {
	lang, _ := language.Get(name)
	saveToYaml(asLanguageConfig(lang), buildYamlFilePath(languageDirPath, name))
}

// GetLanguageConfigFileList returns the list of language configuration files found in language directory
func GetLanguageConfigFileList() (list []string) {
	return listYamlFilesIn(languageDirPath)
}

func loadLanguageConfigs() {
	trace("Loading languages configuration")
	// Loop on all YAML files in language directory
	for _, entry := range GetLanguageConfigFileList() {
		err := language.Register(asLanguage(*loadLanguageConfig(entry)))
		if err != nil {
			trace("Error in ", entry, ": ", err)
		}
	}
}

func loadLanguageConfig(yamlFilename string) *LanguageConfig {
	var languageCfg LanguageConfig
	loadFromYaml(filepath.Join(languageDirPath, yamlFilename), &languageCfg)
	languageCfg.Name = extractNameFromYamlFilename(yamlFilename)
	return &languageCfg
}

func asLanguage(languageCfg LanguageConfig) *language.Language {
	return language.New(
		languageCfg.Name,
		asLanguageToolchains(languageCfg.Toolchains),
		asLanguageFileTreeFilter(languageCfg.SourceFiles),
		asLanguageFileTreeFilter(languageCfg.TestFiles),
	)
}

func asLanguageFileTreeFilter(filesCfg LanguageFileTreeFilterConfig) language.FileTreeFilter {
	return language.FileTreeFilter{
		Directories:  asLanguageDirectoryTable(filesCfg.Directories),
		FilePatterns: asLanguageFilePatternTable(filesCfg.FilePatterns),
	}
}

func asLanguageToolchains(toolchainsCfg LanguageToolchainConfig) language.Toolchains {
	return language.Toolchains{
		Default:    toolchainsCfg.Default,
		Compatible: asLanguageToolchainTable(toolchainsCfg.Compatible),
	}
}

func asLanguageToolchainTable(toolchainsCfg []string) []string {
	return append([]string(nil), toolchainsCfg...)
}

func asLanguageDirectoryTable(directoryTableCfg []string) []string {
	return append([]string(nil), directoryTableCfg...)
}

func asLanguageFilePatternTable(filePatternTableCfg []string) []string {
	return append([]string(nil), filePatternTableCfg...)
}

func resetLanguageConfigs() {
	trace("Resetting languages configuration")
	// Loop on all existing languages
	for _, name := range language.Names() {
		trace("- ", name)
		language.Reset(name)
	}
}

func asLanguageConfig(lang language.LangInterface) LanguageConfig {
	return LanguageConfig{
		Name:        lang.GetName(),
		Toolchains:  asLanguageToolchainsConfig(lang.GetToolchains()),
		SourceFiles: asLanguageFileTreeFilterConfig(lang.GetSrcFileFilter()),
		TestFiles:   asLanguageFileTreeFilterConfig(lang.GetTestFileFilter()),
	}
}

func asLanguageFileTreeFilterConfig(files language.FileTreeFilter) LanguageFileTreeFilterConfig {
	return LanguageFileTreeFilterConfig{
		Directories:  asLanguageDirectoryTableConfig(files.Directories),
		FilePatterns: asLanguageFilePatternTableConfig(files.FilePatterns),
	}
}

func asLanguageToolchainsTableConfig(toolchains []string) []string {
	return append([]string(nil), toolchains...)
}

func asLanguageToolchainsConfig(toolchains language.Toolchains) LanguageToolchainConfig {
	return LanguageToolchainConfig{
		Default:    toolchains.Default,
		Compatible: asLanguageToolchainsTableConfig(toolchains.Compatible),
	}
}

func asLanguageDirectoryTableConfig(directoryTable []string) []string {
	return append([]string(nil), directoryTable...)
}

func asLanguageFilePatternTableConfig(filePatternTable []string) []string {
	return append([]string(nil), filePatternTable...)
}

func initLanguageConfigDirPath() {
	languageDirPath = filepath.Join(configDirPath, languageDir)
}

func createLanguageConfigDir() {
	createConfigSubDir(languageDirPath, "TCR language configuration directory")
}

// GetLanguageConfigDirPath returns the path to the language configuration directory
func GetLanguageConfigDirPath() string {
	return languageDirPath
}

func showLanguageConfigs() {
	trace("Configured languages:")
	entries := GetLanguageConfigFileList()
	if len(entries) == 0 {
		trace("- none (will use built-in languages)")
	}
	for _, entry := range entries {
		loadLanguageConfig(entry).show()
	}
}

func (l LanguageConfig) show() {
	prefix := "language." + l.Name
	l.Toolchains.show(prefix + ".toolchains")
	l.SourceFiles.show(prefix + ".source-files")
	l.TestFiles.show(prefix + ".test-files")
}

func (lt LanguageToolchainConfig) show(prefix string) {
	showConfigValue(prefix+".default", lt.Default)
	showConfigValue(prefix+".compatible-with", lt.Compatible)
}

func (ftf LanguageFileTreeFilterConfig) show(prefix string) {
	showConfigValue(prefix+".directories", ftf.Directories)
	showConfigValue(prefix+".patterns", ftf.FilePatterns)
}

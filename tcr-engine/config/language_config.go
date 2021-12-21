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
	"github.com/murex/tcr/tcr-engine/language"
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
	// LanguageToolchainConfig defines the structure for toolchains related to a language
	LanguageToolchainConfig struct {
		Default    string   `yaml:"default"`
		Compatible []string `yaml:"compatible-with,flow"`
	}

	// LanguageFilesConfig defines the structure for files and directories related to a language
	LanguageFilesConfig struct {
		Directories []string `yaml:"directories,flow"`
		Filters     []string `yaml:"name-filters,flow"`
	}

	// LanguageConfig defines the structure of a language configuration.
	LanguageConfig struct {
		Name        string                  `yaml:"name"`
		Toolchains  LanguageToolchainConfig `yaml:"toolchains"`
		SourceFiles LanguageFilesConfig     `yaml:"source-files"`
		TestFiles   LanguageFilesConfig     `yaml:"test-files"`
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
		lang, _ := language.Get(name)
		saveToYaml(asLanguageConfig(lang), buildYamlFilePath(languageDirPath, name))
	}
}

func loadLanguageConfigs() {
	entries, err := os.ReadDir(languageDirPath)
	if err != nil || len(entries) == 0 {
		// If we cannot open language directory or if it's empty, we don't go any further
		return
	}

	// Loop on all files in language directory
	trace("Loading languages configuration")
	for _, entry := range entries {
		if entry.IsDir() {
			break
		}
		trace("- ", entry.Name())
		var languageCfg LanguageConfig
		loadFromYaml(filepath.Join(languageDirPath, entry.Name()), &languageCfg)
		err := language.Register(asLanguage(languageCfg))
		if err != nil {
			trace("Error in ", entry.Name(), ": ", err)
		}
	}
}

func asLanguage(languageCfg LanguageConfig) language.Language {
	return language.Language{
		Name:       languageCfg.Name,
		Toolchains: asLanguageToolchains(languageCfg.Toolchains),
		SrcFiles:   asLanguageFiles(languageCfg.SourceFiles),
		TestFiles:  asLanguageFiles(languageCfg.TestFiles),
	}
}

func asLanguageFiles(filesCfg LanguageFilesConfig) language.Files {
	return language.Files{
		Directories: asLanguageDirectoryTable(filesCfg.Directories),
		Filters:     asLanguageFilterTable(filesCfg.Filters),
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

func asLanguageFilterTable(fileFilterTableCfg []string) []string {
	return append([]string(nil), fileFilterTableCfg...)
}

func resetLanguageConfigs() {
	trace("Resetting languages configuration")
	// Loop on all existing languages
	for _, name := range language.Names() {
		trace("- ", name)
		language.Reset(name)
	}
}

func asLanguageConfig(lang *language.Language) LanguageConfig {
	return LanguageConfig{
		Name:        lang.GetName(),
		Toolchains:  asLanguageToolchainsConfig(lang.Toolchains),
		SourceFiles: asLanguageFilesConfig(lang.SrcFiles),
		TestFiles:   asLanguageFilesConfig(lang.TestFiles),
	}
}

func asLanguageFilesConfig(files language.Files) LanguageFilesConfig {
	return LanguageFilesConfig{
		Directories: asLanguageDirectoryTableConfig(files.Directories),
		Filters:     asLanguageFilterTableConfig(files.Filters),
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

func asLanguageFilterTableConfig(fileFilterTable []string) []string {
	return append([]string(nil), fileFilterTable...)
}

func initLanguageConfigDirPath() {
	languageDirPath = filepath.Join(configDirPath, languageDir)
}

func createLanguageConfigDir() {
	createConfigSubDir(languageDirPath, "TCR language configuration directory")
}

func showLanguageConfigs() {
	// TODO Implement display of languages configuration
}

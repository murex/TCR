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
	"fmt"
	"github.com/murex/tcr/tcr-engine/language"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_can_save_language_configuration(t *testing.T) {
	// TODO bypass filesystem
	//saveToYaml(lang, "")
}

func Test_convert_language_name_to_config(t *testing.T) {
	lang := language.ALanguage()
	cfg := asLanguageConfig(lang)
	assert.Equal(t, lang.GetName(), cfg.Name)
	assert.Equal(t, cfg.Name, asLanguage(cfg).GetName())
}

func Test_convert_language_default_toolchain_to_config(t *testing.T) {
	lang := language.ALanguage()
	cfg := asLanguageConfig(lang)
	assert.Equal(t, lang.GetToolchains().Default, cfg.Toolchains.Default)
	assert.Equal(t, cfg.Toolchains.Default, asLanguage(cfg).GetToolchains().Default)
}

func Test_convert_language_compatible_toolchains_to_config(t *testing.T) {
	lang := language.ALanguage()
	cfg := asLanguageConfig(lang)
	assert.Equal(t, lang.GetToolchains().Compatible, cfg.Toolchains.Compatible)
	assert.Equal(t, cfg.Toolchains.Compatible, asLanguage(cfg).GetToolchains().Compatible)
}

func Test_convert_language_src_filter_directories_to_config(t *testing.T) {
	lang := language.ALanguage(
		language.WithSrcFiles(
			language.AFileTreeFilter(language.WithDirectories("src-dir1", "src-dir2")),
		),
	)
	cfg := asLanguageConfig(lang)
	assert.Equal(t, lang.GetSrcFileFilter().Directories, cfg.SourceFiles.Directories)
	assert.Equal(t, cfg.SourceFiles.Directories, asLanguage(cfg).GetSrcFileFilter().Directories)
}

func Test_convert_language_src_filter_patterns_to_config(t *testing.T) {
	lang := language.ALanguage(
		language.WithSrcFiles(
			language.AFileTreeFilter(language.WithPatterns("src-pattern1", "src-pattern2")),
		),
	)
	cfg := asLanguageConfig(lang)
	assert.Equal(t, lang.GetSrcFileFilter().FilePatterns, cfg.SourceFiles.FilePatterns)
	assert.Equal(t, cfg.SourceFiles.FilePatterns, asLanguage(cfg).GetSrcFileFilter().FilePatterns)
}

func Test_convert_language_test_filter_directories_to_config(t *testing.T) {
	lang := language.ALanguage(
		language.WithTestFiles(
			language.AFileTreeFilter(language.WithDirectories("test-dir1", "test-dir2")),
		),
	)
	cfg := asLanguageConfig(lang)
	assert.Equal(t, lang.GetTestFileFilter().Directories, cfg.TestFiles.Directories)
	assert.Equal(t, cfg.TestFiles.Directories, asLanguage(cfg).GetTestFileFilter().Directories)
}

func Test_convert_language_test_filter_patterns_to_config(t *testing.T) {
	lang := language.ALanguage(
		language.WithTestFiles(
			language.AFileTreeFilter(language.WithPatterns("test-pattern1", "test-pattern2")),
		),
	)
	cfg := asLanguageConfig(lang)
	assert.Equal(t, lang.GetTestFileFilter().FilePatterns, cfg.TestFiles.FilePatterns)
	assert.Equal(t, cfg.TestFiles.FilePatterns, asLanguage(cfg).GetTestFileFilter().FilePatterns)
}

func Test_show_language_configs_with_no_saved_config(t *testing.T) {
	expected := []string{
		"Configured languages:",
		"- none (will use built-in languages)",
	}
	assertConfigTrace(t, expected,
		func() {
			languageDirPath = ""
			showLanguageConfigs()
		},
	)
}

func Test_reset_language_configs_with_no_saved_config(t *testing.T) {
	expected := []string{
		"Resetting languages configuration",
	}
	for _, builtin := range language.Names() {
		expected = append(expected, "- "+builtin)
	}
	assertConfigTrace(t, expected,
		func() {
			resetLanguageConfigs()
		},
	)
}

func Test_show_language_config(t *testing.T) {
	lang := language.ALanguage()
	cfg := asLanguageConfig(lang)
	prefix := "- language." + cfg.Name
	expected := []string{
		fmt.Sprintf("%v.toolchains.default: %v", prefix, cfg.Toolchains.Default),
		fmt.Sprintf("%v.toolchains.compatible-with: %v", prefix, cfg.Toolchains.Compatible),
		fmt.Sprintf("%v.source-files.directories: %v", prefix, cfg.SourceFiles.Directories),
		fmt.Sprintf("%v.source-files.patterns: %v", prefix, cfg.SourceFiles.FilePatterns),
		fmt.Sprintf("%v.test-files.directories: %v", prefix, cfg.TestFiles.Directories),
		fmt.Sprintf("%v.test-files.patterns: %v", prefix, cfg.TestFiles.FilePatterns),
	}
	assertConfigTrace(t, expected,
		func() {
			cfg.show()
		},
	)
}

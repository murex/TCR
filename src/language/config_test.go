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

package language

import (
	"fmt"
	"github.com/murex/tcr/utils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_convert_language_name_to_config(t *testing.T) {
	lang := ALanguage()
	cfg := asConfig(lang)
	assert.Equal(t, lang.GetName(), cfg.Name)
	assert.Equal(t, cfg.Name, asLanguage(cfg).GetName())
}

func Test_convert_language_default_toolchain_to_config(t *testing.T) {
	lang := ALanguage()
	cfg := asConfig(lang)
	assert.Equal(t, lang.GetToolchains().Default, cfg.Toolchains.Default)
	assert.Equal(t, cfg.Toolchains.Default, asLanguage(cfg).GetToolchains().Default)
}

func Test_convert_language_compatible_toolchains_to_config(t *testing.T) {
	lang := ALanguage()
	cfg := asConfig(lang)
	assert.Equal(t, lang.GetToolchains().Compatible, cfg.Toolchains.Compatible)
	assert.Equal(t, cfg.Toolchains.Compatible, asLanguage(cfg).GetToolchains().Compatible)
}

func Test_convert_language_src_filter_directories_to_config(t *testing.T) {
	lang := ALanguage(
		WithSrcFiles(
			AFileTreeFilter(WithDirectories("src-dir1", "src-dir2")),
		),
	)
	cfg := asConfig(lang)
	assert.Equal(t, lang.GetSrcFileFilter().Directories, cfg.SourceFiles.Directories)
	assert.Equal(t, cfg.SourceFiles.Directories, asLanguage(cfg).GetSrcFileFilter().Directories)
}

func Test_convert_language_src_filter_patterns_to_config(t *testing.T) {
	lang := ALanguage(
		WithSrcFiles(
			AFileTreeFilter(WithPatterns("src-pattern1", "src-pattern2")),
		),
	)
	cfg := asConfig(lang)
	assert.Equal(t, lang.GetSrcFileFilter().FilePatterns, cfg.SourceFiles.FilePatterns)
	assert.Equal(t, cfg.SourceFiles.FilePatterns, asLanguage(cfg).GetSrcFileFilter().FilePatterns)
}

func Test_convert_language_test_filter_directories_to_config(t *testing.T) {
	lang := ALanguage(
		WithTestFiles(
			AFileTreeFilter(WithDirectories("test-dir1", "test-dir2")),
		),
	)
	cfg := asConfig(lang)
	assert.Equal(t, lang.GetTestFileFilter().Directories, cfg.TestFiles.Directories)
	assert.Equal(t, cfg.TestFiles.Directories, asLanguage(cfg).GetTestFileFilter().Directories)
}

func Test_convert_language_test_filter_patterns_to_config(t *testing.T) {
	lang := ALanguage(
		WithTestFiles(
			AFileTreeFilter(WithPatterns("test-pattern1", "test-pattern2")),
		),
	)
	cfg := asConfig(lang)
	assert.Equal(t, lang.GetTestFileFilter().FilePatterns, cfg.TestFiles.FilePatterns)
	assert.Equal(t, cfg.TestFiles.FilePatterns, asLanguage(cfg).GetTestFileFilter().FilePatterns)
}

func Test_show_language_configs_with_no_saved_config(t *testing.T) {
	expected := []string{
		"Configured languages:",
		"- none (will use built-in languages)",
	}
	utils.AssertSimpleTrace(t, expected,
		func() {
			languageDirPath = ""
			ShowConfigs()
		},
	)
}

func Test_reset_language_configs_with_no_saved_config(t *testing.T) {
	expected := []string{
		"Resetting languages configuration",
	}
	for _, builtin := range Names() {
		expected = append(expected, "- "+builtin)
	}
	utils.AssertSimpleTrace(t, expected,
		func() {
			ResetConfigs()
		},
	)
}

func Test_show_language_config(t *testing.T) {
	lang := ALanguage()
	cfg := asConfig(lang)
	prefix := "- language." + cfg.Name
	expected := []string{
		fmt.Sprintf("%v.toolchains.default: %v", prefix, cfg.Toolchains.Default),
		fmt.Sprintf("%v.toolchains.compatible-with: %v", prefix, cfg.Toolchains.Compatible),
		fmt.Sprintf("%v.source-files.directories: %v", prefix, cfg.SourceFiles.Directories),
		fmt.Sprintf("%v.source-files.patterns: %v", prefix, cfg.SourceFiles.FilePatterns),
		fmt.Sprintf("%v.test-files.directories: %v", prefix, cfg.TestFiles.Directories),
		fmt.Sprintf("%v.test-files.patterns: %v", prefix, cfg.TestFiles.FilePatterns),
	}
	utils.AssertSimpleTrace(t, expected,
		func() {
			cfg.show()
		},
	)
}

func Test_save_and_load_a_language_config(t *testing.T) {
	const name = "my-language"
	lang := ALanguage(WithName(name))
	errRegister := Register(lang)
	if errRegister != nil {
		t.Fatal(errRegister)
	}

	// Set up a temporary directory
	appFS = afero.NewOsFs()
	dir, errTempDir := os.MkdirTemp("", "tcr-language")
	if errTempDir != nil {
		t.Fatal(errTempDir)
	}
	defer t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})
	// Prepare language configuration directory
	initConfigDirPath(dir)
	createConfigDir()

	// Save the language configuration file
	saveConfig(name)

	// Load the language configuration file
	yamlConfig := loadConfig(utils.BuildYAMLFilename(name))

	// Check that we get back the same language data
	if assert.NotNil(t, yamlConfig) {
		assert.Equal(t, lang, asLanguage(*yamlConfig))
	}
}

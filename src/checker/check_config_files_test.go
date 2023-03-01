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

package checker

import (
	"github.com/murex/tcr/checker/model"
	"github.com/murex/tcr/params"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_check_config_files(t *testing.T) {
	assertCheckGroupRunner(t,
		checkConfigFiles,
		&checkConfigRunners,
		*params.AParamSet(),
		"configuration files")
}

func Test_check_config_directory(t *testing.T) {
	currentDir, _ := filepath.Abs(".")
	tests := []struct {
		desc     string
		value    string
		expected []model.CheckPoint
	}{
		{
			"not set", "",
			[]model.CheckPoint{
				model.OkCheckPoint("configuration directory parameter is not set explicitly"),
				model.OkCheckPoint("using current directory as configuration directory"),
				model.OkCheckPoint("TCR configuration root directory is ", currentDir),
			},
		},
		{
			"set", ".",
			[]model.CheckPoint{
				model.OkCheckPoint("configuration directory parameter is ."),
				model.OkCheckPoint("TCR configuration root directory is ", currentDir),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithConfigDir(test.value))
			initTestCheckEnv(p)
			assert.Equal(t, test.expected, checkConfigDirectory(p))
		})
	}
}

func Test_check_language_config(t *testing.T) {
	currentDir, _ := filepath.Abs(".")
	tests := []struct {
		desc     string
		dir      string
		fileList []string
		expected []model.CheckPoint
	}{
		{
			"0 language config file", "language", []string{},
			[]model.CheckPoint{
				model.OkCheckPoint("language configuration directory is ", filepath.Join(currentDir, "language")),
				model.WarningCheckPoint("no language configuration file found"),
			},
		},
		{
			"1 language config file", "language", []string{"go.yml"},
			[]model.CheckPoint{
				model.OkCheckPoint("language configuration directory is ", filepath.Join(currentDir, "language")),
				model.OkCheckPoint("language configuration files:"),
				model.OkCheckPoint("- go.yml"),
			},
		},
		{
			"2 language config files", "language", []string{"java.yml", "cpp.yml"},
			[]model.CheckPoint{
				model.OkCheckPoint("language configuration directory is ", filepath.Join(currentDir, "language")),
				model.OkCheckPoint("language configuration files:"),
				model.OkCheckPoint("- java.yml"),
				model.OkCheckPoint("- cpp.yml"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			saved := languageConfigDir
			t.Cleanup(func() { languageConfigDir = saved })
			languageConfigDir = configSubDir{
				name:        "language",
				getDirPath:  func() string { return test.dir },
				getFileList: func() []string { return test.fileList },
			}
			assert.Equal(t, test.expected, checkLanguageConfig(*params.AParamSet()))
		})
	}
}

func Test_check_toolchain_config(t *testing.T) {
	currentDir, _ := filepath.Abs(".")
	tests := []struct {
		desc     string
		dir      string
		fileList []string
		expected []model.CheckPoint
	}{
		{
			"0 toolchain config file", "toolchain", []string{},
			[]model.CheckPoint{
				model.OkCheckPoint("toolchain configuration directory is ", filepath.Join(currentDir, "toolchain")),
				model.WarningCheckPoint("no toolchain configuration file found"),
			},
		},
		{
			"1 toolchain config file", "toolchain", []string{"go-tools.yml"},
			[]model.CheckPoint{
				model.OkCheckPoint("toolchain configuration directory is ", filepath.Join(currentDir, "toolchain")),
				model.OkCheckPoint("toolchain configuration files:"),
				model.OkCheckPoint("- go-tools.yml"),
			},
		},
		{
			"2 toolchain config files", "toolchain", []string{"gradle.yml", "maven.yml"},
			[]model.CheckPoint{
				model.OkCheckPoint("toolchain configuration directory is ", filepath.Join(currentDir, "toolchain")),
				model.OkCheckPoint("toolchain configuration files:"),
				model.OkCheckPoint("- gradle.yml"),
				model.OkCheckPoint("- maven.yml"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			saved := toolchainConfigDir
			t.Cleanup(func() { toolchainConfigDir = saved })
			toolchainConfigDir = configSubDir{
				name:        "toolchain",
				getDirPath:  func() string { return test.dir },
				getFileList: func() []string { return test.fileList },
			}
			assert.Equal(t, test.expected, checkToolchainConfig(*params.AParamSet()))
		})
	}
}

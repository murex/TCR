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
	"errors"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/murex/tcr/checker/model"
	"github.com/murex/tcr/language"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/toolchain"
	"github.com/murex/tcr/toolchain/command"
	"github.com/stretchr/testify/assert"
)

func Test_check_toolchain(t *testing.T) {
	assertCheckGroupRunner(t,
		checkToolchain,
		&checkToolchainRunners,
		*params.AParamSet(),
		"toolchain")
}

func Test_check_toolchain_parameter(t *testing.T) {
	tests := []struct {
		desc     string
		tchn     string
		tchnErr  error
		expected []model.CheckPoint
	}{
		{"not set", "", nil, nil},
		{
			"set and valid", "valid-toolchain", nil,
			[]model.CheckPoint{
				model.OkCheckPoint("toolchain parameter is set to valid-toolchain"),
				model.OkCheckPoint("valid-toolchain toolchain is valid"),
			},
		},
		{
			"set but invalid", "wrong-toolchain", errors.New("toolchain parameter error"),
			[]model.CheckPoint{
				model.OkCheckPoint("toolchain parameter is set to wrong-toolchain"),
				model.ErrorCheckPoint("toolchain parameter error"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithToolchain(test.tchn))
			checkEnv.tchnErr = test.tchnErr
			assert.Equal(t, test.expected, checkToolchainParameter(p))
		})
	}
}

func Test_check_toolchain_language_compatibility(t *testing.T) {
	tests := []struct {
		desc     string
		tchn     string
		lang     language.LangInterface
		langErr  error
		expected []model.CheckPoint
	}{
		{"toolchain param is not set",
			"", nil, nil, nil},
		{
			"unknown language", "my-toolchain",
			nil, errors.New("unknown language"),
			[]model.CheckPoint{
				model.WarningCheckPoint("skipping toolchain/language compatibility check"),
			},
		},
		{
			"all green", "my-toolchain",
			language.ALanguage(
				language.WithName("my-language"),
				language.WithDefaultToolchain("my-toolchain"),
				language.WithCompatibleToolchain("my-toolchain"),
			), nil,
			[]model.CheckPoint{
				model.OkCheckPoint("my-toolchain toolchain is compatible with my-language language"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithToolchain(test.tchn))
			checkEnv.lang = test.lang
			checkEnv.langErr = test.langErr
			assert.Equal(t, test.expected, checkToolchainLanguageCompatibility(p))
		})
	}
}

func Test_check_toolchain_detection(t *testing.T) {
	tests := []struct {
		desc         string
		tchnParam    string
		tchnDetected toolchain.TchnInterface
		tchnErr      error
		lang         language.LangInterface
		langErr      error
		expected     []model.CheckPoint
	}{
		{"toolchain param is set",
			"some-toolchain", nil, nil, nil, nil, nil},
		{
			"unknown language",
			"", nil, nil,
			nil, errors.New("unknown language"),
			[]model.CheckPoint{
				model.OkCheckPoint("toolchain parameter is not set explicitly"),
				model.WarningCheckPoint("language is unknown"),
				model.ErrorCheckPoint("cannot retrieve toolchain from an unknown language"),
			},
		},
		{
			"invalid toolchain",
			"", nil, errors.New("wrong toolchain"),
			language.ALanguage(
				language.WithName("java"),
				language.WithDefaultToolchain("gradle"),
				language.WithCompatibleToolchain("gradle"),
			), nil,
			[]model.CheckPoint{
				model.OkCheckPoint("toolchain parameter is not set explicitly"),
				model.OkCheckPoint("using language's default toolchain"),
				model.ErrorCheckPoint("wrong toolchain"),
			},
		},
		{
			"all green",
			"", toolchain.AToolchain(toolchain.WithName("gradle")), nil,
			language.ALanguage(
				language.WithName("java"),
				language.WithDefaultToolchain("gradle"),
				language.WithCompatibleToolchain("gradle"),
			), nil,
			[]model.CheckPoint{
				model.OkCheckPoint("toolchain parameter is not set explicitly"),
				model.OkCheckPoint("using language's default toolchain"),
				model.OkCheckPoint("default toolchain for java language is gradle"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithToolchain(test.tchnParam))
			checkEnv.tchn = test.tchnDetected
			checkEnv.tchnErr = test.tchnErr
			checkEnv.lang = test.lang
			checkEnv.langErr = test.langErr
			assert.Equal(t, test.expected, checkToolchainDetection(p))
		})
	}
}

func Test_check_toolchain_platform(t *testing.T) {
	tests := []struct {
		desc     string
		tchn     toolchain.TchnInterface
		expected []model.CheckPoint
	}{
		{"with no toolchain", nil, nil},
		{
			"with toolchain", toolchain.AToolchain(),
			[]model.CheckPoint{
				model.OkCheckPoint("local platform is ",
					command.OsName(runtime.GOOS), "/",
					command.ArchName(runtime.GOARCH)),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet()
			checkEnv.tchn = test.tchn
			assert.Equal(t, test.expected, checkToolchainPlatform(p))
		})
	}
}

func Test_check_toolchain_build_command(t *testing.T) {
	tests := []struct {
		desc     string
		tchn     toolchain.TchnInterface
		expected []model.CheckPoint
	}{
		{"with no toolchain", nil, nil},
		{
			"with wrong build command path",
			toolchain.NewFakeToolchain(nil, toolchain.TestStats{}).
				WithCheckCommandAccess(func() (string, error) { return "", errors.New("command not found") }).
				WithBuildCommandPath(func() string { return "/some/path" }).
				WithBuildCommandLine(func() string { return "command arg1 arg2" }),
			[]model.CheckPoint{
				model.OkCheckPoint("build command line: command arg1 arg2"),
				model.ErrorCheckPoint("cannot access build command: /some/path"),
			},
		},
		{
			"with valid build command path",
			toolchain.NewFakeToolchain(nil, toolchain.TestStats{}).
				WithCheckCommandAccess(func() (string, error) { return "/some/path", nil }).
				WithBuildCommandPath(func() string { return "/some/path" }).
				WithBuildCommandLine(func() string { return "command arg1 arg2" }),
			[]model.CheckPoint{
				model.OkCheckPoint("build command line: command arg1 arg2"),
				model.OkCheckPoint("build command path: /some/path"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			checkEnv.tchn = test.tchn
			assert.Equal(t, test.expected, checkToolchainBuildCommand(*params.AParamSet()))
		})
	}
}

func Test_check_toolchain_test_command(t *testing.T) {
	tests := []struct {
		desc     string
		tchn     toolchain.TchnInterface
		expected []model.CheckPoint
	}{
		{"with no toolchain", nil, nil},
		{
			"with wrong test command path",
			toolchain.NewFakeToolchain(nil, toolchain.TestStats{}).
				WithCheckCommandAccess(func() (string, error) { return "", errors.New("command not found") }).
				WithTestCommandPath(func() string { return "/some/path" }).
				WithTestCommandLine(func() string { return "command arg1 arg2" }),
			[]model.CheckPoint{
				model.OkCheckPoint("test command line: command arg1 arg2"),
				model.ErrorCheckPoint("cannot access test command: /some/path"),
			},
		},
		{
			"with valid test command path",
			toolchain.NewFakeToolchain(nil, toolchain.TestStats{}).
				WithCheckCommandAccess(func() (string, error) { return "/some/path", nil }).
				WithTestCommandPath(func() string { return "/some/path" }).
				WithTestCommandLine(func() string { return "command arg1 arg2" }),
			[]model.CheckPoint{
				model.OkCheckPoint("test command line: command arg1 arg2"),
				model.OkCheckPoint("test command path: /some/path"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			checkEnv.tchn = test.tchn
			assert.Equal(t, test.expected, checkToolchainTestCommand(*params.AParamSet()))
		})
	}
}

func Test_check_toolchain_test_result_dir(t *testing.T) {
	workdir, _ := filepath.Abs("/")
	tests := []struct {
		desc     string
		tchn     toolchain.TchnInterface
		expected []model.CheckPoint
	}{
		{"with no toolchain", nil, nil},
		{
			"with empty test result dir",
			toolchain.AToolchain(toolchain.WithTestResultDir("")),
			[]model.CheckPoint{
				model.WarningCheckPoint(
					"test result directory parameter is not set explicitly (default: work directory)"),
				model.OkCheckPoint("test result directory absolute path is ", workdir),
			},
		},
		{
			"with set test result dir",
			toolchain.AToolchain(toolchain.WithTestResultDir("some/path")),
			[]model.CheckPoint{
				model.OkCheckPoint("test result directory parameter is some/path"),
				model.OkCheckPoint("test result directory absolute path is ", filepath.Join(workdir, "some/path")),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			checkEnv.tchn = test.tchn
			_ = toolchain.SetWorkDir(workdir)
			assert.Equal(t, test.expected, checkToolchainTestResultDir(*params.AParamSet()))
		})
	}
}

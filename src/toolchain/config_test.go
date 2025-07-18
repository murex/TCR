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
	"fmt"
	"github.com/murex/tcr/helpers"
	"github.com/murex/tcr/toolchain/command"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_convert_toolchain_name_to_config(t *testing.T) {
	tchn := AToolchain()
	cfg := asConfig(tchn)
	assert.Equal(t, tchn.GetName(), cfg.Name)
	assert.Equal(t, cfg.Name, asToolchain(cfg).GetName())
}

func Test_convert_toolchain_command_os_name_to_config(t *testing.T) {
	cmd := command.ACommand()
	cfg := asCommandConfig(*cmd)
	assert.Equal(t, fmt.Sprint(cmd.Os), fmt.Sprint(cfg.Os))
	assert.Equal(t, fmt.Sprint(cfg.Os), fmt.Sprint(asCommand(cfg).Os))
}

func Test_convert_toolchain_command_arch_name_to_config(t *testing.T) {
	cmd := command.ACommand()
	cfg := asCommandConfig(*cmd)
	assert.Equal(t, fmt.Sprint(cmd.Arch), fmt.Sprint(cfg.Arch))
	assert.Equal(t, fmt.Sprint(cfg.Arch), fmt.Sprint(asCommand(cfg).Arch))
}

func Test_convert_toolchain_command_path_to_config(t *testing.T) {
	cmd := command.ACommand(command.WithPath("some-command-path"))
	cfg := asCommandConfig(*cmd)
	assert.Equal(t, cmd.Path, cfg.Command)
	assert.Equal(t, cfg.Command, asCommand(cfg).Path)
}

func Test_convert_toolchain_command_arguments_to_config(t *testing.T) {
	cmd := command.ACommand(command.WithArgs([]string{"arg1", "arg2"}))
	cfg := asCommandConfig(*cmd)
	assert.Equal(t, fmt.Sprint(cmd.Arguments), fmt.Sprint(cfg.Arguments))
	assert.Equal(t, fmt.Sprint(cfg.Arguments), fmt.Sprint(asCommand(cfg).Arguments))
}

func Test_show_toolchain_configs_with_no_saved_config(t *testing.T) {
	expected := []string{
		"Configured toolchains:",
		"- none (will use built-in toolchains)",
	}
	helpers.AssertSimpleTrace(t, expected,
		func() {
			toolchainDirPath = ""
			ShowConfigs()
		},
	)
}

func Test_reset_toolchain_configs_with_no_saved_config(t *testing.T) {
	expected := []string{
		"Resetting toolchains configuration",
	}
	for _, builtin := range Names() {
		expected = append(expected, "- "+builtin)
	}
	helpers.AssertSimpleTrace(t, expected,
		func() {
			ResetConfigs()
		},
	)
}

func Test_show_toolchain_config(t *testing.T) {
	tchn := AToolchain()
	cfg := asConfig(tchn)
	prefix := "- toolchain." + cfg.Name
	buildCmd := cfg.BuildCommand[0]
	testCmd := cfg.TestCommand[0]
	expected := []string{
		fmt.Sprintf("%v.build.os: %v", prefix, buildCmd.Os),
		fmt.Sprintf("%v.build.arch: %v", prefix, buildCmd.Arch),
		fmt.Sprintf("%v.build.command: %v", prefix, buildCmd.Command),
		fmt.Sprintf("%v.build.args: %v", prefix, buildCmd.Arguments),
		fmt.Sprintf("%v.test.os: %v", prefix, testCmd.Os),
		fmt.Sprintf("%v.test.arch: %v", prefix, testCmd.Arch),
		fmt.Sprintf("%v.test.command: %v", prefix, testCmd.Command),
		fmt.Sprintf("%v.test.args: %v", prefix, testCmd.Arguments),
		fmt.Sprintf("%v.test-result-dir: %v", prefix, tchn.GetTestResultDir()),
	}
	helpers.AssertSimpleTrace(t, expected,
		func() {
			cfg.show()
		},
	)
}

func Test_save_and_load_a_toolchain_config(t *testing.T) {
	const name = "my-toolchain"
	tchn := AToolchain(WithName(name))
	errRegister := Register(tchn)
	if errRegister != nil {
		t.Fatal(errRegister)
	}

	// Set up a temporary directory
	appFS = afero.NewOsFs()
	dir, errTempDir := os.MkdirTemp("", "tcr-toolchain")
	if errTempDir != nil {
		t.Fatal(errTempDir)
	}
	defer t.Cleanup(func() {
		_ = os.RemoveAll(dir)
		delete(registered, name)
	})
	// Prepare toolchain configuration directory
	initConfigDirPath(dir)
	createConfigDir()

	// Save the toolchain configuration file
	saveConfig(name)

	// Load the toolchain configuration file
	yamlConfig := loadConfig(helpers.BuildYAMLFilename(name))

	// Check that we get back the same toolchain data
	if assert.NotNil(t, yamlConfig) {
		assert.Equal(t, tchn, asToolchain(*yamlConfig))
	}
}

func Test_save_and_load_all_toolchain_configs(t *testing.T) {
	// Set up a temporary directory
	appFS = afero.NewOsFs()
	dir, errTempDir := os.MkdirTemp("", "tcr-toolchain")
	if errTempDir != nil {
		t.Fatal(errTempDir)
	}
	defer t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})
	// Prepare toolchain configuration directory
	initConfigDirPath(dir)
	createConfigDir()

	// Make a copy of the registered toolchains to be used as a reference
	expected := make(map[string]TchnInterface)
	for name, tchn := range registered {
		expected[name] = tchn
	}

	// Save all toolchains configuration files. By default these are all built-in toolchains
	SaveConfigs()

	// Empty the map of registered toolchains
	for name := range registered {
		delete(registered, name)
	}

	// Load all toolchains configuration files
	loadConfigs()

	// Check that the toolchains are back in the registered map
	if assert.Equal(t, len(expected), len(registered)) {
		for name, expectedTchn := range expected {
			registeredTchn, found := registered[name]
			if assert.True(t, found) {
				assert.Equal(t, expectedTchn, registeredTchn)
			}
		}
	}
}

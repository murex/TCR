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
	"github.com/murex/tcr/tcr-engine/toolchain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_can_save_toolchain_configuration(t *testing.T) {
	// TODO bypass filesystem
	//saveToYaml(tchn, "")
}

func Test_convert_toolchain_name_to_config(t *testing.T) {
	tchn := toolchain.AToolchain()
	cfg := asToolchainConfig(tchn)
	assert.Equal(t, tchn.GetName(), cfg.Name)
	assert.Equal(t, cfg.Name, asToolchain(cfg).GetName())
}

func Test_convert_toolchain_command_os_name_to_config(t *testing.T) {
	cmd := toolchain.ACommand()
	cfg := asToolchainCommandConfig(*cmd)
	assert.Equal(t, fmt.Sprint(cmd.Os), fmt.Sprint(cfg.Os))
	assert.Equal(t, fmt.Sprint(cfg.Os), fmt.Sprint(asToolchainCommand(cfg).Os))
}

func Test_convert_toolchain_command_arch_name_to_config(t *testing.T) {
	cmd := toolchain.ACommand()
	cfg := asToolchainCommandConfig(*cmd)
	assert.Equal(t, fmt.Sprint(cmd.Arch), fmt.Sprint(cfg.Arch))
	assert.Equal(t, fmt.Sprint(cfg.Arch), fmt.Sprint(asToolchainCommand(cfg).Arch))
}

func Test_convert_toolchain_command_path_to_config(t *testing.T) {
	cmd := toolchain.ACommand(toolchain.WithPath("some-command-path"))
	cfg := asToolchainCommandConfig(*cmd)
	assert.Equal(t, cmd.Path, cfg.Command)
	assert.Equal(t, cfg.Command, asToolchainCommand(cfg).Path)
}

func Test_convert_toolchain_command_arguments_to_config(t *testing.T) {
	cmd := toolchain.ACommand(toolchain.WithArgs([]string{"arg1", "arg2"}))
	cfg := asToolchainCommandConfig(*cmd)
	assert.Equal(t, fmt.Sprint(cmd.Arguments), fmt.Sprint(cfg.Arguments))
	assert.Equal(t, fmt.Sprint(cfg.Arguments), fmt.Sprint(asToolchainCommand(cfg).Arguments))
}

func Test_show_toolchain_configs_with_no_saved_config(t *testing.T) {
	expected := []string{
		"Configured toolchains:",
		"- none (will use built-in toolchains)",
	}
	assertConfigTrace(t, expected,
		func() {
			showToolchainConfigs()
		},
	)
}

func Test_reset_toolchain_configs_with_no_saved_config(t *testing.T) {
	expected := []string{
		"Resetting toolchains configuration",
	}
	for _, builtin := range toolchain.Names() {
		expected = append(expected, "- "+builtin)
	}
	assertConfigTrace(t, expected,
		func() {
			resetToolchainConfigs()
		},
	)
}

func Test_show_toolchain_config(t *testing.T) {
	tchn := toolchain.AToolchain()
	cfg := asToolchainConfig(tchn)
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
	}
	assertConfigTrace(t, expected,
		func() {
			cfg.show()
		},
	)
}

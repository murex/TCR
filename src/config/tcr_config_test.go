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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/murex/tcr/helpers"
	"github.com/murex/tcr/language"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/toolchain"
	"github.com/murex/tcr/variant"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testDataRootDir = "../testdata"
)

var (
	testDataDirJava = filepath.Join(testDataRootDir, "java")
)

func Test_get_config_dir_path_when_not_set(t *testing.T) {
	initConfigDirPath()
	assert.Equal(t, ".tcr", GetConfigDirPath())
}

func Test_get_config_dir_path_when_set(t *testing.T) {
	d, _ := os.MkdirTemp("", "tcr-config-dir")
	defer func() {
		_ = os.RemoveAll(configDirPath)
	}()

	cmd := NewCobraTestCmd()
	cmd.Run = func(cmd *cobra.Command, args []string) {
		InitForTest()
	}
	AddParameters(cmd, d)
	cmd.SetArgs([]string{"--config-dir", d})
	_ = cmd.Execute()

	assert.Equal(t, filepath.Join(d, ".tcr"), GetConfigDirPath())
}

func Test_show_config_value(t *testing.T) {
	key, value := "some-key", "some-value"
	expected := []string{"- " + key + ": " + value}
	helpers.AssertSimpleTrace(t, expected,
		func() {
			helpers.TraceKeyValue(key, value)
		},
	)
}

func Test_reset_tcr_config_with_no_saved_config(t *testing.T) {
	d, _ := os.MkdirTemp("", "tcr-reset-dir")
	defer func() {
		_ = os.RemoveAll(d)
	}()
	expected := []string{
		"Creating TCR configuration directory: " + filepath.Join(d, ".tcr"),
		"No configuration file found",
		"Loading toolchains configuration",
		"Loading languages configuration",
		"Resetting configuration to default values",
	}

	expected = append(expected, "Resetting toolchains configuration")
	for _, builtinTchn := range toolchain.Names() {
		expected = append(expected, "- "+builtinTchn)
	}

	expected = append(expected, "Resetting languages configuration")
	for _, builtinLang := range language.Names() {
		expected = append(expected, "- "+builtinLang)
	}
	expected = append(expected, "Saving configuration: "+filepath.Join(d, ".tcr", "config.yml"))

	expected = append(expected, "Creating TCR toolchain configuration directory: "+filepath.Join(d, ".tcr", "toolchain"))
	expected = append(expected, "Saving toolchains configuration")
	for _, builtinTchn := range toolchain.Names() {
		expected = append(expected, "- "+builtinTchn)
	}

	expected = append(expected, "Creating TCR language configuration directory: "+filepath.Join(d, ".tcr", "language"))
	expected = append(expected, "Saving languages configuration")
	for _, builtinLang := range language.Names() {
		expected = append(expected, "- "+builtinLang)
	}
	helpers.AssertSimpleTrace(t, expected,
		func() {
			cmd := NewCobraTestCmd()
			cmd.Run = func(cmd *cobra.Command, args []string) {
				InitForTest()
				Reset()
			}
			AddParameters(cmd, d)
			cmd.SetArgs([]string{"--config-dir", d})
			_ = cmd.Execute()
		},
	)
}

func Test_save_tcr_config_with_no_saved_config(t *testing.T) {
	d, _ := os.MkdirTemp("", "tcr-save-dir")
	defer func() {
		_ = os.RemoveAll(d)
	}()
	expected := []string{
		"Creating TCR configuration directory: " + filepath.Join(d, ".tcr"),
		"No configuration file found",
		"Loading toolchains configuration",
		"Loading languages configuration",
		"Saving configuration: " + filepath.Join(d, ".tcr", "config.yml"),
	}

	expected = append(expected, "Creating TCR toolchain configuration directory: "+filepath.Join(d, ".tcr", "toolchain"))
	expected = append(expected, "Saving toolchains configuration")
	for _, builtinTchn := range toolchain.Names() {
		expected = append(expected, "- "+builtinTchn)
	}

	expected = append(expected, "Creating TCR language configuration directory: "+filepath.Join(d, ".tcr", "language"))
	expected = append(expected, "Saving languages configuration")
	for _, builtinLang := range language.Names() {
		expected = append(expected, "- "+builtinLang)
	}
	helpers.AssertSimpleTrace(t, expected,
		func() {
			cmd := NewCobraTestCmd()
			cmd.Run = func(cmd *cobra.Command, args []string) {
				InitForTest()
				Save()
			}
			AddParameters(cmd, d)
			cmd.SetArgs([]string{"--config-dir", d})
			_ = cmd.Execute()
		},
	)
}

func Test_init_tcr_config_with_no_config_file(t *testing.T) {
	expected := []string{"No configuration file found"}
	helpers.AssertSimpleTrace(t, expected,
		func() {
			initTCRConfig()
		},
	)
}

func Test_cobra_command_init_with_no_saved_config(t *testing.T) {
	d, _ := os.MkdirTemp("", "tcr-init-dir")
	defer func() {
		_ = os.RemoveAll(d)
	}()
	expected := []string{
		"Creating TCR configuration directory: " + filepath.Join(d, ".tcr"),
		"No configuration file found",
		"Loading toolchains configuration",
		"Loading languages configuration",
	}
	helpers.AssertSimpleTrace(t, expected,
		func() {
			cmd := NewCobraTestCmd()
			cmd.Run = func(cmd *cobra.Command, args []string) {
				InitForTest()
			}
			AddParameters(cmd, d)
			cmd.SetArgs([]string{"--config-dir", d})
			_ = cmd.Execute()
		},
	)
}

func Test_cobra_command_init_with_saved_config(t *testing.T) {
	d := testDataDirJava

	expected := []string{
		"Loading configuration: " + filepath.Join(d, ".tcr", "config.yml"),
		"Loading toolchains configuration",
		"Loading languages configuration",
	}
	helpers.AssertSimpleTrace(t, expected,
		func() {
			cmd := NewCobraTestCmd()
			cmd.Run = func(cmd *cobra.Command, args []string) {
				InitForTest()
			}
			AddParameters(cmd, d)
			cmd.SetArgs([]string{"--config-dir", d})
			_ = cmd.Execute()
		},
	)
}

var testParams params.Params

func NewCobraTestCmd() *cobra.Command {
	return &cobra.Command{
		Use: "test",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			UpdateEngineParams(&testParams)
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}
}

func InitForTest() {
	initConfig(nil)
}

func Test_show_tcr_config_with_default_values(t *testing.T) {
	prefix := "- config"
	expected := []string{
		"TCR configuration:",
		fmt.Sprintf("%v.git.auto-push: %v", prefix, false),
		fmt.Sprintf("%v.git.polling-period: %v", prefix, 2*time.Second),
		fmt.Sprintf("%v.mob-timer.duration: %v", prefix, 5*time.Minute),
		fmt.Sprintf("%v.tcr.language: %v", prefix, ""),
		fmt.Sprintf("%v.tcr.toolchain: %v", prefix, ""),
		fmt.Sprintf("%v.tcr.trace: %v", prefix, "none"),
		fmt.Sprintf("%v.tcr.variant: %v", prefix, variant.Relaxed),
		fmt.Sprintf("%v.vcs.name: %v", prefix, "git"),
	}
	helpers.AssertSimpleTrace(t, expected,
		func() {
			showTCRConfig()
		},
	)
}

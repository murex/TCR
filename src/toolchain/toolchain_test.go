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
	"github.com/murex/tcr/toolchain/command"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func Test_set_work_dir(t *testing.T) {
	testFlags := []struct {
		desc         string
		path         string
		expectError  bool
		expectedPath func() string
	}{
		{
			"with empty path",
			"",
			false,
			func() string { path, _ := os.Getwd(); return path },
		},
		{
			"with existing directory",
			filepath.Join(testDataDirJava),
			false,
			func() string { path, _ := filepath.Abs(testDataDirJava); return path },
		},
		{
			"with non-existing directory",
			filepath.Join(testDataDirJava, "fake-dir"),
			true,
			func() string { return "" },
		},
		{
			"with existing file",
			filepath.Join(testDataDirJava, "pom.xml"),
			true,
			func() string { return "" },
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			err := SetWorkDir(tt.path)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedPath(), GetWorkDir())
		})
	}
}

func Test_get_build_commands(t *testing.T) {
	cmd1 := command.ACommand(command.WithPath("build-cmd-1"))
	cmd2 := command.ACommand(command.WithPath("build-cmd-2"))
	tchn := AToolchain(WithBuildCommand(cmd1), WithBuildCommand(cmd2))
	commands := tchn.GetBuildCommands()
	assert.Contains(t, commands, *cmd1)
	assert.Contains(t, commands, *cmd2)
}

func Test_get_test_commands(t *testing.T) {
	cmd1 := command.ACommand(command.WithPath("test-cmd-1"))
	cmd2 := command.ACommand(command.WithPath("test-cmd-2"))
	tchn := AToolchain(WithTestCommand(cmd1), WithTestCommand(cmd2))
	commands := tchn.GetTestCommands()
	assert.Contains(t, commands, *cmd1)
	assert.Contains(t, commands, *cmd2)
}

func Test_build_command_line(t *testing.T) {
	cmd := command.ACommand(command.WithPath("build-cmd"), command.WithArgs([]string{"arg1", "arg2"}))
	tchn := AToolchain(WithNoBuildCommand(), WithBuildCommand(cmd))
	assert.Equal(t, "build-cmd arg1 arg2", tchn.BuildCommandLine())
}

func Test_test_command_line(t *testing.T) {
	cmd := command.ACommand(command.WithPath("test-cmd"), command.WithArgs([]string{"arg1", "arg2"}))
	tchn := AToolchain(WithNoTestCommand(), WithTestCommand(cmd))
	assert.Equal(t, "test-cmd arg1 arg2", tchn.TestCommandLine())
}

func Test_check_command_access_for_valid_command(t *testing.T) {
	tchn := AToolchain()
	path, err := tchn.CheckCommandAccess("go")
	assert.NoError(t, err)
	assert.NotZero(t, path)
}

func Test_check_command_access_for_invalid_command(t *testing.T) {
	tchn := AToolchain()
	path, err := tchn.CheckCommandAccess("invalid-command")
	assert.Error(t, err)
	assert.Zero(t, path)
}

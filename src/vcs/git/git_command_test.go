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

package git

import (
	"errors"
	"github.com/murex/tcr/vcs/shell"
	"github.com/stretchr/testify/assert"
	"testing"
)

// newGitCommandStub creates a new git shell command stub
func newGitCommandStub() *shell.CommandStub {
	return shell.NewCommandStub(*newGitCommandImpl())
}

// Make sure that shell.NewCommandFunc is back to normal after each test
func teardownTest() {
	shell.NewCommandFunc = shell.NewCommand
}

func Test_is_git_command_available(t *testing.T) {
	tests := []struct {
		desc     string
		inPath   bool
		expected bool
	}{
		{"git command found", true, true},
		{"git command not found", false, false},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newGitCommandStub()
				stub.IsInPathFunc = func() bool {
					return test.inPath
				}
				return stub
			}
			assert.Equal(t, test.expected, IsGitCommandAvailable())
		})
	}
}

func Test_get_git_command_path(t *testing.T) {
	tests := []struct {
		desc     string
		realPath string
		expected string
	}{
		{"git command found", "/some-path/git", "/some-path/git"},
		{"git command not found", "", ""},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newGitCommandStub()
				stub.GetFullPathFunc = func() string {
					return test.realPath
				}
				return stub
			}
			assert.Equal(t, test.expected, GetGitCommandPath())
		})
	}
}

func Test_get_git_command_version(t *testing.T) {
	tests := []struct {
		desc         string
		gitCmdOutput string
		gitCmdError  error
		expected     string
	}{
		{
			"git command found",
			"git version 2.39.1.windows.1" + shell.GetAttributes().EOL,
			nil,
			"2.39.1.windows.1",
		},
		{
			"git command not found",
			"",
			errors.New("git command not found"),
			"unknown",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newGitCommandStub()
				stub.RunFunc = func(params ...string) (output []byte, err error) {
					return []byte(test.gitCmdOutput), test.gitCmdError
				}
				return stub
			}
			assert.Equal(t, test.expected, GetGitCommandVersion())
		})
	}
}

func Test_get_git_username(t *testing.T) {
	tests := []struct {
		desc         string
		gitCmdOutput string
		gitCmdError  error
		expected     string
	}{
		{
			"git command found",
			"John Doe" + shell.GetAttributes().EOL,
			nil,
			"John Doe",
		},
		{
			"git command not found",
			"",
			errors.New("git command not found"),
			"not set",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newGitCommandStub()
				stub.RunFunc = func(params ...string) (output []byte, err error) {
					return []byte(test.gitCmdOutput), test.gitCmdError
				}
				return stub
			}
			assert.Equal(t, test.expected, GetGitUserName())
		})
	}
}

func Test_get_git_config_value(t *testing.T) {
	tests := []struct {
		desc         string
		gitConfigKey string
		gitCmdOutput string
		gitCmdError  error
		expected     string
	}{
		{
			"valid git config key",
			"user.name",
			"John Doe" + shell.GetAttributes().EOL,
			nil,
			"John Doe",
		},
		{
			"invalid git config key",
			"undefined.config.key",
			"" + shell.GetAttributes().EOL,
			nil,
			"not set",
		},
		{
			"git command not found",
			"",
			"",
			errors.New("git command not found"),
			"not set",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newGitCommandStub()
				stub.RunFunc = func(params ...string) (output []byte, err error) {
					return []byte(test.gitCmdOutput), test.gitCmdError
				}
				return stub
			}
			assert.Equal(t, test.expected, getGitConfigValue(test.gitConfigKey))
		})
	}
}

var (
	runGitCmdTests = []struct {
		desc           string
		gitParams      []string
		gitCmdOutput   string
		gitCmdError    error
		expectedOutput string
		expectedError  error
	}{
		{
			"git command found with valid subcommand",
			[]string{"status"},
			"nothing to commit, working tree clean" + shell.GetAttributes().EOL,
			nil,
			"nothing to commit, working tree clean" + shell.GetAttributes().EOL,
			nil,
		},
		{
			"git command found with invalid subcommand",
			[]string{"xxxx"},
			"git: 'xxxx' is not a git command. See 'git --help'." + shell.GetAttributes().EOL,
			errors.New("git invalid subcommand"),
			"git: 'xxxx' is not a git command. See 'git --help'." + shell.GetAttributes().EOL,
			errors.New("git invalid subcommand"),
		},
		{
			"git command not found",
			[]string{"status"},
			"",
			errors.New("git command not found"),
			"",
			errors.New("git command not found"),
		},
	}
)

func Test_run_git_command(t *testing.T) {
	for _, test := range runGitCmdTests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newGitCommandStub()
				stub.RunFunc = func(params ...string) (output []byte, err error) {
					return []byte(test.gitCmdOutput), test.gitCmdError
				}
				return stub
			}
			output, err := runGitCommand(test.gitParams...)
			assert.Equal(t, test.expectedOutput, string(output))
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func Test_trace_git_command(t *testing.T) {
	for _, test := range runGitCmdTests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			var trace string
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newGitCommandStub()
				stub.RunFunc = func(params ...string) (output []byte, err error) {
					return []byte(test.gitCmdOutput), test.gitCmdError
				}
				stub.TraceFunc = func(params ...string) error {
					output, err := stub.Run(params...)
					trace = string(output)
					return err
				}
				return stub
			}
			err := traceGitCommand(test.gitParams...)
			assert.Equal(t, test.expectedOutput, trace)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

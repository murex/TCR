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

package p4

import (
	"errors"
	"github.com/murex/tcr/vcs/shell"
	"github.com/stretchr/testify/assert"
	"testing"
)

// newP4CommandStub creates a new p4 shell command stub
func newP4CommandStub() *shell.CommandStub {
	return shell.NewCommandStub(*newP4CommandImpl())
}

// Make sure that shell.NewCommandFunc is back to normal after each test
func teardownTest() {
	shell.NewCommandFunc = shell.NewCommand
}

func Test_is_p4_command_available(t *testing.T) {
	tests := []struct {
		desc     string
		inPath   bool
		expected bool
	}{
		{"p4 command found", true, true},
		{"p4 command not found", false, false},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newP4CommandStub()
				stub.IsInPathFunc = func() bool {
					return test.inPath
				}
				return stub
			}
			assert.Equal(t, test.expected, IsP4CommandAvailable())
		})
	}
}

func Test_get_p4_command_path(t *testing.T) {
	tests := []struct {
		desc     string
		realPath string
		expected string
	}{
		{"p4 command found", "/some-path/p4", "/some-path/p4"},
		{"p4 command not found", "", ""},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newP4CommandStub()
				stub.GetFullPathFunc = func() string {
					return test.realPath
				}
				return stub
			}
			assert.Equal(t, test.expected, GetP4CommandPath())
		})
	}
}

func Test_get_p4_command_version(t *testing.T) {
	tests := []struct {
		desc        string
		p4CmdOutput string
		p4CmdError  error
		expected    string
	}{
		{
			"p4 command found",
			`
Perforce - The Fast Software Configuration Management System.
Copyright 1995-2022 Perforce Software.  All rights reserved.
This product includes software developed by the OpenSSL Project
for use in the OpenSSL Toolkit (http://www.openssl.org/)
Version of OpenSSL Libraries: OpenSSL 1.1.1q  5 Jul 2022
See 'p4 help [ -l ] legal' for additional license information on
these licenses and others.
Extensions/scripting support built-in.
Parallel sync threading built-in.
Rev. P4/NTX64/2022.2/2369846 (2022/11/14).
`,
			nil,
			"P4/NTX64/2022.2/2369846",
		},
		{
			"p4 command not found",
			"",
			errors.New("p4 command not found"),
			"unknown",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newP4CommandStub()
				stub.RunFunc = func(params ...string) (output []byte, err error) {
					return []byte(test.p4CmdOutput), test.p4CmdError
				}
				return stub
			}
			assert.Equal(t, test.expected, GetP4CommandVersion())
		})
	}
}

func Test_get_p4_username(t *testing.T) {
	tests := []struct {
		desc        string
		p4CmdOutput string
		p4CmdError  error
		expected    string
	}{
		{
			"p4 command found",
			"P4USER=John Doe" + shell.GetAttributes().EOL,
			nil,
			"John Doe",
		},
		{
			"p4 command not found",
			"",
			errors.New("p4 command not found"),
			"not set",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newP4CommandStub()
				stub.RunFunc = func(params ...string) (output []byte, err error) {
					return []byte(test.p4CmdOutput), test.p4CmdError
				}
				return stub
			}
			assert.Equal(t, test.expected, GetP4UserName())
		})
	}
}

func Test_get_p4_config_value(t *testing.T) {
	tests := []struct {
		desc           string
		p4VariableName string
		p4CmdOutput    string
		p4CmdError     error
		expected       string
	}{
		{
			"valid p4 variable name",
			"P4CLIENT",
			"P4CLIENT=client-name" + shell.GetAttributes().EOL,
			nil,
			"client-name",
		},
		{
			"invalid p4 variable name",
			"P4XXXX",
			"" + shell.GetAttributes().EOL,
			nil,
			"not set",
		},
		{
			"p4 command not found",
			"",
			"",
			errors.New("p4 command not found"),
			"not set",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newP4CommandStub()
				stub.RunFunc = func(params ...string) (output []byte, err error) {
					return []byte(test.p4CmdOutput), test.p4CmdError
				}
				return stub
			}
			assert.Equal(t, test.expected, getP4ConfigValue(test.p4VariableName))
		})
	}
}

var (
	runP4CmdTests = []struct {
		desc           string
		p4Params       []string
		p4CmdOutput    string
		p4CmdError     error
		expectedOutput string
		expectedError  error
	}{
		{
			"p4 command found with valid subcommand",
			[]string{"users", "johndoe"},
			"johndoe <john.doe@acme.com> (DOE John) accessed 2023/01/17" + shell.GetAttributes().EOL,
			nil,
			"johndoe <john.doe@acme.com> (DOE John) accessed 2023/01/17" + shell.GetAttributes().EOL,
			nil,
		},
		{
			"p4 command found with invalid subcommand",
			[]string{"xxxx"},
			"Unknown command.  Try 'p4 help' for info." + shell.GetAttributes().EOL,
			errors.New("p4 invalid subcommand"),
			"Unknown command.  Try 'p4 help' for info." + shell.GetAttributes().EOL,
			errors.New("p4 invalid subcommand"),
		},
		{
			"p4 command not found",
			[]string{"users", "johndoe"},
			"",
			errors.New("p4 command not found"),
			"",
			errors.New("p4 command not found"),
		},
	}
)

func Test_run_p4_command(t *testing.T) {
	for _, test := range runP4CmdTests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newP4CommandStub()
				stub.RunFunc = func(params ...string) (output []byte, err error) {
					return []byte(test.p4CmdOutput), test.p4CmdError
				}
				return stub
			}
			output, err := runP4Command(test.p4Params...)
			assert.Equal(t, test.expectedOutput, string(output))
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func Test_trace_p4_command(t *testing.T) {
	for _, test := range runP4CmdTests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			var trace string
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newP4CommandStub()
				stub.RunFunc = func(params ...string) (output []byte, err error) {
					return []byte(test.p4CmdOutput), test.p4CmdError
				}
				stub.TraceFunc = func(params ...string) error {
					output, err := stub.Run(params...)
					trace = string(output)
					return err
				}
				return stub
			}
			err := traceP4Command(test.p4Params...)
			assert.Equal(t, test.expectedOutput, trace)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func Test_run_piped_p4_command(t *testing.T) {
	for _, test := range runP4CmdTests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newP4CommandStub()
				stub.RunAndPipeFunc = func(_ shell.Command, params ...string) (output []byte, err error) {
					// We don't run a real command pipe here (would imply that we run real commands,
					// which is exactly what we want to avoid)
					// instead, we return what the piped to command would return
					return []byte(test.p4CmdOutput), test.p4CmdError
				}
				return stub
			}
			output, err := runPipedP4Command(shell.NewCommandFunc("stubbed"), test.p4Params...)
			assert.Equal(t, test.expectedOutput, string(output))
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func Test_trace_piped_p4_command(t *testing.T) {
	for _, test := range runP4CmdTests {
		t.Run(test.desc, func(t *testing.T) {
			defer teardownTest()
			var trace string
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := newP4CommandStub()
				stub.RunAndPipeFunc = func(_ shell.Command, params ...string) (output []byte, err error) {
					// We don't run a real command pipe here (would imply that we run real commands,
					// which is exactly what we want to avoid)
					// instead, we return what the piped to command would return
					return []byte(test.p4CmdOutput), test.p4CmdError
				}
				stub.TraceAndPipeFunc = func(toCmd shell.Command, params ...string) error {
					output, err := stub.RunAndPipe(toCmd, params...)
					trace = string(output)
					return err
				}
				return stub
			}
			err := tracePipedP4Command(shell.NewCommandFunc("stubbed"), test.p4Params...)
			assert.Equal(t, test.expectedOutput, trace)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

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
	"github.com/murex/tcr/vcs/p4"
	"github.com/murex/tcr/vcs/shell"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_check_p4_environment(t *testing.T) {
	t.Skip("TODO")
}

func Test_check_p4_command(t *testing.T) {
	tests := []struct {
		desc     string
		isInPath bool
		fullPath string
		version  string
		expected []model.CheckPoint
	}{
		{
			"p4 command not found", false, "", "",
			[]model.CheckPoint{
				model.ErrorCheckPoint("p4 command was not found on path"),
			},
		},
		{
			"p4 command found", true, "some-path", "some-version",
			[]model.CheckPoint{
				model.OkCheckPoint("p4 command path is some-path"),
				model.OkCheckPoint("p4 version is some-version"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			defer p4.RestoreP4Command()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := p4.NewP4CommandStub()
				stub.IsInPathFunc = func() bool {
					return test.isInPath
				}
				stub.GetFullPathFunc = func() string {
					return test.fullPath
				}
				stub.RunFunc = func(params ...string) (out []byte, err error) {
					return []byte("Rev. " + test.version + " (0000/00/00)."), nil
				}
				return stub
			}
			initTestCheckEnv(*params.AParamSet())
			assert.Equal(t, test.expected, checkP4Command())
		})
	}
}

func Test_check_p4_config(t *testing.T) {
	tests := []struct {
		desc     string
		username string
		expected []model.CheckPoint
	}{
		{
			"p4 username set", "jane-doe",
			[]model.CheckPoint{
				model.OkCheckPoint("p4 username is jane-doe"),
			},
		},
		{
			"p4 username not set", "not set",
			[]model.CheckPoint{
				model.WarningCheckPoint("p4 username is not set"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			defer p4.RestoreP4Command()
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := p4.NewP4CommandStub()
				stub.RunFunc = func(params ...string) (out []byte, err error) {
					return []byte(test.username), nil
				}
				return stub
			}
			initTestCheckEnv(*params.AParamSet())
			assert.Equal(t, test.expected, checkP4Config())
		})
	}
}

func Test_check_p4_workspace(t *testing.T) {
	t.Skip("TODO")
}

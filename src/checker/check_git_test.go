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
	"testing"

	"github.com/murex/tcr/checker/model"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/vcs"
	"github.com/murex/tcr/vcs/fake"
	"github.com/murex/tcr/vcs/git"
	"github.com/murex/tcr/vcs/shell"
	"github.com/stretchr/testify/assert"
)

func Test_check_git_environment(t *testing.T) {
	assertCheckGroupRunner(t,
		checkGitEnvironment,
		&checkGitRunners,
		*params.AParamSet(params.WithVCS(git.Name)),
		"git environment")
}

func Test_check_git_command(t *testing.T) {
	tests := []struct {
		desc     string
		isInPath bool
		fullPath string
		version  string
		expected []model.CheckPoint
	}{
		{
			"git command not found", false, "", "",
			[]model.CheckPoint{
				model.ErrorCheckPoint("git command was not found on path"),
			},
		},
		{
			"git command found", true, "some-path", "some-version",
			[]model.CheckPoint{
				model.OkCheckPoint("git command path is some-path"),
				model.OkCheckPoint("git version is some-version"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Cleanup(git.RestoreGitCommand)
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := git.NewGitCommandStub()
				stub.IsInPathFunc = func() bool {
					return test.isInPath
				}
				stub.GetFullPathFunc = func() string {
					return test.fullPath
				}
				stub.RunFunc = func(params ...string) (out []byte, err error) {
					return []byte("git version " + test.version), nil
				}
				return stub
			}
			p := *params.AParamSet()
			initTestCheckEnv(p)
			assert.Equal(t, test.expected, checkGitCommand(p))
		})
	}
}

func Test_check_git_config(t *testing.T) {
	tests := []struct {
		desc     string
		username string
		expected []model.CheckPoint
	}{
		{
			"git username not set", "not set",
			[]model.CheckPoint{
				model.WarningCheckPoint("git username is not set"),
			},
		},
		{
			"git username set", "jane-doe",
			[]model.CheckPoint{
				model.OkCheckPoint("git username is jane-doe"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Cleanup(git.RestoreGitCommand)
			shell.NewCommandFunc = func(name string, params ...string) shell.Command {
				stub := git.NewGitCommandStub()
				stub.RunFunc = func(params ...string) (out []byte, err error) {
					return []byte(test.username), nil
				}
				return stub
			}
			p := *params.AParamSet()
			initTestCheckEnv(p)
			assert.Equal(t, test.expected, checkGitConfig(p))
		})
	}
}

func Test_check_git_repository(t *testing.T) {
	tests := []struct {
		desc            string
		sourceTreeError error
		vcsError        error
		expected        []model.CheckPoint
	}{
		{
			"git source tree init failed",
			errors.New("git source tree init failed"), nil,
			[]model.CheckPoint{
				model.ErrorCheckPoint("cannot retrieve git repository information from base directory name"),
			},
		},
		{
			"git error",
			nil, errors.New("some git error"),
			[]model.CheckPoint{
				model.ErrorCheckPoint("some git error"),
			},
		},
		{
			"root branch warning",
			nil, nil,
			[]model.CheckPoint{
				model.OkCheckPoint("git repository root is vcs-fake-root-dir"),
				model.OkCheckPoint("git working branch is vcs-fake-working-branch"),
				model.WarningCheckPoint("running TCR from a root branch is not recommended"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet()
			initTestCheckEnv(p)
			checkEnv.sourceTreeErr = test.sourceTreeError
			checkEnv.vcsErr = test.vcsError
			assert.Equal(t, test.expected, checkGitRepository(p))
		})
	}
}

func Test_check_git_remote(t *testing.T) {
	tests := []struct {
		desc     string
		vcs      vcs.Interface
		expected []model.CheckPoint
	}{
		{
			"VCS not initialized", nil, []model.CheckPoint{},
		},
		{
			"git remote disabled due to undefined remote name",
			fake.NewVCSFake(fake.Settings{RemoteEnabled: false}),
			[]model.CheckPoint{
				model.WarningCheckPoint("git remote not found: vcs-fake-remote-name"),
				model.OkCheckPoint("git remote is disabled: all operations will be done locally"),
			},
		},
		{
			"git remote access not working",
			fake.NewVCSFake(fake.Settings{RemoteEnabled: true, RemoteAccessWorking: false}),
			[]model.CheckPoint{
				model.OkCheckPoint("git remote name is vcs-fake-remote-name"),
				model.ErrorCheckPoint("git remote access does not seem to be working"),
			},
		},
		{
			"git remote access working",
			fake.NewVCSFake(fake.Settings{RemoteEnabled: true, RemoteAccessWorking: true}),
			[]model.CheckPoint{
				model.OkCheckPoint("git remote name is vcs-fake-remote-name"),
				model.OkCheckPoint("git remote access seems to be working"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			checkEnv.vcs = test.vcs
			assert.Equal(t, test.expected, checkGitRemote(*params.AParamSet()))
		})
	}
}

func Test_check_git_auto_push(t *testing.T) {
	tests := []struct {
		desc     string
		value    bool
		expected []model.CheckPoint
	}{
		{"enabled", true, []model.CheckPoint{
			model.OkCheckPoint("git auto-push is turned on: every commit will be pushed to origin"),
		},
		},
		{"disabled", false, []model.CheckPoint{
			model.OkCheckPoint("git auto-push is turned off: commits will only be applied locally"),
		},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithAutoPush(test.value))
			assert.Equal(t, test.expected, checkGitAutoPush(p))
		})
	}
}

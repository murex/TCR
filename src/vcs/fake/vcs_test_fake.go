//go:build test_helper

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

package fake

import (
	"errors"
	"github.com/murex/tcr/vcs"
)

type (
	// Command is the name of a Git command
	Command string

	// Commands is a slice of Command
	Commands []Command
)

func (gc Commands) contains(command Command) bool {
	for _, value := range gc {
		if value == command {
			return true
		}
	}
	return false
}

// List of supported VCS commands
const (
	AddCommand     Command = "add"
	CommitCommand  Command = "commit"
	DiffCommand    Command = "diff"
	LogCommand     Command = "log"
	PullCommand    Command = "pull"
	PushCommand    Command = "push"
	RestoreCommand Command = "restore"
	RevertCommand  Command = "revert"
	StashCommand   Command = "stash"
	UnStashCommand Command = "unStash"
)

type (
	// Settings provide a few ways to tune VCS Fake behaviour
	Settings struct {
		FailingCommands Commands
		ChangedFiles    vcs.FileDiffs
		Logs            vcs.LogItems
	}

	// VCSFake provides a fake implementation of the VCS interface
	VCSFake struct {
		settings    Settings
		pushEnabled bool
		lastCommand Command
	}
)

func (vf *VCSFake) fakeCommand(cmd Command) (err error) {
	vf.lastCommand = cmd
	if vf.settings.FailingCommands.contains(cmd) {
		err = errors.New("vcs " + string(cmd) + " fake error")
	}
	return
}

// NewVCSFake initializes a fake VCS implementation which does nothing
// apart from emulating errors on VCS operations
func NewVCSFake(settings Settings) (*VCSFake, error) {
	return &VCSFake{settings: settings}, nil
}

func (vf *VCSFake) Name() string {
	return "vcs-fake"
}

func (vf *VCSFake) SessionSummary() string {
	return "VCS session \"fake\""
}

// GetLastCommand returns the last command called
func (vf *VCSFake) GetLastCommand() Command {
	return vf.lastCommand
}

// Add does nothing. Returns an error if in the list of failing commands
func (vf *VCSFake) Add(_ ...string) error {
	return vf.fakeCommand(AddCommand)
}

// Commit does nothing. Returns an error if in the list of failing commands
func (vf *VCSFake) Commit(_ bool, _ ...string) error {
	return vf.fakeCommand(CommitCommand)
}

// Restore does nothing. Returns an error if in the list of failing commands
func (vf *VCSFake) Restore(_ string) error {
	return vf.fakeCommand(RestoreCommand)
}

// Push does nothing. Returns an error if in the list of failing commands
func (vf *VCSFake) Push() error {
	return vf.fakeCommand(PushCommand)
}

// Pull does nothing. Returns an error if in the list of failing commands
func (vf *VCSFake) Pull() error {
	return vf.fakeCommand(PullCommand)
}

// Diff returns the list of files modified configured at fake initialization
func (vf *VCSFake) Diff() (_ vcs.FileDiffs, err error) {
	return vf.settings.ChangedFiles, vf.fakeCommand(DiffCommand)
}

// Log returns the list of VCS logs configured at fake initialization
func (vf *VCSFake) Log(msgFilter func(msg string) bool) (logs vcs.LogItems, err error) {
	err = vf.fakeCommand(LogCommand)

	if msgFilter == nil {
		logs = vf.settings.Logs
		return
	}

	for _, log := range vf.settings.Logs {
		if msgFilter(log.Message) {
			logs.Add(log)
		}
	}
	return
}

// Stash does nothing. Returns an error if in the list of failing commands
func (vf *VCSFake) Stash(_ string) error {
	return vf.fakeCommand(StashCommand)
}

// UnStash does nothing. Returns an error if in the list of failing commands
func (vf *VCSFake) UnStash(_ bool) error {
	return vf.fakeCommand(UnStashCommand)
}

// Revert does nothing. Returns an error if in the list of failing commands
func (vf *VCSFake) Revert() error {
	return vf.fakeCommand(RevertCommand)
}

// GetRootDir returns the root directory path
func (vf *VCSFake) GetRootDir() string {
	return ""
}

// GetRemoteName returns the current VCS remote name
func (vf *VCSFake) GetRemoteName() string {
	return ""
}

// GetWorkingBranch returns the current VCS working branch
func (vf *VCSFake) GetWorkingBranch() string {
	return ""
}

// IsOnRootBranch indicates if VCS is currently on its root branch or not
func (vf *VCSFake) IsOnRootBranch() bool {
	return true
}

// EnablePush sets a flag allowing to turn on/off VCS push operations
func (vf *VCSFake) EnablePush(flag bool) {
	vf.pushEnabled = flag
}

// IsPushEnabled indicates if VCS push operations are turned on
func (vf *VCSFake) IsPushEnabled() bool {
	return vf.pushEnabled
}

// IsRemoteEnabled indicates if VCS remote operations are enabled
func (*VCSFake) IsRemoteEnabled() bool {
	return false
}

// CheckRemoteAccess returns true if VCS remote can be accessed
func (*VCSFake) CheckRemoteAccess() bool {
	return true
}

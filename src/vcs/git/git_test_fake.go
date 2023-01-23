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

package git

import (
	"errors"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
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

// List of supported git commands
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
	// FakeSettings provide a few ways to tune Fake behaviour
	FakeSettings struct {
		FailingCommands Commands
		ChangedFiles    vcs.FileDiffs
		Logs            vcs.LogItems
	}

	// Fake provides a fake implementation of the git interface
	Fake struct {
		impl        vcs.Interface
		settings    FakeSettings
		lastCommand Command
	}
)

// inMemoryRepoInit initializes a brand new repository in memory (for use in tests)
func inMemoryRepoInit(_ string) (repo *git.Repository, fs billy.Filesystem, err error) {
	fs = memfs.New()
	repo, err = git.Init(memory.NewStorage(), fs)
	return
}

func (gf *Fake) fakeCommand(cmd Command) (err error) {
	gf.lastCommand = cmd
	if gf.settings.FailingCommands.contains(cmd) {
		err = errors.New("git " + string(cmd) + " fake error")
	}
	return
}

// NewFake initializes a fake git implementation which does nothing
// apart from emulating errors on git operations
func NewFake(settings FakeSettings) (*Fake, error) {
	impl, err := newGitImpl(inMemoryRepoInit, "")
	return &Fake{impl: impl, settings: settings}, err
}

func (gf *Fake) Name() string {
	return "fake-vcs"
}

func (gf *Fake) SessionSummary() string {
	return "VCS session \"fake\""
}

// GetLastCommand returns the last command called
func (gf *Fake) GetLastCommand() Command {
	return gf.lastCommand
}

// Add does nothing. Returns an error if in the list of failing commands
func (gf *Fake) Add(_ ...string) error {
	return gf.fakeCommand(AddCommand)
}

// Commit does nothing. Returns an error if in the list of failing commands
func (gf *Fake) Commit(_ bool, _ ...string) error {
	return gf.fakeCommand(CommitCommand)
}

// Restore does nothing. Returns an error if in the list of failing commands
func (gf *Fake) Restore(_ string) error {
	return gf.fakeCommand(RestoreCommand)
}

// Push does nothing. Returns an error if in the list of failing commands
func (gf *Fake) Push() error {
	return gf.fakeCommand(PushCommand)
}

// Pull does nothing. Returns an error if in the list of failing commands
func (gf *Fake) Pull() error {
	return gf.fakeCommand(PullCommand)
}

// Diff returns the list of files modified configured at fake initialization
func (gf *Fake) Diff() (_ vcs.FileDiffs, err error) {
	return gf.settings.ChangedFiles, gf.fakeCommand(DiffCommand)
}

// Log returns the list of git logs configured at fake initialization
func (gf *Fake) Log(msgFilter func(msg string) bool) (logs vcs.LogItems, err error) {
	err = gf.fakeCommand(LogCommand)

	if msgFilter == nil {
		logs = gf.settings.Logs
		return
	}

	for _, log := range gf.settings.Logs {
		if msgFilter(log.Message) {
			logs.Add(log)
		}
	}
	return
}

// Stash does nothing. Returns an error if in the list of failing commands
func (gf *Fake) Stash(_ string) error {
	return gf.fakeCommand(StashCommand)
}

// UnStash does nothing. Returns an error if in the list of failing commands
func (gf *Fake) UnStash(_ bool) error {
	return gf.fakeCommand(UnStashCommand)
}

// Revert does nothing. Returns an error if in the list of failing commands
func (gf *Fake) Revert() error {
	return gf.fakeCommand(RevertCommand)
}

// GetRootDir returns the root directory path
func (gf *Fake) GetRootDir() string {
	return gf.impl.GetRootDir()
}

// GetRemoteName returns the current git remote name
func (gf *Fake) GetRemoteName() string {
	return gf.impl.GetRemoteName()
}

// GetWorkingBranch returns the current git working branch
func (gf *Fake) GetWorkingBranch() string {
	return gf.impl.GetWorkingBranch()
}

// IsOnRootBranch indicates if git is currently on its root branch or not
func (gf *Fake) IsOnRootBranch() bool {
	return gf.impl.IsOnRootBranch()
}

// EnablePush sets a flag allowing to turn on/off git push operations
func (gf *Fake) EnablePush(flag bool) {
	gf.impl.EnablePush(flag)
}

// IsPushEnabled indicates if git push operations are turned on
func (gf *Fake) IsPushEnabled() bool {
	return gf.impl.IsPushEnabled()
}

// IsRemoteEnabled indicates if git remote operations are enabled
func (gf *Fake) IsRemoteEnabled() bool {
	return gf.impl.IsRemoteEnabled()
}

// CheckRemoteAccess returns true if git remote can be accessed
func (gf *Fake) CheckRemoteAccess() bool {
	return gf.impl.CheckRemoteAccess()
}

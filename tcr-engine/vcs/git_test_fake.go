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

package vcs

import (
	"errors"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

type (
	// GitCommand is the name of a Git command
	GitCommand string

	// GitCommands is a slice of GitCommand
	GitCommands []GitCommand
)

func (gc GitCommands) contains(command GitCommand) bool {
	for _, value := range gc {
		if value == command {
			return true
		}
	}
	return false
}

// List of supported git commands
const (
	AddCommand     GitCommand = "add"
	CommitCommand  GitCommand = "commit"
	DiffCommand    GitCommand = "diff"
	LogCommand     GitCommand = "log"
	PullCommand    GitCommand = "pull"
	PushCommand    GitCommand = "push"
	RestoreCommand GitCommand = "restore"
	RevertCommand  GitCommand = "revert"
	StashCommand   GitCommand = "stash"
	UnStashCommand GitCommand = "unStash"
)

type (
	// GitFakeSettings provide a few ways to tune GitFake behaviour
	GitFakeSettings struct {
		FailingCommands GitCommands
		ChangedFiles    FileDiffs
		Logs            GitLogItems
// GitFake provides a fake implementation of the git interface
type GitFake struct {
	impl            GitInterface
	failingCommands GitCommands
	changedFiles    []FileDiff
}

// inMemoryRepoInit initializes a brand new repository in memory (for use in tests)
func inMemoryRepoInit(_ string) (repo *git.Repository, fs billy.Filesystem, err error) {
	fs = memfs.New()
	repo, err = git.Init(memory.NewStorage(), fs)
	return
}

func (g GitFake) fakeGitCommand(cmd GitCommand) (err error) {
	if g.failingCommands.contains(cmd) {
		err = errors.New("git " + string(cmd) + " fake error")
	}

	// GitFake provides a fake implementation of the git interface
	GitFake struct {
		GitImpl
		settings GitFakeSettings
	}
)

// NewGitFake initializes a fake git implementation which does nothing
// apart from emulating errors on git operations
func NewGitFake(settings GitFakeSettings) GitInterface {
	return &GitFake{settings: settings}
}

func (g GitFake) fakeGitCommand(cmd GitCommand) (err error) {
	if g.settings.FailingCommands.contains(cmd) {
		err = errors.New("git " + string(cmd) + " fake error")
	}
	return
func NewGitFake(failingCommands GitCommands, changedFiles []FileDiff) (GitInterface, error) {
	impl, _ := newGitImpl(inMemoryRepoInit, "")
	return &GitFake{impl: impl, failingCommands: failingCommands, changedFiles: changedFiles}, nil
}

// Add does nothing. Returns an error if in the list of failing commands
func (g *GitFake) Add(_ ...string) error {
	return g.fakeGitCommand(AddCommand)
}

// Commit does nothing. Returns an error if in the list of failing commands
func (g GitFake) Commit(_ bool, _ ...string) error {
	return g.fakeGitCommand(CommitCommand)
}

// Restore does nothing. Returns an error if in the list of failing commands
func (g GitFake) Restore(_ string) error {
	return g.fakeGitCommand(RestoreCommand)
}

// Push does nothing. Returns an error if in the list of failing commands
func (g GitFake) Push() error {
	return g.fakeGitCommand(PushCommand)
}

// Pull does nothing. Returns an error if in the list of failing commands
func (g GitFake) Pull() error {
	return g.fakeGitCommand(PullCommand)
}

// Diff returns the list of files modified configured at fake initialization
func (g GitFake) Diff() (_ FileDiffs, err error) {
	return g.settings.ChangedFiles, g.fakeGitCommand(DiffCommand)
}

// Log returns the list of git logs configured at fake initialization
func (g GitFake) Log(_ func(msg string) bool) (logs GitLogItems, err error) {
	return g.settings.Logs, g.fakeGitCommand(LogCommand)
}

// Stash does nothing. Returns an error if in the list of failing commands
func (g GitFake) Stash(_ string) error {
	return g.fakeGitCommand(StashCommand)
}

// UnStash does nothing. Returns an error if in the list of failing commands
func (g GitFake) UnStash(_ bool) error {
	return g.fakeGitCommand(UnStashCommand)
}

// Revert does nothing. Returns an error if in the list of failing commands
func (g GitFake) Revert() error {
	return g.fakeGitCommand(RevertCommand)
}

// GetRootDir returns the root directory path
func (g GitFake) GetRootDir() string {
	return g.impl.GetRootDir()
}

// GetRemoteName returns the current git remote name
func (g GitFake) GetRemoteName() string {
	return g.impl.GetRemoteName()
}

// GetWorkingBranch returns the current git working branch
func (g GitFake) GetWorkingBranch() string {
	return g.impl.GetWorkingBranch()
}

// EnablePush sets a flag allowing to turn on/off git push operations
func (g GitFake) EnablePush(flag bool) {
	g.impl.EnablePush(flag)
}

// IsPushEnabled indicates if git push operations are turned on
func (g GitFake) IsPushEnabled() bool {
	return g.impl.IsPushEnabled()
}

// IsRemoteEnabled indicates if git remote operations are enabled
func (g GitFake) IsRemoteEnabled() bool {
	return g.impl.IsRemoteEnabled()
}

// CheckRemoteAccess returns true if git remote can be accessed
func (g GitFake) CheckRemoteAccess() bool {
	return g.impl.CheckRemoteAccess()
}

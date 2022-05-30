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
)

// GitCommand is the name of a Git command
type GitCommand string

// GitCommands is a slice of GitCommand
type GitCommands []GitCommand

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
	PullCommand    GitCommand = "pull"
	PushCommand    GitCommand = "push"
	RestoreCommand GitCommand = "restore"
	RevertCommand  GitCommand = "revert"
	StashCommand   GitCommand = "stash"
	UnStashCommand GitCommand = "unStash"
)

// GitFake provides a fake implementation of the git interface
type GitFake struct {
	GitImpl
	failingCommands GitCommands
	changedFiles    []FileDiff
}

func (g GitFake) fakeGitCommand(cmd GitCommand) (err error) {
	if g.failingCommands.contains(cmd) {
		err = errors.New("git " + string(cmd) + " fake error")
	}
	return
}

// NewGitFake initializes a fake git implementation which does nothing
// apart from emulating errors on git operations
func NewGitFake(failingCommands GitCommands, changedFiles []FileDiff) (GitInterface, error) {
	return &GitFake{failingCommands: failingCommands, changedFiles: changedFiles}, nil
}

// Add does nothing. Returns an error if in the list of failing commands
func (g *GitFake) Add(_ ...string) error {
	return g.fakeGitCommand(AddCommand)
}

// Commit does nothing. Returns an error if in the list of failing commands
func (g GitFake) Commit(_ bool, _ string) error {
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
func (g GitFake) Diff() (_ []FileDiff, err error) {
	return g.changedFiles, g.fakeGitCommand(DiffCommand)
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

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

// GitFake provides a fake implementation of the git interface
type GitFake struct {
	GitImpl
	failAdd      bool
	failCommit   bool
	failRestore  bool
	failPush     bool
	failPull     bool
	failDiff     bool
	failStash    bool
	failUnstash  bool
	failRevert   bool
	changedFiles []FileDiff
}

// NewGitFake initializes a fake git implementation which does nothing
// apart from emulating errors on git operations
func NewGitFake(
	failAdd, failCommit, failRestore, failPush, failPull, failDiff, failStash, failUnstash, failRevert bool,
	changedFiles []FileDiff) (GitInterface, error) {
	return &GitFake{
		failAdd:      failAdd,
		failCommit:   failCommit,
		failRestore:  failRestore,
		failPush:     failPush,
		failPull:     failPull,
		failDiff:     failDiff,
		failStash:    failStash,
		failUnstash:  failUnstash,
		failRevert:   failRevert,
		changedFiles: changedFiles,
	}, nil
}

// Add does nothing. Returns an error if failAdd flag is set
func (g *GitFake) Add(_ ...string) error {
	return fakeOperation("add", g.failAdd)
}

// Commit does nothing. Returns an error if failCommit flag is set
func (g GitFake) Commit(_ bool, _ string) error {
	return fakeOperation("commit", g.failCommit)
}

// Restore does nothing. Returns an error if failRestore flag is set
func (g GitFake) Restore(_ string) error {
	return fakeOperation("restore", g.failRestore)
}

// Push does nothing. Returns an error if failPush flag is set
func (g GitFake) Push() error {
	return fakeOperation("push", g.failPush)
}

// Pull does nothing. Returns an error if failPull flag is set
func (g GitFake) Pull() error {
	return fakeOperation("pull", g.failPull)
}

// Diff returns the list of files modified configured at fake initialization
func (g GitFake) Diff() (_ []FileDiff, err error) {
	return g.changedFiles, fakeOperation("diff", g.failDiff)
}

// Stash does nothing. Returns an error if failStash flag is set
func (g GitFake) Stash(_ string) error {
	return fakeOperation("stash", g.failStash)
}

// UnStash does nothing. Returns an error if failUnstash flag is set
func (g GitFake) UnStash(_ bool) error {
	return fakeOperation("unstash", g.failUnstash)
}

// Revert does nothing. Returns an error if failRevert flag is set
func (g GitFake) Revert() error {
	return fakeOperation("revert", g.failRevert)
}

func fakeOperation(operation string, shouldFail bool) (err error) {
	//fmt.Printf("faking git %v operation (failure=%v)\n", operation, shouldFail)
	if shouldFail {
		err = errors.New("git " + operation + " fake error")
	}
	return
}

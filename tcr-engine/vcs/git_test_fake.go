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
	failCommit  bool
	failRestore bool
	failPush    bool
	failPull    bool
}

// NewGitFake initializes a fake git implementation which does nothing
// apart from emulating errors on git operations
func NewGitFake(failCommit, failRestore, failPush, failPull bool) (GitInterface, error) {
	return &GitFake{
		failCommit:  failCommit,
		failRestore: failRestore,
		failPush:    failPush,
		failPull:    failPull,
	}, nil
}

// WorkingBranch returns an empty working branch
func (g GitFake) WorkingBranch() string {
	return ""
}

// Commit restores to last commit.
func (g GitFake) Commit() error {
	return fakeOperation("commit", g.failCommit)
}

// Restore restores to last commit for everything under dir.
func (g GitFake) Restore(_ string) error {
	return fakeOperation("restore", g.failRestore)
}

// Push runs a git push operation.
func (g GitFake) Push() error {
	return fakeOperation("push", g.failPush)
}

// Pull runs a git pull operation.
func (g GitFake) Pull() error {
	return fakeOperation("pull", g.failPull)
}

// EnablePush does nothing
func (g GitFake) EnablePush(_ bool) {
}

// IsPushEnabled always returns false
func (g GitFake) IsPushEnabled() bool {
	return false
}

func fakeOperation(operation string, shouldFail bool) (err error) {
	//fmt.Printf("faking git %v operation (failure=%v)\n", operation, shouldFail)
	if shouldFail {
		err = errors.New("git " + operation + " fake error")
	}
	return
}

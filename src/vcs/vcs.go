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

package vcs

const (
	// DefaultPushEnabled indicates the default state for VCS auto-push option
	DefaultPushEnabled = false
)

// Interface provides the interface that a VCS implementation must satisfy for TCR engine to be
// able to interact with it
type Interface interface {
	GetRootDir() string
	GetRemoteName() string
	GetWorkingBranch() string
	IsOnRootBranch() bool
	Add(paths ...string) error
	Commit(amend bool, messages ...string) error
	Restore(path string) error
	Revert() error
	Push() error
	Pull() error
	Stash(message string) error
	UnStash(keep bool) error
	Diff() (diffs FileDiffs, err error)
	Log(msgFilter func(msg string) bool) (logs LogItems, err error)
	EnablePush(flag bool)
	IsPushEnabled() bool
	IsRemoteEnabled() bool
	CheckRemoteAccess() bool
}

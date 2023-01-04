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
	"github.com/go-git/go-billy/v5"
	"github.com/murex/tcr/vcs"
)

// p4Impl provides the implementation of the Perforce interface
type p4Impl struct {
	baseDir                     string
	rootDir                     string
	filesystem                  billy.Filesystem // TODO - Needed?
	remoteName                  string           // TODO - Needed?
	remoteEnabled               bool             // TODO - Needed?
	workingBranch               string
	workingBranchExistsOnRemote bool // TODO - Needed?
	pushEnabled                 bool // TODO - Needed?
	runP4Function               func(params ...string) (output []byte, err error)
	traceP4Function             func(params ...string) (err error)
}

// New initializes the p4 implementation based on the provided directory from local clone
func New(dir string) (vcs.Interface, error) {
	return newP4Impl(dir)
}

func newP4Impl(dir string) (*p4Impl, error) {
	var p = p4Impl{
		baseDir:         dir,
		pushEnabled:     vcs.DefaultPushEnabled,
		runP4Function:   runP4Command,
		traceP4Function: traceP4Command,
	}

	// TODO p4 -F %clientRoot% -ztag info -> gives workspace root path
	var err error

	// TODO initialization
	//p.repository, p.filesystem, err = initRepo(dir)
	//if err != nil {
	//	return nil, err
	//}
	//
	//p.rootDir = retrieveRootDir(p.filesystem)
	//
	//p.workingBranch, err = retrieveWorkingBranch(p.repository)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if isRemoteDefined(DefaultRemoteName, p.repository) {
	//	p.remoteEnabled = true
	//	p.remoteName = DefaultRemoteName
	//	p.workingBranchExistsOnRemote, err = p.isWorkingBranchOnRemote()
	//}

	return &p, err
}

// GetRootDir returns the root directory path
func (p *p4Impl) GetRootDir() string {
	//TODO implement me
	panic("implement me")
}

// GetRemoteName returns the current p4 "remote name"
func (p *p4Impl) GetRemoteName() string {
	//TODO implement me
	panic("implement me")
}

// IsRemoteEnabled indicates if p4 remote operations are enabled
func (p *p4Impl) IsRemoteEnabled() bool {
	//TODO implement me
	panic("implement me")
}

// GetWorkingBranch returns the current p4 working branch
func (p *p4Impl) GetWorkingBranch() string {
	//TODO implement me
	panic("implement me")
}

// IsOnRootBranch indicates if p4 is currently on its root branch or not.
// Very likely meaningless in the case of p4
func (p *p4Impl) IsOnRootBranch() bool {
	//TODO implement me
	panic("implement me")
}

// Add adds the listed paths to p4 index.
// TODO: p4 add or p4 edit
func (p *p4Impl) Add(paths ...string) error {
	//TODO implement me
	panic("implement me")
}

// Commit commits changes to p4 index.
// TODO: p4 submit
func (p *p4Impl) Commit(amend bool, messages ...string) error {
	//TODO implement me
	panic("implement me")
}

// Restore restores to last commit for the provided path.
// TODO: p4 revert -c
func (p *p4Impl) Restore(path string) error {
	//TODO implement me
	panic("implement me")
}

// Revert runs a p4 revert operation.
// TODO: p4 revert
func (p *p4Impl) Revert() error {
	//TODO implement me
	panic("implement me")
}

// Push runs a p4 push operation.
// TODO: p4 submit?
func (p *p4Impl) Push() error {
	//TODO implement me
	panic("implement me")
}

// Pull runs a p4 pull operation.
// TODO: p4 sync
func (p *p4Impl) Pull() error {
	//TODO implement me
	panic("implement me")
}

// Stash creates a p4 stash.
// TODO: ???
func (p *p4Impl) Stash(message string) error {
	//TODO implement me
	panic("implement me")
}

// UnStash applies a p4 stash. Depending on the keep argument value, either a "stash apply" or a "stash pop"
// command is executed under the hood.
// TODO: ???
func (p *p4Impl) UnStash(keep bool) error {
	//TODO implement me
	panic("implement me")
}

// Diff returns the list of files modified since last commit with diff info for each file
// TODO: p4 diff
func (p *p4Impl) Diff() (diffs vcs.FileDiffs, err error) {
	//TODO implement me
	panic("implement me")
}

// Log returns the list of p4 log items compliant with the provided msgFilter.
// When no msgFilter is provided, returns all p4 log items unfiltered.
// TODO:  p4 changes ./...
func (p *p4Impl) Log(msgFilter func(msg string) bool) (logs vcs.LogItems, err error) {
	//TODO implement me
	panic("implement me")
}

// EnablePush sets a flag allowing to turn on/off p4 push operations
// TODO: ???
func (p *p4Impl) EnablePush(flag bool) {
	//TODO implement me
	panic("implement me")
}

// IsPushEnabled indicates if git push operations are turned on
// TODO: ???
func (p *p4Impl) IsPushEnabled() bool {
	//TODO implement me
	panic("implement me")
}

// CheckRemoteAccess returns true if p4 remote can be accessed.
// TODO: ???
func (p *p4Impl) CheckRemoteAccess() bool {
	//TODO implement me
	panic("implement me")
}

// traceP4 runs a p4 command and traces its output.
// The command is launched from the p4 root directory
func (p *p4Impl) traceP4(args ...string) error {
	return p.traceP4Function(append([]string{"-d", p.GetRootDir()}, args...)...)
}

// runP4 calls p4 command in a separate process and returns its output traces
// The command is launched from the p4 root directory
func (p *p4Impl) runP4(args ...string) (output []byte, err error) {
	return p.runP4Function(append([]string{"-d", p.GetRootDir()}, args...)...)
}

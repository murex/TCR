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
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/vcs"
	"github.com/spf13/afero"
)

// p4Impl provides the implementation of the Perforce interface
type p4Impl struct {
	baseDir         string
	rootDir         string
	filesystem      afero.Fs
	runP4Function   func(params ...string) (output []byte, err error)
	traceP4Function func(params ...string) (err error)
}

// New initializes the p4 implementation based on the provided directory from local clone
func New(dir string) (vcs.Interface, error) {
	return newP4Impl(plainOpen, dir)
}

func newP4Impl(initDepotFs func() afero.Fs, dir string) (*p4Impl, error) {
	var p = p4Impl{
		baseDir:         dir,
		filesystem:      initDepotFs(),
		runP4Function:   runP4Command,
		traceP4Function: traceP4Command,
	}
	err := p.retrieveRootDir()
	return &p, err
}

// plainOpen is the regular function used to open a p4 depot
func plainOpen() afero.Fs {
	return afero.NewOsFs()
}

// retrieveRootDir retrieves the local root directory for the depot's workspace
func (p *p4Impl) retrieveRootDir() error {
	root, err := runP4Command("-F", "%clientRoot%", "-ztag", "info")
	if err != nil {
		p.rootDir = ""
		return err
	}
	p.rootDir = string(root)
	return nil
}

// GetRootDir returns the root directory path
func (p *p4Impl) GetRootDir() string {
	return p.rootDir
}

// GetRemoteName returns the current p4 "remote name"
// TODO: clarify if there is some info that could be used as remote name (server Id maybe?)
func (*p4Impl) GetRemoteName() string {
	// For now, always return an empty string
	return ""
}

// IsRemoteEnabled indicates if p4 remote operations are enabled
// remote is always enabled with p4 due its architecture (all changes occur directly on the server)
func (*p4Impl) IsRemoteEnabled() bool {
	return true
}

// GetWorkingBranch returns the current p4 working branch
// TODO clarify if we need to handle p4 branches
func (*p4Impl) GetWorkingBranch() string {
	// For now, always return an empty string
	return ""
}

// IsOnRootBranch indicates if p4 is currently on its root branch or not.
// TODO clarify if we need to handle p4 branches
func (*p4Impl) IsOnRootBranch() bool {
	// For now, always return false
	return false
}

// Add adds the listed paths to p4 changelist.
func (p *p4Impl) Add(paths ...string) error {
	p4Args := []string{"reconcile", "-A"}
	if len(paths) == 0 {
		p4Args = append(p4Args, ".")
	} else {
		p4Args = append(p4Args, paths...)
	}
	return p.traceP4(p4Args...)
}

// Commit commits changes to p4 index.
// TODO: p4 submit
func (p *p4Impl) Commit(_ bool, messages ...string) error {
	// p4 --field "Description=My pending change" --field "Files=" change -o | p4 change -i
	// TODO - REQUIRED for simple TCR workflow
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

// Push runs a push operation.
// TODO: confirm that it does nothing as already submitted to the server through commit?
func (*p4Impl) Push() error {
	return nil
}

// Pull runs a pull operation ("p4 sync")
func (p *p4Impl) Pull() error {
	return p.traceP4("sync")
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
	// TODO - REQUIRED for simple TCR workflow
	panic("implement me")
}

// Log returns the list of p4 log items compliant with the provided msgFilter.
// When no msgFilter is provided, returns all p4 log items unfiltered.
// TODO:  p4 changes ./...
func (p *p4Impl) Log(msgFilter func(msg string) bool) (logs vcs.LogItems, err error) {
	//TODO implement me
	panic("implement me")
}

// EnablePush sets a flag allowing to turn on/off p4 push operations.
// Auto-push is always on with p4 due its architecture (all changes occur directly on the server)
func (*p4Impl) EnablePush(_ bool) {
	report.PostInfo("Perforce auto-push is always on")
}

// IsPushEnabled indicates if p4 push operations are turned on.
// Auto-push is always on with p4 due its architecture (all changes occur directly on the server)
func (*p4Impl) IsPushEnabled() bool {
	return true
}

// CheckRemoteAccess returns true if p4 remote can be accessed.
func (p *p4Impl) CheckRemoteAccess() bool {
	// TODO check if anything should be done here. Returning true for now
	return true
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

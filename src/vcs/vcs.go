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
	// DefaultAutoPushEnabled provides the default value for auto-push (off by default)
	DefaultAutoPushEnabled = false
	// DefaultTrace provides the default value for VCS commands trace (off by default)
	DefaultTrace = false // VCS trace is off by default
)

var trace bool

func init() {
	trace = DefaultTrace
}

// SetTrace turn on/off trace flag for the VCS package. When trace is on, all calls
// to VCS commands are traced by TCR
func SetTrace(flag bool) {
	trace = flag
}

// GetTrace returns VCS trace status
func GetTrace() bool {
	return trace
}

// Interface provides the interface that a VCS implementation must satisfy for TCR engine to be
// able to interact with it
type Interface interface {
	Name() string
	SessionSummary() string
	GetRootDir() string
	GetRemoteName() string
	GetWorkingBranch() string
	IsOnRootBranch() bool
	Add(paths ...string) error
	Commit(messages ...string) error
	RevertLocal(path string) error
	RollbackLastCommit() error
	Push() error
	Pull() error
	Diff() (diffs FileDiffs, err error)
	Log(msgFilter func(msg string) bool) (logs LogItems, err error)
	EnableAutoPush(flag bool)
	IsAutoPushEnabled() bool
	IsRemoteEnabled() bool
	CheckRemoteAccess() bool
	SupportsEmojis() bool
}

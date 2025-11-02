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

package factory

import (
	"fmt"
	"strings"

	"github.com/murex/tcr/vcs"
	"github.com/murex/tcr/vcs/git"
	"github.com/murex/tcr/vcs/p4"
)

// InitVCS returns the VCS instance of type defined by name, working on the provided directory
var InitVCS func(name string, dir string, remoteName string) (vcs.Interface, error)

func init() {
	// InitVCS is set by default to real VCS implementation factory.
	// This may be overridden when running tests to prevent going through VCS initialization
	// when running tests
	InitVCS = initVCS
}

// UnsupportedVCSError is returned when the provided VCS name is not supported by the factory.
type UnsupportedVCSError struct {
	vcsName string
}

// Error returns the error description
func (e *UnsupportedVCSError) Error() string {
	return fmt.Sprintf("VCS not supported: \"%s\"", e.vcsName)
}

// initVCS returns the VCS instance of type defined by name, working on the provided directory
func initVCS(name string, dir string, remoteName string) (vcs.Interface, error) {
	switch strings.ToLower(name) {
	case git.Name:
		return git.New(dir, remoteName)
	case p4.Name:
		return p4.New(dir)
	default:
		return nil, &UnsupportedVCSError{name}
	}
}

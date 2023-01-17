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

package shell

// Command provides the interface for a shell command.
type Command interface {
	Name() string
	Params() []string
	IsInPath() bool
	GetFullPath() string
	Run(params ...string) (output []byte, err error)
	Trace(params ...string) error
	RunAndPipe(toCmd Command, params ...string) (output []byte, err error)
	TraceAndPipe(toCmd Command, params ...string) error
}

// NewCommandFunc is the constructor used when creating a new command.
// It points by default to the real command implementation constructor.
// It's replaced in most of the tests by a stubbed command constructor allowing to
// bypass real shell commands execution (which are both time-consuming
// and depending on the environment)
var NewCommandFunc = NewCommand

//go:build test_helper

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

// CommandStub is a shell command stub provided for test purpose
type CommandStub struct {
	impl             CommandImpl
	IsInPathFunc     func() bool
	GetFullPathFunc  func() string
	RunFunc          func(params ...string) (output []byte, err error)
	TraceFunc        func(params ...string) error
	RunAndPipeFunc   func(toCmd Command, params ...string) (output []byte, err error)
	TraceAndPipeFunc func(toCmd Command, params ...string) error
}

// NewCommandStub creates a new shell command stub on top of a command implementation.
// CommandStub methods are set by default to use the command's real implementation methods.
// They can be overridden through setting its function attributes.
func NewCommandStub(impl CommandImpl) *CommandStub {
	return &CommandStub{
		impl:             impl,
		IsInPathFunc:     impl.IsInPath,
		GetFullPathFunc:  impl.GetFullPath,
		RunFunc:          impl.Run,
		TraceFunc:        impl.Trace,
		RunAndPipeFunc:   impl.RunAndPipe,
		TraceAndPipeFunc: impl.TraceAndPipe,
	}
}

// Name returns the command name (using the real implementation)
func (stub *CommandStub) Name() string {
	return stub.impl.Name()
}

// Params returns the parameters that the command will run with (using the real implementation)
func (stub *CommandStub) Params() []string {
	return stub.impl.Params()
}

// IsInPath indicates if the command can be found in the path.
func (stub *CommandStub) IsInPath() bool {
	return stub.IsInPathFunc()
}

// GetFullPath returns the full path for this command
func (stub *CommandStub) GetFullPath() string {
	return stub.GetFullPathFunc()
}

// String returns the command as a single string (including additional params if any)
func (stub *CommandStub) String(params ...string) string {
	return stub.impl.String(params...)
}

// Run calls the command with the provided parameters in a separate process and returns its output traces combined
func (stub *CommandStub) Run(params ...string) (output []byte, err error) {
	return stub.RunFunc(params...)
}

// Trace calls the command with the provided parameters and reports its output traces
func (stub *CommandStub) Trace(params ...string) error {
	return stub.TraceFunc(params...)
}

// RunAndPipe calls the command with the provided parameters in a separate process
// and pipes its output to cmd. Returns toCmd's output traces combined
func (stub *CommandStub) RunAndPipe(toCmd Command, params ...string) (output []byte, err error) {
	return stub.RunAndPipeFunc(toCmd, params...)
}

// TraceAndPipe calls the command with the provided parameters in a separate process
// and pipes its output to cmd. Reports toCmd's output traces
func (stub *CommandStub) TraceAndPipe(toCmd Command, params ...string) error {
	return stub.TraceAndPipeFunc(toCmd, params...)
}

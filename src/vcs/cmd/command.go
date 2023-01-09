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

package cmd

import (
	"github.com/codeskyblue/go-sh"
	"github.com/murex/tcr/report"
	"os/exec"
)

// ShellCommand is a command that can be launched from a shell
type ShellCommand struct {
	name   string
	params []string
}

// New creates a new shell command instance
func New(name string, params ...string) *ShellCommand {
	return &ShellCommand{name: name, params: params}
}

// IsInPath indicates if the command can be found in the path
func (sc *ShellCommand) IsInPath() bool {
	_, err := exec.LookPath(sc.name)
	return err == nil
}

// GetFullPath returns the full path for this command
func (sc *ShellCommand) GetFullPath() string {
	path, _ := exec.LookPath(sc.name)
	return path
}

// Run calls the command with the provided parameters in a separate process and returns its output traces combined
func (sc *ShellCommand) Run(params ...string) (output []byte, err error) {
	//report.PostWarning("Command: ", sc.name, " ", append(sc.params, params...))
	return sh.Command(sc.name, append(sc.params, params...)).CombinedOutput()
}

// Trace calls the command with the provided parameters and reports its output traces
func (sc *ShellCommand) Trace(params ...string) error {
	output, err := sc.Run(params...)
	if len(output) > 0 {
		report.PostText(string(output))
	}
	return err
}

// RunAndPipe calls the command with the provided parameters in a separate process
// and pipes its output to cmd. Returns cmd's output traces combined
func (sc *ShellCommand) RunAndPipe(cmd *ShellCommand, params ...string) (output []byte, err error) {
	//report.PostWarning("Command: ", sc.name, " ", append(sc.params, params...), " | ", cmd.name, " ", cmd.params)
	return sh.NewSession().
		Command(sc.name, append(sc.params, params...)).
		Command(cmd.name, cmd.params).
		CombinedOutput()
}

// TraceAndPipe calls the command with the provided parameters in a separate process
// and pipes its output to cmd. Reports cmd's output traces
func (sc *ShellCommand) TraceAndPipe(cmd *ShellCommand, params ...string) error {
	output, err := sc.RunAndPipe(cmd, params...)
	if len(output) > 0 {
		report.PostText(string(output))
	}
	return err
}

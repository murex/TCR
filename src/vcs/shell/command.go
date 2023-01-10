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

package shell

import (
	"github.com/codeskyblue/go-sh"
	"github.com/murex/tcr/report"
	"os/exec"
)

// Command is a command that can be launched from a shell
type Command struct {
	name   string
	params []string
}

// NewCommand creates a new shell command instance
func NewCommand(name string, params ...string) *Command {
	return &Command{name: name, params: params}
}

// IsInPath indicates if the command can be found in the path
func (c *Command) IsInPath() bool {
	_, err := exec.LookPath(c.name)
	return err == nil
}

// GetFullPath returns the full path for this command
func (c *Command) GetFullPath() string {
	path, _ := exec.LookPath(c.name)
	return path
}

// Run calls the command with the provided parameters in a separate process and returns its output traces combined
func (c *Command) Run(params ...string) (output []byte, err error) {
	//report.PostWarning("Command: ", c.name, " ", append(c.params, params...))
	return sh.Command(c.name, append(c.params, params...)).CombinedOutput()
}

// Trace calls the command with the provided parameters and reports its output traces
func (c *Command) Trace(params ...string) error {
	output, err := c.Run(params...)
	if len(output) > 0 {
		report.PostText(string(output))
	}
	return err
}

// RunAndPipe calls the command with the provided parameters in a separate process
// and pipes its output to cmd. Returns toCmd's output traces combined
func (c *Command) RunAndPipe(toCmd *Command, params ...string) (output []byte, err error) {
	//report.PostWarning("Command: ", c.name, " ", append(c.params, params...), " | ", shell.name, " ", shell.params)
	return sh.NewSession().
		Command(c.name, append(c.params, params...)).
		Command(toCmd.name, toCmd.params).
		CombinedOutput()
}

// TraceAndPipe calls the command with the provided parameters in a separate process
// and pipes its output to cmd. Reports toCmd's output traces
func (c *Command) TraceAndPipe(toCmd *Command, params ...string) error {
	output, err := c.RunAndPipe(toCmd, params...)
	if len(output) > 0 {
		report.PostText(string(output))
	}
	return err
}

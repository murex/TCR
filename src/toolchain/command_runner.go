/*
Copyright (c) 2024 Murex

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

package toolchain

import (
	"github.com/codeskyblue/go-sh"
	"github.com/murex/tcr/report"
	"syscall"
)

// CommandRunner is in charge of managing the lifecycle of a command
type CommandRunner struct {
	session *sh.Session
}

// commandRunner singleton instance
var commandRunner = getCommandRunner()

func getCommandRunner() *CommandRunner {
	return &CommandRunner{}
}

// Run launches the execution of the provided command
func (r *CommandRunner) Run(cmd *Command) (result CommandResult) {
	result = CommandResult{Status: CommandStatusUnknown, Output: ""}
	report.PostText(cmd.asCommandLine())

	r.session = sh.NewSession().SetDir(GetWorkDir())
	outputBytes, err := r.session.Command(cmd.Path, cmd.Arguments).CombinedOutput()

	if err == nil {
		result.Status = CommandStatusPass
	} else {
		result.Status = CommandStatusFail
	}

	if outputBytes != nil {
		result.Output = string(outputBytes)
		report.PostText(result.Output)
	}
	r.session = nil
	return result
}

// AbortRunningCommand triggers aborting of any command that is currently running
func (r *CommandRunner) AbortRunningCommand() {
	if r.session == nil {
		report.PostWarning("There is no command running at this time")
		return
	}
	report.PostWarning("Aborting running command")
	r.session.Kill(syscall.SIGKILL)
	r.session = nil
}

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
	"bufio"
	"github.com/murex/tcr/report"
	"os/exec"
)

type (
	// CommandRunner is in charge of managing the lifecycle of a command
	CommandRunner struct {
		command *exec.Cmd
	}

	// CommandStatus is the result status of a Command execution
	CommandStatus string

	// CommandResult contains the result from running a Command
	// - Status
	CommandResult struct {
		Status CommandStatus
		Output string
	}
)

// Failed indicates is a Command failed
func (r CommandResult) Failed() bool {
	return r.Status == CommandStatusFail
}

// Passed indicates is a Command passed
func (r CommandResult) Passed() bool {
	return r.Status == CommandStatusPass
}

// List of possible values for CommandStatus
const (
	CommandStatusPass    CommandStatus = "pass"
	CommandStatusFail    CommandStatus = "fail"
	CommandStatusUnknown CommandStatus = "unknown"
)

// commandRunner singleton instance
var commandRunner = getCommandRunner()

func getCommandRunner() *CommandRunner {
	return &CommandRunner{
		command: nil,
	}
}

// Run launches the execution of the provided command
func (r *CommandRunner) Run(cmd *Command) (result CommandResult) {
	result = CommandResult{Status: CommandStatusUnknown, Output: ""}
	report.PostText(cmd.asCommandLine())

	// Prepare the command
	r.command = exec.Command(cmd.Path, cmd.Arguments...) //nolint:gosec
	r.command.Dir = GetWorkDir()

	// Allow simultaneous trace and capture of command's stdout and stderr
	outReader, _ := r.command.StdoutPipe()
	errReader, _ := r.command.StderrPipe()
	doneOut, doneErr := make(chan bool), make(chan bool)
	r.reportCommandTrace(bufio.NewScanner(outReader), doneOut)
	r.reportCommandTrace(bufio.NewScanner(errReader), doneErr)

	// Start the command asynchronously
	errStart := r.command.Start()
	if errStart != nil {
		report.PostError("Failed to run command: ", errStart.Error())
		r.command = nil
		return result
	}

	// Wait for the command to finish
	errWait := r.command.Wait()
	if errWait != nil {
		result.Status = CommandStatusFail
	} else {
		result.Status = CommandStatusPass
	}

	r.command = nil
	return result
}

func (*CommandRunner) reportCommandTrace(scanner *bufio.Scanner, done chan bool) {
	go func() {
		for scanner.Scan() {
			report.PostText(scanner.Text())
		}
		done <- true
	}()
}

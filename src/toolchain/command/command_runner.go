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

package command

import (
	"bufio"
	"io"
	"os/exec"
	"sync"

	"github.com/murex/tcr/report"
)

type (
	// Runner is in charge of managing the lifecycle of a command
	Runner struct {
		command *exec.Cmd
		// commandMutex is here to enforce that commands run in sequence
		commandMutex sync.Mutex
	}

	// Status is the result status of a Command execution
	Status string

	// Result contains the result from running a Command
	Result struct {
		Status Status
		Output string
	}
)

// Failed indicates is a Command failed
func (r Result) Failed() bool {
	return r.Status == StatusFail
}

// Passed indicates is a Command passed
func (r Result) Passed() bool {
	return r.Status == StatusPass
}

// List of possible values for Status
const (
	StatusPass    Status = "pass"
	StatusFail    Status = "fail"
	StatusUnknown Status = "unknown"
)

// runner singleton instance
var runner = &Runner{
	command: nil,
}

// GetRunner returns the command runner singleton instance
func GetRunner() *Runner {
	return runner
}

// Run launches the execution of the provided command
func (r *Runner) Run(fromDir string, cmd *Command) (result Result) {
	runner.commandMutex.Lock()
	result = Result{Status: StatusUnknown, Output: ""}
	report.PostText(cmd.AsCommandLine())

	// Prepare the command
	r.command = exec.Command(cmd.Path, cmd.Arguments...) //nolint:gosec
	if fromDir != "" {
		r.command.Dir = fromDir
	}

	// Allow simultaneous trace and capture of command's stdout and stderr
	outReader, _ := r.command.StdoutPipe()
	errReader, _ := r.command.StderrPipe()
	r.reportTrace(outReader)
	r.reportTrace(errReader)

	// Start the command asynchronously
	errStart := r.command.Start()
	if errStart != nil {
		report.PostError("Failed to run command: ", errStart.Error())
		// We currently return fail status when command cannot be launched.
		// This is to replicate previous implementation's behaviour where
		// we could not differentiate between failure to launch and failure from execution.
		// We could later on use a different return value in this situation,
		// but we need to ensure first that TCR engine can handle it correctly.
		result.Status = StatusFail
		r.command = nil
		runner.commandMutex.Unlock()
		return result
	}

	// Wait for the command to finish
	errWait := r.command.Wait()
	if errWait != nil {
		result.Status = StatusFail
	} else {
		result.Status = StatusPass
	}

	r.command = nil
	runner.commandMutex.Unlock()
	return result
}

func (*Runner) reportTrace(readCloser io.ReadCloser) {
	scanner := bufio.NewScanner(readCloser)
	go func() {
		for scanner.Scan() {
			report.PostText(scanner.Text())
		}
	}()
}

// AbortRunningCommand triggers aborting of any command that is currently running
func (r *Runner) AbortRunningCommand() bool {
	if r.command == nil || r.command.Process == nil {
		report.PostWarning("There is no command running at this time")
		return false
	}
	report.PostWarning("Aborting command: \"", r.command.String(), "\"")
	_ = r.command.Process.Kill()
	// Calling Kill() may be a bit too brutal (may leave children process alive)
	//_ = r.command.Process.Signal(os.Kill)
	return true
}

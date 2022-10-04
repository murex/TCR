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

package cli

import (
	"bytes"
	"github.com/codeskyblue/go-sh"
	"github.com/murex/tcr/tcr-engine/report"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// sttyCmdDisabled allows to turn on/off the underlying call to stty command.
// It's false by default (regular usage). Most test methods should turn it on
// as a command call is generally useless when running tests, and time-consuming.
var sttyCmdDisabled bool

// tputCmdDisabled allows to turn on/off the underlying call to tput command.
// We try to run tput once, and if not available from the terminal at hand,
// we stop calling it and use default terminal width instead.
var tputCmdDisabled bool

func init() {
	sttyCmdDisabled = false
	tputCmdDisabled = false
}

func readStty(state *bytes.Buffer) (err error) {
	if sttyCmdDisabled {
		return nil
	}
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	cmd.Stdout = state
	return cmd.Run()
}

func setStty(state *bytes.Buffer) (err error) {
	if sttyCmdDisabled {
		return nil
	}
	cmd := exec.Command("stty", state.String()) //nolint:gosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	//report.PostInfo("Command: ", cmd)
	return cmd.Run()
}

// SetRaw changes the terminal state to "raw" mode
func SetRaw() bytes.Buffer {
	var initialState bytes.Buffer
	err := readStty(&initialState)
	//report.PostInfo(initialState.String())
	if err != nil {
		report.PostError("stty -g: ", err)
	}

	cbreakErr := setStty(bytes.NewBufferString("cbreak"))
	if cbreakErr != nil {
		report.PostError("stty cbreak: ", cbreakErr)
	}

	echoErr := setStty(bytes.NewBufferString("-echo"))
	if echoErr != nil {
		report.PostError("stty -echo: ", echoErr)
	}

	return initialState
}

// Restore puts back the terminal state to a "normal" state
func Restore() {
	//func Restore(state *bytes.Buffer)
	// For some unknown reason restoring previous stty state
	// fails on WSL, while working correctly on Windows git bash
	// Still need to test it on macOS and on a non-WSL Linux box
	// Until then we set back echo and -cbreak instead of
	// restoring the previous state
	//report.PostInfo("Restoring stty initial state")
	_ = setStty(bytes.NewBufferString("-cbreak"))
	_ = setStty(bytes.NewBufferString("echo"))
}

// getTerminalColumns returns the terminal's current number of column. If anything goes wrong (for
// example when running from Windows PowerShell), we fall back on a fixed number of columns
func getTerminalColumns() int {
	if tputCmdDisabled {
		return defaultTerminalWidth
	}
	output, err := sh.Command("tput", "cols").Output()
	if err != nil {
		tputCmdDisabled = true
		return defaultTerminalWidth
	}
	columns, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		tputCmdDisabled = true
		return defaultTerminalWidth
	}
	return columns
}

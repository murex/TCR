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
	"github.com/kami-zh/go-capturer"
	"github.com/murex/tcr/tcr-engine/engine"
	"github.com/murex/tcr/tcr-engine/report"
	"github.com/murex/tcr/tcr-engine/role"
	"github.com/murex/tcr/tcr-engine/runmode"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func Test_confirm_with_default_answer_to_yes(t *testing.T) {
	assertConfirmBehaviour(t, []byte{enterKey}, true, true)
}

func Test_confirm_with_default_answer_to_no(t *testing.T) {
	assertConfirmBehaviour(t, []byte{enterKey}, false, false)
}

func Test_confirm_with_a_yes_answer(t *testing.T) {
	assertConfirmBehaviour(t, []byte{'y'}, false, true)
	assertConfirmBehaviour(t, []byte{'Y'}, false, true)
}

func Test_confirm_with_a_no_answer(t *testing.T) {
	assertConfirmBehaviour(t, []byte{'n'}, true, false)
	assertConfirmBehaviour(t, []byte{'N'}, true, false)
}

func Test_confirm_question_with_default_answer_to_no(t *testing.T) {
	assert.Equal(t, "[y/N]", yesOrNoAdvice(false))
}

func Test_confirm_question_with_default_answer_to_yes(t *testing.T) {
	assert.Equal(t, "[Y/n]", yesOrNoAdvice(true))
}

func assertConfirmBehaviour(t *testing.T, input []byte, defaultValue bool, expected bool) {
	stdin := os.Stdin
	stdout := os.Stdout
	// Restore stdin and stdout right after the test.
	defer func() { os.Stdin = stdin; os.Stdout = stdout }()
	// We fake stdin so that we can simulate a key press
	os.Stdin = fakeStdin(t, input)
	// Displayed info on stdout is not used in the test
	os.Stdout = os.NewFile(0, os.DevNull)

	term := New(engine.Params{}, engine.NewTcrEngine())
	assert.Equal(t, expected, term.Confirm("", defaultValue))
}

func fakeStdin(t *testing.T, input []byte) *os.File {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}
	_ = w.Close()
	return r
}

func asCyanTrace(str string) string {
	return "\x1b[36mTCR\x1b[0m \x1b[36m" + str + "\x1b[0m\n"
}

func asCyanTraceWithSeparatorLine(str string) string {
	return asCyanTrace("---------------------------------------------------------------------------") +
		asCyanTrace(str)
}

func asYellowTrace(str string) string {
	return "\x1b[33mTCR\x1b[0m \x1b[33m" + str + "\x1b[0m\n"
}

func asRedTrace(str string) string {
	return "\x1b[31mTCR\x1b[0m \x1b[31m" + str + "\x1b[0m\n"
}

func asGreenTrace(str string) string {
	return "\x1b[32mTCR\x1b[0m \x1b[32m" + str + "\x1b[0m\n"
}

func asNeutralTrace(str string) string {
	return str + "\n"
}

func Test_terminal_tracing_methods(t *testing.T) {
	term := TerminalUI{}
	setLinePrefix("TCR")

	var testFlags = []struct {
		desc     string
		method   func()
		expected string
	}{
		{
			"info method",
			func() {
				term.info("Some info message")
			},
			asCyanTrace("Some info message"),
		},
		{
			"warning method",
			func() {
				term.warning("Some warning message")
			},
			asYellowTrace("Some warning message"),
		},
		{
			"error method",
			func() {
				term.error("Some error message")
			},
			asRedTrace("Some error message"),
		},
		{
			"trace method",
			func() {
				term.trace("Some trace message")
			},
			asNeutralTrace("Some trace message"),
		},
		{
			"title method",
			func() {
				term.title("Some title")
			},
			asCyanTraceWithSeparatorLine("Some title"),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expected, capturer.CaptureStdout(tt.method))
		})
	}
}

func Test_notify_role_starting(t *testing.T) {
	term := TerminalUI{}
	setLinePrefix("TCR")

	var testFlags = []struct {
		currentRole role.Role
		expected    string
	}{
		{
			currentRole: role.Driver{},
			expected:    asCyanTraceWithSeparatorLine("Starting with Driver role. Press ? for options"),
		},
		{
			currentRole: role.Navigator{},
			expected:    asCyanTraceWithSeparatorLine("Starting with Navigator role. Press ? for options"),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.currentRole.Name(), func(t *testing.T) {
			assert.Equal(t, tt.expected, capturer.CaptureStdout(func() {
				term.NotifyRoleStarting(tt.currentRole)
			}))
		})
	}
}

func Test_notify_role_ending(t *testing.T) {
	term := TerminalUI{}
	setLinePrefix("TCR")

	var testFlags = []struct {
		currentRole role.Role
		expected    string
	}{
		{
			currentRole: role.Driver{},
			expected:    asCyanTrace("Ending Driver role"),
		},
		{
			currentRole: role.Navigator{},
			expected:    asCyanTrace("Ending Navigator role"),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.currentRole.Name(), func(t *testing.T) {
			assert.Equal(t, tt.expected, capturer.CaptureStdout(func() {
				term.NotifyRoleEnding(tt.currentRole)
			}))
		})
	}
}

func Test_show_running_mode(t *testing.T) {
	term := TerminalUI{}
	setLinePrefix("TCR")

	var testFlags = []struct {
		currentMode runmode.RunMode
		expected    string
	}{
		{
			currentMode: runmode.Mob{},
			expected:    asCyanTraceWithSeparatorLine("Running in mob mode"),
		},
		{
			currentMode: runmode.Solo{},
			expected:    asCyanTraceWithSeparatorLine("Running in solo mode"),
		},
		{
			currentMode: runmode.OneShot{},
			expected:    asCyanTraceWithSeparatorLine("Running in one-shot mode"),
		},
		{
			currentMode: runmode.Check{},
			expected:    asCyanTraceWithSeparatorLine("Running in check mode"),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.currentMode.Name(), func(t *testing.T) {
			assert.Equal(t, tt.expected, capturer.CaptureStdout(func() {
				term.ShowRunningMode(tt.currentMode)
			}))
		})
	}
}

func Test_terminal_reporting(t *testing.T) {
	term := TerminalUI{}
	setLinePrefix("TCR")

	var testFlags = []struct {
		desc     string
		method   func()
		expected string
	}{
		{
			"PostInfo method",
			func() { report.PostInfo("Some info report") },
			asCyanTrace("Some info report"),
		},
		{
			"PostWarning method",
			func() { report.PostWarning("Some warning report") },
			asYellowTrace("Some warning report"),
		},
		{
			"PostError method",
			func() { report.PostWarning("Some error report") },
			asYellowTrace("Some error report"),
		},
		{
			"PostTitle method",
			func() { report.PostTitle("Some title report") },
			asCyanTraceWithSeparatorLine("Some title report"),
		},
		{
			"PostText method",
			func() { report.PostText("Some text report") },
			asNeutralTrace("Some text report"),
		},
		{
			"PostNotification method",
			func() { report.PostNotification("Some notification report") },
			asGreenTrace("Some notification report"),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expected, capturer.CaptureStdout(func() {
				term.StartReporting()
				time.Sleep(1 * time.Millisecond)
				tt.method()
				time.Sleep(1 * time.Millisecond)
				term.StopReporting()
			}))
		})
	}
}

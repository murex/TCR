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
	"github.com/murex/tcr/desktop"
	"github.com/murex/tcr/engine"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/role"
	"github.com/murex/tcr/runmode"
	"github.com/murex/tcr/vcs/git"
	"github.com/murex/tcr/vcs/p4"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
	"os"
	"strings"
	"testing"
	"time"
)

func init() {
	// turn off the dynamic retriever of terminal width when running tests
	tputCmdDisabled = true
}

func slowTestTag(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

func Test_confirm_answer(t *testing.T) {
	testFlags := []struct {
		desc         string
		input        []byte
		defaultValue bool
		expected     bool
	}{
		{"Enter key with Default Yes", []byte{enterKey}, true, true},
		{"Enter key with Default No", []byte{enterKey}, false, false},
		{"Y key", []byte{'Y'}, false, true},
		{"y key", []byte{'y'}, false, true},
		{"N key", []byte{'N'}, true, false},
		{"n key", []byte{'n'}, true, false},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assertConfirmBehaviour(t, tt.input, tt.defaultValue, tt.expected)
		})
	}
}

func assertConfirmBehaviour(t *testing.T, input []byte, defaultValue bool, expected bool) {
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr
	// Restore stdin, stdout and stderr right after the test.
	defer func() { os.Stdin = stdin; os.Stdout = stdout; os.Stderr = stderr }()
	// We fake stdin so that we can simulate a key press
	os.Stdin = fakeStdin(t, input)
	// Displayed info on stdout and stderr is not used in the test
	os.Stdout = os.NewFile(0, os.DevNull)
	os.Stderr = os.NewFile(0, os.DevNull)

	term := New(params.Params{}, engine.NewTCREngine())
	sttyCmdDisabled = true
	assert.Equal(t, expected, term.Confirm("", defaultValue))
	sttyCmdDisabled = false
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

func Test_confirm_question_with_default_answer_to_no(t *testing.T) {
	assert.Equal(t, "[y/N]", yesOrNoAdvice(false))
}

func Test_confirm_question_with_default_answer_to_yes(t *testing.T) {
	assert.Equal(t, "[Y/n]", yesOrNoAdvice(true))
}

func terminalSetup(p params.Params) (term *TerminalUI, fakeEngine *engine.FakeTCREngine, fakeNotifier *desktop.FakeNotifier) {
	setLinePrefix("TCR")
	fakeEngine = engine.NewFakeTCREngine()
	fakeNotifier = &desktop.FakeNotifier{}
	term = &TerminalUI{params: p, tcr: fakeEngine, desktop: desktop.NewDesktop(fakeNotifier)}
	term.mainMenu = term.initMainMenu()
	term.roleMenu = term.initRoleMenu()
	sttyCmdDisabled = true
	report.Reset()
	term.StartReporting()
	return
}

func terminalTeardown(term TerminalUI) {
	term.StopReporting()
	sttyCmdDisabled = false
	term.MuteDesktopNotifications(false)
}

func asCyanTrace(str string) string {
	return "\x1b[36mTCR\x1b[0m \x1b[36m" + str + "\x1b[0m\n"
}

func asCyanTraceWithSeparatorLine(str string) string {
	return asCyanTrace(strings.Repeat(horizontalLineCharacter, 75)) +
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
	var term TerminalUI

	var testFlags = []struct {
		desc     string
		method   func()
		expected string
	}{
		{
			"info method",
			func() {
				term.ReportInfo(false, "Some info message")
			},
			asCyanTrace("Some info message"),
		},
		{
			"warning method",
			func() {
				term.ReportWarning(false, "Some warning message")
			},
			asYellowTrace("Some warning message"),
		},
		{
			"error method",
			func() {
				term.ReportError(false, "Some error message")
			},
			asRedTrace("Some error message"),
		},
		{
			"trace method",
			func() {
				term.ReportSimple(false, "Some trace message")
			},
			asNeutralTrace("Some trace message"),
		},
		{
			"title method",
			func() {
				term.ReportTitle(false, "Some title")
			},
			asCyanTraceWithSeparatorLine("Some title"),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			term, _, _ := terminalSetup(*params.AParamSet())
			assert.Equal(t, tt.expected, capturer.CaptureStdout(tt.method))
			terminalTeardown(*term)
		})
	}
}

func Test_notify_role_starting(t *testing.T) {
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
			term, _, _ := terminalSetup(*params.AParamSet())
			assert.Equal(t, tt.expected, capturer.CaptureStdout(func() {
				term.NotifyRoleStarting(tt.currentRole)
			}))
			terminalTeardown(*term)
		})
	}
}

func Test_notify_role_ending(t *testing.T) {
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
			term, _, _ := terminalSetup(*params.AParamSet())
			assert.Equal(t, tt.expected, capturer.CaptureStdout(func() {
				term.NotifyRoleEnding(tt.currentRole)
			}))
			terminalTeardown(*term)
		})
	}
}

func Test_list_role_menu_options(t *testing.T) {
	title := "some title"
	var testFlags = []struct {
		currentRole role.Role
		expected    string
	}{
		{
			currentRole: role.Driver{},
			expected: asCyanTraceWithSeparatorLine(title) +
				asCyanTrace("\tT "+menuArrow+" "+timerStatusMenuHelper) +
				asCyanTrace("\tQ "+menuArrow+" "+quitDriverRoleMenuHelper) +
				asCyanTrace("\t? "+menuArrow+" "+optionsMenuHelper),
		},
		{
			currentRole: role.Navigator{},
			expected: asCyanTraceWithSeparatorLine(title) +
				asCyanTrace("\tQ "+menuArrow+" "+quitNavigatorRoleMenuHelper) +
				asCyanTrace("\t? "+menuArrow+" "+optionsMenuHelper),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.currentRole.Name(), func(t *testing.T) {
			term, _, _ := terminalSetup(*params.AParamSet())
			_ = term.runTCR(tt.currentRole)
			assert.Equal(t, tt.expected, capturer.CaptureStdout(func() {
				term.listMenuOptions(term.roleMenu, title)
			}))
			terminalTeardown(*term)
		})
	}
}

func Test_simple_message_methods(t *testing.T) {
	var term *TerminalUI
	var testFlags = []struct {
		desc     string
		method   func()
		expected string
	}{
		{
			desc:     "keyNotRecognizedMessage",
			method:   term.keyNotRecognizedMessage,
			expected: asYellowTrace("Key not recognized. Press ? for available options"),
		},
		{
			desc: "whatShallWeDo",
			method: func() {
				term.whatShallWeDo()
			},
			expected: asCyanTraceWithSeparatorLine("What shall we do?") +
				asCyanTrace("\tD "+menuArrow+" "+enterDriverRoleMenuHelper) +
				asCyanTrace("\tN "+menuArrow+" "+enterNavigatorRoleMenuHelper) +
				asCyanTrace("\tP "+menuArrow+" "+gitAutoPushMenuHelper) +
				asCyanTrace("\tL "+menuArrow+" "+pullMenuHelper) +
				asCyanTrace("\tS "+menuArrow+" "+pushMenuHelper) +
				asCyanTrace("\tQ "+menuArrow+" "+quitMenuHelper) +
				asCyanTrace("\t? "+menuArrow+" "+optionsMenuHelper),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			term, _, _ = terminalSetup(*params.AParamSet())
			assert.Equal(t, tt.expected, capturer.CaptureStdout(tt.method))
			terminalTeardown(*term)
		})
	}
}

func Test_show_running_mode(t *testing.T) {
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
		{
			currentMode: runmode.Log{},
			expected:    asCyanTraceWithSeparatorLine("Running in log mode"),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.currentMode.Name(), func(t *testing.T) {
			term, _, _ := terminalSetup(*params.AParamSet())
			assert.Equal(t, tt.expected, capturer.CaptureStdout(func() {
				term.ShowRunningMode(tt.currentMode)
			}))
			terminalTeardown(*term)
		})
	}
}

func Test_terminal_reporting(t *testing.T) {
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
			func() { report.PostError("Some error report") },
			asRedTrace("Some error report"),
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
			"PostTimerWithEmphasis method",
			func() {
				report.PostTimerWithEmphasis("Some timer with emphasis report")
			},
			asGreenTrace("Some timer with emphasis report"),
		},
		{
			"PostSuccessWithEmphasis method",
			func() {
				report.PostSuccessWithEmphasis("Some success with emphasis report")
			},
			asGreenTrace("Some success with emphasis report"),
		},
		{
			"PostWarningWithEmphasis method",
			func() {
				report.PostWarningWithEmphasis("Some warning with emphasis report")
			},
			asYellowTrace("Some warning with emphasis report"),
		},
		{
			"PostErrorWithEmphasis method",
			func() {
				report.PostErrorWithEmphasis("Some error with emphasis report")
			},
			asRedTrace("Some error with emphasis report"),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expected, capturer.CaptureStdout(func() {
				term, _, _ := terminalSetup(*params.AParamSet())
				time.Sleep(1 * time.Millisecond)
				tt.method()
				terminalTeardown(*term)
			}))
		})
	}
}

func Test_terminal_notification_box_title(t *testing.T) {
	var testFlags = []struct {
		desc     string
		method   func(a ...interface{})
		expected string
	}{
		{
			"timer with emphasis",
			report.PostTimerWithEmphasis,
			"⏳ TCR",
		},
		{
			"success with emphasis",
			report.PostSuccessWithEmphasis,
			"🟢 TCR",
		},
		{
			"warning with emphasis",
			report.PostWarningWithEmphasis,
			"🔶 TCR",
		},
		{
			"error with emphasis",
			report.PostErrorWithEmphasis,
			"🟥 TCR",
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			term, _, fakeNotifier := terminalSetup(*params.AParamSet())
			tt.method(tt.desc)
			time.Sleep(1 * time.Millisecond)
			assert.Equal(t, tt.expected, fakeNotifier.LastTitle)
			terminalTeardown(*term)

		})
	}
}

func Test_show_session_info(t *testing.T) {
	expected := asCyanTraceWithSeparatorLine("Base Directory: fake") +
		asCyanTrace("Work Directory: fake") +
		asCyanTrace("Language=fake, Toolchain=fake") +
		asYellowTrace("VCS \"fake\" is unknown")

	assert.Equal(t, expected, capturer.CaptureStdout(func() {
		term, _, _ := terminalSetup(*params.AParamSet())
		term.ShowSessionInfo()
		terminalTeardown(*term)
	}))
}

func Test_report_vcs_info(t *testing.T) {
	tests := []struct {
		desc     string
		info     engine.SessionInfo
		expected string
	}{
		{
			"VCS not set",
			engine.SessionInfo{VCSName: "", VCSSessionSummary: "", GitAutoPush: false},
			asYellowTrace("VCS \"\" is unknown"),
		},
		{
			"VCS unknown",
			engine.SessionInfo{VCSName: "dummy", VCSSessionSummary: "", GitAutoPush: false},
			asYellowTrace("VCS \"dummy\" is unknown"),
		},
		{
			"git with auto-push on",
			engine.SessionInfo{VCSName: "git", VCSSessionSummary: "git branch \"my-branch\"", GitAutoPush: true},
			asCyanTrace("Running on git branch \"my-branch\" with auto-push enabled"),
		},
		{
			"git with auto-push off",
			engine.SessionInfo{VCSName: "git", VCSSessionSummary: "git branch \"my-branch\"", GitAutoPush: false},
			asCyanTrace("Running on git branch \"my-branch\" with auto-push disabled"),
		},
		{
			"p4",
			engine.SessionInfo{VCSName: "p4", VCSSessionSummary: "p4 client \"my-client\"", GitAutoPush: false},
			asCyanTrace("Running with p4 client \"my-client\""),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert.Equal(t, test.expected, capturer.CaptureStdout(func() {
				term, _, _ := terminalSetup(*params.AParamSet())
				term.reportVCSInfo(test.info)
				terminalTeardown(*term)
			}))
		})
	}
}

func Test_report_commit_message_suffix(t *testing.T) {
	tests := []struct {
		desc     string
		suffix   string
		expected string
	}{
		{
			"not set",
			"",
			"",
		},
		{
			"single-line suffix",
			"simple suffix",
			asCyanTrace("Commit message suffix: \"simple suffix\""),
		},
		{
			"multi-line suffix",
			"line 1\nline 2",
			asCyanTrace("Commit message suffix: \"line 1\nline 2\""),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert.Equal(t, test.expected, capturer.CaptureStdout(func() {
				term, _, _ := terminalSetup(*params.AParamSet())
				term.reportMessageSuffix(test.suffix)
				terminalTeardown(*term)
			}))
		})
	}
}

func Test_main_menu(t *testing.T) {
	slowTestTag(t)
	testFlags := []struct {
		desc     string
		vcsName  string
		input1   []byte
		input2   []byte
		expected []engine.TCRCall
	}{
		{
			"Enter key has no action", git.Name, []byte{enterKey}, nil,
			engine.NoTCRCall,
		},
		{
			"? key has no action on TCR", git.Name, []byte{'?'}, nil,
			engine.NoTCRCall,
		},
		{
			"Q key has no action on TCR", git.Name, []byte{'q'}, []byte{'Q'},
			engine.NoTCRCall,
		},
		{
			"T key has no action in main menu", git.Name, []byte{'t'}, []byte{'T'},
			engine.NoTCRCall,
		},
		{
			"P key is actionable with git", git.Name, []byte{'p'}, []byte{'P'},
			[]engine.TCRCall{
				engine.TCRCallToggleAutoPush,
				engine.TCRCallGetSessionInfo,
			},
		},
		{
			"P key has no action with p4", p4.Name, []byte{'p'}, []byte{'P'},
			engine.NoTCRCall,
		},
		{
			"L key is actionable with git", git.Name, []byte{'l'}, []byte{'L'},
			[]engine.TCRCall{
				engine.TCRCallVCSPull,
			},
		},
		{
			"L key has no action with p4", p4.Name, []byte{'l'}, []byte{'L'},
			engine.NoTCRCall,
		},
		{
			"S key is actionable with git", git.Name, []byte{'s'}, []byte{'S'},
			[]engine.TCRCall{
				engine.TCRCallVCSPush,
			},
		},
		{
			"S key has no action with p4", p4.Name, []byte{'s'}, []byte{'S'},
			engine.NoTCRCall,
		},
		{
			"Y key has no action with git", git.Name, []byte{'y'}, []byte{'Y'},
			engine.NoTCRCall,
		},
		{
			"Y key is actionable with p4", p4.Name, []byte{'y'}, []byte{'Y'},
			[]engine.TCRCall{
				engine.TCRCallVCSPull,
			},
		},
		{
			"D+Q keys", git.Name, []byte{'d', 'q'}, []byte{'D', 'Q'},
			[]engine.TCRCall{
				engine.TCRCallRunAsDriver,
				engine.TCRCallStop,
			},
		},
		{
			"N+Q keys", git.Name, []byte{'n', 'q'}, []byte{'N', 'Q'},
			[]engine.TCRCall{
				engine.TCRCallRunAsNavigator,
				engine.TCRCallStop,
			},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			for _, input := range [][]byte{tt.input1, tt.input2} {
				if input != nil {
					assertMainMenuActions(t, tt.vcsName, input, tt.expected)
				}
			}
		})
	}
}

func assertMainMenuActions(t *testing.T, vcsName string, input []byte, expected []engine.TCRCall) {
	t.Helper()
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr
	// Restore stdin, stdout and stderr right after the test.
	defer func() { os.Stdin = stdin; os.Stdout = stdout; os.Stderr = stderr }()
	// We fake stdin so that we can simulate a key press
	// We always add a 'q' at the end to make sure we get out of the infinite loop
	os.Stdin = fakeStdin(t, append(input, 'q'))
	// Displayed info on stdout and stderr is not used in the test
	os.Stdout = os.NewFile(0, os.DevNull)
	os.Stderr = os.NewFile(0, os.DevNull)

	term, fakeEngine, _ := terminalSetup(*params.AParamSet(params.WithVCS(vcsName)))
	term.enterMainMenu()
	assert.Equal(t, append(expected, engine.TCRCallQuit), fakeEngine.GetCallHistory())
	terminalTeardown(*term)
}

func Test_driver_menu(t *testing.T) {
	testFlags := []struct {
		desc     string
		input    []byte
		expected []engine.TCRCall
	}{
		{
			"Enter key has no action", []byte{enterKey},
			engine.NoTCRCall,
		},
		{
			"? key has no action on TCR", []byte{'?'},
			engine.NoTCRCall,
		},
		{
			"Q key has no action on TCR", []byte{'q', 'Q'},
			engine.NoTCRCall,
		},
		{
			"T key triggers reporting timer status", []byte{'t', 'T'},
			[]engine.TCRCall{engine.TCRCallReportMobTimerStatus},
		},
		{
			"P key has no action", []byte{'p', 'P'},
			engine.NoTCRCall,
		},
		{
			"D key has no action", []byte{'d', 'D'},
			engine.NoTCRCall,
		},
		{
			"N key has no action", []byte{'n', 'N'},
			engine.NoTCRCall,
		},
		{
			"L key has no action", []byte{'l', 'L'},
			engine.NoTCRCall,
		},
		{
			"S key has no action", []byte{'s', 'S'},
			engine.NoTCRCall,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			for _, input := range tt.input {
				assertStartAsActions(t, role.Driver{}, input,
					append([]engine.TCRCall{engine.TCRCallRunAsDriver}, tt.expected...))
			}
		})
	}
}

func Test_navigator_menu(t *testing.T) {
	testFlags := []struct {
		desc     string
		input    []byte
		expected []engine.TCRCall
	}{
		{
			"Enter key has no action", []byte{enterKey},
			engine.NoTCRCall,
		},
		{
			"? key has no action on TCR", []byte{'?'},
			engine.NoTCRCall,
		},
		{
			"Q key has no action on TCR", []byte{'q', 'Q'},
			engine.NoTCRCall,
		},
		{
			"T key has no action", []byte{'t', 'T'},
			engine.NoTCRCall,
		},
		{
			"P key has no action", []byte{'p', 'P'},
			engine.NoTCRCall,
		},
		{
			"L key has no action", []byte{'l', 'L'},
			engine.NoTCRCall,
		},
		{
			"S key has no action", []byte{'s', 'S'},
			engine.NoTCRCall,
		},
		{
			"D key has no action", []byte{'d', 'D'},
			engine.NoTCRCall,
		},
		{
			"N key has no action", []byte{'n', 'N'},
			engine.NoTCRCall,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			for _, input := range tt.input {
				assertStartAsActions(t, role.Navigator{}, input,
					append([]engine.TCRCall{engine.TCRCallRunAsNavigator}, tt.expected...))
			}
		})
	}
}

func assertStartAsActions(t *testing.T, r role.Role, input byte, expected []engine.TCRCall) {
	t.Helper()
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr
	// Restore stdin, stdout and stderr right after the test.
	defer func() { os.Stdin = stdin; os.Stdout = stdout; os.Stderr = stderr }()
	// We fake stdin so that we can simulate a key press
	// We always add a 'q' at the end to make sure we get out of the infinite loop
	os.Stdin = fakeStdin(t, []byte{input, 'q'})
	// Displayed info on stdout and stderr is not used in the test
	os.Stdout = os.NewFile(0, os.DevNull)
	os.Stderr = os.NewFile(0, os.DevNull)

	term, fakeEngine, _ := terminalSetup(*params.AParamSet())
	term.enterRole(r)
	assert.Equal(t, append(expected, engine.TCRCallStop), fakeEngine.GetCallHistory())
	terminalTeardown(*term)
}

func Test_start_terminal(t *testing.T) {
	testFlags := []struct {
		desc     string
		mode     runmode.RunMode
		input    []byte
		expected []engine.TCRCall
	}{
		{
			"solo mode", runmode.Solo{}, []byte{'q'},
			[]engine.TCRCall{
				engine.TCRCallRunAsDriver,
				engine.TCRCallStop,
				engine.TCRCallQuit,
			},
		},
		{
			"mob mode", runmode.Mob{}, []byte{'q'},
			[]engine.TCRCall{
				engine.TCRCallQuit,
			},
		},
		{
			"one-shot mode", runmode.OneShot{}, []byte{},
			[]engine.TCRCall{
				engine.TCRCallRunTcrCycle,
				engine.TCRCallQuit,
			},
		},
		{
			"check mode", runmode.Check{}, []byte{},
			[]engine.TCRCall{
				engine.TCRCallRunCheck,
				engine.TCRCallQuit,
			},
		},
		{
			"log mode", runmode.Log{}, []byte{},
			[]engine.TCRCall{
				engine.TCRCallPrintLog,
				engine.TCRCallQuit,
			},
		},
		{
			"stats mode", runmode.Stats{}, []byte{},
			[]engine.TCRCall{
				engine.TCRCallPrintStats,
				engine.TCRCallQuit,
			},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assertStartTerminal(t, tt.mode, tt.input, tt.expected)
		})
	}
}

func assertStartTerminal(t *testing.T, mode runmode.RunMode, input []byte, expected []engine.TCRCall) {
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr
	// Restore stdin, stdout and stderr right after the test.
	defer func() { os.Stdin = stdin; os.Stdout = stdout; os.Stderr = stderr }()
	// We fake stdin so that we can simulate a key press
	os.Stdin = fakeStdin(t, input)
	// Displayed info on stdout and stderr is not used in the test
	os.Stdout = os.NewFile(0, os.DevNull)
	os.Stderr = os.NewFile(0, os.DevNull)

	term, fakeEngine, _ := terminalSetup(*params.AParamSet(params.WithRunMode(mode)))
	term.Start()
	terminalTeardown(*term)
	assert.Equal(t, expected, fakeEngine.GetCallHistory())
}

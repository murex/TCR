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
	"fmt"
	"github.com/murex/tcr/desktop"
	"github.com/murex/tcr/engine"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/role"
	"github.com/murex/tcr/runmode"
	"github.com/murex/tcr/settings"
	"github.com/murex/tcr/ui"
	"os"
)

// TerminalUI is the user interface implementation when using the Command Line Interface
type TerminalUI struct {
	reportingChannel chan bool
	tcr              engine.TCRInterface
	params           params.Params
	desktop          *desktop.Desktop
}

const (
	enterKey  = 0x0a
	escapeKey = 0x1b
)

const (
	pullMenuHelper              = "Pull from remote"
	pushMenuHelper              = "Push to remote"
	driverRoleMenuHelper        = "Driver role"
	navigatorRoleMenuHelper     = "Navigator role"
	autoPushMenuHelper          = "Turn on/off VCS auto-push"
	quitMenuHelper              = "Quit"
	optionsMenuHelper           = "List available options"
	timerStatusMenuHelper       = "Timer status"
	quitDriverRoleMenuHelper    = "Quit Driver role"
	quitNavigatorRoleMenuHelper = "Quit Navigator role"
)

// New creates a new instance of terminal
func New(p params.Params, tcr engine.TCRInterface) ui.UserInterface {
	setLinePrefix("[" + settings.ApplicationName + "]")
	var term = TerminalUI{params: p, tcr: tcr, desktop: desktop.NewDesktop(nil)}
	term.MuteDesktopNotifications(false)
	term.StartReporting()
	StartInterruptHandler()
	return &term
}

// StartReporting tells the terminal to start reporting information
func (term *TerminalUI) StartReporting() {
	term.reportingChannel = report.Subscribe(term)
}

// MuteDesktopNotifications allows preventing desktop Notification popups from being displayed.
// Used for test automation at the moment. Could be turned into a feature later if there is need for it.
func (term *TerminalUI) MuteDesktopNotifications(muted bool) {
	if muted {
		term.desktop.MuteNotifications()
	} else {
		term.desktop.UnmuteNotifications()
	}
}

// StopReporting tells the terminal to stop reporting information
func (term *TerminalUI) StopReporting() {
	if term.reportingChannel != nil {
		report.Unsubscribe(term.reportingChannel)
	}
}

// NotifyRoleStarting tells the user that TCR engine is starting with the provided role
func (term *TerminalUI) NotifyRoleStarting(r role.Role) {
	term.ReportTitle(false, "Starting with ", r.LongName(), ". Press ? for options")
}

// NotifyRoleEnding tells the user that TCR engine is ending the provided role
func (term *TerminalUI) NotifyRoleEnding(r role.Role) {
	term.ReportInfo(false, "Ending ", r.LongName())
}

// ReportSimple reports simple messages
func (*TerminalUI) ReportSimple(_ bool, a ...interface{}) {
	printUntouched(a...)
}

// ReportInfo reports info messages
func (*TerminalUI) ReportInfo(_ bool, a ...interface{}) {
	printInCyan(a...)
}

// ReportTitle reports title messages
func (*TerminalUI) ReportTitle(_ bool, a ...interface{}) {
	printHorizontalLine()
	printInCyan(a...)
}

// ReportTimer reports timer messages
func (term *TerminalUI) ReportTimer(emphasis bool, a ...interface{}) {
	printInGreen(a...)
	term.notifyOnEmphasis(emphasis, "⏳", a...)
}

// ReportSuccess reports success messages
func (term *TerminalUI) ReportSuccess(emphasis bool, a ...interface{}) {
	printInGreen(a...)
	term.notifyOnEmphasis(emphasis, "🟢", a...)
}

// ReportWarning reports warning messages
func (term *TerminalUI) ReportWarning(emphasis bool, a ...interface{}) {
	printInYellow(a...)
	term.notifyOnEmphasis(emphasis, "🔶", a...)
}

// ReportError reports error messages
func (term *TerminalUI) ReportError(emphasis bool, a ...interface{}) {
	printInRed(a...)
	term.notifyOnEmphasis(emphasis, "🟥", a...)
}

func (term *TerminalUI) notifyOnEmphasis(emphasis bool, emoji string, a ...interface{}) {
	if emphasis {
		err := term.desktop.ShowNotification(desktop.NormalLevel, emoji+" "+settings.ApplicationName, fmt.Sprint(a...))
		if err != nil {
			term.ReportWarning(false, "Failed to show desktop notification: ", err.Error())
		}
	}
}

func (term *TerminalUI) mainMenu() {
	term.whatShallWeDo()

	keyboardInput := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(keyboardInput)
		if err != nil {
			term.ReportWarning(false, "Something went wrong while reading from stdin: ", err)
		}

		switch keyboardInput[0] {
		case '?':
			term.listMainMenuOptions("Available Options:")
		case 'd', 'D':
			term.startAs(role.Driver{})
			term.whatShallWeDo()
		case 'n', 'N':
			term.startAs(role.Navigator{})
			term.whatShallWeDo()
		case 'p', 'P':
			term.tcr.ToggleAutoPush()
			term.ShowSessionInfo()
			term.whatShallWeDo()
		case 'l', 'L':
			term.vcsPull()
			term.whatShallWeDo()
		case 's', 'S':
			term.vcsPush()
			term.whatShallWeDo()
		case 'q', 'Q':
			Restore()
			term.tcr.Quit()
			// The return statement below is never reached when running the application
			// due to the call to os.Exit() done in tcr.Quit(). We still want to keep
			// it here for running the tests, so that we're able to get out of the infinite loop
			// even when tcr.Quit() is faked.
			return
		case enterKey:
			// We ignore enter key press
			continue
		default:
			term.ReportWarning(false, "No action is mapped to shortcut '", string(keyboardInput), "'")
			term.listMainMenuOptions("Please choose one of the following:")
		}
	}
}

func (term *TerminalUI) vcsPull() {
	term.tcr.VCSPull()
}

func (term *TerminalUI) vcsPush() {
	term.tcr.VCSPush()
}

func (term *TerminalUI) whatShallWeDo() {
	term.listMainMenuOptions("What shall we do?")
}

func (term *TerminalUI) startAs(r role.Role) {
	// We ask TCR engine to start...
	switch r {
	case role.Navigator{}:
		term.tcr.RunAsNavigator()
	case role.Driver{}:
		term.tcr.RunAsDriver()
	default:
		term.ReportWarning(false, "No action defined for ", r.LongName())
	}

	// ...Until the user decides to stop
	keyboardInput := make([]byte, 1)
	for stopRequest := false; !stopRequest; {
		if _, err := os.Stdin.Read(keyboardInput); err != nil {
			term.ReportWarning(false, "Something went wrong while reading from stdin: ", err)
		}
		switch keyboardInput[0] {
		case '?':
			term.listRoleMenuOptions(r, "Available Options:")
		case 'q', 'Q', escapeKey:
			term.ReportWarning(false, "OK, I heard you")
			stopRequest = true
			term.tcr.Stop()
		case 't', 'T':
			term.showTimerStatus()
		case enterKey:
			// We ignore enter key press
			continue
		default:
			term.keyNotRecognizedMessage()
		}
	}
}

func (term *TerminalUI) keyNotRecognizedMessage() {
	term.ReportWarning(false, "Key not recognized. Press ? for available options")
}

func (term *TerminalUI) showTimerStatus() {
	if settings.EnableMobTimer {
		if r := term.tcr.GetCurrentRole(); r != nil && r.RunsWithTimer() {
			term.tcr.ReportMobTimerStatus()
		} else {
			term.keyNotRecognizedMessage()
		}
	}
}

// ShowRunningMode shows the current running mode
func (term *TerminalUI) ShowRunningMode(mode runmode.RunMode) {
	term.ReportTitle(false, "Running in ", mode.Name(), " mode")
}

// ShowSessionInfo shows main information related to the current TCR session
func (term *TerminalUI) ShowSessionInfo() {
	info := term.tcr.GetSessionInfo()

	term.ReportTitle(false, "Base Directory: ", info.BaseDir)
	term.ReportInfo(false, "Work Directory: ", info.WorkDir)
	term.ReportInfo(false, "Language=", info.LanguageName, ", Toolchain=", info.ToolchainName)

	autoPush := "disabled"
	if info.GitAutoPush {
		autoPush = "enabled"
	}
	term.ReportInfo(false, "Running on ", info.VCSSessionSummary, " with auto-push ", autoPush)
}

// Confirm asks the user for confirmation
func (term *TerminalUI) Confirm(message string, defaultAnswer bool) bool {
	_ = SetRaw()
	defer Restore()

	term.ReportWarning(false, message)
	term.ReportWarning(false, "Do you want to proceed? ", yesOrNoAdvice(defaultAnswer))

	keyboardInput := make([]byte, 1)
	for {
		_, _ = os.Stdin.Read(keyboardInput)
		switch keyboardInput[0] {
		case 'y', 'Y':
			return true
		case 'n', 'N':
			return false
		case enterKey:
			return defaultAnswer
		}
	}
}

func yesOrNoAdvice(defaultAnswer bool) string {
	if defaultAnswer {
		return "[Y/n]"
	}
	return "[y/N]"
}

// Start runs the terminal session
func (term *TerminalUI) Start() {
	term.initTCREngine()

	if term.params.Mode.IsInteractive() {
		_ = SetRaw()
		defer Restore()
	}

	switch term.params.Mode {
	case runmode.Solo{}:
		// When running TCR in solo mode, there's no selection menu:
		// we directly enter driver mode, and quit when done
		term.startAs(role.Driver{})
		Restore()
		term.tcr.Quit()
	case runmode.Mob{}:
		// When running TCR in mob mode, every participant
		// is given the possibility to switch between
		// driver and navigator modes
		term.mainMenu()
	case runmode.OneShot{}:
		// When running TCR in one-shot mode, there's no selection menu:
		// we directly ask TCR engine to run one cycle and quit when done
		term.tcr.RunTCRCycle()
		term.tcr.Quit()
	case runmode.Check{}:
		// When running TCR in check mode, there's no selection menu:
		// we directly ask TCR engine to run a check and quit when done
		term.tcr.RunCheck(term.params)
		term.tcr.Quit()
	case runmode.Log{}:
		// When running TCR in log mode, there's no selection menu:
		// we directly ask TCR engine to print the commit history and quit when done
		term.tcr.PrintLog(term.params)
		term.tcr.Quit()
	case runmode.Stats{}:
		// When running TCR in stats mode, there's no selection menu:
		// we directly ask TCR engine to print the stats and quit when done
		term.tcr.PrintStats(term.params)
		term.tcr.Quit()
	default:
		term.ReportError(false, "Unknown run mode: ", term.params.Mode)
	}
}

func (term *TerminalUI) initTCREngine() {
	term.tcr.Init(term, term.params)
}

func (term *TerminalUI) printMenuOption(shortcut byte, description ...interface{}) {
	term.ReportInfo(false, append([]interface{}{"\t", string(shortcut), " -> "}, description...)...)
}

func (term *TerminalUI) listMainMenuOptions(title string) {
	term.ReportTitle(false, title)
	term.printMenuOption('D', driverRoleMenuHelper)
	term.printMenuOption('N', navigatorRoleMenuHelper)
	term.printMenuOption('P', autoPushMenuHelper)
	term.printMenuOption('L', pullMenuHelper)
	term.printMenuOption('S', pushMenuHelper)
	term.printMenuOption('Q', quitMenuHelper)
	term.printMenuOption('?', optionsMenuHelper)
}

func (term *TerminalUI) listRoleMenuOptions(r role.Role, title string) {
	term.ReportTitle(false, title)
	if settings.EnableMobTimer && r != nil && r.RunsWithTimer() {
		term.printMenuOption('T', timerStatusMenuHelper)
	}
	term.printMenuOption('Q', "Quit ", r.LongName())
	term.printMenuOption('?', optionsMenuHelper)
}

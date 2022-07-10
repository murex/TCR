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
	"github.com/murex/tcr/tcr-cli/desktop"
	"github.com/murex/tcr/tcr-engine/engine"
	"github.com/murex/tcr/tcr-engine/params"
	"github.com/murex/tcr/tcr-engine/report"
	"github.com/murex/tcr/tcr-engine/role"
	"github.com/murex/tcr/tcr-engine/runmode"
	"github.com/murex/tcr/tcr-engine/settings"
	"github.com/murex/tcr/tcr-engine/ui"
	"os"
)

// TerminalUI is the user interface implementation when using the Command Line Interface
type TerminalUI struct {
	reportingChannel chan bool
	tcr              engine.TcrInterface
	params           params.Params
}

const (
	enterKey  = 0x0a
	escapeKey = 0x1b
)

// New creates a new instance of terminal
func New(p params.Params, tcr engine.TcrInterface) ui.UserInterface {
	setLinePrefix("[" + settings.ApplicationName + "]")
	var term = TerminalUI{params: p, tcr: tcr}
	term.MuteDesktopNotifications(false)
	term.StartReporting()
	return &term
}

// StartReporting tells the terminal to start reporting information
func (term *TerminalUI) StartReporting() {
	term.reportingChannel = report.Subscribe(func(msg report.Message) {
		switch msg.Type {
		case report.Normal:
			term.trace(msg.Text)
		case report.Title:
			term.title(msg.Text)
		case report.Info:
			term.info(msg.Text)
		case report.Warning:
			term.warning(msg.Text)
		case report.Error:
			term.error(msg.Text)
		case report.Notification:
			term.notification(msg.Text)
		}
	})
}

// MuteDesktopNotifications allows preventing desktop notification popups from being displayed.
// Used for test automation at the moment. Could be turned into a feature later if there is need for it.
func (term *TerminalUI) MuteDesktopNotifications(muted bool) {
	if muted {
		desktop.MuteNotifications()
	} else {
		desktop.UnmuteNotifications()
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
	term.title("Starting with ", r.LongName(), ". Press ? for options")
}

// NotifyRoleEnding tells the user that TCR engine is ending the provided role
func (term *TerminalUI) NotifyRoleEnding(r role.Role) {
	term.info("Ending ", r.LongName())
}

func (term *TerminalUI) info(a ...interface{}) {
	printInCyan(a...)
}

func (term *TerminalUI) title(a ...interface{}) {
	printHorizontalLine()
	printInCyan(a...)
}

func (term *TerminalUI) warning(a ...interface{}) {
	printInYellow(a...)
}

func (term *TerminalUI) error(a ...interface{}) {
	printInRed(a...)
}

func (term *TerminalUI) notification(a ...interface{}) {
	printInGreen(a...)
	err := desktop.ShowNotification(desktop.NormalLevel, settings.ApplicationName, fmt.Sprint(a...))
	if err != nil {
		term.warning("Failed to show desktop notification: ", err.Error())
	}
}

func (term *TerminalUI) trace(a ...interface{}) {
	printUntouched(a...)
}

func (term *TerminalUI) mainMenu() {
	term.whatShallWeDo()

	keyboardInput := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(keyboardInput)
		if err != nil {
			term.warning("Something went wrong while reading from stdin: ", err)
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
			term.warning("No action is mapped to shortcut '", string(keyboardInput), "'")
			term.listMainMenuOptions("Please choose one of the following:")
		}
	}
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
		term.warning("No action defined for ", r.LongName())
	}

	// ...Until the user decides to stop
	keyboardInput := make([]byte, 1)
	for stopRequest := false; !stopRequest; {
		if _, err := os.Stdin.Read(keyboardInput); err != nil {
			term.warning("Something went wrong while reading from stdin: ", err)
		}
		switch keyboardInput[0] {
		case '?':
			term.listRoleMenuOptions(r, "Available Options:")
		case 'q', 'Q', escapeKey:
			term.warning("OK, I heard you")
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
	term.warning("Key not recognized. Press ? for available options")
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
	term.title("Running in ", mode.Name(), " mode")
}

// ShowSessionInfo shows main information related to the current TCR session
func (term *TerminalUI) ShowSessionInfo() {
	info := term.tcr.GetSessionInfo()

	term.title("Base Directory: ", info.BaseDir)
	term.info("Work Directory: ", info.WorkDir)
	term.info("Language=", info.LanguageName, ", Toolchain=", info.ToolchainName)

	autoPush := "disabled"
	if info.AutoPush {
		autoPush = "enabled"
	}
	term.info(
		"Running on git branch \"", info.BranchName,
		"\" with auto-push ", autoPush)
}

// Confirm asks the user for confirmation
func (term *TerminalUI) Confirm(message string, defaultAnswer bool) bool {

	_ = SetRaw()
	defer Restore()

	term.warning(message)
	term.warning("Do you want to proceed? ", yesOrNoAdvice(defaultAnswer))

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
	term.initTcrEngine()

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
		term.error("Unknown run mode: ", term.params.Mode)
	}
}

func (term *TerminalUI) initTcrEngine() {
	term.tcr.Init(term, term.params)
}

func (term *TerminalUI) printMenuOption(shortcut byte, description ...interface{}) {
	term.info(append([]interface{}{"\t", string(shortcut), " -> "}, description...)...)
}

func (term *TerminalUI) listMainMenuOptions(title string) {
	term.title(title)
	term.printMenuOption('D', role.Driver{}.LongName())
	term.printMenuOption('N', role.Navigator{}.LongName())
	term.printMenuOption('P', "Turn on/off git auto-push")
	term.printMenuOption('Q', "Quit")
	term.printMenuOption('?', "List available options")
}

func (term *TerminalUI) listRoleMenuOptions(r role.Role, title string) {
	term.title(title)
	if settings.EnableMobTimer && r != nil && r.RunsWithTimer() {
		term.printMenuOption('T', "Timer status")
	}
	term.printMenuOption('Q', "Quit ", r.LongName())
	term.printMenuOption('?', "List available options")
}

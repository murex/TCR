/*
Copyright (c) 2021 Murex

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
	"github.com/murex/tcr/tcr-engine/report"
	"github.com/murex/tcr/tcr-engine/role"
	"github.com/murex/tcr/tcr-engine/runmode"
	"github.com/murex/tcr/tcr-engine/settings"
	"github.com/murex/tcr/tcr-engine/stty"
	"github.com/murex/tcr/tcr-engine/ui"
	"os"
	"strings"
)

// Terminal is the user interface implementation when using the Command Line Interface
type Terminal struct {
	reportingChannel chan bool
	params           engine.Params
}

const (
	enterKey  = 0x0a
	escapeKey = 0x1b
)

// New creates a new instance of terminal
func New(p engine.Params) ui.UserInterface {
	setLinePrefix("[" + settings.ApplicationName + "]")
	var term = Terminal{params: p}
	term.StartReporting()
	return &term
}

// StartReporting tells the terminal to start reporting information
func (term *Terminal) StartReporting() {
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

// StopReporting tells the terminal to stop reporting information
func (term *Terminal) StopReporting() {
	report.Unsubscribe(term.reportingChannel)
}

// NotifyRoleStarting tells the user that TCR engine is starting with the provided role
func (term *Terminal) NotifyRoleStarting(r role.Role) {
	term.title("Starting with ", strings.Title(r.Name()), " role. Press ? for options")
}

// NotifyRoleEnding tells the user that TCR engine is ending the provided role
func (term *Terminal) NotifyRoleEnding(r role.Role) {
	term.info("Ending ", strings.Title(r.Name()), " role")
}

func (term *Terminal) info(a ...interface{}) {
	printInCyan(a...)
}

func (term *Terminal) title(a ...interface{}) {
	printHorizontalLine()
	printInCyan(a...)
}

func (term *Terminal) warning(a ...interface{}) {
	printInYellow(a...)
}

func (term *Terminal) error(a ...interface{}) {
	printInRed(a...)
}

func (term *Terminal) notification(a ...interface{}) {
	printInGreen(a...)
	desktop.ShowNotification(desktop.NormalLevel, settings.ApplicationName, fmt.Sprint(a...))
}

func (term *Terminal) trace(a ...interface{}) {
	printUntouched(a...)
}

func (term *Terminal) mainMenu() {
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
			engine.ToggleAutoPush()
			term.ShowSessionInfo()
			term.whatShallWeDo()
		case 'q', 'Q':
			stty.Restore()
			engine.Quit()
		case enterKey:
			// We ignore enter key press
			continue
		default:
			term.warning("No action is mapped to shortcut '", string(keyboardInput), "'")
			term.listMainMenuOptions("Please choose one of the following:")
		}
	}
}

func (term *Terminal) whatShallWeDo() {
	term.listMainMenuOptions("What shall we do?")
}

func (term *Terminal) startAs(r role.Role) {

	// We ask TCR engine to start...
	switch r {
	case role.Navigator{}:
		engine.RunAsNavigator()
	case role.Driver{}:
		engine.RunAsDriver()
	default:
		term.warning("No action defined for role ", r.Name())
	}

	// ...Until the user decides to stop
	keyboardInput := make([]byte, 1)
	for stopRequest := false; !stopRequest; {
		if _, err := os.Stdin.Read(keyboardInput); err != nil {
			term.warning("Something went wrong while reading from stdin: ", err)
		}
		switch keyboardInput[0] {
		case '?':
			term.listRoleMenuOptions("Available Options:")
		case 'q', 'Q', escapeKey:
			term.warning("OK, I heard you")
			stopRequest = true
			engine.Stop()
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

func (term *Terminal) keyNotRecognizedMessage() {
	term.warning("Key not recognized. Press ? for available options")
}

func (term *Terminal) showTimerStatus() {
	if settings.EnableMobTimer {
		if r := engine.GetCurrentRole(); r != nil && r.RunsWithTimer() {
			engine.ReportMobTimerStatus()
		} else {
			term.keyNotRecognizedMessage()
		}
	}

}

// ShowRunningMode shows the current running mode
func (term *Terminal) ShowRunningMode(mode runmode.RunMode) {
	term.title("Running in ", mode.Name(), " mode")
}

// ShowSessionInfo shows main information related to the current TCR session
func (term *Terminal) ShowSessionInfo() {
	d, l, t, ap, b := engine.GetSessionInfo()

	term.title("Working Directory: ", d)
	term.info("Language=", l, ", Toolchain=", t)

	autoPush := "disabled"
	if ap {
		autoPush = "enabled"
	}
	term.info(
		"Running on git branch \"", b,
		"\" with auto-push ", autoPush)
}

// Confirm asks the user for confirmation
func (term *Terminal) Confirm(message string, defaultAnswer bool) bool {

	_ = stty.SetRaw()
	defer stty.Restore()

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
func (term *Terminal) Start() {

	term.initTcrEngine()

	_ = stty.SetRaw()
	defer stty.Restore()

	switch term.params.Mode {
	case runmode.Solo{}:
		// When running TCR in solo mode, there's no
		// selection menu: we directly enter driver mode, and quit when done
		term.startAs(role.Driver{})
		engine.Quit()
	case runmode.Mob{}:
		// When running TCR in mob mode, every participant
		// is given the possibility to switch between
		// driver and navigator modes
		term.mainMenu()
	default:
		term.error("Unknown run mode: ", term.params.Mode)
	}
}

func (term *Terminal) initTcrEngine() {
	engine.Init(term, term.params)
}

func (term *Terminal) printMenuOption(shortcut byte, description ...interface{}) {
	term.info(append([]interface{}{"\t", string(shortcut), " -> "}, description...)...)
}

func (term *Terminal) listMainMenuOptions(title string) {
	term.title(title)
	term.printMenuOption('D', strings.Title(role.Driver{}.Name()), " role")
	term.printMenuOption('N', strings.Title(role.Navigator{}.Name()), " role")
	term.printMenuOption('P', "Turn on/off git auto-push")
	term.printMenuOption('Q', "Quit")
	term.printMenuOption('?', "List available options")
}

func (term *Terminal) listRoleMenuOptions(title string) {
	term.title(title)
	r := engine.GetCurrentRole()
	if settings.EnableMobTimer && r != nil && r.RunsWithTimer() {
		term.printMenuOption('T', "Timer status")
	}
	term.printMenuOption('Q', "Quit ", r.Name(), " role")
	term.printMenuOption('?', "List available options")
}

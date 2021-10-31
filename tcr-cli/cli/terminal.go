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
	"github.com/murex/tcr-cli/desktop"
	"github.com/murex/tcr-engine/engine"
	"github.com/murex/tcr-engine/report"
	"github.com/murex/tcr-engine/role"
	"github.com/murex/tcr-engine/runmode"
	"github.com/murex/tcr-engine/stty"
	"github.com/murex/tcr-engine/ui"
	"os"
	"strings"
)

// Terminal is the user interface implementation when using the Command Line Interface
type Terminal struct {
	reportingChannel chan bool
	params           engine.Params
}

const (
	tcrLinePrefix = "[TCR]"

	enterKey  = 0x0a
	escapeKey = 0x1b
)

// New creates a new instance of terminal
func New(p engine.Params) ui.UserInterface {
	setLinePrefix(tcrLinePrefix)
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
		case report.Event:
			term.event(msg.Text)
			desktop.ShowNotification(desktop.NormalLevel, "TCR", msg.Text)
		}
	})
}

// StopReporting tells the terminal to stop reporting information
func (term *Terminal) StopReporting() {
	report.Unsubscribe(term.reportingChannel)
}

// NotifyRoleStarting tells the user that TCR engine is starting with the provided role
func (term *Terminal) NotifyRoleStarting(r role.Role) {
	term.title("Starting as a ", strings.Title(r.Name()), ". Press ESC or Q when done")
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

func (term *Terminal) event(a ...interface{}) {
	printInGreen(a...)
}

func (term *Terminal) trace(a ...interface{}) {
	printUntouched(a...)
}

func (term *Terminal) mainMenu() {
	term.printOptionsMenu()

	keyboardInput := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(keyboardInput)
		if err != nil {
			term.warning("Something went wrong while reading from stdin: ", err)
		}

		switch keyboardInput[0] {
		case 'd', 'D':
			term.startAs(role.Driver{})
		case 'n', 'N':
			term.startAs(role.Navigator{})
		case 'p', 'P':
			engine.ToggleAutoPush()
			term.ShowSessionInfo()
		case 'q', 'Q':
			stty.Restore()
			engine.Quit()
		case enterKey:
			// We ignore enter key press
			continue
		default:
			term.warning("No action is mapped to shortcut '",
				string(keyboardInput), "'")
		}
		term.printOptionsMenu()
	}
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
		_, err := os.Stdin.Read(keyboardInput)
		if err != nil {
			term.warning("Something went wrong while reading from stdin: ", err)
		}
		switch keyboardInput[0] {
		case 'q', 'Q', escapeKey:
			term.warning("OK, I heard you")
			stopRequest = true
			engine.Stop()
		case enterKey:
			// We ignore enter key press
			continue
		default:
			term.warning("Key not recognized. Press ESC or Q to leave ", r.Name(), " role")
		}
	}
}

// ShowRunningMode shows the current running mode
func (term *Terminal) ShowRunningMode(mode runmode.RunMode) {
	term.title("Running in ", mode.Name(), " mode")
}

func (term *Terminal) printOptionsMenu() {
	term.title("What shall we do?")
	term.info("\tD -> ", strings.Title(role.Driver{}.Name()), " role")
	term.info("\tN -> ", strings.Title(role.Navigator{}.Name()), " role")
	term.info("\tP -> Turn on/off git auto-push")
	term.info("\tQ -> Quit")
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
		// selection menu: we directly enter driver mode
		term.startAs(role.Driver{})
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

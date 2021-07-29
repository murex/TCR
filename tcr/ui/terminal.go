package ui

import (
	"github.com/mengdaming/tcr/stty"
	"github.com/mengdaming/tcr/tcr"
	"github.com/mengdaming/tcr/tcr/engine"
	"github.com/mengdaming/tcr/tcr/role"
	"github.com/mengdaming/tcr/trace"
	"strings"

	"os"
)

type Terminal struct {
}

const (
	tcrLinePrefix = "[TCR]"
)

func NewTerminal() tcr.UserInterface {
	trace.SetLinePrefix(tcrLinePrefix)
	var term = Terminal{}
	return &term
}

func (term *Terminal) NotifyRoleStarting(role tcr.Role) {
	term.horizontalLine()
	term.Info("Starting as a ", strings.Title(role.Name()),
		". Press ESC to return to the main menu")
}

func (term *Terminal) NotifyRoleEnding(role tcr.Role) {
	term.Info("Leaving ", strings.Title(role.Name()), " role")
}

func (term *Terminal) Info(a ...interface{}) {
	trace.Info(a...)
}

func (term *Terminal) Warning(a ...interface{}) {
	trace.Warning(a...)
}

func (term *Terminal) Error(a ...interface{}) {
	trace.Error(a...)
}

func (term *Terminal) Trace(a ...interface{}) {
	trace.Echo(a...)
}

func (term *Terminal) horizontalLine() {
	trace.HorizontalLine()
}

func (term *Terminal) WaitForAction() {
	term.printOptionsMenu()

	_ = stty.SetRaw()
	defer stty.Restore()

	keyboardInput := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(keyboardInput)
		if err != nil {
			term.Warning("Something went wrong while reading from stdin: ", err)
		}

		switch keyboardInput[0] {
		case 'd', 'D':
			term.enterRoleMenu(role.Driver{})
		case 'n', 'N':
			term.enterRoleMenu(role.Navigator{})
		case 'p', 'P':
			engine.ToggleAutoPush()
			term.ShowSessionInfo()
		case 'q', 'Q':
			stty.Restore()
			engine.Quit()
		default:
			term.Warning("No action is mapped to shortcut '",
				string(keyboardInput), "'")
		}
		term.printOptionsMenu()
	}
}

func (term *Terminal) enterRoleMenu(r role.Role) {

	keyboardInput := make([]byte, 1)
	interrupt := make(chan bool)

	switch r {
	case role.Navigator{}:
		go engine.RunAsNavigator(interrupt)
	case role.Driver{}:
		go engine.RunAsDriver(interrupt)
	default:
		term.Warning("No action defined for role ", r.Name())
	}

	for {
		var done = false
		_, err := os.Stdin.Read(keyboardInput)
		if err != nil {
			term.Warning("Something went wrong while reading from stdin: ", err)
		}
		switch keyboardInput[0] {
		// 27 is for ESC key, at least on Windows
		case 27:
			term.Warning("OK, heading back to the main menu")
			interrupt <- true
			done = true
		default:
			term.Warning("Key not recognized. Press ESC to return to the main menu")
		}
		if done {
			break
		}
	}
}

func (term *Terminal) ShowRunningMode(mode tcr.WorkMode) {
	term.horizontalLine()
	term.Info("Running in ", mode, " mode")
}

func (term *Terminal) printOptionsMenu() {
	term.horizontalLine()
	term.Info("What shall we do?")
	term.Info("\tD -> Driver role")
	term.Info("\tN -> Navigator role")
	term.Info("\tP -> Turn on/off git auto-push")
	term.Info("\tQ -> Quit")
}

func (term *Terminal) ShowSessionInfo() {
	d, l, t, ap, b := engine.GetSessionInfo()

	term.horizontalLine()
	term.Info("Working Directory: ", d)
	term.Info("Language=", l, ", Toolchain=", t)

	autoPush := "disabled"
	if ap {
		autoPush = "enabled"
	}
	term.Info(
		"Running on git branch \"", b,
		"\" with auto-push ", autoPush)
}

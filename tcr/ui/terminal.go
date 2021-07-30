package ui

import (
	"github.com/mengdaming/tcr/stty"
	"github.com/mengdaming/tcr/tcr"
	"github.com/mengdaming/tcr/tcr/engine"
	"github.com/mengdaming/tcr/tcr/role"
	"github.com/mengdaming/tcr/trace"
	"os"
	"strings"
)

type Terminal struct {
}

const (
	tcrLinePrefix = "[TCR]"

	enterKey  = 0x0a
	escapeKey = 0x1b
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
			term.startAs(role.Driver{})
		case 'n', 'N':
			term.startAs(role.Navigator{})
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

func (term *Terminal) startAs(r role.Role) {

	// We ask TCR engine to start...
	stopEngine := make(chan bool)
	switch r {
	case role.Navigator{}:
		go engine.RunAsNavigator(stopEngine)
	case role.Driver{}:
		go engine.RunAsDriver(stopEngine)
	default:
		term.Warning("No action defined for role ", r.Name())
	}

	// ...Until the user decides to stop and go back to the main menu
	keyboardInput := make([]byte, 1)
	for stopRequest := false; !stopRequest; {
		_, err := os.Stdin.Read(keyboardInput)
		if err != nil {
			term.Warning("Something went wrong while reading from stdin: ", err)
		}
		switch keyboardInput[0] {
		case escapeKey:
			term.Warning("OK, heading back to the main menu")
			stopRequest = true
			stopEngine <- true
		default:
			term.Warning("Key not recognized. Press ESC to return to the main menu")
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
	term.Info("\tD -> ", strings.Title(role.Driver{}.Name()), " role")
	term.Info("\tN -> ", strings.Title(role.Navigator{}.Name()), " role")
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

func (term *Terminal) Confirm(message string, defaultAnswer bool) bool {

	_ = stty.SetRaw()
	defer stty.Restore()

	term.Warning(message)
	term.Warning("Do you want to proceed? ", yesOrNoAdvice(defaultAnswer))

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
	if defaultAnswer == true {
		return "[Y/n]"
	} else {
		return "[y/N]"
	}
}

package ui

import (
	"github.com/mengdaming/tcr/stty"
	"github.com/mengdaming/tcr/tcr"
	"github.com/mengdaming/tcr/tcr/engine"
	"github.com/mengdaming/tcr/trace"
	"strings"

	"os"
)

type Terminal struct {
}

func NewTerminal() tcr.UserInterface {
	var term = Terminal{}
	return &term
}

func (term *Terminal) NotifyRoleStarting(role tcr.Role) {
	term.horizontalLine()
	term.Info("Starting as a ", strings.Title(string(role)),
		". Press CTRL-C to go back to the main menu")
}

func (term *Terminal) NotifyRoleEnding(role tcr.Role) {
	term.Info("Leaving ", strings.Title(string(role)), " role")
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
			engine.RunAsDriver()
		case 'n', 'N':
			engine.RunAsNavigator()
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

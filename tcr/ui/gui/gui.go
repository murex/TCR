package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mengdaming/tcr/tcr/engine"
	"github.com/mengdaming/tcr/tcr/role"
	"github.com/mengdaming/tcr/tcr/runmode"
	"github.com/mengdaming/tcr/tcr/ui"
	"github.com/mengdaming/tcr/tcr/ui/cli"
	"github.com/mengdaming/tcr/trace"
	"image/color"
)

type GUI struct {
	app       fyne.App
	win       fyne.Window
	directory widget.Label
	language  widget.Label
	toolchain widget.Label
	branch    widget.Label
	autoPush  widget.Label
}

const (
	linePrefix = "[TCR-GUI]"
)

// TODO We re-route temporarily messages to the terminal
var term ui.UserInterface

func New() ui.UserInterface {
	term = cli.New()
	trace.SetLinePrefix(linePrefix)
	var gui = GUI{}
	gui.initApp()
	return &gui
}

func (gui *GUI) RunInMode(mode runmode.RunMode) {
	// TODO setup according to mode value
	gui.win.ShowAndRun()
}

func (gui *GUI) ShowRunningMode(mode runmode.RunMode) {
	// TODO Replace with GU-specific implementation
	term.ShowRunningMode(mode)
}

func (gui *GUI) NotifyRoleStarting(r role.Role) {
	// TODO Replace with GU-specific implementation
	term.NotifyRoleStarting(r)
}

func (gui *GUI) NotifyRoleEnding(r role.Role) {
	// TODO Replace with GU-specific implementation
	term.NotifyRoleEnding(r)
}

func (gui *GUI) ShowSessionInfo() {
	d, l, t, ap, b := engine.GetSessionInfo()

	gui.directory.Text = fmt.Sprintf("Directory: %v", d)
	gui.language.Text = fmt.Sprintf("Language: %v", l)
	gui.toolchain.Text = fmt.Sprintf("Toolchain: %v", t)
	gui.branch.Text = fmt.Sprintf("Branch: %v", b)
	if ap {
		gui.autoPush.Text = "Auto-Push: enabled"
	} else {
		gui.autoPush.Text = "Auto-Push: disabled"
	}
}

func (gui *GUI) Info(a ...interface{}) {
	// TODO Replace with GU-specific implementation
	term.Info(a...)
}

func (gui *GUI) Warning(a ...interface{}) {
	// TODO Replace with GU-specific implementation
	term.Warning(a...)
}

func (gui *GUI) Error(a ...interface{}) {
	// TODO Replace with GU-specific implementation
	term.Error(a...)
}

func (gui *GUI) Trace(a ...interface{}) {
	// TODO Replace with GU-specific implementation
	term.Trace(a...)
}

func (gui *GUI) Confirm(message string, def bool) bool {
	// TODO Replace with GU-specific implementation
	return term.Confirm(message, def)
}

func (gui *GUI) initApp() {
	gui.app = app.New()
	gui.win = gui.app.NewWindow("TCR")
	gui.win.Resize(fyne.NewSize(400, 600))

	top := canvas.NewText("top bar", color.White)
	left := canvas.NewText("left", color.White)
	right := canvas.NewText("right", color.White)
	middle := canvas.NewText("middle", color.White)

	// Session Information container

	gui.directory = *widget.NewLabel("Directory")
	gui.language = *widget.NewLabel("Language")
	gui.toolchain = *widget.NewLabel("Toolchain")
	gui.branch = *widget.NewLabel("Branch")
	gui.autoPush = *widget.NewLabel("Auto-Push")
	sessionInfo := container.NewVBox(
		container.NewHBox(
			&gui.directory,
		),
		container.NewHBox(
			&gui.language,
			&gui.toolchain,
			&gui.branch,
			&gui.autoPush,
		),
	)

	// Top level container

	content := container.New(layout.NewBorderLayout(top,
		sessionInfo, left, right),
		top, left, middle, right, sessionInfo)

	gui.win.SetContent(content)
}

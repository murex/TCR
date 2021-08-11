package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/mengdaming/tcr/tcr/engine"
	"github.com/mengdaming/tcr/tcr/report"
	"github.com/mengdaming/tcr/tcr/role"
	"github.com/mengdaming/tcr/tcr/runmode"
	"github.com/mengdaming/tcr/tcr/ui"
	"github.com/mengdaming/tcr/tcr/ui/cli"
	"image/color"
	"strings"
)

const (
	defaultWidth  = 400
	defaultHeight = 800
)

var (
	redColor    = color.RGBA{R: 255, G: 0, B: 0}
	cyanColor   = color.RGBA{R: 0, G: 139, B: 139}
	yellowColor = color.RGBA{R: 255, G: 255, B: 0}
	whiteColor  = color.RGBA{R: 255, G: 255, B: 255}
)

type GUI struct {
	term         ui.UserInterface
	reporting    chan bool
	app          fyne.App
	win          fyne.Window
	actionBar    ActionBar
	traceArea    *TraceArea
	sessionPanel *SessionPanel
}

func New() ui.UserInterface {
	var gui = GUI{}
	// Until the GUI is able to report, we rely on the terminal to report information
	gui.term = cli.New()
	gui.initApp()
	report.PostInfo("Opening TCR GUI")

	gui.term.StopReporting()
	gui.StartReporting()
	return &gui
}

func (gui *GUI) StartReporting() {
	gui.reporting = report.Subscribe(func(msg report.Message) {
		switch msg.Type {
		case report.Normal:
			gui.trace(msg.Text)
		case report.Title:
			gui.title(msg.Text)
		case report.Info:
			gui.info(msg.Text)
		case report.Warning:
			gui.warning(msg.Text)
		case report.Error:
			gui.error(msg.Text)
		}
	})
}

func (gui *GUI) StopReporting() {
	report.Unsubscribe(gui.reporting)
}

func (gui *GUI) Start(mode runmode.RunMode) {
	gui.confirmRootBranch()
	gui.ShowRunningMode(mode)
	gui.win.ShowAndRun()
}

func (gui *GUI) ShowRunningMode(mode runmode.RunMode) {
	gui.sessionPanel.setMode(mode)
	gui.adjustActionBar(mode)
}

func (gui *GUI) NotifyRoleStarting(r role.Role) {
	report.PostTitle("Starting as a ", strings.Title(r.Name()))
}

func (gui *GUI) NotifyRoleEnding(r role.Role) {
	report.PostInfo("Ending ", strings.Title(r.Name()), " role")
}

func (gui *GUI) ShowSessionInfo() {
	gui.sessionPanel.setSessionInfo()
}

func (gui *GUI) info(a ...interface{}) {
	gui.traceArea.printText(cyanColor, a...)
}

func (gui *GUI) title(a ...interface{}) {
	gui.traceArea.printHeader(a...)
}

func (gui *GUI) warning(a ...interface{}) {
	gui.traceArea.printText(yellowColor, a...)
}

func (gui *GUI) error(a ...interface{}) {
	gui.traceArea.printText(redColor, a...)
}

func (gui *GUI) trace(a ...interface{}) {
	gui.traceArea.printText(whiteColor, a...)
}

type confirmationInfo struct {
	required      bool
	title         string
	message       string
	defaultAnswer bool
}

var rootBranchConfirmation confirmationInfo

func (gui *GUI) Confirm(message string, def bool) bool {
	// We need to defer the confirmation dialog until the window is displayed
	gui.warning(message)
	rootBranchConfirmation = confirmationInfo{
		required:      true,
		title:         message,
		message:       "Are you sure you want to continue?",
		defaultAnswer: def,
	}
	return true
}

func (gui *GUI) confirmRootBranch() {
	if rootBranchConfirmation.required {
		// TODO See if there is a way to change the default button selection in fyne confirmation dialog
		dialog.ShowConfirm(
			rootBranchConfirmation.title,
			rootBranchConfirmation.message,
			func(response bool) {
				if response == false {
					gui.quit()
				}
			},
			gui.win,
		)
	}
}

func (gui *GUI) quit() {
	gui.StopReporting()
	gui.term.StartReporting()
	engine.Quit()
}

func (gui *GUI) initApp() {
	gui.app = app.New()
	// TODO Add a TCR-Specific icon
	gui.app.SetIcon(theme.FyneLogo())
	gui.win = gui.app.NewWindow("TCR")
	gui.win.Resize(fyne.NewSize(defaultWidth, defaultHeight))
	gui.win.SetCloseIntercept(func() {
		gui.quit()
		gui.win.Close()
	})
	gui.win.CenterOnScreen()

	gui.actionBar = NewMobActionBar()
	gui.traceArea = NewTraceArea()
	gui.sessionPanel = NewSessionPanel()

	topLevel := container.New(layout.NewBorderLayout(
		gui.sessionPanel.container, gui.actionBar.getContainer(), nil, nil),
		gui.sessionPanel.container, gui.actionBar.getContainer(), gui.traceArea.container)

	gui.win.SetContent(topLevel)
}

func (gui *GUI) adjustActionBar(mode runmode.RunMode) {
	// TODO replace containers
	switch mode {
	case runmode.Mob{}:
		gui.actionBar = NewMobActionBar()
	case runmode.Solo{}:
		gui.actionBar = NewSoloActionBar()
	}
}

package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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
	rbConfirm    DeferredConfirmDialog
	topLevel     *fyne.Container
	layout       fyne.Layout
	runMode      runmode.RunMode
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

func (gui *GUI) Start(_ runmode.RunMode) {
	gui.rbConfirm.showIfNeeded()
	gui.win.ShowAndRun()
}

func (gui *GUI) ShowRunningMode(mode runmode.RunMode) {
	gui.setRunMode(mode)
}

func (gui *GUI) NotifyRoleStarting(r role.Role) {
	report.PostTitle("Starting as a ", strings.Title(r.Name()))
}

func (gui *GUI) NotifyRoleEnding(r role.Role) {
	report.PostInfo("Ending ", strings.Title(r.Name()), " role")
}

func (gui *GUI) ShowSessionInfo() {
	d, l, t, ap, b := engine.GetSessionInfo()
	gui.win.SetTitle(fmt.Sprintf("TCR - %v", d))
	gui.sessionPanel.setLanguage(l)
	gui.sessionPanel.setToolchain(t)
	gui.sessionPanel.setBranch(b)
	gui.sessionPanel.setGitAutoPush(ap)
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

func (gui *GUI) quit() {
	gui.StopReporting()
	gui.term.StartReporting()
	engine.Quit()
}

func (gui *GUI) initApp() {
	gui.app = app.New()
	icon, _ := fyne.LoadResourceFromPath("Icon.png")
	gui.app.SetIcon(icon)
	gui.win = gui.app.NewWindow("TCR")
	gui.win.Resize(fyne.NewSize(defaultWidth, defaultHeight))
	gui.win.SetCloseIntercept(func() {
		gui.quit()
		gui.win.Close()
	})
	gui.win.CenterOnScreen()

	gui.actionBar = NewActionBar()
	gui.traceArea = NewTraceArea()
	gui.sessionPanel = gui.NewSessionPanel()

	gui.layout = layout.NewBorderLayout(
		gui.sessionPanel.container, gui.actionBar.container, nil, nil)
	gui.topLevel = container.New(gui.layout,
		gui.sessionPanel.container,
		gui.actionBar.container,
		gui.traceArea.container,
	)
	gui.win.SetContent(gui.topLevel)
}

func (gui *GUI) setRunMode(mode runmode.RunMode) {
	if mode != gui.runMode {
		gui.runMode = mode
		report.PostInfo(fmt.Sprintf("Run mode set to %v", gui.runMode.Name()))
		gui.sessionPanel.setRunMode(gui.runMode)
		gui.actionBar.setRunMode(gui.runMode)
	}
}

func (gui *GUI) getRunMode() runmode.RunMode {
	return gui.runMode
}

func (gui *GUI) Confirm(message string, def bool) bool {
	gui.warning(message)
	// We need to defer showing the confirmation dialog until the window is displayed
	gui.rbConfirm = NewDeferredConfirmDialog(message, def, gui.quit, gui.win)
	return true
}

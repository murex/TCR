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

package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/murex/tcr/tcr-cli/cli"
	"github.com/murex/tcr/tcr-engine/engine"
	"github.com/murex/tcr/tcr-engine/report"
	"github.com/murex/tcr/tcr-engine/role"
	"github.com/murex/tcr/tcr-engine/runmode"
	"github.com/murex/tcr/tcr-engine/settings"
	"github.com/murex/tcr/tcr-engine/ui"
	"image/color"
	"strings"
)

const (
	defaultWidth  = 400
	defaultHeight = 800
)

var (
	redColor   = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	cyanColor  = color.RGBA{R: 0, G: 139, B: 139, A: 255}
	greenColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	//yellowColor = color.RGBA{R: 255, G: 255, B: 0, A: 255}
	orangeColor = color.RGBA{R: 255, G: 165, B: 0, A: 255}
	//whiteColor  = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	grayColor = color.RGBA{R: 128, G: 128, B: 128, A: 255}
	//blackColor  = color.RGBA{R: 0, G: 0, B: 0, A: 255}
)

// GUI is the user interface implementation when using the TCR with its graphical user interface
type GUI struct {
	term          ui.UserInterface
	reporting     chan bool
	app           fyne.App
	win           fyne.Window
	actionBar     ActionBar
	traceArea     *TraceArea
	sessionPanel  *SessionPanel
	rbConfirm     DeferredConfirmDialog
	baseDirDialog BaseDirSelectionDialog
	topLevel      *fyne.Container
	layout        fyne.Layout
	runMode       runmode.RunMode
	params        engine.Params
}

// New creates a new instance of graphical user interface
func New(p engine.Params) ui.UserInterface {
	var gui = GUI{params: p}
	// Until the GUI is able to report, we rely on the terminal to report information
	gui.term = cli.New(p)
	report.PostInfo("Opening ", settings.ApplicationName, " GUI")

	gui.initApp()
	gui.initBaseDirSelectionDialog()

	gui.StartReporting()
	return &gui
}

// StartReporting tells the GUI to start reporting information
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
		case report.Notification:
			gui.notification(msg.Text)
		}
	})
}

// StopReporting tells the GUI to stop reporting information
func (gui *GUI) StopReporting() {
	report.Unsubscribe(gui.reporting)
}

// Start runs the GUI session
func (gui *GUI) Start() {
	if isBaseDirDefined(gui.params.BaseDir) {
		gui.initTcrEngine(gui.params.BaseDir)
	} else {
		gui.baseDirDialog.show()
	}
	gui.rbConfirm.showIfNeeded()
	gui.win.ShowAndRun()
}

func isBaseDirDefined(dir string) bool {
	return dir != ""
}

// ShowRunningMode shows the current running mode
func (gui *GUI) ShowRunningMode(mode runmode.RunMode) {
	gui.setRunMode(mode)
}

// NotifyRoleStarting tells the user that TCR engine is starting with the provided role
func (gui *GUI) NotifyRoleStarting(r role.Role) {
	report.PostTitle("Starting as a ", strings.Title(r.Name()))
	gui.sessionPanel.disableActions()
}

// NotifyRoleEnding tells the user that TCR engine is ending the provided role
func (gui *GUI) NotifyRoleEnding(r role.Role) {
	report.PostInfo("Ending ", strings.Title(r.Name()), " role")
	gui.sessionPanel.enableActions()
}

// ShowSessionInfo shows main information related to the current TCR session
func (gui *GUI) ShowSessionInfo() {
	d, l, t, ap, b := engine.GetSessionInfo()
	gui.win.SetTitle(fmt.Sprintf("%v - %v", settings.ApplicationName, d))
	gui.sessionPanel.setLanguage(l)
	gui.sessionPanel.setToolchain(t)
	gui.sessionPanel.setBranch(b)
	gui.sessionPanel.setGitAutoPush(ap)
}

func (gui *GUI) info(a ...interface{}) {
	gui.traceArea.printText(cyanColor, false, a...)
}

func (gui *GUI) title(a ...interface{}) {
	gui.traceArea.printHeader(a...)
}

func (gui *GUI) warning(a ...interface{}) {
	gui.traceArea.printText(orangeColor, false, a...)
}

func (gui *GUI) error(a ...interface{}) {
	gui.traceArea.printText(redColor, false, a...)
}

func (gui *GUI) notification(a ...interface{}) {
	gui.traceArea.printText(greenColor, false, a...)
	gui.app.SendNotification(fyne.NewNotification(settings.ApplicationName, fmt.Sprint(a...)))
}

func (gui *GUI) trace(a ...interface{}) {
	gui.traceArea.printText(grayColor, true, a...)
}

func (gui *GUI) quit(message string) {
	gui.StopReporting()
	gui.term.StartReporting()
	if message != "" {
		report.PostInfo(message)
	}
	engine.Quit()
}

func (gui *GUI) initApp() {
	gui.app = app.NewWithID("TCR")
	icon, _ := fyne.LoadResourceFromPath("Icon.png")
	gui.app.SetIcon(icon)
	gui.win = gui.app.NewWindow(settings.ApplicationName)
	gui.win.Resize(fyne.NewSize(defaultWidth, defaultHeight))
	gui.win.SetCloseIntercept(func() {
		gui.quit("Closing application")
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
	gui.BuildMainMenu()

}

func (gui *GUI) setRunMode(mode runmode.RunMode) {
	if mode != gui.runMode {
		gui.runMode = mode
		report.PostInfo(fmt.Sprintf("Run mode set to %v", gui.runMode.Name()))
		gui.sessionPanel.setRunMode(gui.runMode)
		gui.actionBar.setRunMode(gui.runMode)
		engine.SetRunMode(mode)
	}
}

// Confirm asks the user for confirmation through a popup confirmation window
func (gui *GUI) Confirm(message string, def bool) bool {
	gui.warning(message)
	// We need to defer showing the confirmation dialog until the window is displayed
	gui.rbConfirm = NewDeferredConfirmDialog(message, def, gui.quit, gui.win)
	gui.rbConfirm.showIfNeeded()
	return true
}

func (gui *GUI) initBaseDirSelectionDialog() {
	gui.baseDirDialog = NewBaseDirSelectionDialog(gui.initTcrEngine, gui.win)
}

func (gui *GUI) initTcrEngine(baseDir string) {
	if baseDir == "" {
		gui.quit("Operation cancelled")
	}
	gui.params.BaseDir = baseDir
	engine.Init(gui, gui.params)
	gui.term.StopReporting()
}

func (gui *GUI) showTimerStatus() {
	if settings.EnableMobTimer {
		engine.ReportMobTimerStatus()
	}
}

// TODO Add possibility to tune mob timer duration from the GUI

// TODO Add possibility to disable/enable mob timer from the GUI

// TODO Replace timer trace with a visual component such as a progress bar or even a clock

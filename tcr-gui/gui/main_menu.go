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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/murex/tcr/tcr-engine/settings"
)

// BuildMainMenu creates the application's main menu
func (gui *GUI) BuildMainMenu() {
	tcrMenu := gui.buildTcrMenu()
	timerMenu := gui.buildTimerMenu()
	helpMenu := gui.buildHelpMenu()
	gui.win.SetMainMenu(fyne.NewMainMenu(tcrMenu, timerMenu, helpMenu))
}

func (gui *GUI) buildTcrMenu() *fyne.Menu {
	quitMenuItem := fyne.NewMenuItem("Quit",
		func() {
			gui.quit("Closing application")
		},
	)
	return fyne.NewMenu("TCR", quitMenuItem)
}

func (gui *GUI) buildTimerMenu() *fyne.Menu {
	statusMenuItem := fyne.NewMenuItem("Status",
		func() {
			gui.showTimerStatus()
		},
	)
	return fyne.NewMenu("Timer", statusMenuItem)
}

func (gui *GUI) buildHelpMenu() *fyne.Menu {
	aboutMenuItem := fyne.NewMenuItem("About",
		func() {
			gui.NewAboutDialog().Show()
		},
	)
	return fyne.NewMenu("Help", aboutMenuItem)
}

// NewAboutDialog creates a dialog window with "About TCR" contents
func (gui *GUI) NewAboutDialog() dialog.Dialog {
	c := container.New(layout.NewFormLayout())
	for _, buildInfo := range settings.GetBuildInfo() {
		c.Add(widget.NewLabelWithStyle(buildInfo.Label, fyne.TextAlignTrailing, fyne.TextStyle{Italic: true}))
		c.Add(widget.NewLabelWithStyle(buildInfo.Value, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	}
	return dialog.NewCustom("About "+settings.ApplicationName, "Dismiss", c, gui.win)
}

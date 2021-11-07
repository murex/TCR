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
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/murex/tcr-engine/settings"
)

// BuildMainMenu creates the application's main menu
func (gui *GUI) BuildMainMenu() {

	tcrMenu := fyne.NewMenu("TCR")

	// --------------------------------

	buildInfoMenuItem := fyne.NewMenuItem("Build Info",
		func() {
			b := new(bytes.Buffer)
			for key, value := range settings.GetBuildInfo() {
				_, _ = fmt.Fprintf(b, "%s:\t\"%s\"\n", key, value)
			}
			// TODO Replace dialog window with a custom window with nicer formatting of information
			w := dialog.NewInformation("Build Information", b.String(), gui.win)
			//w := dialog.NewCustom("Build Information", "Dismiss",
			//	container.New(layout.NewFormLayout(), widget.NewLabel("xxx"), widget.NewLabel("yyy")), gui.win)
			w.Show()
		},
	)

	helpMenu := fyne.NewMenu("Help", buildInfoMenuItem)

	// --------------------------------

	gui.win.SetMainMenu(fyne.NewMainMenu(tcrMenu, helpMenu))
}

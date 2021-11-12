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
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/murex/tcr/tcr-engine/runmode"
)

// ActionBar is the area containing the buttons that the user can use to start/stop TCR engine.
// The content of the action bar can change depending on the current run mode
type ActionBar struct {
	currentMode  runmode.RunMode
	buttonPanels map[runmode.RunMode]ActionBarButtonPanel
	container    *fyne.Container
}

// ActionBarButtonPanel provides the interface that any action bar button panel must implement
// so that it can be used inside the action bar
type ActionBarButtonPanel interface {
	getContainer() *fyne.Container
	updateButtonsState(running bool)
}

// NewActionBar creates an action bar panel containing the action buttons
func NewActionBar() ActionBar {
	var ab = ActionBar{}

	// The 2 mode-specific button panels
	ab.buttonPanels = make(map[runmode.RunMode]ActionBarButtonPanel)
	ab.buttonPanels[runmode.Solo{}] = NewSoloActionBarButtonPanel()
	ab.buttonPanels[runmode.Mob{}] = NewMobActionBarButtonPanel()

	ab.container = container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel(""), // Empty label is used as a placeholder until run mode is set
		layout.NewSpacer(),
	)
	return ab
}

func (ab *ActionBar) setRunMode(mode runmode.RunMode) {
	if mode != ab.currentMode {
		ab.currentMode = mode
		bp := ab.buttonPanels[mode]
		ab.container.Objects[1] = bp.getContainer()
		ab.container.Refresh()
	}
}

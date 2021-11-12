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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/murex/tcr/tcr-engine/engine"
)

// MobActionBarButtonPanel is the action bar panel containing the buttons when running in mob
type MobActionBarButtonPanel struct {
	startNavigatorButton *widget.Button
	startDriverButton    *widget.Button
	stopButton           *widget.Button
	container            *fyne.Container
}

func (ab *MobActionBarButtonPanel) getContainer() *fyne.Container {
	return ab.container
}

// NewMobActionBarButtonPanel creates the action bar panel containing the buttons when running in mob
func NewMobActionBarButtonPanel() ActionBarButtonPanel {
	var ab = MobActionBarButtonPanel{}

	ab.startDriverButton = widget.NewButtonWithIcon("Start as Driver",
		theme.MediaPlayIcon(),
		func() {
			ab.updateButtonsState(true)
			engine.RunAsDriver()
		},
	)
	ab.startNavigatorButton = widget.NewButtonWithIcon("Start as Navigator",
		theme.MediaPlayIcon(),
		func() {
			ab.updateButtonsState(true)
			engine.RunAsNavigator()
		},
	)
	ab.stopButton = widget.NewButtonWithIcon("Stop",
		theme.MediaStopIcon(),
		func() {
			ab.updateButtonsState(false)
			engine.Stop()
		},
	)

	// Initial state
	ab.updateButtonsState(false)

	ab.container = container.NewHBox(
		layout.NewSpacer(),
		ab.startDriverButton,
		ab.startNavigatorButton,
		layout.NewSpacer(),
		ab.stopButton,
		layout.NewSpacer(),
	)

	return &ab
}

func (ab *MobActionBarButtonPanel) updateButtonsState(running bool) {
	if running {
		ab.startDriverButton.Disable()
		ab.startNavigatorButton.Disable()
		ab.stopButton.Enable()
	} else {
		ab.startDriverButton.Enable()
		ab.startNavigatorButton.Enable()
		ab.stopButton.Disable()
	}
}

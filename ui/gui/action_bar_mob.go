package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/mengdaming/tcr/engine"
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

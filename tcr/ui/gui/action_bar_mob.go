package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/mengdaming/tcr/tcr/engine"
)

type MobActionBar struct {
	startNavigatorButton *widget.Button
	startDriverButton    *widget.Button
	stopButton           *widget.Button
	container            *fyne.Container
}

func (ab *MobActionBar) getContainer() *fyne.Container {
	return ab.container
}

func NewMobActionBar() ActionBar {
	var ab = MobActionBar{}

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

	ab.container = container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(
			layout.NewSpacer(),
			ab.startDriverButton,
			ab.startNavigatorButton,
			layout.NewSpacer(),
			ab.stopButton,
			layout.NewSpacer(),
		),
	)

	return &ab
}

func (ab *MobActionBar) updateButtonsState(running bool) {
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

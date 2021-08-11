package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/mengdaming/tcr/tcr/engine"
)

type SoloActionBar struct {
	startButton *widget.Button
	stopButton  *widget.Button
	container   *fyne.Container
}

func (ab *SoloActionBar) getContainer() *fyne.Container {
	return ab.container
}

func NewSoloActionBar() ActionBar {
	var ab = SoloActionBar{}

	ab.startButton = widget.NewButtonWithIcon("Start",
		theme.MediaPlayIcon(),
		func() {
			ab.updateButtonsState(true)
			engine.RunAsDriver()
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
			ab.startButton,
			layout.NewSpacer(),
			ab.stopButton,
			layout.NewSpacer(),
		),
	)

	return &ab
}

func (ab *SoloActionBar) updateButtonsState(running bool) {
	if running {
		ab.startButton.Disable()
		ab.stopButton.Enable()
	} else {
		ab.startButton.Enable()
		ab.stopButton.Disable()
	}
}

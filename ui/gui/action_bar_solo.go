package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/mengdaming/tcr/engine"
)

type SoloActionBarButtonPanel struct {
	startButton *widget.Button
	stopButton  *widget.Button
	container   *fyne.Container
}

func (ab *SoloActionBarButtonPanel) getContainer() *fyne.Container {
	return ab.container
}

func NewSoloActionBarButtonPanel() ActionBarButtonPanel {
	var ab = SoloActionBarButtonPanel{}

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

	ab.container = container.NewHBox(
		layout.NewSpacer(),
		ab.startButton,
		layout.NewSpacer(),
		ab.stopButton,
		layout.NewSpacer(),
	)

	return &ab
}

func (ab *SoloActionBarButtonPanel) updateButtonsState(running bool) {
	if running {
		ab.startButton.Disable()
		ab.stopButton.Enable()
	} else {
		ab.startButton.Enable()
		ab.stopButton.Disable()
	}
}

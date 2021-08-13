package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mengdaming/tcr/tcr/runmode"
)

type ActionBar struct {
	currentMode  runmode.RunMode
	buttonPanels map[runmode.RunMode]ActionBarButtonPanel
	container    *fyne.Container
}

type ActionBarButtonPanel interface {
	getContainer() *fyne.Container
	updateButtonsState(running bool)
}

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

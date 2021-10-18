package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/murex/tcr-engine/runmode"
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

package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/murex/tcr-engine/engine"
	"github.com/murex/tcr-engine/runmode"
)

// SessionPanel is the panel showing all information related to the current TCR session
type SessionPanel struct {
	languageLabel  *widget.Label
	toolchainLabel *widget.Label
	branchLabel    *widget.Label
	modeLabel      *widget.Label
	modeSelect     *widget.Select
	autoPushToggle *widget.Check
	container      *fyne.Container
}

// NewSessionPanel creates a new instance of session information panel
func (gui *GUI) NewSessionPanel() *SessionPanel {
	sp := SessionPanel{}

	sp.modeLabel = widget.NewLabel("Mode")
	sp.modeSelect = widget.NewSelect(
		runmode.Names(),
		func(selected string) {
			var newMode = runmode.Map()[selected]
			gui.setRunMode(newMode)
		},
	)

	sp.languageLabel = widget.NewLabel("Language")
	sp.toolchainLabel = widget.NewLabel("Toolchain")
	sp.branchLabel = widget.NewLabel("Branch")

	sp.autoPushToggle = widget.NewCheck("Auto-Push",
		func(checked bool) {
			engine.SetAutoPush(checked)
		},
	)

	sp.container = container.NewVBox(
		container.NewHBox(
			sp.modeLabel,
			sp.modeSelect,
			widget.NewSeparator(),
			sp.languageLabel,
			widget.NewSeparator(),
			sp.toolchainLabel,
			widget.NewSeparator(),
			sp.branchLabel,
			widget.NewSeparator(),
			sp.autoPushToggle,
		),
		widget.NewSeparator(),
	)
	return &sp
}

func (sp *SessionPanel) setRunMode(mode runmode.RunMode) {
	sp.modeSelect.SetSelected(mode.Name())
	sp.autoPushToggle.SetChecked(mode.AutoPushDefault())
}

func (sp *SessionPanel) setLanguage(language string) {
	sp.languageLabel.SetText(fmt.Sprintf("Language: %v", language))
}

func (sp *SessionPanel) setToolchain(toolchain string) {
	sp.toolchainLabel.SetText(fmt.Sprintf("Toolchain: %v", toolchain))
}

func (sp *SessionPanel) setBranch(branch string) {
	sp.branchLabel.SetText(fmt.Sprintf("Branch: %v", branch))
}

func (sp *SessionPanel) setGitAutoPush(autoPush bool) {
	sp.autoPushToggle.SetChecked(autoPush)
}

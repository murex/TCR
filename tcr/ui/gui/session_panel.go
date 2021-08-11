package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mengdaming/tcr/tcr/engine"
	"github.com/mengdaming/tcr/tcr/report"
	"github.com/mengdaming/tcr/tcr/runmode"
)

type SessionPanel struct {
	directoryLabel *widget.Label
	languageLabel  *widget.Label
	toolchainLabel *widget.Label
	branchLabel    *widget.Label
	modeLabel      *widget.Label
	modeSelect     *widget.Select
	autoPushToggle *widget.Check
	container      *fyne.Container
}

func NewSessionPanel() *SessionPanel {
	sp := SessionPanel{}
	sp.directoryLabel = widget.NewLabel("Directory")
	sp.modeLabel = widget.NewLabel("Mode")

	// TODO retrieve the list from runmode package
	modeOptions := []string{runmode.Mob{}.Name(), runmode.Solo{}.Name()}
	sp.modeSelect = widget.NewSelect(modeOptions, func(s string) {
		report.PostWarning(fmt.Sprintf("Switching to %v mode", s))
	})

	sp.languageLabel = widget.NewLabel("Language")
	sp.toolchainLabel = widget.NewLabel("Toolchain")
	sp.branchLabel = widget.NewLabel("Branch")
	sp.autoPushToggle = widget.NewCheck("Auto-Push", func(checked bool) {
		engine.ToggleAutoPush()
		autoPushStr := "off"
		if checked {
			autoPushStr = "on"
		}
		report.PostInfo(fmt.Sprintf("Git auto-push is turned %v", autoPushStr))
	})

	sp.container = container.NewVBox(
		container.NewHBox(
			sp.directoryLabel,
		),
		widget.NewSeparator(),
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

func (sp *SessionPanel) setMode(mode runmode.RunMode) {
	//sp.modeLabel.SetText(fmt.Sprintf("Mode: %v", mode.Name()))
	sp.modeSelect.SetSelected(mode.Name())
	//sp.modeSelect.Refresh()
}

func (sp *SessionPanel) setSessionInfo() {
	d, l, t, ap, b := engine.GetSessionInfo()

	sp.directoryLabel.SetText(fmt.Sprintf("Directory: %v", d))
	sp.languageLabel.SetText(fmt.Sprintf("Language: %v", l))
	sp.toolchainLabel.SetText(fmt.Sprintf("Toolchain: %v", t))
	sp.branchLabel.SetText(fmt.Sprintf("Branch: %v", b))
	sp.autoPushToggle.SetChecked(ap)
}

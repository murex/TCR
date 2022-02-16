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
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/murex/tcr/tcr-engine/engine"
	"github.com/murex/tcr/tcr-engine/runmode"
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
		runmode.InteractiveModes(),
		func(selected string) {
			var newMode = runmode.Map()[selected]
			gui.setRunMode(newMode)
		},
	)

	sp.languageLabel = widget.NewLabel("Language")
	sp.toolchainLabel = widget.NewLabel("Toolchain")
	sp.branchLabel = widget.NewLabel("BranchName")

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
	sp.branchLabel.SetText(fmt.Sprintf("BranchName: %v", branch))
}

func (sp *SessionPanel) setGitAutoPush(autoPush bool) {
	sp.autoPushToggle.SetChecked(autoPush)
}

func (sp *SessionPanel) disableActions() {
	sp.modeSelect.Disable()
	sp.autoPushToggle.Disable()
}

func (sp *SessionPanel) enableActions() {
	sp.modeSelect.Enable()
	sp.autoPushToggle.Enable()
}

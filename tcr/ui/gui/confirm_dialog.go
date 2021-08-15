package gui

import (
	"fyne.io/fyne/v2/dialog"
)

type confirmationInfo struct {
	required      bool
	title         string
	message       string
	defaultAnswer bool
}

var rootBranchConfirmation confirmationInfo

func (gui *GUI) Confirm(message string, def bool) bool {
	// We need to defer the confirmation dialog until the window is displayed
	gui.warning(message)
	rootBranchConfirmation = confirmationInfo{
		required:      true,
		title:         message,
		message:       "Are you sure you want to continue?",
		defaultAnswer: def,
	}
	return true
}

func (gui *GUI) confirmRootBranch() {
	if rootBranchConfirmation.required {
		c := dialog.NewConfirm(
			rootBranchConfirmation.title,
			rootBranchConfirmation.message,
			func(response bool) {
				if response != rootBranchConfirmation.defaultAnswer {
					gui.quit()
				}
			},
			gui.win,
		)
		// This is a workaround to allow default selection of "No" answer.
		// Fyne ConfirmDialog don't allow this by default
		if rootBranchConfirmation.defaultAnswer == false {
			c.SetConfirmText("No")
			c.SetDismissText("Yes")
		}
		c.Show()
	}
}

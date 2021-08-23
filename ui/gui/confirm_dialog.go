package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type DeferredConfirmDialog struct {
	confirmDialog *dialog.ConfirmDialog
}

func NewDeferredConfirmDialog(message string, defaultSelected bool, cbAction func(), parent fyne.Window) DeferredConfirmDialog {
	cd := DeferredConfirmDialog{}
	cd.confirmDialog = dialog.NewConfirm(
		message, "Are you sure you want to continue?",
		func(response bool) {
			if response != defaultSelected {
				cbAction()
			}
		},
		parent,
	)
	if !defaultSelected {
		cd.confirmDialog.SetConfirmText("No")
		cd.confirmDialog.SetDismissText("Yes")
	}
	return cd
}

func (dlg *DeferredConfirmDialog) showIfNeeded() {
	if dlg.confirmDialog != nil {
		dlg.confirmDialog.Show()
	}
}

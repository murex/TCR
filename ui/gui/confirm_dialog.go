package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// DeferredConfirmDialog is a customized ConfirmDialog for which display is deferred until the application
// main window is displayed on the screen. It also adds the possibility to indicate the button to be selected
// by default in the dialog window
type DeferredConfirmDialog struct {
	confirmDialog *dialog.ConfirmDialog
}

// NewDeferredConfirmDialog creates a new instance of deferred confirmation dialog
func NewDeferredConfirmDialog(message string, defaultSelected bool, cbAction func(info string), parent fyne.Window) DeferredConfirmDialog {
	cd := DeferredConfirmDialog{}
	cd.confirmDialog = dialog.NewConfirm(
		message, "Are you sure you want to continue?",
		func(response bool) {
			if response != defaultSelected {
				cbAction("Ok, let's stop here then")
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

package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// BaseDirSelectionDialog allows selecting base working directory. The dialog shows only when
// there is no valid directory defined at startup
type BaseDirSelectionDialog struct {
	selectionDialog *dialog.FileDialog
}

// NewBaseDirSelectionDialog crates an instance of base dir selection dialog
func NewBaseDirSelectionDialog(cbAction func(baseDir string), parent fyne.Window) BaseDirSelectionDialog {
	dlg := BaseDirSelectionDialog{}

	dlg.selectionDialog = dialog.NewFolderOpen(func(uri fyne.ListableURI, e error) {
		var selected string
		if uri != nil {
			selected = uri.Path()
		}
		cbAction(selected)
	}, parent)

	return dlg
}

func (dlg BaseDirSelectionDialog) show() {
	dlg.selectionDialog.Show()
}

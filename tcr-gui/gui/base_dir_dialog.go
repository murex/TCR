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

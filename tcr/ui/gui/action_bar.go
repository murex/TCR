package gui

import "fyne.io/fyne/v2"

type ActionBar interface {
	getContainer() *fyne.Container
	updateButtonsState(running bool)
}

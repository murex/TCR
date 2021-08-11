package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strings"
)

type TraceArea struct {
	vbox      *fyne.Container
	container *container.Scroll
}

func NewTraceArea() *TraceArea {
	ta := TraceArea{}
	ta.vbox = container.NewVBox(
		widget.NewLabelWithStyle("Welcome to TCR!",
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		))

	ta.container = container.NewVScroll(ta.vbox)

	return &ta
}

func (ta *TraceArea) printLine() {
	ta.vbox.Add(widget.NewSeparator())
	ta.scrollToBottom()
}

func (ta *TraceArea) printText(col color.Color, a ...interface{}) {
	for _, s := range strings.Split(fmt.Sprint(a...), "\n") {
		ta.vbox.Add(canvas.NewText(s, col))
	}
	ta.scrollToBottom()
}

func (ta *TraceArea) printHeader(a ...interface{}) {
	ta.printLine()
	ta.vbox.Add(widget.NewLabelWithStyle(
		fmt.Sprint(a...),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	))
	ta.printLine()
	ta.scrollToBottom()
}

func (ta *TraceArea) scrollToBottom() {
	// The ScrollToTop() call below is some kind of workaround to ensure
	// that the UI indeed refreshes and scrolls to bottom when ScrollToBottom() is called
	ta.container.ScrollToTop()
	ta.container.ScrollToBottom()
}


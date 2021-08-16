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

const (
	traceMaxSize = 200
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

func (ta *TraceArea) printSeparator() {
	ta.print(widget.NewSeparator())
}

func (ta *TraceArea) printText(col color.Color, a ...interface{}) {
	for _, s := range strings.Split(fmt.Sprint(a...), "\n") {
		ta.print(canvas.NewText(s, col))
	}
}

func (ta *TraceArea) printHeader(a ...interface{}) {
	ta.printSeparator()
	ta.print(widget.NewLabelWithStyle(
		fmt.Sprint(a...),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	))
	ta.printSeparator()
}

func (ta *TraceArea) print(object fyne.CanvasObject) {
	// We cap the number of lines in trace to prevent performance decrease as trace piles up
	if len(ta.vbox.Objects) >= traceMaxSize {
		ta.vbox.Objects = ta.vbox.Objects[1:]
	}
	ta.vbox.Add(object)
	ta.scrollToBottom()
}

func (ta *TraceArea) scrollToBottom() {
	ta.container.ScrollToBottom()
}

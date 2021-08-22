package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"math"
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
	ta.container = container.NewScroll(ta.vbox)
	return &ta
}

func (ta *TraceArea) printSeparator() {
	ta.print(widget.NewSeparator())
}

func splitStringFixedLength(str string, length int) []string {
	strLength := len(str)
	nbChunks := int(math.Ceil(float64(strLength) / float64(length)))
	chunks := make([]string, nbChunks)
	var start, stop int
	for i := 0; i < nbChunks; i += 1 {
		start = i * length
		stop = start + length
		if stop > strLength {
			stop = strLength
		}
		chunks[i] = str[start : stop]
	}
	return chunks
}

func (ta *TraceArea) printText(col color.Color, a ...interface{}) {
	str := fmt.Sprint(a...)
	const maxCharsPerLine = 80
	for _, line := range strings.Split(strings.TrimRight(str, "\n"), "\n") {
		for _, chunk := range splitStringFixedLength(line, maxCharsPerLine) {
			ta.print(canvas.NewText(chunk, col))
		}
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

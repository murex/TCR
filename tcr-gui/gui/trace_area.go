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
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"math"
	"strings"
)

const (
	maxNumberOfLines = 200
	maxCharsPerLine  = 80
)

// TraceArea is the area of the GUI in charge of printing traces
type TraceArea struct {
	vbox      *fyne.Container
	container *container.Scroll
}

// NewTraceArea initializes the trace area widget
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
	for i := 0; i < nbChunks; i++ {
		start = i * length
		stop = start + length
		if stop > strLength {
			stop = strLength
		}
		chunks[i] = str[start:stop]
	}
	return chunks
}

func (ta *TraceArea) printText(col color.Color, monospace bool, a ...interface{}) {
	str := fmt.Sprint(a...)
	for _, line := range strings.Split(strings.TrimRight(str, "\n"), "\n") {
		for _, chunk := range splitStringFixedLength(line, maxCharsPerLine) {
			txt := canvas.NewText(chunk, col)
			txt.TextStyle = fyne.TextStyle{Monospace: monospace}
			ta.print(txt)
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
	if len(ta.vbox.Objects) >= maxNumberOfLines {
		ta.vbox.Objects = ta.vbox.Objects[1:]
	}
	ta.vbox.Add(object)
	ta.scrollToBottom()
}

func (ta *TraceArea) scrollToBottom() {
	ta.container.ScrollToBottom()
}

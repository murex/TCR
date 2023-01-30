/*
Copyright (c) 2022 Murex

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

package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
	"strings"
	"testing"
)

const (
	ansiEscape   = "\x1b[0m"
	ansiRedFg    = "\x1b[31m"
	ansiGreenFg  = "\x1b[32m"
	ansiYellowFg = "\x1b[33m"
	ansiCyanFg   = "\x1b[36m"
	newline      = "\n"
)

func Test_change_line_prefix(t *testing.T) {
	prefix := "PREFIX"
	msg := "Message"
	setLinePrefix(prefix)
	out := capturer.CaptureStdout(func() {
		printInCyan(msg)
	})
	assert.Contains(t, out, prefix)
}

func Test_print_untouched_function_does_not_alter_data(t *testing.T) {
	msg := "Dummy Message"
	out := capturer.CaptureStdout(func() {
		printUntouched(msg)
	})
	assert.Equal(t, msg+"\n", out)
}

func assertPrintInColor(t *testing.T, printInColorFunc func(a ...interface{}), ansiCode string) {
	msg := "Some message in color"
	assertPrintFormatting(t, func() { printInColorFunc(msg) }, ansiCode, "TCR", msg)
}

func assertPrintFormatting(t *testing.T, printFunc func(), ansiCode string, prefix string, msg string) {
	setLinePrefix(prefix)
	out := capturer.CaptureStdout(printFunc)
	expected := ansiCode + prefix + ansiEscape + " " + ansiCode + msg + ansiEscape + newline
	assert.Equal(t, expected, out)
}

func Test_print_in_cyan_function_formatting(t *testing.T) {
	assertPrintInColor(t, printInCyan, ansiCyanFg)
}

func Test_print_in_green_function_formatting(t *testing.T) {
	assertPrintInColor(t, printInGreen, ansiGreenFg)
}

func Test_print_in_yellow_function_formatting(t *testing.T) {
	assertPrintInColor(t, printInYellow, ansiYellowFg)
}

func Test_print_in_red_function_formatting(t *testing.T) {
	assertPrintInColor(t, printInRed, ansiRedFg)
}

func Test_print_horizontal_line_when_terminal_width_is_unknown(t *testing.T) {
	assertPrintFormatting(t, printHorizontalLine, ansiCyanFg, "TCR",
		strings.Repeat(horizontalLineCharacter, 75))
}

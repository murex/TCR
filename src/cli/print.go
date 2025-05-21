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
	"fmt"
	"github.com/logrusorgru/aurora"
	"strings"
)

const (
	horizontalLineCharacter = "─" // Character used for printing horizontal lines
	defaultTerminalWidth    = 80  // Default terminal width if current terminal is not recognized
)

var (
	colorizer  = aurora.NewAurora(true)
	linePrefix = ""
)

func setLinePrefix(value string) {
	linePrefix = value
}

func printPrefixedAndColored(fgColor aurora.Color, message string) {
	setupTerminal()
	_, _ = fmt.Println(
		colorizer.Colorize(linePrefix, fgColor),
		colorizer.Colorize(message, fgColor))
}

func printInCyan(a ...any) {
	printPrefixedAndColored(aurora.CyanFg, fmt.Sprint(a...))
}

func printInYellow(a ...any) {
	printPrefixedAndColored(aurora.YellowFg, fmt.Sprint(a...))
}

func printInGreen(a ...any) {
	printPrefixedAndColored(aurora.GreenFg, fmt.Sprint(a...))
}

func printInRed(a ...any) {
	printPrefixedAndColored(aurora.RedFg, fmt.Sprint(a...))
}

func printUntouched(a ...any) {
	_, _ = fmt.Println(a...)
}

func printHorizontalLine() {
	lineWidth := max(getTerminalColumns()-len(linePrefix)-2, 0)
	printInCyan(strings.Repeat(horizontalLineCharacter, lineWidth))
}

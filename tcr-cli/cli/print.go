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

package cli

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/logrusorgru/aurora"
	"strconv"
	"strings"
)

const (
	defaultLinePrefix = ""

	horizontalLineCharacter = "-"

	// Default terminal width if current terminal is not recognized
	defaultTerminalWidth = 80
)

var (
	colorizer  = aurora.NewAurora(true)
	linePrefix = defaultLinePrefix
)

func setLinePrefix(value string) {
	linePrefix = value
}

func printPrefixedAndColored(fgColor aurora.Color, message string) {
	setupConsole()
	fmt.Println(
		colorizer.Colorize(linePrefix, fgColor),
		colorizer.Colorize(message, fgColor))
}

func printInCyan(a ...interface{}) {
	printPrefixedAndColored(aurora.CyanFg, fmt.Sprint(a...))
}

func printInYellow(a ...interface{}) {
	printPrefixedAndColored(aurora.YellowFg, fmt.Sprint(a...))
}

func printInGreen(a ...interface{}) {
	printPrefixedAndColored(aurora.GreenFg, fmt.Sprint(a...))
}

func printInRed(a ...interface{}) {
	printPrefixedAndColored(aurora.RedFg, fmt.Sprint(a...))
}

func printUntouched(a ...interface{}) {
	fmt.Println(a...)
}

func printHorizontalLine() {
	termWidth := getTerminalColumns()
	prefixWidth := len(linePrefix) + 1
	lineWidth := termWidth - prefixWidth - 1
	if lineWidth < 0 {
		lineWidth = 0
	}
	horizontalLine := strings.Repeat(horizontalLineCharacter, lineWidth)
	printInCyan(horizontalLine)
}

func getTerminalColumns() int {
	output, err := sh.Command("tput", "cols").Output()
	if err != nil {
		return defaultTerminalWidth
	}

	columns, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return defaultTerminalWidth
	}
	return columns
}

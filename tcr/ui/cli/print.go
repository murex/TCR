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
	colorizer      = aurora.NewAurora(true)
	linePrefix     = defaultLinePrefix
)

func setLinePrefix(value string) {
	linePrefix = value
}

func printPrefixedAndColored(fgColor aurora.Color, message string) {
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


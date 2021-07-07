package trace

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/logrusorgru/aurora"
	"os"
	"strconv"
	"strings"
)

const (
	linePrefix = "[TCR]"

	horizontalLineCharacter = "-"

	// Default terminal width if current terminal is not recognized
	defaultTerminalWidth = 80
)

var (
	colorizer = aurora.NewAurora(true)
)

func printColoredLine(fgColor aurora.Color, message string) {
	fmt.Println(
		colorizer.Colorize(linePrefix, fgColor),
		colorizer.Colorize(message, fgColor))
}

func Info(a ...interface{}) {
	printColoredLine(aurora.CyanFg, fmt.Sprint(a...))
}

func Warning(a ...interface{}) {
	printColoredLine(aurora.YellowFg, fmt.Sprint(a...))
}

func Error(a ...interface{}) {
	printColoredLine(aurora.RedFg, fmt.Sprint(a...))
	os.Exit(1)
}

func HorizontalLine() {
	termWidth := getTermColumns()
	prefixWidth := len(linePrefix) + 1
	lineWidth := termWidth - prefixWidth - 1
	if lineWidth < 0 {
		lineWidth = 0
	}
	horizontalLine := strings.Repeat(horizontalLineCharacter, lineWidth)
	Info(horizontalLine)
}

func getTermColumns() int {
	output, err := sh.Command("tput", "cols").Output()
	if err != nil {
		return defaultTerminalWidth
	}

	termColumns, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return defaultTerminalWidth
	}
	return termColumns
}

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
	colorizer      = aurora.NewAurora(true)
	exitFunction   = os.Exit
	exitReturnCode = 0
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
	exitFunction(1)
}

func Echo(a ...interface{}) {
	fmt.Println(a...)
}

func HorizontalLine() {
	termWidth := getTerminalColumns()
	prefixWidth := len(linePrefix) + 1
	lineWidth := termWidth - prefixWidth - 1
	if lineWidth < 0 {
		lineWidth = 0
	}
	horizontalLine := strings.Repeat(horizontalLineCharacter, lineWidth)
	Info(horizontalLine)
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

// --------------------------------------------------------------------------
// For testability purpose: When running in test mode, we bypass the os.Exit()
// call, so that we can check the return code if needed
// These functions should not be used when running in production mode

func SetTestMode() {
	Warning("Running in test mode")
	exitFunction = func(rc int) {
		exitReturnCode = rc
	}
}

func GetExitReturnCode() int {
	return exitReturnCode
}
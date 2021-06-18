package trace

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"os"
	"strconv"
	"strings"
)

const (
	linePrefix = "[TCR]"

	// Default terminal width if current terminal is not recognized
	terminalDefaultWidth = 80

	colorReset  = 0
	colorRed    = 31
	colorYellow = 33
	colorCyan   = 36
)

func printColoredLine(message string, colorCode int) {
	colored := fmt.Sprintf("\x1b[%dm%s %s\x1b[%dm", colorCode, linePrefix, message, colorReset)
	fmt.Println(colored)
}

func Info(message string) {
	// 36 = Cyan
	printColoredLine(message, colorCyan)
}

func Warning(message string) {
	// 33 = Yellow
	printColoredLine(message, colorYellow)
}

func Error(message string) {
	// 31 = Red
	printColoredLine(message, colorRed)
	os.Exit(1)
}

func HorizontalLine() {
	termWidth := getTermColumns()
	prefixWidth := len(linePrefix) + 1
	lineWidth := termWidth - prefixWidth - 1
	if lineWidth < 0 {
		lineWidth = 0
	}
	horizontalLine := strings.Repeat("-", lineWidth)
	Info(horizontalLine)
}

func getTermColumns() int {
	output, err1 := sh.Command("tput", "cols").Output()
	if err1 != nil {
		return terminalDefaultWidth
	}
	termColumns, err2 := strconv.Atoi(strings.TrimSpace(string(output)))
	if err2 != nil {
		return terminalDefaultWidth
	}
	return termColumns
}

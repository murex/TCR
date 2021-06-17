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
	termColumns, err := getTermColumns()
	if err {
		termColumns = 40 // Default width
	}

	prefixSize := len(linePrefix) + 1
	repeated := termColumns - prefixSize
	horizontalLine := strings.Repeat("-", repeated)
	Info(horizontalLine)
	message := fmt.Sprintf("Width: %v", termColumns)
	Info(message)
}

func getTermColumns() (int, bool) {
	termColumns, err := cols()
	return termColumns, err
}

func cols() (int, bool) {
	output, err := sh.Command("tput", "cols").Output()
	cols, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, true
	}
	return cols, false
}

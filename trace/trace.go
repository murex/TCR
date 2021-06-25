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

	horizontalLineCharacter = "-"

	// Default terminal width if current terminal is not recognized
	terminalDefaultWidth = 80

	ansiColorReset  = 0
	ansiColorRed    = 31
	ansiColorYellow = 33
	ansiColorCyan   = 36
)

func printColoredLine(message string, colorCode int) {
	colored := fmt.Sprintf("\x1b[%dm%s %s\x1b[%dm", colorCode, linePrefix, message, ansiColorReset)
	fmt.Println(colored)
}

func Info(a ...interface{}) {
	message := fmt.Sprint(a)
	printColoredLine(message, ansiColorCyan)
}

func Warning(a ...interface{}) {
	message := fmt.Sprint(a)
	printColoredLine(message, ansiColorYellow)
}

func Error(a ...interface{}) {
	message := fmt.Sprint(a)
	printColoredLine(message, ansiColorRed)
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

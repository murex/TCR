package trace

import (
	"fmt"
	"github.com/buger/goterm"
	"os"
)

func printColoredLine(message string, colorCode int) {
	linePrefix := "[TCR]"
	colored := fmt.Sprintf("\x1b[%dm%s %s\x1b[0m", colorCode, linePrefix, message)
	fmt.Println(colored)
}

func Info(message string) {
	// 36 = Cyan
	printColoredLine(message, 36)
}

func Warning(message string) {
	// 33 = Yellow
	printColoredLine(message,33)
}

func Error(message string) {
	// 31 = Red
	printColoredLine(message, 31)
	os.Exit(1)
}

func HorizontalLine() {
	termColumns := goterm.Width()
	message := fmt.Sprintf("Width: %v",termColumns)
	Info(message)
}
//tcr_horizontal_line() {
//term_columns=$(tput cols)
//repeated=$(("${term_columns}" - 7))
//line=$(head -c "${repeated}" </dev/zero | tr '\0' '-')
//tcr_info "$line"
//}

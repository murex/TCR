package tcr

import (
	"github.com/mengdaming/tcr/stty"
	"github.com/mengdaming/tcr/trace"

	"os"
)


func mobMainMenu() {
	printOptionsMenu()

	_ = stty.SetRaw()
	defer stty.Restore()

	keyboardInput := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(keyboardInput)
		if err != nil {
			trace.Warning("Something went wrong while reading from stdin: ", err)
		}

		switch keyboardInput[0] {
		case 'd', 'D':
			runAsDriver()
		case 'n', 'N':
			runAsNavigator()
		case 'p', 'P':
			toggleAutoPush()
			printTCRHeader()
		case 'q', 'Q':
			stty.Restore()
			quit()
		default:
			trace.Warning("No action is mapped to shortcut '",
				string(keyboardInput), "'" )
		}
		printOptionsMenu()
	}
}

func printOptionsMenu() {
	trace.HorizontalLine()
	trace.Info("What shall we do?")
	trace.Info("\tD -> Driver mode")
	trace.Info("\tN -> Navigator mode")
	trace.Info("\tP -> Turn on/off git auto-push")
	trace.Info("\tQ -> Quit")
}


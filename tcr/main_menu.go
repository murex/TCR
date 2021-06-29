package tcr

import (
	"github.com/mengdaming/tcr/stty"
	"github.com/mengdaming/tcr/trace"
	"os"
)

//var initialSttyState bytes.Buffer

func mainMenu() {
	printOptionsMenu()

	_ = stty.SetRaw()
	defer stty.Restore()

	keyboardInput := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(keyboardInput)
		if err != nil {
			trace.Warning("Something went wrong while reading from stdin: ", err)
		}
		//trace.Info("Read character: ", keyboardInput)

		switch keyboardInput[0] {
		case 'd', 'D':
			runAsDriver()
		case 'n', 'N':
			runAsNavigator()
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

//+build !windows

package tcr

import (
	"bytes"
	"github.com/mengdaming/tcr/stty"
	"github.com/mengdaming/tcr/trace"
	"os"
)

var initialSttyState bytes.Buffer

func mainMenu() {
	printOptionsMenu()

	initialSttyState = stty.SetRaw()
	defer stty.Restore(&initialSttyState)

	b := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(b)
		if err != nil {
			trace.Warning(err)
		}
		//trace.Info("Read character: ", b)

		switch b[0] {
		case 'd', 'D':
			runAsDriver()
		case 'n', 'N':
			runAsNavigator()
		case 'q', 'Q':
			stty.Restore(&initialSttyState)
			quit()
		}
		printOptionsMenu()
	}
}

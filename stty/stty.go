package stty

import (
	"bytes"
	"github.com/mengdaming/tcr/trace"
	"os"
	"os/exec"
)

func readStty(state *bytes.Buffer) (err error) {
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	cmd.Stdout = state
	return cmd.Run()
}

func setStty(state *bytes.Buffer) (err error) {
	cmd := exec.Command("stty", state.String())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func SetRaw() bytes.Buffer {
	var initialState bytes.Buffer
	err := readStty(&initialState)
	if err != nil {
		trace.Error(err)
	}

	cbreakErr := setStty(bytes.NewBufferString("cbreak"))
	if cbreakErr != nil {
		trace.Error("stty cbreak: ", cbreakErr)
	}

	echoErr := setStty(bytes.NewBufferString("-echo"))
	if echoErr != nil {
		trace.Error("stty -echo: ", cbreakErr)
	}

	return initialState
}

func Restore(state *bytes.Buffer) {
	//trace.Info("Restoring stty initial state")
	err := setStty(state)
	if err != nil {
		trace.Error(err)
	}
}

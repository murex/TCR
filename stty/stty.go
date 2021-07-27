package stty

import (
	"github.com/mengdaming/tcr/trace"

	"bytes"
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
	//trace.Info("Command: ", cmd)
	return cmd.Run()
}

func SetRaw() bytes.Buffer {
	var initialState bytes.Buffer
	err := readStty(&initialState)
	//trace.Info(initialState.String())
	if err != nil {
		trace.Error("stty -g: ", err)
	}

	cbreakErr := setStty(bytes.NewBufferString("cbreak"))
	if cbreakErr != nil {
		trace.Error("stty cbreak: ", cbreakErr)
	}

	echoErr := setStty(bytes.NewBufferString("-echo"))
	if echoErr != nil {
		trace.Error("stty -echo: ", echoErr)
	}

	return initialState
}

//func Restore(state *bytes.Buffer) {
func Restore() {
	// For some unknown reason restoring previous stty state
	// fails on WSL, while working correctly on Windows git bash
	// Still need to test it on MacOS and on a non-WSL Linux box
	// Until then we set back echo and -cbreak instead of
	// restoring the previous state
	//trace.Info("Restoring stty initial state")
	_ = setStty(bytes.NewBufferString("-cbreak"))
	_ = setStty(bytes.NewBufferString("echo"))
}



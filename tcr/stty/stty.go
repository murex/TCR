package stty

import (
	"bytes"
	"github.com/mengdaming/tcr/tcr/report"
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
	//report.PostInfo("Command: ", cmd)
	return cmd.Run()
}

func SetRaw() bytes.Buffer {
	var initialState bytes.Buffer
	err := readStty(&initialState)
	//report.PostInfo(initialState.String())
	if err != nil {
		report.PostError("stty -g: ", err)
	}

	cbreakErr := setStty(bytes.NewBufferString("cbreak"))
	if cbreakErr != nil {
		report.PostError("stty cbreak: ", cbreakErr)
	}

	echoErr := setStty(bytes.NewBufferString("-echo"))
	if echoErr != nil {
		report.PostError("stty -echo: ", echoErr)
	}

	return initialState
}

func Restore() {
	//func Restore(state *bytes.Buffer) {
	// For some unknown reason restoring previous stty state
	// fails on WSL, while working correctly on Windows git bash
	// Still need to test it on MacOS and on a non-WSL Linux box
	// Until then we set back echo and -cbreak instead of
	// restoring the previous state
	//report.PostInfo("Restoring stty initial state")
	_ = setStty(bytes.NewBufferString("-cbreak"))
	_ = setStty(bytes.NewBufferString("echo"))
}



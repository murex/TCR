package sandbox

import (
	"fmt"
	"github.com/go-cmd/cmd"
	"time"
)

func CmdSandbox() {
	tryGoCmdAsync()
}

func tryGoCmdAsync() {
	// Start a long-running process, capture stdout and stderr
	findCmd := cmd.NewCmd("find", ".", "-name", "\\*.go")
	statusChan := findCmd.Start() // non-blocking

	ticker := time.NewTicker(2 * time.Second)

	// Print last line of stdout every 2s
	go func() {
		for range ticker.C {
			status := findCmd.Status()
			n := len(status.Stdout)
			fmt.Println(status.Stdout[n-1])
		}
	}()

	// Stop command after 1 hour
	go func() {
		<-time.After(1 * time.Hour)
		err := findCmd.Stop()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}()

	// Check if command is done
	select {
	case finalStatus := <-statusChan:
		fmt.Println("Final Status:", finalStatus)
		// done
	default:
		// no, still running
	}

	// Block waiting for command to exit, be stopped, or be killed
	finalStatus := <-statusChan
	fmt.Println("Final Status:", finalStatus)
}

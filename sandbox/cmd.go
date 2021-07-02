package sandbox

import (
	"fmt"
	"github.com/go-cmd/cmd"
	"golang.org/x/net/context"
	"gopkg.in/tomb.v2"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func CmdSandbox() {
	//tryGoCmdAsync()
	//tryLookPath()
	//tryContext()
	tryTombV2()
}

func tryGoCmdAsync() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	// Start a long-running process, capture stdout and stderr
//	findCmd := cmd.NewCmd("find", ".", "-name", "\\*.go")
	findCmd := cmd.NewCmd("ls", "-a")
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

func tryLookPath() {
	command := "ls"
	goExecPath, err := exec.LookPath(command)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println(command, "executable:", goExecPath)
	}
}

func tryContext() {
	c1, cancel := context.WithCancel(context.Background())

	exitCh := make(chan struct{})
	go func(ctx context.Context) {
		for {
			fmt.Println("In loop. Press ^C to stop.")
			// Do something useful in a real usecase.
			// Here we just sleep for this example.
			time.Sleep(time.Second)

			select {
			case <-ctx.Done():
				fmt.Println("received done, exiting in 500 milliseconds")
				time.Sleep(500 * time.Millisecond)
				exitCh <- struct{}{}
				return
			default:
			}
		}
	}(c1)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		select {
		case <-signalCh:
			cancel()
			return
		}
	}()
	<-exitCh
}

// ------------------------------------------------

func tryTombV2() {
	var t tomb.Tomb
	t.Go(func() error {
		for {
			time.Sleep(2 * time.Second)
			select {
			case <-t.Dying():
				fmt.Println("terminate the goroutine")
				return nil
			default:
				r, err := http.Get("https://mcorbin.fr")
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("%d\n", r.StatusCode)
			}
		}
	})
	t.Go(func() error {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		s := <-sig
		fmt.Printf("received signal %s\n", s)
		t.Kill(nil)
		return nil
	})
	err := t.Wait()
	if err != nil {
		fmt.Println(err)
	}
}
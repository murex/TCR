package sandbox

import (
	"bufio"
	"fmt"
	"github.com/daspoet/gowinkey"
	"github.com/mengdaming/tcr/trace"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func KeyboardSandbox() {

	trace.HorizontalLine()
	trace.Info("Experimenting with keyboard input (Windows OS)")

	tryGoWinKey()
	//tryScanner()

}

func tryTerm() {
	if ! terminal.IsTerminal(0) {
		b, _ := ioutil.ReadAll(os.Stdin)
		fmt.Println(string(b))
	} else {
		fmt.Println("no piped data")
	}
}

func tryTermNoEof() {
	r := strings.NewReader("some io.Reader stream to be read\n")

	buf := make([]byte, 1)
	if _, err := io.ReadAtLeast(os.Stdin, buf, 1); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf)

	// buffer smaller than minimal read size.
	shortBuf := make([]byte, 3)
	if _, err := io.ReadAtLeast(r, shortBuf, 1); err != nil {
		fmt.Println("error:", err)
	}

	// minimal read size bigger than io.Reader stream
	longBuf := make([]byte, 64)
	if _, err := io.ReadAtLeast(r, longBuf, 1); err != nil {
		fmt.Println("error:", err)
	}

}

func tryGoWinKey() {
	events, _ := gowinkey.ListenSelective(gowinkey.VK_W, gowinkey.VK_A, gowinkey.VK_S, gowinkey.VK_D)

	//	timer := time.AfterFunc(time.Second * 5, stopFn)
	//	defer timer.Stop()

	for e := range events {
		switch e.Type {
		case gowinkey.KeyPressed:
			fmt.Println("pressed", e)
		case gowinkey.KeyReleased:
			fmt.Println("released", e)
		}
	}
}

func tryScanner() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println("You typed:", scanner.Text())
	}

	if scanner.Err() != nil {
		// handle error
	}
}


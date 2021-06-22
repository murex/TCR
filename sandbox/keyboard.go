//+build !windows

package sandbox

import (
	"fmt"
	"github.com/mengdaming/tcr/trace"
	"github.com/pkg/term"
	"github.com/tj/go-terminput"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"os"
)

func KeyboardSandbox() {

		trace.HorizontalLine()
		trace.Info("Experimenting with keyboard input (non-Windows OS)")

		tryTermInput()

}

func tryTerm() {
	if ! terminal.IsTerminal(0) {
		b, _ := ioutil.ReadAll(os.Stdin)
		fmt.Println(string(b))
	} else {
		fmt.Println("no piped data")
	}
}


func tryTermInput() {
	t, err := term.Open("/dev/tty")
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	t.SetRaw()
	defer t.Restore()

	fmt.Printf("Type something, use 'q' to exit.\r\n")

	for {
		e, err := terminput.Read(t)
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}

		if e.Key() == terminput.KeyEscape || e.Rune() == 'q' {
			break
		}

		fmt.Printf("%s — shift=%v ctrl=%v alt=%v meta=%v\r\n", e.String(), e.Shift(), e.Ctrl(), e.Alt(), e.Meta())
	}
}



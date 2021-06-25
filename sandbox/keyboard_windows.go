package sandbox

import (
	"bufio"
	"fmt"
	"github.com/containerd/console"
	"github.com/daspoet/gowinkey"
	"github.com/janmir/go-winput"
	"github.com/mengdaming/tcr/trace"
	"github.com/vence722/inputhook"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/term"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func KeyboardSandbox() {

	trace.HorizontalLine()
	trace.Info("Experimenting with keyboard input (Windows OS)")

	//tryGoWinKey()
	//tryGoWinput()
	//tryReadRune()
	//tryInputHook()
	//tryBufioScanner()
	//tryWithMakeRaw()
	//tryConsole()
	tryGetTermDimWithSttyCommand()

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

func tryGoWinput() {
	input := winput.New()
	input.Type("Hello from Winput!")
	ok := input.HotKey(winput.HotKeySelectAll)
	trace.Info("Type status: ", ok)
}

func tryReadRune() {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println(err)
	}

	// prints the unicode code point of the character
	fmt.Println("Result: ", char)
}

func tryInputHook() {
	inputhook.HookKeyboard(hookCallback)
	ch := make(chan bool)
	<-ch
}

func hookCallback(keyEvent int, keyCode int) {
	fmt.Println("keyEvent:", keyEvent)
	fmt.Println("keyCode:", keyCode)
}

func tryBufioScanner() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func tryWithMakeRaw() {
	// fd 0 is stdin
	var state, err = term.MakeRaw(0)
	if err != nil {
		log.Fatalln("setting stdin to raw:", err)
	}
	defer func() {
		if err := term.Restore(0, state); err != nil {
			log.Println("warning, failed to restore terminal:", err)
		}
	}()

	in := bufio.NewReader(os.Stdin)
	for {
		r, _, err := in.ReadRune()
		if err != nil {
			log.Println("stdin:", err)
			break
		}
		fmt.Printf("read rune %q\r\n", r)
		if r == 'q' {
			break
		}
	}
}

func tryConsole() {
	current := console.Current()
	defer current.Reset()

	if err := current.SetRaw(); err != nil {
	}
	ws, _ := current.Size()
	current.Resize(ws)

}

func tryGetTermDimWithSttyCommand() () {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	var termDim []byte
	var err error
	if termDim, err = cmd.Output(); err != nil {
		return
	}
	var width, height int
	_,_ = fmt.Sscan(string(termDim), &height, &width)
	fmt.Println("Height", height, "Width", width)
	return
}


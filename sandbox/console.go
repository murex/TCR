package sandbox

import (
	"fmt"
	"github.com/DrSmithFr/go-console/pkg/input/argument"
	"github.com/DrSmithFr/go-console/pkg/input/option"
	"github.com/DrSmithFr/go-console/pkg/style"
)

func ConsoleSandbox() {
	tryGoConsole()
}

func tryGoConsole() {
	io := style.
		NewConsoleStyler().
		AddInputArgument(
			argument.
				New("name", argument.REQUIRED),
		).
		AddInputOption(
			option.
				New("foo", option.NONE).
				SetShortcut("f"),
		).
		ParseInput().
		ValidateInput()

	// enable stylish errors
	defer io.HandleRuntimeException()

	//
	// Do what you want with console and args !
	//

	io.Title("Starting console")

	io.TextArray([]string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	})

	io.Note(
		fmt.Sprintf(
			"name argument value '%s'", "toto",
	//		io.GetInput().GetArgument("name"),
		),
	)

	if option.DEFINED == io.GetInput().GetOption("foo") {
		io.Success("foo option is set")
	}

	//panic("this error will be stylish!")
}
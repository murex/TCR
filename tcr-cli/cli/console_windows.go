package cli

import (
	"golang.org/x/sys/windows"
	"os"
)

func setupConsole() {
	// This enforces that ANSI color codes are properly interpreted when running on Windows OS
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	_ = windows.GetConsoleMode(stdout, &originalMode)
	_ = windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}

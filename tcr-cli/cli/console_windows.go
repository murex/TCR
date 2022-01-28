/*
Copyright (c) 2022 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package cli

import (
	"golang.org/x/sys/windows"
	"os"
)

func setupConsole() {
	// Cf. https://docs.microsoft.com/en-us/windows/console/setconsolemode
	setupWindowsConsoleStdout()
	setupWindowsConsoleStdin()
}

func setupWindowsConsoleStdin() {
	// This enforces that we do not wait for enter key and do not have keystrokes echo-ed on console output
	stdin := windows.Handle(os.Stdin.Fd())
	var originalMode uint32
	_ = windows.GetConsoleMode(stdin, &originalMode)
	_ = windows.SetConsoleMode(stdin, originalMode&^
		uint32(windows.ENABLE_LINE_INPUT)&^
		uint32(windows.ENABLE_ECHO_INPUT))
}

func setupWindowsConsoleStdout() {
	// This enforces that ANSI color codes are properly interpreted when running on Windows OS
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32
	_ = windows.GetConsoleMode(stdout, &originalMode)
	_ = windows.SetConsoleMode(stdout, originalMode|
		windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}

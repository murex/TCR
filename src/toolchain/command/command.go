/*
Copyright (c) 2024 Murex

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

package command

import (
	"errors"
	"runtime"
	"strings"
)

type (
	// OsName is the name of a supported operating system
	OsName string

	// ArchName is the name of a supported architecture
	ArchName string

	// Command is a command that can be run by a toolchain.
	// It contains 2 filters (Os and Arch) allowing to restrict it to specific OS(s)/Architecture(s).
	// - Path is the path to the command to be run.
	// - Arguments is the arguments to be passed to the command when executed.
	Command struct {
		Os        []OsName
		Arch      []ArchName
		Path      string
		Arguments []string
	}
)

// List of possible values for OsName
const (
	OsDarwin  OsName = "darwin"
	OsLinux   OsName = "linux"
	OsWindows OsName = "windows"
)

// List of possible values for ArchName
const (
	Arch386   ArchName = "386"
	ArchAmd64 ArchName = "amd64"
	ArchArm64 ArchName = "arm64"
)

// GetAllOsNames return the list of all supported OS Names
func GetAllOsNames() []OsName {
	return []OsName{OsDarwin, OsLinux, OsWindows}
}

// GetAllArchNames return the list of all supported OS Architectures
func GetAllArchNames() []ArchName {
	return []ArchName{Arch386, ArchAmd64, ArchArm64}
}

func (command Command) runsOnLocalMachine() bool {
	return command.runsOnPlatform(OsName(runtime.GOOS), ArchName(runtime.GOARCH))
}

func (command Command) runsOnPlatform(osName OsName, archName ArchName) bool {
	return command.runsWithOs(osName) && command.runsWithArch(archName)
}

func (command Command) runsWithOs(osName OsName) bool {
	for _, o := range command.Os {
		if o == osName {
			return true
		}
	}
	return false
}

func (command Command) runsWithArch(archName ArchName) bool {
	for _, a := range command.Arch {
		if a == archName {
			return true
		}
	}
	return false
}

func (command Command) check() error {
	if err := command.checkPath(); err != nil {
		return err
	}
	if err := command.checkOsTable(); err != nil {
		return err
	}
	return command.checkArchTable()
}

func (command Command) checkPath() error {
	if command.Path == "" {
		return errors.New("command path is empty")
	}
	return nil
}

func (command Command) checkOsTable() error {
	if command.Os == nil {
		return errors.New("command's OS list is empty")
	}
	for _, osName := range command.Os {
		if osName == "" {
			return errors.New("a command OS name is empty")
		}
	}
	return nil
}

func (command Command) checkArchTable() error {
	if command.Arch == nil {
		return errors.New("command's architecture list is empty")
	}
	for _, archName := range command.Arch {
		if archName == "" {
			return errors.New("a command architecture name is empty")
		}
	}
	return nil
}

// AsCommandLine returns a string formatted like if the command was run from the command line
func (command Command) AsCommandLine() string {
	return command.Path + " " + strings.Join(command.Arguments, " ")
}

// FindCommand retrieves out of a list of commands the first that is compatible with provided OS and Architecture
func FindCommand(commands []Command, osName OsName, archName ArchName) *Command {
	for _, cmd := range commands {
		if cmd.runsOnPlatform(osName, archName) {
			return &cmd
		}
	}
	return nil
}

// FindCompatibleCommand retrieves out of a list of commands one that is compatible
// with local machine's OS and Architecture
func FindCompatibleCommand(commands []Command) *Command {
	for _, command := range commands {
		if command.runsOnLocalMachine() {
			return &command
		}
	}
	return nil
}

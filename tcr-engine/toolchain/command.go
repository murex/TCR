/*
Copyright (c) 2021 Murex

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

package toolchain

import (
	"errors"
	"github.com/codeskyblue/go-sh"
	"github.com/murex/tcr/tcr-engine/report"
	"os"
	"path/filepath"
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
	OsDarwin  = "darwin"
	OsLinux   = "linux"
	OsWindows = "windows"
)

// List of possible values for OsArch
const (
	Arch386   = "386"
	ArchAmd64 = "amd64"
	ArchArm64 = "arm64"
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

func (command Command) runsOnPlatform(os OsName, arch ArchName) bool {
	return command.runsWithOs(os) && command.runsWithArch(arch)
}

func (command Command) runsWithOs(os OsName) bool {
	for _, osName := range command.Os {
		if osName == os {
			return true
		}
	}
	return false
}

func (command Command) runsWithArch(arch ArchName) bool {
	for _, archName := range command.Arch {
		if archName == arch {
			return true
		}
	}
	return false
}

func (command Command) run() error {
	report.PostText(command.asCommandLine())
	output, err := sh.Command(adjustCommandPath(command.Path), command.Arguments).CombinedOutput()
	if output != nil {
		report.PostText(string(output))
	}
	return err
}

func (command Command) check() error {
	if err := command.checkPath(); err != nil {
		return err
	}
	if err := command.checkOsTable(); err != nil {
		return err
	}
	if err := command.checkArchTable(); err != nil {
		return err
	}
	return nil
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

func (command Command) asCommandLine() string {
	return adjustCommandPath(command.Path) + " " + strings.Join(command.Arguments, " ")
}

func findCommand(commands []Command, os OsName, arch ArchName) *Command {
	for _, cmd := range commands {
		if cmd.runsOnPlatform(os, arch) {
			return &cmd
		}
	}
	return nil
}

func findCompatibleCommand(commands []Command) *Command {
	for _, command := range commands {
		if command.runsOnLocalMachine() {
			return &command
		}
	}
	return nil
}

func adjustCommandPath(cmdPath string) string {
	// If this is an absolute path, we return it after cleaning it up
	if filepath.IsAbs(cmdPath) {
		return filepath.Clean(cmdPath)
	}
	// If not, we check if it can be a relative path from the current directory.
	// If the file is found, we return it
	wd, _ := os.Getwd()
	pathFromWorkingDir := filepath.Join(wd, cmdPath)
	info, err := os.Stat(pathFromWorkingDir)
	if err == nil && !info.IsDir() {
		return pathFromWorkingDir
	}
	// As a last resort, we assume it's available in the $PATH
	return filepath.Clean(cmdPath)
}

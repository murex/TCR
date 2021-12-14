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
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/murex/tcr/tcr-engine/report"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// TchnInterface provides the interface that any toolchain needs to implement so that it can be used by TcR
type TchnInterface interface {
	GetName() string
	RunBuild() error
	RunTests() error
	BuildCommandName() string
	BuildCommandArgs() []string
	TestCommandName() string
	TestCommandArgs() []string
	reset()
}

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

	// Toolchain defines the structure of a toolchain.
	// - Name is the name of the toolchain, it must be unique in the list of available toolchains
	// - BuildCommands is a table of commands that can be called when running the build. The first one
	// matching the current OS and configuration will be the one to be called.
	// - TestCommands is a table of commands that can be called when running the tests. The first one
	// matching the current OS and configuration will be the one to be called.
	Toolchain struct {
		Name          string
		BuildCommands []Command
		TestCommands  []Command
	}
)

// List of possible values for OsName
const (
	OsDarwin  = "darwin"
	OsLinux   = "linux"
	OsWindows = "windows"
)

// GetAllOsNames return the list of all supported OS Names
func GetAllOsNames() []OsName {
	return []OsName{OsDarwin, OsLinux, OsWindows}
}

// List of possible values for OsArch
const (
	Arch386   = "386"
	ArchAmd64 = "amd64"
	ArchArm64 = "arm64"
)

// GetAllArchNames return the list of all supported OS Architectures
func GetAllArchNames() []ArchName {
	return []ArchName{Arch386, ArchAmd64, ArchArm64}
}

var (
	builtInToolchains   = make(map[string]Toolchain)
	supportedToolchains = make(map[string]Toolchain)
)

func addSupportedToolchain(tchn Toolchain) {
	supportedToolchains[strings.ToLower(tchn.GetName())] = tchn
}

func isSupported(name string) bool {
	_, found := supportedToolchains[strings.ToLower(name)]
	return found
}

// GetToolchain returns the toolchain instance with the provided name
// The toolchain name is case insensitive.
func GetToolchain(name string) (*Toolchain, error) {
	if name == "" {
		return nil, errors.New("toolchain name not provided")
	}
	tchn, found := supportedToolchains[strings.ToLower(name)]
	if found {
		return &tchn, nil
	}
	return nil, errors.New(fmt.Sprint("toolchain not supported: ", name))
}

// Names returns the list of available toolchain names
func Names() []string {
	var names []string
	for _, tchn := range supportedToolchains {
		names = append(names, tchn.Name)
	}
	return names
}

// Reset resets the toolchain with the provided name to its default values
func Reset(name string) {
	//_, found := supportedToolchains[strings.ToLower(name)]
	//if found {
	// TODO
	//tchn.reset()
	//}
}

func runBuild(toolchain TchnInterface) error {
	return runCommand(toolchain.BuildCommandName(), toolchain.BuildCommandArgs())
}

func runTests(toolchain TchnInterface) error {
	return runCommand(toolchain.TestCommandName(), toolchain.TestCommandArgs())
}

func runCommand(cmdPath string, cmdArgs []string) error {
	output, err := sh.Command(tuneCommandPath(cmdPath), cmdArgs).CombinedOutput()
	if output != nil {
		report.PostText(string(output))
	}
	return err
}

func tuneCommandPath(cmdPath string) string {
	// TODO handle different types of paths (relative, absolute, no path)
	wd, _ := os.Getwd()
	return filepath.Join(wd, cmdPath)
}

func isBuiltIn(name string) bool {
	_, found := builtInToolchains[strings.ToLower(name)]
	return found
}

func addBuiltInToolchain(tchn Toolchain) error {
	if tchn.Name == "" {
		return errors.New("toolchain name cannot be an empty string")
	}
	builtInToolchains[strings.ToLower(tchn.Name)] = tchn
	addSupportedToolchain(tchn)
	return nil
}

func (tchn Toolchain) reset() {
	//TODO implement me
}

// GetName provides the name of the toolchain
func (tchn Toolchain) GetName() string {
	return tchn.Name
}

// RunBuild runs the build with this toolchain
func (tchn Toolchain) RunBuild() error {
	return runBuild(tchn)
}

// RunTests runs the tests with this toolchain
func (tchn Toolchain) RunTests() error {
	return runTests(tchn)
}

// BuildCommandName returns the build command name for this toolchain
func (tchn Toolchain) BuildCommandName() string {
	var cmd = findCompatibleCommand(tchn.BuildCommands)
	return cmd.Path
}

// BuildCommandArgs returns a table with the list of build command arguments for this toolchain
func (tchn Toolchain) BuildCommandArgs() []string {
	var cmd = findCompatibleCommand(tchn.BuildCommands)
	return cmd.Arguments
}

// TestCommandName returns the test command name for this toolchain
func (tchn Toolchain) TestCommandName() string {
	var cmd = findCompatibleCommand(tchn.TestCommands)
	return cmd.Path
}

// TestCommandArgs returns a table with the list of test command arguments for this toolchain
func (tchn Toolchain) TestCommandArgs() []string {
	var cmd = findCompatibleCommand(tchn.TestCommands)
	return cmd.Arguments
}

func findCompatibleCommand(commands []Command) *Command {
	for _, command := range commands {
		if runsOnLocalMachine(command) {
			return &command
		}
	}
	return nil
}

func runsOnLocalMachine(command Command) bool {
	return runsWithLocalOs(command) && runsWithLocalArch(command)
}

func runsWithLocalOs(command Command) bool {
	return runsWithOs(command, runtime.GOOS)
}

func runsWithOs(command Command, os string) bool {
	for _, osName := range command.Os {
		if string(osName) == os {
			return true
		}
	}
	return false
}

func runsWithLocalArch(command Command) bool {
	return runsWithArch(command, runtime.GOARCH)
}

func runsWithArch(command Command, arch string) bool {
	for _, archName := range command.Arch {
		if string(archName) == arch {
			return true
		}
	}
	return false
}

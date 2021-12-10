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
	"strings"
)

// Toolchain provides the interface that any toolchain needs to implement so that it can be used by TcR
type Toolchain interface {
	Name() string
	RunBuild() error
	RunTests() error
	BuildCommandName() string
	BuildCommandArgs() []string
	TestCommandName() string
	TestCommandArgs() []string
	reset()
}

// OsName is the name of a supported operating system
type OsName string

// ArchName is the name of a supported architecture
type ArchName string

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
	supportedToolchains = make(map[string]Toolchain)
)

func init() {
	addSupportedToolchain(GradleToolchain{})
	addSupportedToolchain(MavenToolchain{})
	addSupportedToolchain(CmakeToolchain{})
}

func addSupportedToolchain(tchn Toolchain) {
	supportedToolchains[strings.ToLower(tchn.Name())] = tchn
}

func isSupported(name string) bool {
	_, found := supportedToolchains[strings.ToLower(name)]
	return found
}

// GetToolchain returns the toolchain instance with the provided name
func GetToolchain(name string) (Toolchain, error) {
	tchn, found := supportedToolchains[strings.ToLower(name)]
	if found {
		return tchn, nil
	}
	return nil, errors.New(fmt.Sprint("toolchain not supported: ", name))
}

// Names returns the list of available toolchain names
func Names() []string {
	var names []string
	for _, tchn := range supportedToolchains {
		names = append(names, tchn.Name())
	}
	return names
}

// Reset resets the toolchain with the provided name to its default values
func Reset(name string) {
	tchn, found := supportedToolchains[strings.ToLower(name)]
	if found {
		tchn.reset()
	}
}

// New creates a new toolchain instance with the provided name.
// The toolchain name is case insensitive.
func New(name string) (Toolchain, error) {
	if name != "" {
		return GetToolchain(name)
	}
	return nil, errors.New("toolchain name not provided")
}

func runBuild(toolchain Toolchain) error {
	return runCommand(toolchain.BuildCommandName(), toolchain.BuildCommandArgs())
}

func runTests(toolchain Toolchain) error {
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

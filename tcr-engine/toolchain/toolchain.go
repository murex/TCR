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
	wd, _ := os.Getwd()
	buildCommandPath := filepath.Join(wd, toolchain.BuildCommandName())
	output, err := sh.Command(
		buildCommandPath,
		toolchain.BuildCommandArgs()).CombinedOutput()
	if output != nil {
		report.PostText(string(output))
	}
	return err
}

func runTests(tchn Toolchain) error {
	wd, _ := os.Getwd()
	testCommandPath := filepath.Join(wd, tchn.TestCommandName())
	output, err := sh.Command(
		testCommandPath,
		tchn.TestCommandArgs()).CombinedOutput()
	if output != nil {
		report.PostText(string(output))
	}
	return err
}

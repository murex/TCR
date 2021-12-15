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
	"strings"
)

type (
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

// Get returns the toolchain instance with the provided name
// The toolchain name is case insensitive.
func Get(name string) (*Toolchain, error) {
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

func isBuiltIn(name string) bool {
	_, found := builtInToolchains[strings.ToLower(name)]
	return found
}

func addBuiltIn(tchn Toolchain) error {
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
	return findCompatibleCommand(tchn.BuildCommands).run()
}

// RunTests runs the tests with this toolchain
func (tchn Toolchain) RunTests() error {
	return findCompatibleCommand(tchn.TestCommands).run()
}

// buildCommandName returns the build command name for this toolchain
func (tchn Toolchain) buildCommandName() string {
	return findCompatibleCommand(tchn.BuildCommands).Path
}

// buildCommandArgs returns a table with the list of build command arguments for this toolchain
func (tchn Toolchain) buildCommandArgs() []string {
	return findCompatibleCommand(tchn.BuildCommands).Arguments
}

// testCommandName returns the test command name for this toolchain
func (tchn Toolchain) testCommandName() string {
	return findCompatibleCommand(tchn.TestCommands).Path
}

// testCommandArgs returns a table with the list of test command arguments for this toolchain
func (tchn Toolchain) testCommandArgs() []string {
	return findCompatibleCommand(tchn.TestCommands).Arguments
}

func (tchn Toolchain) supportsPlatform(os OsName, arch ArchName) bool {
	return tchn.findBuildCommandFor(os, arch) != nil && tchn.findTestCommandFor(os, arch) != nil
}

func (tchn Toolchain) findBuildCommandFor(os OsName, arch ArchName) *Command {
	return findCommand(tchn.BuildCommands, os, arch)
}

func (tchn Toolchain) findTestCommandFor(os OsName, arch ArchName) *Command {
	return findCommand(tchn.TestCommands, os, arch)
}

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
	builtIn    = make(map[string]Toolchain)
	registered = make(map[string]Toolchain)
)

func register(tchn Toolchain) {
	registered[strings.ToLower(tchn.GetName())] = tchn
}

func isSupported(name string) bool {
	_, found := registered[strings.ToLower(name)]
	return found
}

// Get returns the toolchain instance with the provided name
// The toolchain name is case insensitive.
func Get(name string) (*Toolchain, error) {
	if name == "" {
		return nil, errors.New("toolchain name not provided")
	}
	tchn, found := registered[strings.ToLower(name)]
	if found {
		return &tchn, nil
	}
	return nil, errors.New(fmt.Sprint("toolchain not supported: ", name))
}

// Names returns the list of available toolchain names
func Names() []string {
	var names []string
	for _, tchn := range registered {
		names = append(names, tchn.Name)
	}
	return names
}

// Reset resets the toolchain with the provided name to its default values
func Reset(name string) {
	_, found := registered[strings.ToLower(name)]
	if found && isBuiltIn(name) {
		register(*getBuiltIn(name))
	}
}

func getBuiltIn(name string) *Toolchain {
	var builtIn, _ = builtIn[strings.ToLower(name)]
	return &builtIn
}

func isBuiltIn(name string) bool {
	_, found := builtIn[strings.ToLower(name)]
	return found
}

func addBuiltIn(tchn Toolchain) error {
	if tchn.Name == "" {
		return errors.New("toolchain name cannot be an empty string")
	}
	builtIn[strings.ToLower(tchn.Name)] = tchn
	register(tchn)
	return nil
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

// buildCommandPath returns the build command name for this toolchain
func (tchn Toolchain) buildCommandPath() string {
	return findCompatibleCommand(tchn.BuildCommands).Path
}

// buildCommandArgs returns a table with the list of build command arguments for this toolchain
func (tchn Toolchain) buildCommandArgs() []string {
	return findCompatibleCommand(tchn.BuildCommands).Arguments
}

// testCommandPath returns the test command name for this toolchain
func (tchn Toolchain) testCommandPath() string {
	return findCompatibleCommand(tchn.TestCommands).Path
}

// testCommandArgs returns a table with the list of test command arguments for this toolchain
func (tchn Toolchain) testCommandArgs() []string {
	return findCompatibleCommand(tchn.TestCommands).Arguments
}

func (tchn Toolchain) runsOnPlatform(os OsName, arch ArchName) bool {
	return tchn.findBuildCommandFor(os, arch) != nil && tchn.findTestCommandFor(os, arch) != nil
}

func (tchn Toolchain) findBuildCommandFor(os OsName, arch ArchName) *Command {
	return findCommand(tchn.BuildCommands, os, arch)
}

func (tchn Toolchain) findTestCommandFor(os OsName, arch ArchName) *Command {
	return findCommand(tchn.TestCommands, os, arch)
}

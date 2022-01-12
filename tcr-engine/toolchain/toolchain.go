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
	"sort"
	"strings"
)

type (
	// Toolchain defines the data structure of a toolchain.
	// - name is the name of the toolchain, it must be unique in the list of available toolchains
	// - buildCommands is a table of commands that can be called when running the build. The first one
	// matching the current OS and configuration will be the one to be called.
	// - testCommands is a table of commands that can be called when running the tests. The first one
	// matching the current OS and configuration will be the one to be called.
	Toolchain struct {
		name          string
		buildCommands []Command
		testCommands  []Command
	}

	// TchnInterface provides the interface for interacting with a toolchain
	TchnInterface interface {
		GetName() string
		GetBuildCommands() []Command
		GetTestCommands() []Command
		RunBuild() error
		RunTests() error
		checkName() error
		buildCommandPath() string
		buildCommandArgs() []string
		checkBuildCommand() error
		testCommandPath() string
		testCommandArgs() []string
		checkTestCommand() error
		runsOnPlatform(os OsName, arch ArchName) bool
	}
)

var (
	builtIn    = make(map[string]TchnInterface)
	registered = make(map[string]TchnInterface)
)

// New creates a new Toolchain instance with the provided name, buildCommands and testCommands
func New(name string, buildCommands, testCommands []Command) *Toolchain {
	return &Toolchain{
		name:          name,
		buildCommands: buildCommands,
		testCommands:  testCommands,
	}
}

// Register adds the provided toolchain to the list of supported toolchains
func Register(tchn TchnInterface) error {
	if err := tchn.checkName(); err != nil {
		return err
	}
	if err := tchn.checkBuildCommand(); err != nil {
		return err
	}
	if err := tchn.checkTestCommand(); err != nil {
		return err
	}
	registered[strings.ToLower(tchn.GetName())] = tchn
	return nil
}

func isSupported(name string) bool {
	_, found := registered[strings.ToLower(name)]
	return found
}

func (tchn Toolchain) checkName() error {
	if tchn.name == "" {
		return errors.New("toolchain name is empty")
	}
	return nil
}

func (tchn Toolchain) checkBuildCommand() error {
	if tchn.buildCommands == nil {
		return errors.New("toolchain has no build command")
	}
	return nil
}

func (tchn Toolchain) checkTestCommand() error {
	if tchn.testCommands == nil {
		return errors.New("toolchain has no test command")
	}
	return nil
}

// Get returns the toolchain instance with the provided name
// The toolchain name is case insensitive.
func Get(name string) (TchnInterface, error) {
	if name == "" {
		return nil, errors.New("toolchain name not provided")
	}
	tchn, found := registered[strings.ToLower(name)]
	if found {
		return tchn, nil
	}
	return nil, errors.New(fmt.Sprint("toolchain not supported: ", name))
}

// Names returns the list of available toolchain names sorted alphabetically
func Names() []string {
	var names []string
	for _, tchn := range registered {
		names = append(names, tchn.GetName())
	}
	sort.Strings(names)
	return names
}

// Reset resets the toolchain with the provided name to its default values
func Reset(name string) {
	_, found := registered[strings.ToLower(name)]
	if found && isBuiltIn(name) {
		_ = Register(*getBuiltIn(name))
	}
}

func getBuiltIn(name string) *TchnInterface {
	var builtIn, _ = builtIn[strings.ToLower(name)]
	return &builtIn
}

func isBuiltIn(name string) bool {
	_, found := builtIn[strings.ToLower(name)]
	return found
}

func addBuiltIn(tchn TchnInterface) error {
	if tchn.GetName() == "" {
		return errors.New("toolchain name cannot be an empty string")
	}
	builtIn[strings.ToLower(tchn.GetName())] = tchn
	return Register(tchn)
}

// GetName provides the name of the toolchain
func (tchn Toolchain) GetName() string {
	return tchn.name
}

// GetBuildCommands returns the toolchain's build commands
func (tchn Toolchain) GetBuildCommands() []Command {
	return tchn.buildCommands
}

// GetTestCommands returns the toolchain's test commands
func (tchn Toolchain) GetTestCommands() []Command {
	return tchn.testCommands
}

// RunBuild runs the build with this toolchain
func (tchn Toolchain) RunBuild() error {
	return findCompatibleCommand(tchn.buildCommands).run()
}

// RunTests runs the tests with this toolchain
func (tchn Toolchain) RunTests() error {
	return findCompatibleCommand(tchn.testCommands).run()
}

// buildCommandPath returns the build command name for this toolchain
func (tchn Toolchain) buildCommandPath() string {
	return findCompatibleCommand(tchn.buildCommands).Path
}

// buildCommandArgs returns a table with the list of build command arguments for this toolchain
func (tchn Toolchain) buildCommandArgs() []string {
	return findCompatibleCommand(tchn.buildCommands).Arguments
}

// testCommandPath returns the test command name for this toolchain
func (tchn Toolchain) testCommandPath() string {
	return findCompatibleCommand(tchn.testCommands).Path
}

// testCommandArgs returns a table with the list of test command arguments for this toolchain
func (tchn Toolchain) testCommandArgs() []string {
	return findCompatibleCommand(tchn.testCommands).Arguments
}

func (tchn Toolchain) runsOnPlatform(os OsName, arch ArchName) bool {
	return tchn.findBuildCommandFor(os, arch) != nil && tchn.findTestCommandFor(os, arch) != nil
}

func (tchn Toolchain) findBuildCommandFor(os OsName, arch ArchName) *Command {
	return findCommand(tchn.buildCommands, os, arch)
}

func (tchn Toolchain) findTestCommandFor(os OsName, arch ArchName) *Command {
	return findCommand(tchn.testCommands, os, arch)
}

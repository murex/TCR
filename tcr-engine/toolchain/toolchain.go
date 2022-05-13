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
	"os"
	"path/filepath"
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
		RunTests() (TestResults, error)
		checkName() error
		BuildCommandLine() string
		BuildCommandPath() string
		BuildCommandArgs() []string
		checkBuildCommand() error
		TestCommandLine() string
		TestCommandPath() string
		TestCommandArgs() []string
		checkTestCommand() error
		runsOnPlatform(os OsName, arch ArchName) bool
		CheckCommandAccess(cmdPath string) (string, error)
	}
)

var workDir string

// SetWorkDir sets the work directory from which toolchain commands will be launched
func SetWorkDir(dir string) (err error) {
	workDir, err = dirAbsPath(dir)
	return
}

// dirAbsPath returns the absolute path for the provided directory.
// Returns an error if the directory cannot be accessed or is not a directory
func dirAbsPath(dir string) (string, error) {
	absPath, err := filepath.Abs(dir)
	if err == nil {
		info, err := os.Stat(absPath)
		if err != nil {
			return "", errors.New("cannot access " + absPath)
		}
		if !info.IsDir() {
			return "", errors.New(absPath + " exists but is not a directory")
		}
	}
	return absPath, nil
}

// GetWorkDir returns the work directory from which toolchain commands will be launched
func GetWorkDir() string {
	return workDir
}

// New creates a new Toolchain instance with the provided name, buildCommands and testCommands
func New(name string, buildCommands, testCommands []Command) *Toolchain {
	return &Toolchain{
		name:          name,
		buildCommands: buildCommands,
		testCommands:  testCommands,
	}
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
	_, err := findCompatibleCommand(tchn.buildCommands).run()
	return err
}

// RunTests runs the tests with this toolchain
func (tchn Toolchain) RunTests() (testResults TestResults, err error) {
	var testOutput string
	testOutput, err = findCompatibleCommand(tchn.testCommands).run()
	testResults = extractTestResults(testOutput)
	return
}

// BuildCommandPath returns the build command path for this toolchain
func (tchn Toolchain) BuildCommandPath() string {
	return findCompatibleCommand(tchn.buildCommands).Path
}

// BuildCommandArgs returns a table with the list of build command arguments for this toolchain
func (tchn Toolchain) BuildCommandArgs() []string {
	return findCompatibleCommand(tchn.buildCommands).Arguments
}

// BuildCommandLine returns the toolchain's build command line as a string
func (tchn Toolchain) BuildCommandLine() string {
	return findCompatibleCommand(tchn.buildCommands).asCommandLine()
}

// TestCommandPath returns the test command path for this toolchain
func (tchn Toolchain) TestCommandPath() string {
	return findCompatibleCommand(tchn.testCommands).Path
}

// TestCommandArgs returns a table with the list of test command arguments for this toolchain
func (tchn Toolchain) TestCommandArgs() []string {
	return findCompatibleCommand(tchn.testCommands).Arguments
}

// TestCommandLine returns the toolchain's test command line as a string
func (tchn Toolchain) TestCommandLine() string {
	return findCompatibleCommand(tchn.testCommands).asCommandLine()
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

// CheckCommandAccess verifies if the provided command path can be accessed. Returns the path as
// an absolute command path if found. Returns an empty path otherwise, together with the corresponding error
func (tchn Toolchain) CheckCommandAccess(cmdPath string) (string, error) {
	return checkCommandPath(cmdPath)
}

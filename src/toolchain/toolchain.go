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

package toolchain

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/murex/tcr/report"
	"github.com/murex/tcr/toolchain/command"
	"github.com/murex/tcr/xunit"
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
		buildCommands []command.Command
		testCommands  []command.Command
		testResultDir string
	}

	// TestCommandResult is a Result enriched with test Stats
	TestCommandResult struct {
		command.Result
		Stats TestStats
	}

	// TchnInterface provides the interface for interacting with a toolchain
	TchnInterface interface {
		GetName() string
		GetBuildCommands() []command.Command
		GetTestCommands() []command.Command
		GetTestResultDir() string
		GetTestResultPath() string
		RunBuild() command.Result
		RunTests() TestCommandResult
		checkName() error
		BuildCommandLine() string
		BuildCommandPath() string
		BuildCommandArgs() []string
		checkBuildCommand() error
		TestCommandLine() string
		TestCommandPath() string
		TestCommandArgs() []string
		checkTestCommand() error
		runsOnPlatform(osName command.OsName, archName command.ArchName) bool
		CheckCommandAccess(cmdPath string) (string, error)
		AbortExecution() bool
	}
)

var workDir string

// SetWorkDir sets the work directory from which toolchain commands will be launched
func SetWorkDir(dir string) (err error) {
	workDir, err = dirAbsPath(dir)
	return err
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
func New(name string, buildCommands, testCommands []command.Command, testResultDir string) *Toolchain {
	return &Toolchain{
		name:          name,
		buildCommands: buildCommands,
		testCommands:  testCommands,
		testResultDir: testResultDir,
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
func (tchn Toolchain) GetBuildCommands() []command.Command {
	return tchn.buildCommands
}

// GetTestCommands returns the toolchain's test commands
func (tchn Toolchain) GetTestCommands() []command.Command {
	return tchn.testCommands
}

// RunBuild runs the build with this toolchain
func (tchn Toolchain) RunBuild() command.Result {
	cmd := command.FindCompatibleCommand(tchn.buildCommands)
	return command.GetRunner().Run(GetWorkDir(), cmd)
}

// RunTests runs the tests with this toolchain
func (tchn Toolchain) RunTests() TestCommandResult {
	cmd := command.FindCompatibleCommand(tchn.testCommands)
	result := command.GetRunner().Run(GetWorkDir(), cmd)
	testStats, _ := tchn.parseTestReport()
	return TestCommandResult{result, testStats}
}

// AbortExecution asks the toolchain to abort any command currently executing
func (Toolchain) AbortExecution() bool {
	return command.GetRunner().AbortRunningCommand()
}

// BuildCommandPath returns the build command path for this toolchain
func (tchn Toolchain) BuildCommandPath() string {
	return command.FindCompatibleCommand(tchn.buildCommands).Path
}

// BuildCommandArgs returns a table with the list of build command arguments for this toolchain
func (tchn Toolchain) BuildCommandArgs() []string {
	return command.FindCompatibleCommand(tchn.buildCommands).Arguments
}

// BuildCommandLine returns the toolchain's build command line as a string
func (tchn Toolchain) BuildCommandLine() string {
	return command.FindCompatibleCommand(tchn.buildCommands).AsCommandLine()
}

// TestCommandPath returns the test command path for this toolchain
func (tchn Toolchain) TestCommandPath() string {
	return command.FindCompatibleCommand(tchn.testCommands).Path
}

// TestCommandArgs returns a table with the list of test command arguments for this toolchain
func (tchn Toolchain) TestCommandArgs() []string {
	return command.FindCompatibleCommand(tchn.testCommands).Arguments
}

// TestCommandLine returns the toolchain's test command line as a string
func (tchn Toolchain) TestCommandLine() string {
	return command.FindCompatibleCommand(tchn.testCommands).AsCommandLine()
}

func (tchn Toolchain) runsOnPlatform(osName command.OsName, archName command.ArchName) bool {
	return tchn.findBuildCommandFor(osName, archName) != nil && tchn.findTestCommandFor(osName, archName) != nil
}

func (tchn Toolchain) findBuildCommandFor(osName command.OsName, archName command.ArchName) *command.Command {
	return command.FindCommand(tchn.buildCommands, osName, archName)
}

func (tchn Toolchain) findTestCommandFor(osName command.OsName, archName command.ArchName) *command.Command {
	return command.FindCommand(tchn.testCommands, osName, archName)
}

// CheckCommandAccess verifies if the provided command path can be accessed. Returns the path as
// an absolute command path if found. Returns an empty path otherwise, together with the corresponding error
func (Toolchain) CheckCommandAccess(cmdPath string) (string, error) {
	return checkCommandPath(cmdPath)
}

func (tchn Toolchain) parseTestReport() (TestStats, error) {
	parser := xunit.NewParser()
	err := parser.ParseDir(tchn.GetTestResultPath())
	if err != nil {
		report.PostWarning(err)
		return TestStats{}, err
	}
	return NewTestStats(
		parser.Stats.Run,
		parser.Stats.Passed,
		parser.Stats.Failed,
		parser.Stats.Skipped,
		parser.Stats.InError,
		parser.Stats.Duration,
	), nil
}

// GetTestResultPath provides the absolute path to the test result directory
func (tchn Toolchain) GetTestResultPath() string {
	return filepath.Join(workDir, tchn.GetTestResultDir())
}

// GetTestResultDir returns the directory where to retrieve test results (in xUnit format)
func (tchn Toolchain) GetTestResultDir() string {
	return tchn.testResultDir
}

func adjustCommandPath(cmdPath string) string {
	// If this is an absolute path, we return it after cleaning it up
	if filepath.IsAbs(cmdPath) {
		return filepath.Clean(cmdPath)
	}
	// If not, we check if it can be a relative path from the work directory.
	// If the file is found, we return it
	pathFromWorkDir := filepath.Join(GetWorkDir(), cmdPath)
	info, err := os.Stat(pathFromWorkDir)
	if err == nil && !info.IsDir() {
		return pathFromWorkDir
	}
	// As a last resort, we assume it's available in the $PATH
	return filepath.Clean(cmdPath)
}

func checkCommandPath(cmdPath string) (string, error) {
	if cmdPath == "" {
		return "", errors.New("command path is empty")
	}
	return exec.LookPath(adjustCommandPath(cmdPath))
}

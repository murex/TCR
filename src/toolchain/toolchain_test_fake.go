//go:build test_helper

/*
Copyright (c) 2023 Murex

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

import "github.com/murex/tcr/toolchain/command"

type commandFunc func() string
type checkCommandFunc func() (string, error)

type (
	// Operation is the name of a toolchain operation
	Operation string
	// Operations is a slice of Operation
	Operations []Operation
)

// List of supported toolchain operations
const (
	BuildOperation Operation = "build"
	TestOperation  Operation = "test"
	Never          Operation = ""
)

func (operations Operations) contains(operation Operation) bool {
	for _, op := range operations {
		if op == operation {
			return true
		}
	}
	return false
}

// FakeToolchain is a toolchain fake allowing to emulate failing operations (build and fail)
// responses, but without actually calling a real command.
type FakeToolchain struct {
	Toolchain
	failingOperations  Operations
	testStats          TestStats
	buildCommandPath   commandFunc
	testCommandPath    commandFunc
	buildCommandLine   commandFunc
	testCommandLine    commandFunc
	checkCommandAccess checkCommandFunc
}

// NewFakeToolchain creates a FakeToolchain instance
func NewFakeToolchain(failingOperations Operations, testStats TestStats) *FakeToolchain {
	return &FakeToolchain{
		Toolchain: Toolchain{
			name:          "fake-toolchain",
			buildCommands: nil,
			testCommands:  nil,
		},
		failingOperations:  failingOperations,
		testStats:          testStats,
		buildCommandPath:   func() string { return "" },
		testCommandPath:    func() string { return "" },
		buildCommandLine:   func() string { return "" },
		testCommandLine:    func() string { return "" },
		checkCommandAccess: func() (string, error) { return "", nil },
	}
}

// CheckCommandAccess verifies if the provided command path can be accessed (faked)
func (ft *FakeToolchain) CheckCommandAccess(_ string) (string, error) {
	return ft.checkCommandAccess()
}

// WithCheckCommandAccess allows to change the behaviour of CheckCommandAccess() method
func (ft *FakeToolchain) WithCheckCommandAccess(f checkCommandFunc) *FakeToolchain {
	ft.checkCommandAccess = f
	return ft
}

// BuildCommandPath returns the build command path for this toolchain (faked)
func (ft *FakeToolchain) BuildCommandPath() string {
	return ft.buildCommandPath()
}

// WithBuildCommandPath allows to change the behaviour of BuildCommandPath() method
func (ft *FakeToolchain) WithBuildCommandPath(f commandFunc) *FakeToolchain {
	ft.buildCommandPath = f
	return ft
}

// TestCommandPath returns the test command path for this toolchain (faked)
func (ft *FakeToolchain) TestCommandPath() string {
	return ft.testCommandPath()
}

// WithTestCommandPath allows to change the behaviour of TestCommandPath() method
func (ft *FakeToolchain) WithTestCommandPath(f commandFunc) *FakeToolchain {
	ft.testCommandPath = f
	return ft
}

// BuildCommandLine returns the toolchain's build command line as a string (faked)
func (ft *FakeToolchain) BuildCommandLine() string {
	return ft.buildCommandLine()
}

// WithBuildCommandLine allows to change the behaviour of BuildCommandLine() method
func (ft *FakeToolchain) WithBuildCommandLine(f commandFunc) *FakeToolchain {
	ft.buildCommandLine = f
	return ft
}

// TestCommandLine returns the toolchain's test command line as a string (faked)
func (ft *FakeToolchain) TestCommandLine() string {
	return ft.testCommandLine()
}

// WithTestCommandLine allows to change the behaviour of TestCommandLine() method
func (ft *FakeToolchain) WithTestCommandLine(f commandFunc) *FakeToolchain {
	ft.testCommandLine = f
	return ft
}

func (*FakeToolchain) checkBuildCommand() error {
	return nil
}

func (*FakeToolchain) checkTestCommand() error {
	return nil
}

// RunBuild returns an error if build is part of failingOperations, nil otherwise.
// This method does not call any real command
func (ft *FakeToolchain) RunBuild() command.Result {
	return ft.fakeOperation(BuildOperation)
}

// RunTests returns an error if test is part of failingOperations, nil otherwise.
// This method does not call any real command
func (ft *FakeToolchain) RunTests() TestCommandResult {
	return TestCommandResult{ft.fakeOperation(TestOperation), ft.testStats}
}

func (ft *FakeToolchain) fakeOperation(operation Operation) (result command.Result) {
	if ft.failingOperations.contains(operation) {
		result = command.Result{
			Status: command.StatusFail,
			Output: "toolchain " + string(operation) + " fake error",
		}
	} else {
		result = command.Result{
			Status: command.StatusPass,
			Output: "",
		}
	}
	return
}

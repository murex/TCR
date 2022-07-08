//go:build test_helper

/*
Copyright (c) 2022 Murex

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

import "errors"

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
	failingOperations Operations
	testStats         TestStats
}

// NewFakeToolchain creates a FakeToolchain instance
func NewFakeToolchain(failingOperations Operations, testStats TestStats) *FakeToolchain {
	return &FakeToolchain{
		Toolchain: Toolchain{
			name:          "fake-toolchain",
			buildCommands: nil,
			testCommands:  nil,
		},
		failingOperations: failingOperations,
		testStats:         testStats,
	}
}

func (tchn FakeToolchain) checkBuildCommand() error {
	return nil
}

func (tchn FakeToolchain) checkTestCommand() error {
	return nil
}

// RunBuild returns an error if build is part of failingOperations, nil otherwise.
// This method does not call any real command
func (tchn FakeToolchain) RunBuild() error {
	return tchn.fakeOperation(BuildOperation)
}

// RunTests returns an error if test is part of failingOperations, nil otherwise.
// This method does not call any real command
func (tchn FakeToolchain) RunTests() (TestStats, error) {
	return tchn.testStats, tchn.fakeOperation(TestOperation)
}

func (tchn FakeToolchain) fakeOperation(operation Operation) (err error) {
	if tchn.failingOperations.contains(operation) {
		err = errors.New("toolchain " + string(operation) + " fake error")
	}
	return
}

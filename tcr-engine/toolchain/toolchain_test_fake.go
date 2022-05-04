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

// FakeToolchain is a toolchain fake with 2 flags failingBuild and failingTest
// which determine how the toolchain's build and test command will respond (but without
// actually calling a real command)
type FakeToolchain struct {
	Toolchain
	failingBuild bool
	failingTest  bool
	testOutput   string
}

// NewFakeToolchain creates a FakeToolchain instance
func NewFakeToolchain(failingBuild, failingTest bool, testOutput string) *FakeToolchain {
	return &FakeToolchain{
		Toolchain: Toolchain{
			name:          "fake-toolchain",
			buildCommands: nil,
			testCommands:  nil,
		},
		failingBuild: failingBuild,
		failingTest:  failingTest,
		testOutput:   testOutput,
	}
}

func (tchn FakeToolchain) checkBuildCommand() error {
	return nil
}

func (tchn FakeToolchain) checkTestCommand() error {
	return nil
}

// RunBuild returns an error if failingBuild is true, nil otherwise. THis method does not
// call any real command
func (tchn FakeToolchain) RunBuild() error {
	return runCommandStub(tchn.failingBuild)
}

// RunTests returns an error if failingTest is true, nil otherwise. THis method does not
// call any real command
func (tchn FakeToolchain) RunTests() (string, error) {
	return tchn.testOutput, runCommandStub(tchn.failingTest)
}

func runCommandStub(shouldFail bool) (err error) {
	if shouldFail {
		err = errors.New("fake toolchain failing command")
	}
	return
}

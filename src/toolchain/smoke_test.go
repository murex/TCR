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
	"github.com/murex/tcr/toolchain/command"
	"github.com/murex/tcr/utils"
	"testing"
)

const (
	toolchainName = "gradle-wrapper"
)

func Test_toolchain_returns_error_when_build_command_is_not_found(t *testing.T) {
	// Command does not exist in testDataRootDir
	assertErrorWhenBuildFails(t, toolchainName, testDataRootDir)
}

func Test_toolchain_returns_error_when_build_command_fails(t *testing.T) {
	tchn, _ := Get(toolchainName)
	// Add a built-in toolchain from a valid one, but with invalid build command arguments.
	// This allows to actually run the build command and get an execution error out of it
	failingName := tchn.GetName() + "-failing-build"
	var failingCommands []command.Command
	for _, cmd := range tchn.GetBuildCommands() {
		failingCommands = append(failingCommands, command.Command{
			Os:        cmd.Os,
			Arch:      cmd.Arch,
			Path:      cmd.Path,
			Arguments: []string{"-invalid-argument"},
		})
	}
	_ = addBuiltIn(New(failingName, failingCommands, tchn.GetTestCommands(), tchn.GetTestResultDir()))
	assertErrorWhenBuildFails(t, failingName, testDataDirJava)
}

func Test_toolchain_returns_ok_when_build_passes(t *testing.T) {
	utils.SlowTestTag(t)
	assertNoErrorWhenBuildPasses(t, toolchainName, testDataDirJava)
}

func Test_toolchain_returns_error_when_test_command_is_not_found(t *testing.T) {
	// Command does not exist in testDataRootDir
	assertErrorWhenTestFails(t, toolchainName, testDataRootDir)
}

func Test_toolchain_returns_error_when_test_command_fails(t *testing.T) {
	// Add a built-in toolchain from a valid one, but with invalid test command arguments.
	// This allows to actually run the test command and get an execution error out of it
	tchn, _ := Get(toolchainName)
	failingName := tchn.GetName() + "-failing-test"
	var failingCommands []command.Command
	for _, cmd := range tchn.GetTestCommands() {
		failingCommands = append(failingCommands, command.Command{
			Os:        cmd.Os,
			Arch:      cmd.Arch,
			Path:      cmd.Path,
			Arguments: []string{"-invalid-argument"},
		})
	}
	_ = addBuiltIn(New(failingName, tchn.GetBuildCommands(), failingCommands, tchn.GetTestResultDir()))
	assertErrorWhenTestFails(t, failingName, testDataDirJava)
}

func Test_toolchain_returns_ok_when_tests_pass(t *testing.T) {
	utils.SlowTestTag(t)
	assertNoErrorWhenTestPasses(t, toolchainName, testDataDirJava)
}

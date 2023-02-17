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

import (
	"github.com/murex/tcr/utils"
	"testing"
)

const (
	gradleToolchainName        = "gradle"
	gradleWrapperToolchainName = "gradle-wrapper"
)

const gradleCommandPath = "gradle"

func Test_gradle_and_gradle_wrapper_are_built_in_toolchains(t *testing.T) {
	assertIsABuiltInToolchain(t, gradleToolchainName)
	assertIsABuiltInToolchain(t, gradleWrapperToolchainName)
}

func Test_gradle_toolchains_are_supported(t *testing.T) {
	assertIsSupported(t, gradleToolchainName)
	assertIsSupported(t, gradleWrapperToolchainName)
}

func Test_gradle_toolchains_are_registered(t *testing.T) {
	assertIsRegistered(t, gradleToolchainName)
	assertIsRegistered(t, gradleWrapperToolchainName)
}

func Test_gradle_toolchain_names_are_case_insensitive(t *testing.T) {
	assertNameIsNotCaseSensitive(t, gradleToolchainName)
	assertNameIsNotCaseSensitive(t, gradleWrapperToolchainName)
}

func Test_gradle_toolchains_initialization(t *testing.T) {
	assertToolchainInitialization(t, gradleToolchainName)
	assertToolchainInitialization(t, gradleWrapperToolchainName)
}

func Test_gradle_toolchains_name(t *testing.T) {
	assertToolchainName(t, gradleToolchainName)
	assertToolchainName(t, gradleWrapperToolchainName)
}

func Test_gradle_toolchains_build_command_path(t *testing.T) {
	assertBuildCommandPath(t, gradleCommandPath, gradleToolchainName)
	assertBuildCommandPath(t, gradleWrapperCommandPath, gradleWrapperToolchainName)
}

func Test_gradle_toolchains_build_command_args(t *testing.T) {
	buildCommandArgs := []string{"build", "testClasses", "-x", "test"}
	assertBuildCommandArgs(t, buildCommandArgs, gradleToolchainName)
	assertBuildCommandArgs(t, buildCommandArgs, gradleWrapperToolchainName)
}

func Test_gradle_wrapper_toolchain_returns_error_when_build_fails(t *testing.T) {
	assertErrorWhenBuildFails(t, gradleWrapperToolchainName, testDataRootDir)
}

func Test_gradle_wrapper_toolchain_returns_ok_when_build_passes(t *testing.T) {
	utils.SlowTestTag(t)
	assertNoErrorWhenBuildPasses(t, gradleWrapperToolchainName, testDataDirJava)
}

func Test_gradle_toolchains_test_command_path(t *testing.T) {
	assertTestCommandPath(t, gradleCommandPath, gradleToolchainName)
	assertTestCommandPath(t, gradleWrapperCommandPath, gradleWrapperToolchainName)
}

func Test_gradle_toolchains_test_command_args(t *testing.T) {
	testCommandArgs := []string{"test"}
	assertTestCommandArgs(t, testCommandArgs, gradleToolchainName)
	assertTestCommandArgs(t, testCommandArgs, gradleWrapperToolchainName)
}

func Test_gradle_wrapper_toolchain_returns_error_when_tests_fail(t *testing.T) {
	assertErrorWhenTestFails(t, gradleWrapperToolchainName, testDataRootDir)
}

func Test_gradle_wrapper_toolchain_returns_ok_when_tests_pass(t *testing.T) {
	utils.SlowTestTag(t)
	assertNoErrorWhenTestPasses(t, gradleWrapperToolchainName, testDataDirJava)
}

func Test_gradle_toolchains_supported_platforms(t *testing.T) {
	assertRunsOnAllOsWithAmd64(t, gradleToolchainName)
	assertRunsOnAllOsWithAmd64(t, gradleWrapperToolchainName)
}

func Test_gradle_toolchains_test_result_dir(t *testing.T) {
	const testResultDir = "build/test-results/test"
	assertTestResultDir(t, testResultDir, gradleToolchainName)
	assertTestResultDir(t, testResultDir, gradleWrapperToolchainName)
}

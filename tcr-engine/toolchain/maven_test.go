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
	"testing"
)

const (
	mavenToolchainName        = "maven"
	mavenWrapperToolchainName = "maven-wrapper"
)

const mavenCommandPath = "mvn"

func Test_maven_and_maven_wrapper_are_built_in_toolchains(t *testing.T) {
	assertIsABuiltInToolchain(t, mavenToolchainName)
	assertIsABuiltInToolchain(t, mavenWrapperToolchainName)
}

func Test_maven_toolchains_are_supported(t *testing.T) {
	assertIsSupported(t, mavenToolchainName)
	assertIsSupported(t, mavenWrapperToolchainName)
}

func Test_maven_toolchains_are_registered(t *testing.T) {
	assertIsRegistered(t, mavenToolchainName)
	assertIsRegistered(t, mavenWrapperToolchainName)
}

func Test_maven_toolchain_names_are_case_insensitive(t *testing.T) {
	assertNameIsNotCaseSensitive(t, mavenToolchainName)
	assertNameIsNotCaseSensitive(t, mavenWrapperToolchainName)
}

func Test_maven_toolchains_initialization(t *testing.T) {
	assertToolchainInitialization(t, mavenToolchainName)
	assertToolchainInitialization(t, mavenWrapperToolchainName)
}

func Test_maven_toolchains_name(t *testing.T) {
	assertToolchainName(t, mavenToolchainName)
	assertToolchainName(t, mavenWrapperToolchainName)
}

func Test_maven_toolchains_build_command_path(t *testing.T) {
	assertBuildCommandPath(t, mavenCommandPath, mavenToolchainName)
	assertBuildCommandPath(t, mavenWrapperCommandPath, mavenWrapperToolchainName)
}

func Test_maven_toolchains_build_command_args(t *testing.T) {
	buildCommandArgs := []string{"test-compile"}
	assertBuildCommandArgs(t, buildCommandArgs, mavenToolchainName)
	assertBuildCommandArgs(t, buildCommandArgs, mavenWrapperToolchainName)
}

func Test_maven_wrapper_toolchain_returns_error_when_build_fails(t *testing.T) {
	assertErrorWhenBuildFails(t, mavenWrapperToolchainName, testDataRootDir)
}

func Test_maven_wrapper_toolchain_returns_ok_when_build_passes(t *testing.T) {
	assertNoErrorWhenBuildPasses(t, mavenWrapperToolchainName, testDataDirJava)
}

func Test_maven_toolchains_test_command_path(t *testing.T) {
	assertTestCommandPath(t, mavenCommandPath, mavenToolchainName)
	assertTestCommandPath(t, mavenWrapperCommandPath, mavenWrapperToolchainName)
}

func Test_maven_toolchains_test_command_args(t *testing.T) {
	testCommandArgs := []string{"test"}
	assertTestCommandArgs(t, testCommandArgs, mavenToolchainName)
	assertTestCommandArgs(t, testCommandArgs, mavenWrapperToolchainName)
}

func Test_maven_wrapper_toolchain_returns_error_when_tests_fail(t *testing.T) {
	assertErrorWhenTestFails(t, mavenWrapperToolchainName, testDataRootDir)
}

func Test_maven_wrapper_toolchain_returns_ok_when_tests_pass(t *testing.T) {
	assertNoErrorWhenTestPasses(t, mavenWrapperToolchainName, testDataDirJava)
}

func Test_maven_toolchains_supported_platforms(t *testing.T) {
	assertRunsOnAllOsWithAmd64(t, mavenToolchainName)
	assertRunsOnAllOsWithAmd64(t, mavenWrapperToolchainName)
}

func Test_maven_toolchains_test_result_dir(t *testing.T) {
	const testResultDir = "target/surefire-reports"
	assertTestResultDir(t, testResultDir, mavenToolchainName)
	assertTestResultDir(t, testResultDir, mavenWrapperToolchainName)
}

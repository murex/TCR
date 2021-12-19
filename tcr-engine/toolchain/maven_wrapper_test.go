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
	"testing"
)

const (
	mavenWrapperToolchainName = "maven-wrapper"
)

func Test_maven_wrapper_is_a_built_in_toolchain(t *testing.T) {
	assertIsABuiltInToolchain(t, mavenWrapperToolchainName)
}

func Test_maven_wrapper_toolchain_is_supported(t *testing.T) {
	assertIsSupported(t, mavenWrapperToolchainName)
}

func Test_maven_wrapper_toolchain_name_is_case_insensitive(t *testing.T) {
	assertNameIsNotCaseSensitive(t, mavenWrapperToolchainName)
}

func Test_maven_wrapper_toolchain_initialization(t *testing.T) {
	assertToolchainInitialization(t, mavenWrapperToolchainName)
}

func Test_maven_wrapper_toolchain_name(t *testing.T) {
	assertToolchainName(t, mavenWrapperToolchainName)
}

func Test_maven_wrapper_toolchain_build_command_args(t *testing.T) {
	assertBuildCommandArgs(t, []string{"test-compile"}, mavenWrapperToolchainName)
}

func Test_maven_wrapper_toolchain_returns_error_when_build_fails(t *testing.T) {
	assertErrorWhenBuildFails(t, mavenWrapperToolchainName, testDataRootDir)
}

func Test_maven_wrapper_toolchain_returns_ok_when_build_passes(t *testing.T) {
	assertNoErrorWhenBuildPasses(t, mavenWrapperToolchainName, testDataDirJava)
}

func Test_maven_wrapper_toolchain_test_command_args(t *testing.T) {
	assertTestCommandArgs(t, []string{"test"}, mavenWrapperToolchainName)
}

func Test_maven_wrapper_toolchain_returns_error_when_tests_fail(t *testing.T) {
	assertErrorWhenTestFails(t, mavenWrapperToolchainName, testDataRootDir)
}

func Test_maven_wrapper_toolchain_returns_ok_when_tests_pass(t *testing.T) {
	assertNoErrorWhenTestPasses(t, mavenWrapperToolchainName, testDataDirJava)
}

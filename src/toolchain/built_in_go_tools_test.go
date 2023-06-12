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
	goToolchainName = "go-tools"
)

func Test_go_tools_is_a_built_in_toolchain(t *testing.T) {
	assertIsABuiltInToolchain(t, goToolchainName)
}

func Test_go_tools_toolchain_is_supported(t *testing.T) {
	assertIsSupported(t, goToolchainName)
}

func Test_go_tools_toolchain_is_registered(t *testing.T) {
	assertIsRegistered(t, goToolchainName)
}

func Test_go_tools_toolchain_name_is_case_insensitive(t *testing.T) {
	assertNameIsNotCaseSensitive(t, goToolchainName)
}

func Test_go_tools_toolchain_initialization(t *testing.T) {
	assertToolchainInitialization(t, goToolchainName)
}

func Test_go_tools_toolchain_name(t *testing.T) {
	assertToolchainName(t, goToolchainName)
}

func Test_go_tools_toolchain_build_command_path(t *testing.T) {
	assertBuildCommandPath(t, goToolchainName, "go")
}

func Test_go_tools_toolchain_build_command_args(t *testing.T) {
	assertBuildCommandArgs(t, goToolchainName, []string{"test", "-count=0", "./..."})
}

func Test_go_tools_toolchain_test_command_path(t *testing.T) {
	assertTestCommandPath(t, goToolchainName, "go")
}

func Test_go_tools_toolchain_test_command_args(t *testing.T) {
	assertTestCommandArgs(t, goToolchainName, []string{"test", "-short", "./..."})
}

func Test_go_tools_toolchain_supported_platforms(t *testing.T) {
	assertRunsOnAllOsWithAmd64(t, goToolchainName)
}

func Test_go_tools_test_result_dir(t *testing.T) {
	assertTestResultDir(t, goToolchainName, ".")
}

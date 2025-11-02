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

import (
	"testing"

	. "github.com/murex/tcr/toolchain/built_in_test_data"
)

func Test_is_a_built_in_toolchain(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertIsABuiltInToolchain(t, builtIn.Name)
		})
	}
}

func Test_built_in_toolchain_is_supported(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertIsSupported(t, builtIn.Name)
		})
	}
}

func Test_built_in_toolchain_is_registered(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertIsRegistered(t, builtIn.Name)
		})
	}
}

func Test_built_in_toolchain_name_is_case_insensitive(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertNameIsNotCaseSensitive(t, builtIn.Name)
		})
	}
}

func Test_built_in_toolchain_initialization(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertToolchainInitialization(t, builtIn.Name)
		})
	}
}

func Test_built_in_toolchain_name(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertToolchainName(t, builtIn.Name)
		})
	}
}

func Test_built_in_toolchain_build_command_path(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertBuildCommandPath(t, builtIn.Name, builtIn.BuildCommandPath)
		})
	}
}

func Test_built_in_toolchain_build_command_args(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertBuildCommandArgs(t, builtIn.Name, builtIn.BuildCommandArgs)
		})
	}
}

func Test_built_in_toolchain_test_command_path(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertTestCommandPath(t, builtIn.Name, builtIn.TestCommandPath)
		})
	}
}

func Test_built_in_toolchain_test_command_args(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertTestCommandArgs(t, builtIn.Name, builtIn.TestCommandArgs)
		})
	}
}

func Test_built_in_toolchain_supported_platforms(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertRunsOnAllOsWithAmd64(t, builtIn.Name)
		})
	}
}

func Test_built_in_test_result_dir(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertTestResultDir(t, builtIn.Name, builtIn.TestResultDir)
		})
	}
}

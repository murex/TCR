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
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const (
	gradleWrapperToolchainName = "gradle-wrapper"
)

func Test_gradle_wrapper_is_a_built_in_toolchain(t *testing.T) {
	assert.True(t, isBuiltIn(gradleWrapperToolchainName))
}

func Test_gradle_wrapper_toolchain_is_supported(t *testing.T) {
	assert.True(t, isSupported(gradleWrapperToolchainName))
}

func Test_gradle_wrapper_toolchain_name_is_case_insensitive(t *testing.T) {
	assert.True(t, isSupported(gradleWrapperToolchainName))
	assert.True(t, isSupported(strings.ToUpper(gradleWrapperToolchainName)))
	assert.True(t, isSupported(strings.ToLower(gradleWrapperToolchainName)))
	assert.True(t, isSupported(strings.Title(gradleWrapperToolchainName)))
}

func Test_gradle_wrapper_toolchain_initialization(t *testing.T) {
	toolchain, err := Get(gradleWrapperToolchainName)
	assert.NoError(t, err)
	assert.Equal(t, gradleWrapperToolchainName, toolchain.GetName())
}

func Test_gradle_wrapper_toolchain_name(t *testing.T) {
	toolchain, _ := Get(gradleWrapperToolchainName)
	assert.Equal(t, gradleWrapperToolchainName, toolchain.GetName())
}

func Test_gradle_wrapper_toolchain_build_command_args(t *testing.T) {
	toolchain, _ := Get(gradleWrapperToolchainName)
	assert.Equal(t, []string{"build", "-x", "test"}, toolchain.buildCommandArgs())
}

func Test_gradle_wrapper_toolchain_returns_error_when_build_fails(t *testing.T) {
	toolchain, _ := Get(gradleWrapperToolchainName)
	runFromDir(t, testDataRootDir,
		func(t *testing.T) {
			assert.Error(t, toolchain.RunBuild())
		})
}

func Test_gradle_wrapper_toolchain_returns_ok_when_build_passes(t *testing.T) {
	toolchain, _ := Get(gradleWrapperToolchainName)
	runFromDir(t, testDataDirJava,
		func(t *testing.T) {
			assert.NoError(t, toolchain.RunBuild())
		})
}

func Test_gradle_wrapper_toolchain_test_command_args(t *testing.T) {
	toolchain, _ := Get(gradleWrapperToolchainName)
	assert.Equal(t, []string{"test"}, toolchain.testCommandArgs())
}

func Test_gradle_wrapper_toolchain_returns_error_when_tests_fail(t *testing.T) {
	toolchain, _ := Get(gradleWrapperToolchainName)
	runFromDir(t, testDataRootDir,
		func(t *testing.T) {
			assert.Error(t, toolchain.RunTests())
		})
}

func Test_gradle_wrapper_toolchain_returns_ok_when_tests_pass(t *testing.T) {
	toolchain, _ := Get(gradleWrapperToolchainName)
	runFromDir(t, testDataDirJava,
		func(t *testing.T) {
			assert.NoError(t, toolchain.RunTests())
		})
}

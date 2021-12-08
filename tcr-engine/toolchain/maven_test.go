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
	"testing"
)

func Test_maven_toolchain_is_supported(t *testing.T) {
	assert.True(t, isSupported("maven"))
	assert.True(t, isSupported("Maven"))
	assert.True(t, isSupported("MAVEN"))
}

func Test_get_maven_toolchain_instance(t *testing.T) {
	toolchain, err := getToolchain("maven")
	assert.Equal(t, MavenToolchain{}, toolchain)
	assert.Zero(t, err)
}

func Test_maven_toolchain_initialization(t *testing.T) {
	tchn, err := New("maven")
	assert.Equal(t, MavenToolchain{}, tchn)
	assert.Zero(t, err)
}

func Test_maven_toolchain_name(t *testing.T) {
	assert.Equal(t, "maven", MavenToolchain{}.Name())
}

func Test_maven_toolchain_build_command_name(t *testing.T) {
	assert.Equal(t, "mvnw", MavenToolchain{}.BuildCommandName())
}

func Test_maven_toolchain_build_command_args(t *testing.T) {
	assert.Equal(t, []string{"test-compile"}, MavenToolchain{}.BuildCommandArgs())
}

func Test_maven_toolchain_returns_error_when_build_fails(t *testing.T) {
	runFromDir(t, testDataRootDir,
		func(t *testing.T) {
			assert.NotZero(t, MavenToolchain{}.RunBuild())
		})
}

func Test_maven_toolchain_returns_ok_when_build_passes(t *testing.T) {
	runFromDir(t, testDataDirJava,
		func(t *testing.T) {
			assert.Zero(t, MavenToolchain{}.RunBuild())
		})
}

func Test_maven_toolchain_test_command_name(t *testing.T) {
	assert.Equal(t, "mvnw", MavenToolchain{}.TestCommandName())
}

func Test_maven_toolchain_test_command_args(t *testing.T) {
	assert.Equal(t, []string{"test"}, MavenToolchain{}.TestCommandArgs())
}

func Test_maven_toolchain_returns_error_when_tests_fail(t *testing.T) {
	runFromDir(t, testDataRootDir,
		func(t *testing.T) {
			assert.NotZero(t, MavenToolchain{}.RunTests())
		})
}

func Test_maven_toolchain_returns_ok_when_tests_pass(t *testing.T) {
	runFromDir(t, testDataDirJava,
		func(t *testing.T) {
			assert.Zero(t, MavenToolchain{}.RunTests())
		})
}

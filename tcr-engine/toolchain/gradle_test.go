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
	"github.com/murex/tcr/tcr-engine/language"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_gradle_toolchain_initialization(t *testing.T) {
	tchn, err := New("gradle", language.Java{})
	assert.Equal(t, GradleToolchain{}, tchn)
	assert.Zero(t, err)
}

func Test_gradle_toolchain_name(t *testing.T) {
	assert.Equal(t, "gradle", GradleToolchain{}.Name())
}

func Test_gradle_toolchain_build_command_name(t *testing.T) {
	assert.Equal(t, "gradlew", GradleToolchain{}.buildCommandName())
}

func Test_gradle_toolchain_build_command_args(t *testing.T) {
	assert.Equal(t, []string{"build", "-x", "test"}, GradleToolchain{}.buildCommandArgs())
}

func Test_gradle_toolchain_returns_error_when_build_fails(t *testing.T) {
	runFromDir(t, testDataRootDir,
		func(t *testing.T) {
			assert.NotZero(t, GradleToolchain{}.RunBuild())
		})
}

func Test_gradle_toolchain_returns_ok_when_build_passes(t *testing.T) {
	runFromDir(t, testLanguageRootDir(language.Java{}),
		func(t *testing.T) {
			assert.Zero(t, GradleToolchain{}.RunBuild())
		})
}

func Test_gradle_toolchain_test_command_name(t *testing.T) {
	assert.Equal(t, "gradlew", GradleToolchain{}.testCommandName())
}

func Test_gradle_toolchain_test_command_args(t *testing.T) {
	assert.Equal(t, []string{"test"}, GradleToolchain{}.testCommandArgs())
}

func Test_gradle_toolchain_returns_error_when_tests_fail(t *testing.T) {
	runFromDir(t, testDataRootDir,
		func(t *testing.T) {
			assert.NotZero(t, GradleToolchain{}.RunTests())
		})
}

func Test_gradle_toolchain_returns_ok_when_tests_pass(t *testing.T) {
	runFromDir(t, testLanguageRootDir(language.Java{}),
		func(t *testing.T) {
			assert.Zero(t, GradleToolchain{}.RunTests())
		})
}

func Test_gradle_toolchain_supports_java(t *testing.T) {
	assert.True(t, GradleToolchain{}.supports(language.Java{}))
}

func Test_gradle_toolchain_does_not_support_cpp(t *testing.T) {
	assert.False(t, GradleToolchain{}.supports(language.Cpp{}))
}

func Test_gradle_toolchain_language_compatibility(t *testing.T) {
	var comp bool
	var err error

	comp, err = verifyCompatibility(GradleToolchain{}, language.Java{})
	assert.True(t, comp)
	assert.Zero(t, err)

	comp, err = verifyCompatibility(GradleToolchain{}, language.Cpp{})
	assert.False(t, comp)
	assert.NotZero(t, err)
}

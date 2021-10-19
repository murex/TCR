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
	"github.com/murex/tcr-engine/language"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_maven_toolchain_initialization(t *testing.T) {
	tchn, err := New("maven", language.Java{})
	assert.Equal(t, MavenToolchain{}, tchn)
	assert.Zero(t, err)
}

func Test_maven_toolchain_name(t *testing.T) {
	assert.Equal(t, "maven", MavenToolchain{}.Name())
}

func Test_maven_toolchain_build_command_name(t *testing.T) {
	assert.Equal(t, "mvnw", MavenToolchain{}.buildCommandName())
}

func Test_maven_toolchain_build_command_args(t *testing.T) {
	assert.Equal(t, []string{"test-compile"}, MavenToolchain{}.buildCommandArgs())
}

func Test_maven_toolchain_returns_error_when_build_fails(t *testing.T) {
	runFromDir(t, testDataRootDir,
		func(t *testing.T) {
			assert.NotZero(t, MavenToolchain{}.RunBuild())
		})
}

func Test_maven_toolchain_returns_ok_when_build_passes(t *testing.T) {
	runFromDir(t, testLanguageRootDir(language.Java{}),
		func(t *testing.T) {
			assert.Zero(t, MavenToolchain{}.RunBuild())
		})
}

func Test_maven_toolchain_test_command_name(t *testing.T) {
	assert.Equal(t, "mvnw", MavenToolchain{}.testCommandName())
}

func Test_maven_toolchain_test_command_args(t *testing.T) {
	assert.Equal(t, []string{"test"}, MavenToolchain{}.testCommandArgs())
}

func Test_maven_toolchain_returns_error_when_tests_fail(t *testing.T) {
	runFromDir(t, testDataRootDir,
		func(t *testing.T) {
			assert.NotZero(t, MavenToolchain{}.RunTests())
		})
}

func Test_maven_toolchain_returns_ok_when_tests_pass(t *testing.T) {
	runFromDir(t, testLanguageRootDir(language.Java{}),
		func(t *testing.T) {
			assert.Zero(t, MavenToolchain{}.RunTests())
		})
}

func Test_maven_toolchain_supports_java(t *testing.T) {
	assert.True(t, MavenToolchain{}.supports(language.Java{}))
}

func Test_maven_toolchain_does_not_support_cpp(t *testing.T) {
	assert.False(t, MavenToolchain{}.supports(language.Cpp{}))
}

func Test_maven_toolchain_language_compatibility(t *testing.T) {
	var comp bool
	var err error

	comp, err = verifyCompatibility(MavenToolchain{}, language.Java{})
	assert.True(t, comp)
	assert.Zero(t, err)

	comp, err = verifyCompatibility(MavenToolchain{}, language.Cpp{})
	assert.False(t, comp)
	assert.NotZero(t, err)
}

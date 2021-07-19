package toolchain

import (
	"github.com/mengdaming/tcr/tcr/language"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_gradle_toolchain_initialization(t *testing.T) {
	assert.Equal(t, GradleToolchain{}, NewToolchain("gradle", language.Java{}))
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
	runFromDir(t, testKataRootDir,
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
	runFromDir(t, testKataRootDir,
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
	assert.True(t, verifyCompatibility(GradleToolchain{}, language.Java{}))
	assert.False(t, verifyCompatibility(GradleToolchain{}, language.Cpp{}))
}

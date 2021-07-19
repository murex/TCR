package toolchain

import (
	"github.com/mengdaming/tcr/tcr/language"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_maven_toolchain_initialization(t *testing.T) {
	assert.Equal(t, MavenToolchain{}, NewToolchain("maven", language.Java{}))
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
	runFromDir(t, testKataRootDir,
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
	runFromDir(t, testKataRootDir,
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
	assert.True(t, verifyCompatibility(MavenToolchain{}, language.Java{}))
	assert.False(t, verifyCompatibility(MavenToolchain{}, language.Cpp{}))
}

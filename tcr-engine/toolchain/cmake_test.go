package toolchain

import (
	"github.com/murex/tcr-engine/language"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_cmake_toolchain_initialization(t *testing.T) {
	tchn, err := New("cmake", language.Cpp{})
	assert.Equal(t, CmakeToolchain{}, tchn)
	assert.Zero(t, err)
}

func Test_cmake_toolchain_name(t *testing.T) {
	assert.Equal(t, "cmake", CmakeToolchain{}.Name())
}

func Test_cmake_toolchain_build_command_args(t *testing.T) {
	assert.Equal(t, []string{
		"--build", "build",
		"--config", "Debug",
	}, CmakeToolchain{}.buildCommandArgs())
}

func Test_cmake_toolchain_returns_error_when_build_fails(t *testing.T) {
	runFromDir(t, testDataRootDir,
		func(t *testing.T) {
			assert.NotZero(t, CmakeToolchain{}.RunBuild())
		})
}

// TODO Figure out a way to provide a cmake wrapper
//func test_cmake_toolchain_returns_ok_when_build_passes(t *testing.T) {
//	runFromDir(t, testLanguageRootDir(language.Cpp{}),
//		func(t *testing.T) {
//			assert.Zero(t, CmakeToolchain{}.RunBuild())
//		})
//}

func Test_cmake_toolchain_test_command_args(t *testing.T) {
	assert.Equal(t, []string{
		"--output-on-failure",
		"--test-dir", "build",
		"--build-config", "Debug",
	}, CmakeToolchain{}.testCommandArgs())
}

func Test_cmake_toolchain_returns_error_when_tests_fail(t *testing.T) {
	runFromDir(t, testDataRootDir,
		func(t *testing.T) {
			assert.NotZero(t, CmakeToolchain{}.RunTests())
		})
}

// TODO Figure out a way to provide a cmake wrapper
//func Test_cmake_toolchain_returns_ok_when_tests_pass(t *testing.T) {
//	runFromDir(t, testLanguageRootDir(language.Cpp{}),
//		func(t *testing.T) {
//			assert.Zero(t, CmakeToolchain{}.RunTests())
//		})
//}

func Test_cmake_toolchain_supports_cpp(t *testing.T) {
	assert.True(t, CmakeToolchain{}.supports(language.Cpp{}))
}

func Test_cmake_toolchain_does_not_support_java(t *testing.T) {
	assert.False(t, CmakeToolchain{}.supports(language.Java{}))
}

func Test_cmake_toolchain_language_compatibility(t *testing.T) {
	var comp bool
	var err error

	comp, err = verifyCompatibility(CmakeToolchain{}, language.Cpp{})
	assert.True(t, comp)
	assert.Zero(t, err)

	comp, err = verifyCompatibility(CmakeToolchain{}, language.Java{})
	assert.False(t, comp)
	assert.NotZero(t, err)
}

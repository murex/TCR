package toolchain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_cmake_toolchain_build_command_name(t *testing.T) {
	expected := "build/cmake/cmake-macos-universal/CMake.app/Contents/bin/cmake"
	assert.Equal(t, expected, CmakeToolchain{}.buildCommandName())
}

func Test_cmake_toolchain_test_command_name(t *testing.T) {
	expected := "build/cmake/cmake-macos-universal/CMake.app/Contents/bin/ctest"
	assert.Equal(t, expected, CmakeToolchain{}.testCommandName())
}

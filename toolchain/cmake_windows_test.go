package toolchain

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_cmake_toolchain_build_command_name(t *testing.T) {
	expected := filepath.Join("build", "cmake", "cmake-windows-x86_64", "bin", "cmake.exe")
	assert.Equal(t, expected, CmakeToolchain{}.buildCommandName())
}

func Test_cmake_toolchain_test_command_name(t *testing.T) {
	expected := filepath.Join("build", "cmake", "cmake-windows-x86_64", "bin", "ctest.exe")
	assert.Equal(t, expected, CmakeToolchain{}.testCommandName())
}

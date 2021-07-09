package tcr

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_cmake_toolchain_build_command_name(t *testing.T) {
	expected := filepath.Join(".", "cmake", "cmake-win64-x64", "bin", "cmake.exe")
	assert.Equal(t, expected, CmakeToolchain{}.buildCommandName())
}

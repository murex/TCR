package toolchain

import (
	"path/filepath"
)

func (tchn CmakeToolchain) buildCommandName() string {
	return filepath.Join( "build", "cmake", "cmake-windows-x86_64", "bin", "cmake.exe")
}

func (tchn CmakeToolchain) testCommandName() string {
	return filepath.Join( "build", "cmake", "cmake-windows-x86_64", "bin", "ctest.exe")
}
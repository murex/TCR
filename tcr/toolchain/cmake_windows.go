package toolchain

import "path/filepath"

func (tchn CmakeToolchain) buildCommandName() string {
	return filepath.Join( "build", "cmake", "cmake-win64-x64", "bin", "cmake.exe")
}

func (tchn CmakeToolchain) testCommandName() string {
	return filepath.Join( "build", "cmake", "cmake-win64-x64", "bin", "ctest.exe")
}
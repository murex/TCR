package tcr

import "path/filepath"

func (toolchain CmakeToolchain) buildCommandName() string {
	return filepath.Join( "build", "cmake", "cmake-win64-x64", "bin", "cmake.exe")
}

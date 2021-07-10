package tcr

import "path/filepath"

func (toolchain CmakeToolchain) buildCommandName() string {
	return filepath.Join("cmake", "cmake-macos-universal",
		"CMake.app", "Contents", "bin", "cmake")
}

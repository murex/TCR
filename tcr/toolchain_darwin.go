package tcr

func (toolchain CmakeToolchain) buildCommandName() string {
	return "build/cmake/cmake-macos-universal/CMake.app/Contents/bin/cmake"
}

func (toolchain CmakeToolchain) testCommandName() string {
	return "build/cmake/cmake-macos-universal/CMake.app/Contents/bin/ctest"
}

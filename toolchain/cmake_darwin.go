package toolchain

func (tchn CmakeToolchain) buildCommandName() string {
	return "build/cmake/cmake-macos-universal/CMake.app/Contents/bin/cmake"
}

func (tchn CmakeToolchain) testCommandName() string {
	return "build/cmake/cmake-macos-universal/CMake.app/Contents/bin/ctest"
}

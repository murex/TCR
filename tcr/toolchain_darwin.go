package tcr

func (toolchain CmakeToolchain) buildCommandName() string {
	return "./cmake/cmake-macos-universal/CMake.app/Contents/bin/cmake"
}

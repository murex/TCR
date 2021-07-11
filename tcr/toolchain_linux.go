package tcr

func (toolchain CmakeToolchain) buildCommandName() string {
	return "build/cmake/cmake-Linux-x86_64/bin/cmake"
}

func (toolchain CmakeToolchain) testCommandName() string {
	return "build/cmake/cmake-Linux-x86_64/bin/ctest"
}

package toolchain

func (tchn CmakeToolchain) buildCommandName() string {
	return "build/cmake/cmake-Linux-x86_64/bin/cmake"
}

func (tchn CmakeToolchain) testCommandName() string {
	return "build/cmake/cmake-Linux-x86_64/bin/ctest"
}

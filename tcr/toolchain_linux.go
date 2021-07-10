package tcr

func (toolchain CmakeToolchain) buildCommandName() string {
	return filepath.Join("cmake", "cmake-Linux-x86_64", "bin", "cmake")
}

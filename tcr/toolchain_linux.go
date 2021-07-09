package tcr

func (toolchain CmakeToolchain) buildCommandName() string {
	// TODO OS-Specific command path
	return "cmake"
}


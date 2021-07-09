package tcr

type Toolchain interface {
	name() string
	runBuild() error
	runTests() error
}

// ========================================================================

type GradleToolchain struct {
}

func (toolchain GradleToolchain) name() string {
	return "gradle"
}

func (toolchain GradleToolchain) runBuild() error {
	// TODO
	return nil
}

func (toolchain GradleToolchain) runTests() error {
	// TODO
	return nil
}

func (toolchain GradleToolchain) buildCommandName() string {
	return "./gradlew"
}

func (toolchain GradleToolchain) buildCommandArgs() []string {
	return []string{"-x", "test"}
}

// ========================================================================

type CmakeToolchain struct{}

func (toolchain CmakeToolchain) name() string {
	return "cmake"
}

func (toolchain CmakeToolchain) runBuild() error {
	// TODO
	return nil
}

func (toolchain CmakeToolchain) runTests() error {
	// TODO
	return nil
}

package tcr

import (
	"github.com/mengdaming/tcr/trace"
)

func NewToolchain(name string) Toolchain {
	switch name {
	case GradleToolchain{}.name():
		return GradleToolchain{}
	case CmakeToolchain{}.name():
		return CmakeToolchain{}
	default:
		// TODO check toolchain / language compatibility
		trace.Error("Toolchain \"", name, "\" not supported")
	}
	return nil
}

// ========================================================================

type Toolchain interface {
	name() string
	runBuild() error
	runTests() error
}

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

// ========================================================================

func (toolchain GradleToolchain) buildCommandArgs() []string {
	return []string{"-x", "test"}
}

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

func (toolchain CmakeToolchain) buildCommandArgs() []string {
	return []string{"--build", ".", "--config", "Debug"}
}

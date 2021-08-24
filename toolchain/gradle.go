package toolchain

import (
	"github.com/mengdaming/tcr/language"
)

// GradleToolchain is the toolchain implementation for Maven
type GradleToolchain struct {
}

// Name provides the name of the toolchain
func (tchn GradleToolchain) Name() string {
	return "gradle"
}

// RunBuild runs the build with this toolchain
func (tchn GradleToolchain) RunBuild() error {
	return runBuild(tchn)
}

// RunTests runs the tests with this toolchain
func (tchn GradleToolchain) RunTests() error {
	return runTests(tchn)
}

func (tchn GradleToolchain) buildCommandName() string {
	return "gradlew"
}

func (tchn GradleToolchain) buildCommandArgs() []string {
	return []string{"build", "-x", "test"}
}

func (tchn GradleToolchain) testCommandName() string {
	return "gradlew"
}

func (tchn GradleToolchain) testCommandArgs() []string {
	return []string{"test"}
}

func (tchn GradleToolchain) supports(lang language.Language) bool {
	return lang == language.Java{}
}

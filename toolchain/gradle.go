package toolchain

import (
	"github.com/mengdaming/tcr/language"
)

type GradleToolchain struct {
}

func (tchn GradleToolchain) Name() string {
	return "gradle"
}

func (tchn GradleToolchain) RunBuild() error {
	return runBuild(tchn)
}

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

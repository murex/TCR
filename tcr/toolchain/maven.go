package toolchain

import (
	"github.com/mengdaming/tcr/tcr/language"
)

type MavenToolchain struct {
}

func (tchn MavenToolchain) Name() string {
	return "maven"
}

func (tchn MavenToolchain) RunBuild() error {
	return runBuild(tchn)
}

func (tchn MavenToolchain) RunTests() error {
	return runTests(tchn)
}

func (tchn MavenToolchain) buildCommandName() string {
	return "mvnw"
}

func (tchn MavenToolchain) buildCommandArgs() []string {
	return []string{"test-compile"}
}

func (tchn MavenToolchain) testCommandName() string {
	return "mvnw"
}

func (tchn MavenToolchain) testCommandArgs() []string {
	return []string{"test"}
}

func (tchn MavenToolchain) supports(lang language.Language) bool {
	return lang == language.Java{}
}

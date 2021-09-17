package toolchain

import (
	"github.com/mengdaming/tcr-engine/language"
)

// MavenToolchain is the toolchain implementation for Maven
type MavenToolchain struct {
}

// Name provides the name of the toolchain
func (tchn MavenToolchain) Name() string {
	return "maven"
}

// RunBuild runs the build with this toolchain
func (tchn MavenToolchain) RunBuild() error {
	return runBuild(tchn)
}

// RunTests runs the tests with this toolchain
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

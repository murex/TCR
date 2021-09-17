package toolchain

import (
	"github.com/mengdaming/tcr-engine/language"
)

// CmakeToolchain is the toolchain implementation for CMake
type CmakeToolchain struct{}

// Name provides the name of the toolchain
func (tchn CmakeToolchain) Name() string {
	return "cmake"
}

// RunBuild runs the build with this toolchain
func (tchn CmakeToolchain) RunBuild() error {
	return runBuild(tchn)
}

// RunTests runs the tests with this toolchain
func (tchn CmakeToolchain) RunTests() error {
	return runTests(tchn)
}

func (tchn CmakeToolchain) buildCommandArgs() []string {
	return []string{"--build", "build", "--config", "Debug"}
}

func (tchn CmakeToolchain) testCommandArgs() []string {
	// Important: This (--test-dir option) requires to use cmake 3.20 version or higher
	return []string{"--output-on-failure", "--test-dir", "build", "--build-config", "Debug"}
}

func (tchn CmakeToolchain) supports(lang language.Language) bool {
	return lang == language.Cpp{}
}

package toolchain

import "github.com/mengdaming/tcr/tcr/language"

type CmakeToolchain struct{}

func (tchn CmakeToolchain) Name() string {
	return "cmake"
}

func (tchn CmakeToolchain) RunBuild() error {
	return runBuild(tchn)
}

func (tchn CmakeToolchain) RunTests() error {
	return runTests(tchn)
}

func (tchn CmakeToolchain) buildCommandArgs() []string {
	return []string{"--build", ".", "--config", "Debug"}
}

func (tchn CmakeToolchain) testCommandArgs() []string {
	return []string{"--output-on-failure", "-C", "Debug"}
}

func (tchn CmakeToolchain) supports(lang language.Language) bool {
	return lang == language.Cpp{}
}

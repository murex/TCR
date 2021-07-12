package tcr

import (
	"github.com/codeskyblue/go-sh"
	"github.com/mengdaming/tcr/trace"
	"os"
	"path/filepath"
)

type Toolchain interface {
	name() string
	runBuild() error
	runTests() error
	buildCommandName() string
	buildCommandArgs() []string
	testCommandName() string
	testCommandArgs() []string
	supports(language Language) bool
}

func NewToolchain(name string, language Language) Toolchain {
	switch name {
	case GradleToolchain{}.name():
		return GradleToolchain{}
	case MavenToolchain{}.name():
		return MavenToolchain{}
	case CmakeToolchain{}.name():
		return CmakeToolchain{}
	case "":
		return defaultToolchain(language)
	default:
		trace.Error("Toolchain \"", name, "\" not supported")
	}
	return nil
}

func defaultToolchain(language Language) Toolchain {
	switch language {
	case JavaLanguage{}:
		return GradleToolchain{}
	case CppLanguage{}:
		return CmakeToolchain{}
	default:
		trace.Error("No supported toolchain for language ", language.name())
	}
	return nil
}

func checkToolchainAndLanguageCompatibility(toolchain Toolchain, language Language) {
	if !toolchain.supports(language) {
		trace.Error("Toolchain ", toolchain.name(),
			" does not support language ", language.name())
	}
}

func runBuild(toolchain Toolchain) error {
	wd, _ := os.Getwd()
	buildCommandPath := filepath.Join(wd, toolchain.buildCommandName())
	//trace.Info(buildCommandPath)
	output, err := sh.Command(
		buildCommandPath,
		toolchain.buildCommandArgs()).Output()
	if output != nil {
		trace.Transparent(string(output))
	}
	if err != nil {
		trace.Warning(err)
	}
	return err
}

// Gradle ========================================================================

func runTests(toolchain Toolchain) error {
	wd, _ := os.Getwd()
	testCommandPath := filepath.Join(wd, toolchain.testCommandName())
	//trace.Info(testCommandPath)
	output, err := sh.Command(
		testCommandPath,
		toolchain.testCommandArgs()).Output()
	if output != nil {
		trace.Transparent(string(output))
	}
	if err != nil {
		trace.Warning(err)
	}
	return err
}

type GradleToolchain struct {
}

func (toolchain GradleToolchain) name() string {
	return "gradle"
}

func (toolchain GradleToolchain) runBuild() error {
	return runBuild(toolchain)
}

func (toolchain GradleToolchain) runTests() error {
	return runTests(toolchain)
}

func (toolchain GradleToolchain) buildCommandName() string {
	return "gradlew"
}

func (toolchain GradleToolchain) buildCommandArgs() []string {
	return []string{"build", "-x", "test"}
}

func (toolchain GradleToolchain) testCommandName() string {
	return "gradlew"
}

func (toolchain GradleToolchain) testCommandArgs() []string {
	return []string{"test"}
}

// Cmake ========================================================================

func (toolchain GradleToolchain) supports(language Language) bool {
	return language == JavaLanguage{}
}

type CmakeToolchain struct{}

func (toolchain CmakeToolchain) name() string {
	return "cmake"
}

func (toolchain CmakeToolchain) runBuild() error {
	return runBuild(toolchain)
}

func (toolchain CmakeToolchain) runTests() error {
	return runTests(toolchain)
}

func (toolchain CmakeToolchain) buildCommandArgs() []string {
	return []string{"--build", ".", "--config", "Debug"}
}

func (toolchain CmakeToolchain) testCommandArgs() []string {
	return []string{"--output-on-failure", "-C", "Debug"}
}

// Maven ========================================================================

func (toolchain CmakeToolchain) supports(language Language) bool {
	return language == CppLanguage{}
}

type MavenToolchain struct {
}

func (toolchain MavenToolchain) name() string {
	return "maven"
}

func (toolchain MavenToolchain) runBuild() error {
	return runBuild(toolchain)
}

func (toolchain MavenToolchain) runTests() error {
	return runTests(toolchain)
}

func (toolchain MavenToolchain) buildCommandName() string {
	return "mvnw"
}

func (toolchain MavenToolchain) buildCommandArgs() []string {
	return []string{"test-compile"}
}

func (toolchain MavenToolchain) testCommandName() string {
	return "mvnw"
}

func (toolchain MavenToolchain) testCommandArgs() []string {
	return []string{"test"}
}

func (toolchain MavenToolchain) supports(language Language) bool {
	return language == JavaLanguage{}
}


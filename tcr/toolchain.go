package tcr

import (
	"github.com/codeskyblue/go-sh"
	"github.com/mengdaming/tcr/tcr/language"
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
	supports(lang language.Language) bool
}

func NewToolchain(name string, lang language.Language) Toolchain {
	var toolchain Toolchain = nil
	switch name {
	case GradleToolchain{}.name():
		toolchain = GradleToolchain{}
	case MavenToolchain{}.name():
		toolchain = MavenToolchain{}
	case CmakeToolchain{}.name():
		toolchain = CmakeToolchain{}
	case "":
		toolchain = defaultToolchain(lang)
	default:
		trace.Error("Toolchain \"", name, "\" not supported")
		return nil
	}

	if !verifyCompatibility(toolchain, lang) {
		return nil
	}
	return toolchain
}

func defaultToolchain(lang language.Language) Toolchain {
	switch lang {
	case language.Java{}:
		return GradleToolchain{}
	case language.Cpp{}:
		return CmakeToolchain{}
	default:
		trace.Error("No supported toolchain for language ", lang.Name())
	}
	return nil
}

func verifyCompatibility(toolchain Toolchain, lang language.Language) bool {
	if toolchain == nil || lang == nil {
		return false
	}
	if !toolchain.supports(lang) {
		trace.Error("Toolchain ", toolchain.name(),
			" does not support language ", lang.Name())
		return false
	}
	return true
}

func runBuild(toolchain Toolchain) error {
	wd, _ := os.Getwd()
	buildCommandPath := filepath.Join(wd, toolchain.buildCommandName())
	//trace.Info(buildCommandPath)
	output, err := sh.Command(
		buildCommandPath,
		toolchain.buildCommandArgs()).Output()
	if output != nil {
		trace.Echo(string(output))
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
		trace.Echo(string(output))
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

func (toolchain GradleToolchain) supports(lang language.Language) bool {
	return lang == language.Java{}
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

func (toolchain CmakeToolchain) supports(lang language.Language) bool {
	return lang == language.Cpp{}
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

func (toolchain MavenToolchain) supports(lang language.Language) bool {
	return lang == language.Java{}
}

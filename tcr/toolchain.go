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
}

func NewToolchain(name string) Toolchain {
	switch name {
	case GradleToolchain{}.name():
		return GradleToolchain{}
	case MavenToolchain{}.name():
		return MavenToolchain{}
	case CmakeToolchain{}.name():
		return CmakeToolchain{}
	default:
		// TODO check toolchain / language compatibility
		trace.Error("Toolchain \"", name, "\" not supported")
	}
	return nil
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

// Gradle ========================================================================

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

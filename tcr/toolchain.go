package tcr

import (
	"github.com/codeskyblue/go-sh"
	"github.com/mengdaming/tcr/trace"
	"time"
)

type Toolchain interface {
	name() string
	runBuild() error
	runTests() error
	buildCommandName() string
	buildCommandArgs() []string
}

func NewToolchain(name string) Toolchain {
	// TODO add maven
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

func runBuild(toolchain Toolchain) error {
	output, err := sh.Command(
		toolchain.buildCommandName(),
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

type GradleToolchain struct {
}

func (toolchain GradleToolchain) name() string {
	return "gradle"
}

func (toolchain GradleToolchain) runBuild() error {
	return runBuild(toolchain)
}

func (toolchain GradleToolchain) runTests() error {
	// TODO Call gradle test
	time.Sleep(1 * time.Second)
	return nil
}

func (toolchain GradleToolchain) buildCommandName() string {
	return "./gradlew"
}


func (toolchain GradleToolchain) buildCommandArgs() []string {
	return []string{"build", "-x", "test"}
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
	// TODO Call cmake test
	time.Sleep(1 * time.Second)
	return nil
}

func (toolchain CmakeToolchain) buildCommandArgs() []string {
	return []string{"--build", ".", "--config", "Debug"}
}

// Maven ========================================================================

// TODO Add maven toolchain
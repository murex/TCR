package toolchain

import (
	"github.com/mengdaming/tcr/tcr/language"
	"github.com/mengdaming/tcr/trace"

	"github.com/codeskyblue/go-sh"

	"os"
	"path/filepath"
)

type Toolchain interface {
	Name() string
	RunBuild() error
	RunTests() error
	buildCommandName() string
	buildCommandArgs() []string
	testCommandName() string
	testCommandArgs() []string
	supports(lang language.Language) bool
}

func NewToolchain(name string, lang language.Language) Toolchain {
	var toolchain Toolchain = nil
	switch name {
	case GradleToolchain{}.Name():
		toolchain = GradleToolchain{}
	case MavenToolchain{}.Name():
		toolchain = MavenToolchain{}
	case CmakeToolchain{}.Name():
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
		trace.Error("Toolchain ", toolchain.Name(),
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

func runTests(tchn Toolchain) error {
	wd, _ := os.Getwd()
	testCommandPath := filepath.Join(wd, tchn.testCommandName())
	//trace.Info(testCommandPath)
	output, err := sh.Command(
		testCommandPath,
		tchn.testCommandArgs()).Output()
	if output != nil {
		trace.Echo(string(output))
	}
	return err
}

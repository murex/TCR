/*
Copyright (c) 2021 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package toolchain

import (
	"errors"
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/murex/tcr/tcr-engine/language"
	"github.com/murex/tcr/tcr-engine/report"
	"os"
	"path/filepath"
	"strings"
)

// Toolchain provides the interface that any toolchain needs to implement so that it can be used by TcR
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

var (
	supportedToolchains = make(map[string]Toolchain)
	defaultToolchains   = make(map[language.Language]Toolchain)
)

func init() {
	// List of supported toolchains
	addSupportedToolchain(GradleToolchain{})
	addSupportedToolchain(MavenToolchain{})
	addSupportedToolchain(CmakeToolchain{})

	// Default toolchain for each language
	setDefaultToolchain(language.Java{}, GradleToolchain{})
	setDefaultToolchain(language.Cpp{}, CmakeToolchain{})
}

func addSupportedToolchain(tchn Toolchain) {
	supportedToolchains[strings.ToLower(tchn.Name())] = tchn
}

func isSupported(name string) bool {
	_, found := supportedToolchains[strings.ToLower(name)]
	return found
}

func getToolchain(name string) (Toolchain, error) {
	tchn, found := supportedToolchains[strings.ToLower(name)]
	if found {
		return tchn, nil
	}
	return nil, errors.New(fmt.Sprint("Toolchain not supported: ", name))
}

func setDefaultToolchain(lang language.Language, tchn Toolchain) {
	defaultToolchains[lang] = tchn
}

func getDefaultToolchain(lang language.Language) (Toolchain, error) {
	if lang == nil {
		return nil, errors.New(fmt.Sprint("Language is not defined"))
	}
	tchn, found := defaultToolchains[lang]
	if found {
		return tchn, nil
	}
	return nil, errors.New(fmt.Sprint("No supported toolchain for ", lang.Name(), " language"))
}

// New creates a new toolchain instance with the provided name and for the provided language.
// When toolchain name is not provided (e.g. empty string), we fallback on the default toolchain
// for the provided language.
// The toolchain name is case insensitive.
// This function also verifies that toolchain and language are compatible with each other.
func New(name string, lang language.Language) (Toolchain, error) {
	var toolchain Toolchain
	var err error

	// We first retrieve the toolchain
	if name != "" {
		toolchain, err = getToolchain(name)
	} else {
		toolchain, err = getDefaultToolchain(lang)
	}
	if err != nil {
		return nil, err
	}

	// Then check language compatibility
	comp, err := verifyCompatibility(toolchain, lang)
	if !comp || err != nil {
		return nil, err
	}
	return toolchain, nil
}

func verifyCompatibility(toolchain Toolchain, lang language.Language) (bool, error) {
	if toolchain == nil || lang == nil {
		return false, errors.New("toolchain and/or language is unknown")
	}
	if !toolchain.supports(lang) {
		return false, fmt.Errorf(
			"%v toolchain does not support %v language",
			toolchain.Name(), lang.Name(),
		)
	}
	return true, nil
}

func runBuild(toolchain Toolchain) error {
	wd, _ := os.Getwd()
	buildCommandPath := filepath.Join(wd, toolchain.buildCommandName())
	output, err := sh.Command(
		buildCommandPath,
		toolchain.buildCommandArgs()).CombinedOutput()
	if output != nil {
		report.PostText(string(output))
		//fmt.Println(string(output))
	}
	return err
}

func runTests(tchn Toolchain) error {
	wd, _ := os.Getwd()
	testCommandPath := filepath.Join(wd, tchn.testCommandName())
	output, err := sh.Command(
		testCommandPath,
		tchn.testCommandArgs()).CombinedOutput()
	if output != nil {
		report.PostText(string(output))
		//fmt.Println(string(output))
	}
	return err
}

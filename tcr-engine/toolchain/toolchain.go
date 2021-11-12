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

// New creates a new toolchain instance with the provided name and for the provided language
func New(name string, lang language.Language) (Toolchain, error) {
	var toolchain Toolchain
	var err error

	switch name {
	case GradleToolchain{}.Name():
		toolchain = GradleToolchain{}
	case MavenToolchain{}.Name():
		toolchain = MavenToolchain{}
	case CmakeToolchain{}.Name():
		toolchain = CmakeToolchain{}
	case "":
		toolchain, err = defaultToolchain(lang)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprint("toolchain \"", name, "\" not supported"))
	}

	comp, err := verifyCompatibility(toolchain, lang)
	if !comp || err != nil {
		return nil, err
	}
	return toolchain, nil
}

func defaultToolchain(lang language.Language) (Toolchain, error) {
	switch lang {
	case language.Java{}:
		return GradleToolchain{}, nil
	case language.Cpp{}:
		return CmakeToolchain{}, nil
	default:
		return nil, errors.New(fmt.Sprint("no supported toolchain for ", lang.Name(), " language"))
	}
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

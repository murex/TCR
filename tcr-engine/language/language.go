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

package language

import (
	"errors"
	"fmt"
	"github.com/murex/tcr/tcr-engine/toolchain"
	"path/filepath"
	"strings"
)

// Language is the interface that any supported language implementation must comply with
// in order to be used by TCR engine
type Language interface {
	Name() string
	SrcDirs() []string
	TestDirs() []string
	IsSrcFile(filename string) bool
	defaultToolchain() toolchain.Toolchain
	worksWithToolchain(t toolchain.Toolchain) bool
	GetToolchain(t string) (toolchain.Toolchain, error)
}

var (
	supportedLanguages = make(map[string]Language)
)

func init() {
	addSupportedLanguage(Java{})
	addSupportedLanguage(Cpp{})
}

func addSupportedLanguage(lang Language) {
	supportedLanguages[strings.ToLower(lang.Name())] = lang
}

func isSupported(name string) bool {
	_, found := supportedLanguages[strings.ToLower(name)]
	return found
}

func getLanguage(name string) (Language, error) {
	language, found := supportedLanguages[strings.ToLower(name)]
	if found {
		return language, nil
	}
	return nil, errors.New(fmt.Sprint("language not supported: ", name))
}

// New returns the language to be used in current session. If no value is provided
// for language (e.g. empty string), we try to detect the language based on the directory name.
// Both name and baseDir are case insensitive
func New(name string, baseDir string) (Language, error) {
	if name != "" {
		return getLanguage(name)
	}
	return detectLanguage(baseDir)

}

// detectLanguage is used to identify the language used in the provided directory. The current implementation
// simply looks at the name of the directory and checks if it matches with one of the supported languages
func detectLanguage(baseDir string) (Language, error) {
	return getLanguage(filepath.Base(baseDir))
}

// DirsToWatch returns the list of directories that TCR engine needs to watch for the provided language
func DirsToWatch(baseDir string, lang Language) []string {
	dirList := append(lang.SrcDirs(), lang.TestDirs()...)
	for i := 0; i < len(dirList); i++ {
		dirList[i] = filepath.Join(baseDir, dirList[i])
	}
	//report.PostInfo(dirList)
	return dirList
}

func getToolchain(lang Language, toolchainName string) (toolchain.Toolchain, error) {
	var tchn toolchain.Toolchain
	var err error

	// We first retrieve the toolchain
	if toolchainName != "" {
		tchn, err = toolchain.New(toolchainName)
		if err != nil {
			return nil, err
		}
	} else {
		// If no toolchain is specified, we use the default toolchain for this language
		tchn = lang.defaultToolchain()
	}

	// Then we check language/toolchain compatibility
	comp, err := verifyCompatibility(lang, tchn)
	if !comp || err != nil {
		return nil, err
	}
	return tchn, nil
}

func verifyCompatibility(lang Language, toolchain toolchain.Toolchain) (bool, error) {
	if toolchain == nil || lang == nil {
		return false, errors.New("toolchain and/or language is unknown")
	}
	if !lang.worksWithToolchain(toolchain) {
		return false, fmt.Errorf(
			"%v toolchain does not support %v language",
			toolchain.Name(), lang.Name(),
		)
	}
	return true, nil
}

//func setDefaultToolchain(lang language.Language, tchn Toolchain) {
//	defaultToolchains[lang] = tchn
//}
//
//func getDefaultToolchain(lang language.Language) (Toolchain, error) {
//	if lang == nil {
//		return nil, errors.New("language is not defined")
//	}
//	tchn, found := defaultToolchains[lang]
//	if found {
//		return tchn, nil
//	}
//	return nil, errors.New(fmt.Sprint("no supported toolchain for ", lang.Name(), " language"))
//}

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
	"sort"
	"strings"
)

type (
	// Toolchains defines the structure for toolchains related to a language
	Toolchains struct {
		Default    string
		Compatible []string
	}

	// Language defines the data structure of a language.
	// - Name is the name of the language, it must be unique in the list of available languages
	Language struct {
		Name       string
		Toolchains Toolchains
		SrcFiles   FileTreeFilter
		TestFiles  FileTreeFilter
		baseDir    string
	}
)

var (
	builtIn    = make(map[string]Language)
	registered = make(map[string]Language)
)

// Register adds the provided language to the list of supported languages
func Register(lang Language) error {
	if err := lang.checkName(); err != nil {
		return err
	}
	if err := lang.checkCompatibleToolchains(); err != nil {
		return err
	}
	if err := lang.checkDefaultToolchain(); err != nil {
		return err
	}
	registered[strings.ToLower(lang.GetName())] = lang
	return nil
}

func isSupported(name string) bool {
	_, found := registered[strings.ToLower(name)]
	return found
}

// Get returns the language instance with the provided name.
// The language name is case-insensitive.
// This method does not guarantee that the returned language instance can
// be used out of the box for file filtering operations as it does not
// enforce that baseDir is set. Prefer GetLanguage in this case.
func Get(name string) (*Language, error) {
	if name == "" {
		return nil, errors.New("language name not provided")
	}
	lang, found := registered[strings.ToLower(name)]
	if found {
		return &lang, nil
	}
	return nil, errors.New(fmt.Sprint("language not supported: ", name))
}

// Names returns the list of available language names sorted alphabetically
func Names() []string {
	var names []string
	for _, lang := range registered {
		names = append(names, lang.Name)
	}
	sort.Strings(names)
	return names
}

// Reset resets the language with the provided name to its default values
func Reset(name string) {
	_, found := registered[strings.ToLower(name)]
	if found && isBuiltIn(name) {
		_ = Register(*getBuiltIn(name))
	}
}

func getBuiltIn(name string) *Language {
	var builtIn, _ = builtIn[strings.ToLower(name)]
	return &builtIn
}

func isBuiltIn(name string) bool {
	_, found := builtIn[strings.ToLower(name)]
	return found
}

func addBuiltIn(lang Language) error {
	if lang.Name == "" {
		return errors.New("language name cannot be an empty string")
	}
	builtIn[strings.ToLower(lang.Name)] = lang
	return Register(lang)
}

// GetName provides the name of the toolchain
func (lang Language) GetName() string {
	return lang.Name
}

// GetLanguage returns the language to be used in current session. If no value is provided
// for language (e.g. empty string), we try to detect the language based on the directory name.
// Both name and baseDir are case-insensitive
func GetLanguage(name string, baseDir string) (lang *Language, err error) {
	if name != "" {
		lang, err = getRegisteredLanguage(name)
	} else {
		lang, err = detectLanguageFromDirName(baseDir)
	}
	lang.setBaseDir(baseDir)
	return lang, err
}

func getRegisteredLanguage(name string) (*Language, error) {
	language, found := registered[strings.ToLower(name)]
	if found {
		return &language, nil
	}
	return nil, errors.New(fmt.Sprint("language not supported: ", name))
}

// detectLanguageFromDirName is used to identify the language used in the provided directory. The current implementation
// simply looks at the name of the directory and checks if it matches with one of the supported languages
func detectLanguageFromDirName(baseDir string) (*Language, error) {
	return getRegisteredLanguage(filepath.Base(baseDir))
}

func (lang Language) checkName() error {
	if lang.Name == "" {
		return errors.New("language name is empty")
	}
	return nil
}

func (lang Language) checkCompatibleToolchains() error {
	if lang.Toolchains.Compatible == nil {
		return errors.New("language has no compatible toolchain defined")
	}
	return nil
}

func (lang Language) checkDefaultToolchain() error {
	if lang.Toolchains.Default == "" {
		return errors.New("language has no default toolchain defined")
	} else if !lang.worksWithToolchain(lang.Toolchains.Default) {
		return errors.New("language's default toolchain " +
			lang.Toolchains.Default + " is not listed in compatible toolchains list")
	}
	return nil
}

// IsSrcFile returns true if the provided filePath is recognized as a source file for this language
func (lang Language) IsSrcFile(filepath string) bool {
	// test files take precedence over source files in case of overlapping (such as with go language)
	if lang.IsTestFile(filepath) {
		return false
	}
	return lang.SrcFiles.matches(filepath, lang.baseDir)
}

// IsTestFile returns true if the provided filePath is recognized as a test file for this language
func (lang Language) IsTestFile(filepath string) bool {
	return lang.TestFiles.matches(filepath, lang.baseDir)
}

// IsLanguageFile returns true if the provided filePath is recognized as either a source or a test file for this language
func (lang Language) IsLanguageFile(filename string) bool {
	return lang.IsSrcFile(filename) || lang.IsTestFile(filename)
}

// DirsToWatch returns the list of directories that TCR engine needs to watch for this language
func (lang Language) DirsToWatch(baseDir string) (dirs []string) {
	// First we concatenate the 2 lists
	concat := append(lang.SrcFiles.Directories, lang.TestFiles.Directories...)

	// Then we use a map to remove duplicates
	unique := make(map[string]int)
	for _, dir := range concat {
		unique[dir] = 1
	}

	// Finally, we prefix each item with baseDir
	for dir := range unique {
		dirs = append(dirs, filepath.Join(baseDir, toLocalPath(dir)))
	}
	//report.PostInfo(dirs)
	return dirs
}

// GetToolchain returns the toolchain instance for this language.
// - If toolchainName is provided and is compatible with this language, it will be returned.
// - If toolchainName is provided but is not compatible with this language, an error is returned.
// - If toolchainName is not provided, the language's default toolchain is returned.
func (lang Language) GetToolchain(toolchainName string) (tchn *toolchain.Toolchain, err error) {
	// We first retrieve the toolchain
	if toolchainName != "" {
		// If toolchain is specified, we try to get it
		tchn, err = toolchain.Get(toolchainName)
		if err != nil {
			return nil, err
		}
	} else {
		// If no toolchain is specified, we use the default toolchain for this language
		tchn, err = toolchain.Get(lang.Toolchains.Default)
		if err != nil {
			return nil, err
		}
	}

	// Then we check language/toolchain compatibility
	comp, err := lang.verifyCompatibility(tchn)
	if !comp || err != nil {
		return nil, err
	}
	return tchn, nil
}

func (lang Language) verifyCompatibility(tchn *toolchain.Toolchain) (bool, error) {
	if tchn == nil {
		return false, errors.New("toolchain is unknown")
	}
	if !lang.worksWithToolchain(tchn.GetName()) {
		return false, fmt.Errorf(
			"%v toolchain is not compatible with %v language",
			tchn.GetName(), lang.GetName(),
		)
	}
	return true, nil
}

func (lang Language) worksWithToolchain(toolchainName string) bool {
	for _, compatible := range lang.Toolchains.Compatible {
		if compatible == toolchainName {
			return true
		}
	}
	return false
}

// SrcDirs returns the list of subdirectories that may contain source files for this language
func (lang Language) SrcDirs() []string {
	return lang.SrcFiles.Directories
}

//func (lang Language) isInSrcTree(path string) bool {
//	return lang.SrcFiles.isInFileTree(path, lang.baseDir)
//}
//
//func (lang Language) isInTestTree(path string) bool {
//	return lang.TestFiles.isInFileTree(path, lang.baseDir)
//}

func (lang *Language) setBaseDir(dir string) {
	lang.baseDir, _ = filepath.Abs(dir)
}

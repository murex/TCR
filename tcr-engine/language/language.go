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
)

type (
	// Toolchains defines the structure for toolchains related to a language
	Toolchains struct {
		Default    string
		Compatible []string
	}

	// Language defines the data structure of a language.
	// - name is the name of the language, it must be unique in the list of available languages
	Language struct {
		name           string
		toolchains     Toolchains
		srcFileFilter  FileTreeFilter
		testFileFilter FileTreeFilter
		baseDir        string
	}

	// LangInterface provides the interface for interacting with a language
	LangInterface interface {
		GetName() string
		GetToolchains() Toolchains
		GetSrcFileFilter() FileTreeFilter
		GetTestFileFilter() FileTreeFilter
		GetToolchain(toolchainName string) (toolchain.TchnInterface, error)
		DirsToWatch(baseDir string) []string
		IsSrcFile(aPath string) bool
		IsTestFile(aPath string) bool
		IsLanguageFile(filename string) bool
		AllSrcFiles() ([]string, error)
		AllTestFiles() ([]string, error)
		checkName() error
		checkCompatibleToolchains() error
		checkDefaultToolchain() error
		setBaseDir(dir string)
		worksWithToolchain(toolchainName string) bool
	}
)

// New creates a new Language instance with the provided name, toolchains, srcFiles and testFiles
func New(name string, toolchains Toolchains, srcFiles FileTreeFilter, testFiles FileTreeFilter) *Language {
	return &Language{
		name:           name,
		toolchains:     toolchains,
		srcFileFilter:  srcFiles,
		testFileFilter: testFiles,
	}
}

// GetName provides the name of the toolchain
func (lang *Language) GetName() string {
	return lang.name
}

func (lang *Language) checkName() error {
	if lang.GetName() == "" {
		return errors.New("language name is empty")
	}
	return nil
}

func (lang *Language) checkCompatibleToolchains() error {
	if lang.toolchains.Compatible == nil {
		return errors.New("language has no compatible toolchain defined")
	}
	return nil
}

func (lang *Language) checkDefaultToolchain() error {
	if lang.toolchains.Default == "" {
		return errors.New("language has no default toolchain defined")
	} else if !lang.worksWithToolchain(lang.toolchains.Default) {
		return errors.New("language's default toolchain " +
			lang.toolchains.Default + " is not listed in compatible toolchains list")
	}
	return nil
}

// GetSrcFileFilter provides the language's list of filters for source files
func (lang *Language) GetSrcFileFilter() FileTreeFilter {
	return lang.srcFileFilter
}

// GetTestFileFilter provides the language's list of filters for test files
func (lang *Language) GetTestFileFilter() FileTreeFilter {
	return lang.testFileFilter
}

// IsSrcFile returns true if the provided filePath is recognized as a source file for this language
func (lang *Language) IsSrcFile(aPath string) bool {
	// test files take precedence over source files in case of overlapping (such as with go language)
	if lang.IsTestFile(aPath) {
		return false
	}
	return lang.GetSrcFileFilter().matches(aPath, lang.baseDir)
}

// IsTestFile returns true if the provided filePath is recognized as a test file for this language
func (lang *Language) IsTestFile(aPath string) bool {
	return lang.GetTestFileFilter().matches(aPath, lang.baseDir)
}

// IsLanguageFile returns true if the provided filePath is recognized as either a source
// or a test file for this language
func (lang *Language) IsLanguageFile(aPath string) bool {
	return lang.IsSrcFile(aPath) || lang.IsTestFile(aPath)
}

// DirsToWatch returns the list of directories that TCR engine needs to watch for this language
func (lang *Language) DirsToWatch(baseDir string) (dirs []string) {
	// First we concatenate the 2 lists
	concat := append(lang.GetSrcFileFilter().Directories, lang.GetTestFileFilter().Directories...)

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

// GetToolchains returns the toolchains setup instance for this language.
func (lang *Language) GetToolchains() Toolchains {
	return lang.toolchains
}

// GetToolchain returns the toolchain instance for this language.
// - If toolchainName is provided and is compatible with this language, it will be returned.
// - If toolchainName is provided but is not compatible with this language, an error is returned.
// - If toolchainName is not provided, the language's default toolchain is returned.
func (lang *Language) GetToolchain(toolchainName string) (tchn toolchain.TchnInterface, err error) {
	// We first retrieve the toolchain
	if toolchainName != "" {
		// If toolchain is specified, we try to get it
		tchn, err = toolchain.GetToolchain(toolchainName)
		if err != nil {
			return nil, err
		}
	} else {
		// If no toolchain is specified, we use the default toolchain for this language
		tchn, err = toolchain.GetToolchain(lang.GetToolchains().Default)
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

func (lang *Language) verifyCompatibility(tchn toolchain.TchnInterface) (bool, error) {
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

func (lang *Language) worksWithToolchain(toolchainName string) bool {
	for _, compatible := range lang.GetToolchains().Compatible {
		if compatible == toolchainName {
			return true
		}
	}
	return false
}

// AllSrcFiles returns the list of source files for this language.
// If there is an overlap between source and test files patterns, test files
// are excluded from the returned list
func (lang *Language) AllSrcFiles() (result []string, err error) {
	var srcFiles, testFiles []string

	testFilesMap := make(map[string]bool)
	testFiles, err = lang.allMatchingTestFiles()
	if err != nil {
		return nil, err
	}
	for _, path := range testFiles {
		testFilesMap[path] = true
	}

	srcFiles, err = lang.allMatchingSrcFiles()
	if err != nil {
		return nil, err
	}
	for _, path := range srcFiles {
		if !testFilesMap[path] {
			result = append(result, path)
		}
	}
	return result, nil
}

// AllTestFiles returns the list of test files for this language.
func (lang *Language) AllTestFiles() (result []string, err error) {
	return lang.allMatchingTestFiles()
}

// allMatchingSrcFiles returns the list of source files matching for this language
func (lang *Language) allMatchingSrcFiles() ([]string, error) {
	return lang.GetSrcFileFilter().findAllMatchingFiles(lang.baseDir)
}

// allMatchingTestFiles returns the list of test files matching for this language
func (lang *Language) allMatchingTestFiles() ([]string, error) {
	return lang.GetTestFileFilter().findAllMatchingFiles(lang.baseDir)
}

func (lang *Language) setBaseDir(dir string) {
	// Warning (for tests only): filepath.Abs() does not work with MemMapFs on Windows
	lang.baseDir, _ = filepath.Abs(dir)
}

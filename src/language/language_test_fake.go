//go:build test_helper

/*
Copyright (c) 2023 Murex

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
	"strings"

	"github.com/murex/tcr/toolchain"
)

type allFilesFunc func() ([]string, error)

type (
	// FakeLanguage is a Language fake that makes sure there is no filesystem
	// access. It allows to overwrite AllSrcFiles() and AllTestFiles() behavior
	FakeLanguage struct {
		lang         Language
		allSrcFiles  allFilesFunc
		allTestFiles allFilesFunc
	}
)

// NewFakeLanguage creates a new Language instance compatible with the provided toolchain
func NewFakeLanguage(toolchainName string) *FakeLanguage {
	return &FakeLanguage{
		lang: Language{
			name: "fake-language",
			toolchains: Toolchains{
				Default:    toolchainName,
				Compatible: []string{toolchainName},
			},
		},
		// Default behaviour for AllSrcFiles()
		allSrcFiles: func() ([]string, error) {
			return []string{"fake-src-file1", "fake-src-file2"}, nil
		},
		// Default behaviour for AllTestFiles()
		allTestFiles: func() ([]string, error) {
			return []string{"fake-test-file1", "fake-test-file2"}, nil
		},
	}
}

// WithAllSrcFiles allows to change the behaviour of AllSrcFiles() method
func (fl *FakeLanguage) WithAllSrcFiles(f allFilesFunc) *FakeLanguage {
	fl.allSrcFiles = f
	return fl
}

// WithAllTestFiles allows to change the behaviour of AllTestFiles() method
func (fl *FakeLanguage) WithAllTestFiles(f allFilesFunc) *FakeLanguage {
	fl.allTestFiles = f
	return fl
}

// AllSrcFiles returns the list of source files for this language.
// Always returns a list of 2 fake filenames.
func (fl *FakeLanguage) AllSrcFiles() (result []string, err error) {
	return fl.allSrcFiles()
}

// AllTestFiles returns the list of test files for this language.
// Always returns a list of 2 fake filenames.
func (fl *FakeLanguage) AllTestFiles() (result []string, err error) {
	return fl.allTestFiles()
}

// IsSrcFile returns true if the provided filePath is recognized as a source file for this language
func (fl *FakeLanguage) IsSrcFile(name string) bool {
	return !strings.Contains(name, "test")
}

// GetName uses real Language behaviour
func (fl *FakeLanguage) GetName() string {
	return fl.lang.GetName()
}

// GetToolchains uses real Language behaviour
func (fl *FakeLanguage) GetToolchains() Toolchains {
	return fl.lang.GetToolchains()
}

// GetSrcFileFilter uses real Language behaviour
func (fl *FakeLanguage) GetSrcFileFilter() FileTreeFilter {
	return fl.lang.GetSrcFileFilter()
}

// GetTestFileFilter uses real Language behaviour
func (fl *FakeLanguage) GetTestFileFilter() FileTreeFilter {
	return fl.lang.GetTestFileFilter()
}

// GetToolchain uses real Language behaviour
func (fl *FakeLanguage) GetToolchain(toolchainName string) (toolchain.TchnInterface, error) {
	return fl.lang.GetToolchain(toolchainName)
}

// DirsToWatch uses real Language behaviour
func (fl *FakeLanguage) DirsToWatch(baseDir string) []string {
	return fl.lang.DirsToWatch(baseDir)
}

// IsTestFile uses real Language behaviour
func (fl *FakeLanguage) IsTestFile(aPath string) bool {
	return fl.lang.IsTestFile(aPath)
}

// IsLanguageFile uses real Language behaviour
func (fl *FakeLanguage) IsLanguageFile(filename string) bool {
	return fl.lang.IsLanguageFile(filename)
}

func (fl *FakeLanguage) checkName() error {
	return fl.lang.checkName()
}

func (fl *FakeLanguage) checkCompatibleToolchains() error {
	return fl.lang.checkCompatibleToolchains()
}

func (fl *FakeLanguage) checkDefaultToolchain() error {
	return fl.lang.checkDefaultToolchain()
}

func (fl *FakeLanguage) setBaseDir(dir string) {
	fl.lang.setBaseDir(dir)
}

func (fl *FakeLanguage) worksWithToolchain(toolchainName string) bool {
	return fl.lang.worksWithToolchain(toolchainName)
}

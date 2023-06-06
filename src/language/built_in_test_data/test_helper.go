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

package built_in_test_data //nolint:revive

import (
	"path/filepath"
	"strings"
)

// BuiltInTestData provides test data for a built-in language.
type BuiltInTestData struct {
	Name                 string
	Extensions           []string
	DefaultToolchain     string
	CompatibleToolchains []string
	DirsToWatch          []string
	SrcMatchersDefs      []MatcherDef
	TestMatcherDefs      []MatcherDef
}

// MatcherDef defines the pattern for building a file path matcher.
type MatcherDef struct {
	ParentDir string
	FileBase  string
}

// BuiltInTests is the slice containing all built-in language test data.
var BuiltInTests []BuiltInTestData

var knownFileExtensions []string
var knownToolchains = make(map[string]bool)

func addBuiltIn(data BuiltInTestData) {
	BuiltInTests = append(BuiltInTests, data)
	knownFileExtensions = append(knownFileExtensions, data.Extensions...)
	for _, t := range data.CompatibleToolchains {
		knownToolchains[t] = true
	}
}

// AllKnownFileExtensionsBut provides the list of all file extensions known for built-in languages
// with the exception of the provided parameters.
func AllKnownFileExtensionsBut(ext ...string) (out []string) {
	for _, e := range knownFileExtensions {
		if !contains(ext, e) {
			out = append(out, e)
		}
	}
	return out
}

// AllKnownToolchainsBut provides the list of all toolchains known for built-in languages
// with the exception of the provided parameters.
func AllKnownToolchainsBut(toolchainNames ...string) (out []string) {
	for k := range knownToolchains {
		if !contains(toolchainNames, k) {
			out = append(out, k)
		}
	}
	return out
}

func contains(items []string, searched string) bool {
	for _, item := range items {
		if searched == item {
			return true
		}
	}
	return false
}

// FilePathMatcher defines the pattern for building a file path matcher.
type FilePathMatcher struct {
	FilePath   string
	IsSrcFile  bool
	IsTestFile bool
}

// ShouldMatchSrc returns a FilePathMatcher for src files.
func ShouldMatchSrc(filePath string) FilePathMatcher {
	return FilePathMatcher{FilePath: filePath, IsSrcFile: true, IsTestFile: false}
}

// ShouldMatchTest returns a FilePathMatcher for test files.
func ShouldMatchTest(filePath string) FilePathMatcher {
	return FilePathMatcher{FilePath: filePath, IsSrcFile: false, IsTestFile: true}
}

// ShouldNotMatch returns a FilePathMatcher for neither src or test files.
func ShouldNotMatch(filePath string) FilePathMatcher {
	return FilePathMatcher{FilePath: filePath, IsSrcFile: false, IsTestFile: false}
}

// BuildFilePathMatchers is a convenience method building a set of matching tests with the provided
// dir, fileBaseName and fileExt.
// Among other things, it checks that extension matching is case-insensitive, that temporary files with
// "~" or ".swp" are excluded, and that matching should pass with files in subdirectories from parentDir
func BuildFilePathMatchers(matcher func(string) FilePathMatcher, def MatcherDef, fileExt string) []FilePathMatcher {
	var fileBasePath = filepath.Join(def.ParentDir, def.FileBase)
	var subDirPath = filepath.Join(def.ParentDir, "subDir", def.FileBase)
	return []FilePathMatcher{
		ShouldNotMatch(fileBasePath),
		matcher(fileBasePath + strings.ToLower(fileExt)),
		matcher(fileBasePath + strings.ToUpper(fileExt)),
		ShouldNotMatch(subDirPath),
		matcher(subDirPath + fileExt),
		ShouldNotMatch(fileBasePath + fileExt + "~"),
		ShouldNotMatch(fileBasePath + fileExt + ".swp"),
	}
}

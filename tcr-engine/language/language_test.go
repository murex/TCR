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
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_dirs_to_watch_should_contain_both_source_and_test_dirs(t *testing.T) {
	const srcDir, testDir = "src-dir", "test-dir"
	lang := ALanguage(
		WithSrcFiles(AFileTreeFilter(WithDirectory(srcDir))),
		WithTestFiles(AFileTreeFilter(WithDirectory(testDir))),
	)
	dirs := lang.DirsToWatch("")
	assert.Contains(t, dirs, srcDir)
	assert.Contains(t, dirs, testDir)
}

func Test_dirs_to_watch_should_be_prefixed_with_workdir_path(t *testing.T) {
	const srcDir, testDir = "src-dir", "test-dir"
	baseDir, _ := os.Getwd()
	lang := ALanguage(
		WithSrcFiles(AFileTreeFilter(WithDirectory(srcDir))),
		WithTestFiles(AFileTreeFilter(WithDirectory(testDir))),
	)
	dirs := lang.DirsToWatch(baseDir)
	assert.Contains(t, dirs, filepath.Join(baseDir, srcDir))
	assert.Contains(t, dirs, filepath.Join(baseDir, testDir))
}

func Test_dirs_to_watch_should_not_have_duplicates(t *testing.T) {
	const dir = "dir"
	baseDir, _ := os.Getwd()
	lang := ALanguage(
		WithSrcFiles(AFileTreeFilter(WithDirectory(dir), WithDirectory(dir))),
		WithTestFiles(AFileTreeFilter(WithDirectory(dir), WithDirectory(dir))),
	)
	assert.Equal(t, 1, len(lang.DirsToWatch(baseDir)))
}

// Assert utility functions for language type

func assertDefaultToolchain(t *testing.T, expected string, name string) {
	lang, _ := Get(name)
	assert.Equal(t, expected, lang.GetToolchains().Default)
}

func assertListOfDirsToWatch(t *testing.T, expected []string, name string) {
	dirList := getBuiltIn(name).DirsToWatch("")
	for _, dir := range expected {
		assert.Contains(t, dirList, toLocalPath(dir))
	}
}

func assertCompatibleToolchains(t *testing.T, expected []string, name string) {
	lang := getBuiltIn(name)
	for _, toolchain := range expected {
		assert.True(t, lang.worksWithToolchain(toolchain))
	}
}

func assertIncompatibleToolchains(t *testing.T, expected []string, name string) {
	lang := getBuiltIn(name)
	for _, toolchain := range expected {
		assert.False(t, lang.worksWithToolchain(toolchain))
	}
}

type filePathMatcher struct {
	filePath   string
	isSrcFile  bool
	isTestFile bool
}

func shouldMatchSrc(filePath string) filePathMatcher {
	return filePathMatcher{filePath: filePath, isSrcFile: true, isTestFile: false}
}

func shouldMatchTest(filePath string) filePathMatcher {
	return filePathMatcher{filePath: filePath, isSrcFile: false, isTestFile: true}
}
func shouldNotMatch(filePath string) filePathMatcher {
	return filePathMatcher{filePath: filePath, isSrcFile: false, isTestFile: false}
}

// buildFilePathMatchers is a convenience method building a set of matching tests with the provided
// dir, fileBaseName and fileExt.
// Among other things, it checks that extension matching is case-insensitive, that temporary files with
// "~" or ".swp" are excluded, and that matching should pass with files in subdirectories from parentDir
func buildFilePathMatchers(matcher func(string) filePathMatcher, parentDir string, fileBase string, fileExt string) []filePathMatcher {
	var fileBasePath = filepath.Join(parentDir, fileBase)
	var subDirPath = filepath.Join(parentDir, "subDir", fileBase)
	return []filePathMatcher{
		shouldNotMatch(fileBasePath),
		matcher(fileBasePath + strings.ToLower(fileExt)),
		matcher(fileBasePath + strings.ToUpper(fileExt)),
		shouldNotMatch(subDirPath),
		matcher(subDirPath + fileExt),
		shouldNotMatch(fileBasePath + fileExt + "~"),
		shouldNotMatch(fileBasePath + fileExt + ".swp"),
	}
}

func assertFilePathsMatching(t *testing.T, matchers []filePathMatcher, name string) {
	lang, _ := GetLanguage(name, "")
	for _, matcher := range matchers {
		assert.Equal(t, matcher.isSrcFile, lang.IsSrcFile(matcher.filePath),
			"Should %v be a source file?", matcher.filePath)
		assert.Equal(t, matcher.isTestFile, lang.IsTestFile(matcher.filePath),
			"Should %v be a test file?", matcher.filePath)
	}
}

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

//func WithBaseDir(path string) func(lang *Language) {
//	return func(lang *Language) { lang.setBaseDir(path) }
//}

// =========================================================================================

func Test_does_not_support_empty_language_name(t *testing.T) {
	assert.False(t, isSupported(""))
}

func Test_does_not_support_unregistered_language_name(t *testing.T) {
	assert.False(t, isSupported("unregistered-language"))
}

func Test_unrecognized_language_name(t *testing.T) {
	lang, err := Get("dummy-language")
	assert.Error(t, err)
	assert.Zero(t, lang)
}

func Test_can_add_a_built_in_language(t *testing.T) {
	const name = "new-built-in-language"
	assert.False(t, isBuiltIn(name))
	assert.NoError(t, addBuiltIn(ALanguage(WithName(name))))
	assert.True(t, isBuiltIn(name))
}

func Test_cannot_add_a_built_in_language_with_no_name(t *testing.T) {
	assert.Error(t, addBuiltIn(ALanguage(WithName(""))))
}

func Test_language_name_is_case_insensitive(t *testing.T) {
	const name = "miXeD-CasE"
	_ = Register(ALanguage(WithName(name)))
	assertNameIsNotCaseSensitive(t, name)
}

func Test_can_register_a_language(t *testing.T) {
	const name = "new-language"
	assert.False(t, isSupported(name))
	assert.NoError(t, Register(ALanguage(
		WithName(name),
		WithDefaultToolchain("default-toolchain"),
		WithCompatibleToolchain("default-toolchain"),
	)))
	assert.True(t, isSupported(name))
}

func Test_cannot_register_a_language_with_no_name(t *testing.T) {
	assert.Error(t, Register(ALanguage(WithName(""))))
}

func Test_cannot_register_a_language_with_no_compatible_toolchain(t *testing.T) {
	const name = "no-compatible-toolchain"
	assert.Error(t, Register(ALanguage(WithName(name), WithNoCompatibleToolchain())))
	assert.False(t, isSupported(name))
}

func Test_cannot_register_a_language_with_no_default_toolchain(t *testing.T) {
	const name = "no-default-toolchain"
	assert.Error(t, Register(ALanguage(WithName(name), WithNoDefaultToolchain())))
	assert.False(t, isSupported(name))
}

func Test_cannot_register_a_language_with_default_toolchain_not_compatible(t *testing.T) {
	const name = "default-toolchain-not-compatible"
	assert.Error(t, Register(ALanguage(
		WithName(name),
		WithDefaultToolchain("toolchain1"),
		WithCompatibleToolchain("toolchain2"),
	)))
	assert.False(t, isSupported(name))
}

func Test_does_not_detect_language_from_a_dir_name_not_matching_a_known_language(t *testing.T) {
	lang, err := detectLanguageFromDirName("dummy")
	assert.Error(t, err)
	assert.Zero(t, lang)
}

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

func assertIsABuiltInLanguage(t *testing.T, name string) {
	assert.True(t, isBuiltIn(name))
}

func assertIsSupported(t *testing.T, name string) {
	assert.True(t, isSupported(name))
}

func assertNameIsNotCaseSensitive(t *testing.T, name string) {
	assert.True(t, isSupported(name))
	assert.True(t, isSupported(strings.ToUpper(name)))
	assert.True(t, isSupported(strings.ToLower(name)))
	assert.True(t, isSupported(strings.Title(name)))
}

func assertLanguageInitialization(t *testing.T, name string) {
	lang, err := Get(name)
	assert.NoError(t, err)
	assert.Equal(t, name, lang.GetName())
}

func assertLanguageName(t *testing.T, name string) {
	lang, _ := Get(name)
	assert.Equal(t, name, lang.GetName())
}

func assertDefaultToolchain(t *testing.T, expected string, name string) {
	lang, _ := Get(name)
	assert.Equal(t, expected, lang.GetToolchains().Default)
}

func assertFallbacksOnDirNameIfLanguageIsNotSpecified(t *testing.T, dirName string) {
	lang, err := GetLanguage("", filepath.Join("some", "path", dirName))
	assert.NoError(t, err)
	assert.NotEmpty(t, lang)
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

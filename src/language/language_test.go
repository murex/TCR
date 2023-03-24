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
	"github.com/murex/tcr/toolchain"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var languageFileExtensions []string

func registerLanguageFileExtensionsForTests(ext ...string) {
	languageFileExtensions = append(languageFileExtensions, ext...)
}

func allLanguageFileExtensionsBut(ext ...string) (out []string) {
	for _, e := range languageFileExtensions {
		if !contains(ext, e) {
			out = append(out, e)
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

var knownToolchains = make(map[string]bool)

func registerToolchainsForTests(toolchainNames ...string) {
	for _, t := range toolchainNames {
		knownToolchains[t] = true
	}
}

func allKnownToolchainsBut(toolchainNames ...string) (out []string) {
	for k := range knownToolchains {
		if !contains(toolchainNames, k) {
			out = append(out, k)
		}
	}
	return out
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

func Test_a_file_with_no_name_is_not_a_language_file(t *testing.T) {
	lang := ALanguage()
	assert.False(t, lang.IsLanguageFile(""))
}

func Test_a_matching_source_file_is_a_language_file(t *testing.T) {
	const dir = "dir"
	lang := ALanguage(
		WithSrcFiles(AFileTreeFilter(WithDirectory(dir), WithPattern(".*\\.ext"))),
		WithTestFiles(AFileTreeFilter(WithDirectory(dir), WithClosedPattern())),
	)
	assert.True(t, lang.IsLanguageFile(filepath.Join(dir, "some-file.ext")))
}

func Test_a_matching_test_file_is_a_language_file(t *testing.T) {
	const dir = "dir"
	lang := ALanguage(
		WithSrcFiles(AFileTreeFilter(WithDirectory(dir), WithClosedPattern())),
		WithTestFiles(AFileTreeFilter(WithDirectory(dir), WithPattern(".*\\.ext"))),
	)
	assert.True(t, lang.IsLanguageFile(filepath.Join(dir, "some-file.ext")))
}

func Test_a_file_not_matching_src_or_test_is_not_a_language_file(t *testing.T) {
	const dir = "dir"
	lang := ALanguage(
		WithSrcFiles(AFileTreeFilter(WithDirectory(dir), WithClosedPattern())),
		WithTestFiles(AFileTreeFilter(WithDirectory(dir), WithClosedPattern())),
	)
	assert.False(t, lang.IsLanguageFile(filepath.Join(dir, "some-file.ext")))
}

func Test_get_toolchain_with_unregistered_toolchain(t *testing.T) {
	lang := ALanguage(
		WithDefaultToolchain("some-toolchain"),
		WithCompatibleToolchain("some-toolchain"),
	)

	actual, err := lang.GetToolchain("some-toolchain")
	assert.Zero(t, actual)
	assert.Error(t, err)
	assert.Equal(t, errors.New("toolchain not supported: some-toolchain"), err)
}

func Test_get_toolchain_with_empty_toolchain_name_and_a_default_toolchain(t *testing.T) {
	lang := ALanguage(
		WithDefaultToolchain("some-toolchain"),
		WithCompatibleToolchain("some-toolchain"),
	)
	_ = toolchain.Register(
		toolchain.AToolchain(toolchain.WithName("some-toolchain")))
	actual, err := lang.GetToolchain("")
	toolchain.Unregister("some-toolchain")

	assert.Zero(t, err)
	assert.Equal(t, lang.toolchains.Default, actual.GetName())
}

func Test_get_toolchain_with_empty_toolchain_name_and_no_default_toolchain(t *testing.T) {
	lang := ALanguage(
		WithDefaultToolchain(""),
		WithCompatibleToolchain(""),
	)
	actual, err := lang.GetToolchain("")

	assert.Zero(t, actual)
	assert.Error(t, err)
	assert.Equal(t, errors.New("toolchain name not provided"), err)
}

func Test_get_toolchain_with_non_compatible_toolchain(t *testing.T) {
	lang := ALanguage(
		WithDefaultToolchain("some-toolchain"),
		WithCompatibleToolchain("some-toolchain"),
	)
	_ = toolchain.Register(
		toolchain.AToolchain(toolchain.WithName("other-toolchain")))
	actual, err := lang.GetToolchain("other-toolchain")
	toolchain.Unregister("other-toolchain")

	assert.Zero(t, actual)
	assert.Error(t, err)
	assert.Equal(t, errors.New("other-toolchain toolchain is not compatible with default-language language"), err)
}

func Test_retrieve_language_files(t *testing.T) {
	appFs = afero.NewMemMapFs()
	baseDir := filepath.Join("base-dir")
	srcDir := filepath.Join(baseDir, "src")
	_ = appFs.MkdirAll(srcDir, os.ModeDir)
	srcFile1 := filepath.Join(srcDir, "file1.ext")
	_ = afero.WriteFile(appFs, srcFile1, []byte("some contents"), 0644)
	srcFile2 := filepath.Join(srcDir, "file2.ext")
	_ = afero.WriteFile(appFs, srcFile2, []byte("some contents"), 0644)
	testDir := filepath.Join(baseDir, "test")
	_ = appFs.MkdirAll(testDir, os.ModeDir)
	testFile1 := filepath.Join(testDir, "file1.ext")
	_ = afero.WriteFile(appFs, testFile1, []byte("some contents"), 0644)
	testFile2 := filepath.Join(testDir, "file2.ext")
	_ = afero.WriteFile(appFs, testFile2, []byte("some contents"), 0644)

	lang := ALanguage(
		WithBaseDir(baseDir),
		WithSrcFiles(AFileTreeFilter(WithDirectory("src"), WithPattern(".*\\.ext"))),
		WithTestFiles(AFileTreeFilter(WithDirectory("test"), WithPattern(".*\\.ext"))),
	)
	srcFiles, errSrc := lang.AllSrcFiles()
	assert.NoError(t, errSrc)
	assert.Contains(t, srcFiles, srcFile1)
	assert.Contains(t, srcFiles, srcFile2)
	assert.NotContains(t, srcFiles, testFile1)
	assert.NotContains(t, srcFiles, testFile2)

	testFiles, errTest := lang.AllTestFiles()
	assert.NoError(t, errTest)
	assert.NotContains(t, testFiles, srcFile1)
	assert.NotContains(t, testFiles, srcFile2)
	assert.Contains(t, testFiles, testFile1)
	assert.Contains(t, testFiles, testFile2)
}

// Assert utility functions for language type

func assertDefaultToolchain(t *testing.T, languageName string, toolchainName string) {
	t.Helper()
	lang, _ := Get(languageName)
	assert.Equal(t, toolchainName, lang.GetToolchains().Default)
}

func assertListOfDirsToWatch(t *testing.T, languageName string, dirs ...string) {
	t.Helper()
	dirList := getBuiltIn(languageName).DirsToWatch("")
	for _, dir := range dirs {
		t.Run(dir, func(t *testing.T) {
			assert.Contains(t, dirList, toLocalPath(dir))
		})
	}
}

func assertCompatibleToolchains(t *testing.T, languageName string, toolchainNames ...string) {
	t.Helper()
	lang := getBuiltIn(languageName)
	for _, tchn := range toolchainNames {
		t.Run(tchn, func(t *testing.T) {
			assert.True(t, lang.worksWithToolchain(tchn))
		})
	}
}

func assertIncompatibleToolchains(t *testing.T, languageName string, toolchainNames ...string) {
	t.Helper()
	lang := getBuiltIn(languageName)
	for _, tchn := range toolchainNames {
		t.Run(tchn, func(t *testing.T) {
			assert.False(t, lang.worksWithToolchain(tchn))
		})
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

func assertFilePathsMatching(t *testing.T, languageName string, matchers ...filePathMatcher) {
	t.Helper()
	lang, _ := GetLanguage(languageName, "")
	for _, matcher := range matchers {
		assert.Equal(t, matcher.isSrcFile, lang.IsSrcFile(matcher.filePath),
			"Should %v be a source file?", matcher.filePath)
		assert.Equal(t, matcher.isTestFile, lang.IsTestFile(matcher.filePath),
			"Should %v be a test file?", matcher.filePath)
	}
}

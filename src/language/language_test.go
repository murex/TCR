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
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/murex/tcr/toolchain"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
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
	appFS = afero.NewMemMapFs()
	baseDir := filepath.Join("base-dir")
	srcDir := filepath.Join(baseDir, "src")
	_ = appFS.MkdirAll(srcDir, os.ModeDir)
	srcFile1 := filepath.Join(srcDir, "file1.ext")
	_ = afero.WriteFile(appFS, srcFile1, []byte("some contents"), 0644)
	srcFile2 := filepath.Join(srcDir, "file2.ext")
	_ = afero.WriteFile(appFS, srcFile2, []byte("some contents"), 0644)
	testDir := filepath.Join(baseDir, "test")
	_ = appFS.MkdirAll(testDir, os.ModeDir)
	testFile1 := filepath.Join(testDir, "file1.ext")
	_ = afero.WriteFile(appFS, testFile1, []byte("some contents"), 0644)
	testFile2 := filepath.Join(testDir, "file2.ext")
	_ = afero.WriteFile(appFS, testFile2, []byte("some contents"), 0644)

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

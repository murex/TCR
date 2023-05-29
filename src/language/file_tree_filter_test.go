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
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func Test_convert_back_slashed_path_to_slashed_path(t *testing.T) {
	var input = "some\\path\\with\\backslash"
	var expected = "some/path/with/backslash"
	assert.Equal(t, expected, toSlashedPath(input))
}

func Test_convert_slashed_path_to_slashed_path(t *testing.T) {
	var input = "some/path/with/slash"
	var expected = "some/path/with/slash"
	assert.Equal(t, expected, toSlashedPath(input))
}

func Test_file_tree_filter_with_no_directory(t *testing.T) {
	filter := AFileTreeFilter(WithNoDirectory(), WithOpenPattern())
	assert.True(t, filter.matches("some_file.ext", ""))
	assert.True(t, filter.matches("some_dir/some_file.ext", ""))
}

func Test_file_tree_filter_with_no_pattern(t *testing.T) {
	filter := AFileTreeFilter(WithNoPattern())
	assert.True(t, filter.matches("some_file.ext", ""))
	assert.True(t, filter.matches("some_dir/some_file.ext", ""))
}

func Test_file_tree_filter_with_open_pattern(t *testing.T) {
	filter := AFileTreeFilter(WithOpenPattern())
	assert.True(t, filter.matches("some_file.ext", ""))
	assert.True(t, filter.matches("some_dir/some_file.ext", ""))
}

func Test_file_tree_filter_with_closed_pattern(t *testing.T) {
	filter := AFileTreeFilter(WithClosedPattern())
	assert.False(t, filter.matches("some_file.ext", ""))
	assert.False(t, filter.matches("some_dir/some_file.ext", ""))
}

func Test_file_path_is_in_file_tree(t *testing.T) {
	const matchDir = "some_dir"
	baseDir, _ := os.Getwd()
	filter := AFileTreeFilter(WithDirectory(matchDir), WithOpenPattern())
	for _, dir := range []string{"", ".", "./x", "x", "x/y", "x/y/z"} {
		okPath := filepath.Join(baseDir, matchDir, dir, "some_file")
		assert.True(t, filter.isInFileTree(okPath, baseDir), okPath)
		koPath := filepath.Join(baseDir, dir, "some_file")
		assert.False(t, filter.isInFileTree(koPath, baseDir), koPath)
	}
}

func Test_file_tree_filter_with_one_file_pattern(t *testing.T) {
	filter := AFileTreeFilter(WithPattern(".*\\.ext"))
	assert.True(t, filter.matches("base.ext", ""))
	assert.False(t, filter.matches("base.other_ext", ""))
}

func Test_file_tree_filter_with_multiple_file_patterns(t *testing.T) {
	filter := AFileTreeFilter(WithDirectory(""), WithPattern(".*\\.ext1"), WithPattern(".*\\.ext2"))
	assert.True(t, filter.matches("base.ext1", ""))
	assert.True(t, filter.matches("base.ext2", ""))
	assert.False(t, filter.matches("base.ext3", ""))
}

func Test_find_all_matching_files(t *testing.T) {
	appFS = afero.NewMemMapFs()
	baseDir := filepath.Join("base-dir")
	srcDir := filepath.Join(baseDir, "src")
	_ = appFS.MkdirAll(srcDir, os.ModeDir)
	testDir := filepath.Join(baseDir, "test")
	_ = appFS.MkdirAll(testDir, os.ModeDir)
	matching := filepath.Join(srcDir, "file.ext")
	_ = afero.WriteFile(appFS, matching, []byte("some contents"), 0644)
	nonMatchingName := filepath.Join(srcDir, "file.other-ext")
	_ = afero.WriteFile(appFS, nonMatchingName, []byte("some contents"), 0644)
	nonMatchingDir := filepath.Join(testDir, "file.ext")
	_ = afero.WriteFile(appFS, nonMatchingDir, []byte("some contents"), 0644)

	filter := AFileTreeFilter(WithDirectory("src"), WithPattern(".*\\.ext"))
	files, err := filter.findAllMatchingFiles(baseDir)
	assert.NoError(t, err)
	assert.Contains(t, files, matching)
	assert.NotContains(t, files, nonMatchingName)
	assert.NotContains(t, files, nonMatchingDir)
}

func Test_find_all_matching_files_with_wrong_base_dir(t *testing.T) {
	appFS = afero.NewMemMapFs()
	baseDir := filepath.Join("base-dir")
	_ = appFS.MkdirAll(baseDir, os.ModeDir)

	filter := AFileTreeFilter(WithDirectory("src"), WithPattern(".*\\.ext"))
	_, err := filter.findAllMatchingFiles("wrong-dir")
	assert.Error(t, err)
}

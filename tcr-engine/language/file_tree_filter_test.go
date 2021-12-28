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
	"testing"
)

// aFileTreeFilter is a test data builder for type FileTreeFilter
func aFileTreeFilter(builders ...func(filter *FileTreeFilter)) *FileTreeFilter {
	filter := &FileTreeFilter{
		Directories: []string{},
		Filters:     []string{},
	}

	for _, build := range builders {
		build(filter)
	}
	return filter
}

func withDirectory(dirName string) func(filter *FileTreeFilter) {
	return func(filter *FileTreeFilter) {
		filter.Directories = append(filter.Directories, dirName)
	}
}

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

func Test_file_path_is_in_file_tree(t *testing.T) {
	const srcDir = "dir"
	baseDir, _ := os.Getwd()
	filter := aFileTreeFilter(withDirectory(srcDir))
	for _, dir := range []string{"", ".", "./x", "x", "x/y", "x/y/z"} {
		okPath := filepath.Join(baseDir, srcDir, dir)
		assert.True(t, filter.isInFileTree(okPath, baseDir))
		koPath := filepath.Join(baseDir, dir)
		assert.False(t, filter.isInFileTree(koPath, baseDir))
	}
}

/*
Copyright (c) 2022 Murex

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

package vcs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_compute_changed_lines(t *testing.T) {
	assert.Equal(t, 6, NewFileDiff("", 4, 2).ChangedLines())
}

func Test_total_number_of_changed_lines(t *testing.T) {
	var diffs = FileDiffs{
		NewFileDiff("dir1/File1", 1, 2),
		NewFileDiff("dir1/File2", 3, 4),
		NewFileDiff("dir2/File3", 5, 6),
		NewFileDiff("dir2/File4", 7, 8),
	}

	testFlags := []struct {
		desc      string
		predicate func(filepath string) bool
		expected  int
	}{
		{
			"no predicate", nil, 36},
		{
			"all files",
			func(filepath string) bool { return true },
			36,
		},
		{
			"no files",
			func(filepath string) bool { return false },
			0,
		},
		{
			"file1 only",
			func(filepath string) bool {
				return filepath == "dir1/File1"
			},
			3,
		},
		{
			"dir1 only",
			func(filepath string) bool {
				return strings.Index(filepath, "dir1/") == 0
			},
			10,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expected, diffs.ChangedLines(tt.predicate))
		})
	}
}

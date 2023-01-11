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

package p4

import (
	"errors"
	"fmt"
	"github.com/murex/tcr/vcs"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/charmap"
	"testing"
)

func Test_convert_line(t *testing.T) {
	tests := []struct {
		desc     string
		input    rune
		encoding *charmap.Charmap
		expected string
	}{
		{"UTC-8 smiley", '🙂', nil, "f09f9982"},
		{"ISO8859-1 smiley", '🙂', charmap.ISO8859_1, "c3b0c29fc299c282"},
		{"WINDOWS-1252 smiley", '🙂', charmap.Windows1252, "c3b0c5b8e284a2e2809a"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			//t.Logf("%x - %x", test.input, convertLine(test.encoding, string(test.input)))
			assert.Equal(t, test.expected, fmt.Sprintf("%x", convertLine(test.encoding, string(test.input))))
		})
	}
}

func Test_p4_diff_command(t *testing.T) {
	testFlags := []struct {
		desc         string
		p4DiffOutput string
		p4DiffError  error
		expectError  bool
		expectedArgs []string
		expectedDiff vcs.FileDiffs
	}{
		{"p4 command arguments",
			"",
			errors.New("p4 diff error"),
			true,
			[]string{
				"diff", "-f", "-Od", "-dw", "-ds"},
			nil,
		},
		{"p4 error",
			"",
			errors.New("p4 diff error"),
			true,
			nil,
			nil,
		},
		{"0 file changed",
			"",
			nil,
			false,
			nil,
			nil,
		},
		{"1 file changed",
			"==== //some-depot/some-file.txt#1 - C:\\some-path\\some-file.txt ====\r\n" +
				"add 1 chunks 1 lines\r\n" +
				"deleted 1 chunks 1 lines\r\n" +
				"changed 0 chunks 0 / 0 lines\r\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: "C:\\some-path\\some-file.txt", AddedLines: 1, RemovedLines: 1},
			},
		},
		{"2 files changed",
			"==== //some-depot/file1.txt#1 - C:\\some-path\\file1.txt ====\r\n" +
				"add 1 chunks 1 lines\r\n" +
				"deleted 1 chunks 1 lines\r\n" +
				"changed 1 chunks 1 / 1 lines\r\n" +
				"==== //some-depot/file2.txt#1 - C:\\some-path\\file2.txt ====\r\n" +
				"add 1 chunks 1 lines\r\n" +
				"deleted 1 chunks 1 lines\r\n" +
				"changed 1 chunks 1 / 1 lines\r\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: "C:\\some-path\\file1.txt", AddedLines: 1 + 1, RemovedLines: 1 + 1},
				{Path: "C:\\some-path\\file2.txt", AddedLines: 1 + 1, RemovedLines: 1 + 1},
			},
		},
		{"file changed in sub-directory",
			"==== //some-depot/some-dir/some-file.txt#1 - C:\\some-path\\some-dir\\some-file.txt ====\r\n" +
				"add 1 chunks 1 lines\r\n" +
				"deleted 1 chunks 1 lines\r\n" +
				"changed 0 chunks 0 / 0 lines\r\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: "C:\\some-path\\some-dir\\some-file.txt", AddedLines: 1, RemovedLines: 1},
			},
		},
		{"1 file changed with added lines only",
			"==== //some-depot/some-file.txt#1 - C:\\some-path\\some-file.txt ====\r\n" +
				"add 1 chunks 15 lines\r\n" +
				"deleted 0 chunks 0 lines\r\n" +
				"changed 0 chunks 0 / 0 lines\r\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: "C:\\some-path\\some-file.txt", AddedLines: 15, RemovedLines: 0},
			},
		},
		{"1 file changed with removed lines only",
			"==== //some-depot/some-file.txt#1 - C:\\some-path\\some-file.txt ====\r\n" +
				"add 0 chunks 0 lines\r\n" +
				"deleted 0 chunks 7 lines\r\n" +
				"changed 0 chunks 0 / 0 lines\r\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: "C:\\some-path\\some-file.txt", AddedLines: 0, RemovedLines: 7},
			},
		},
		{"1 file changed with changed lines only",
			"==== //some-depot/some-file.txt#1 - C:\\some-path\\some-file.txt ====\r\n" +
				"add 0 chunks 0 lines\r\n" +
				"deleted 0 chunks 0 lines\r\n" +
				"changed 1 chunks 2 / 5 lines\r\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: "C:\\some-path\\some-file.txt", AddedLines: 5, RemovedLines: 2},
			},
		},
		{"1 file changed with added, deleted and changed lines",
			"==== //some-depot/some-file.txt#1 - C:\\some-dir\\some-file.txt ====\r\n" +
				"add 2 chunks 10 lines\r\n" +
				"deleted 1 chunks 5 lines\r\n" +
				"changed 3 chunks 15 / 2 lines\r\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: "C:\\some-dir\\some-file.txt", AddedLines: 10 + 2, RemovedLines: 5 + 15},
			},
		},
		{"noise in output trace",
			"Usage: diff [ -d<flags> -f -m max -Od -s<flag> -t ] [files...]\r\n" +
				"Invalid option: -xxx.\r\n",
			nil,
			true,
			nil,
			nil,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			p, _ := newP4Impl(inMemoryDepotInit, "")
			p.rootDir = ""
			p.runP4Function = func(args ...string) (output []byte, err error) {
				actualArgs = args[2:]
				return []byte(tt.p4DiffOutput), tt.p4DiffError
			}
			fileDiffs, err := p.Diff()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.expectedArgs != nil {
				assert.Equal(t, tt.expectedArgs, actualArgs)
			}
			assert.Equal(t, tt.expectedDiff, fileDiffs)
		})
	}
}

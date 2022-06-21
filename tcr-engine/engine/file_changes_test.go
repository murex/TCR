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

package engine

import (
	"github.com/murex/tcr/tcr-engine/language"
	"github.com/murex/tcr/tcr-engine/vcs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_compute_number_of_changed_lines(t *testing.T) {
	testFlags := []struct {
		desc         string
		files        vcs.FileDiffs
		expectedSrc  int
		expectedTest int
	}{
		{"0 file", nil, 0, 0},
		{
			"1 src file - 0 test file",
			vcs.FileDiffs{
				vcs.NewFileDiff("src/main/Hello.java", 1, 2),
			},
			3, 0,
		},
		{
			"2 src files - 0 test file",
			vcs.FileDiffs{
				vcs.NewFileDiff("src/main/Hello1.java", 1, 2),
				vcs.NewFileDiff("src/main/Hello2.java", 3, 4),
			},
			10, 0,
		},
		{
			"1 src file - 1 test file",
			vcs.FileDiffs{
				vcs.NewFileDiff("src/main/Hello.java", 1, 2),
				vcs.NewFileDiff("src/test/HelloTest.java", 3, 4),
			},
			3, 7,
		},
		{
			"0 src file - 1 test file",
			vcs.FileDiffs{
				vcs.NewFileDiff("src/test/HelloTest.java", 1, 2),
			},
			0, 3,
		},
		{
			"0 src file - 2 test files",
			vcs.FileDiffs{
				vcs.NewFileDiff("src/test/HelloTest1.java", 1, 2),
				vcs.NewFileDiff("src/test/HelloTest2.java", 3, 4),
			},
			0, 10,
		},
	}

	lang, _ := language.Get("java")

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expectedSrc, computeSrcLinesChanged(lang, tt.files))
			assert.Equal(t, tt.expectedTest, computeTestLinesChanged(lang, tt.files))
		})
	}
}

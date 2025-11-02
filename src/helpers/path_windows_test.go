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

package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_is_sub_path_of_windows(t *testing.T) {
	tests := []struct {
		desc     string
		subPath  string
		refPath  string
		expected bool
	}{
		{"direct sub-dir", "C:\\user\\bob", "C:\\user", true},
		{"direct not quite sub-dir", "C:\\usex\\bob", "C:\\user", false},
		{"deep sub-dir", "C:\\user\\bob\\deep\\dir", "C:\\user", true},
		{"direct sub-dir with trailing on sub", "C:\\user\\bob\\", "C:\\user", true},
		{"direct sub-dir with trailing on ref", "C:\\user\\bob", "C:\\user\\", true},
		{"backslash sub-dir and slash ref-dir", "C:\\user\\bob", "C:/user", true},
		{"slash ref-dir and backslash ref-dir", "C:/user/bob", "C:\\user", true},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert.Equal(t, test.expected, IsSubPathOf(test.subPath, test.refPath),
				fmt.Sprintf("%s vs %s", test.subPath, test.refPath))
		})
	}
}

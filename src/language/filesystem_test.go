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
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func Test_existing_dirs_in(t *testing.T) {
	baseDir := filepath.Join("base-dir")
	existing := filepath.Join(baseDir, "existing")
	missing := filepath.Join(baseDir, "missing")
	existingFile := filepath.Join(baseDir, "existing-file")

	tests := []struct {
		desc        string
		input       []string
		expected    []string
		expectedErr error
	}{
		{
			"one existing dir",
			[]string{existing},
			[]string{existing},
			nil,
		},
		{
			"one missing dir",
			[]string{missing},
			nil,
			nil,
		},
		{
			"one existing and one missing dir",
			[]string{existing, missing},
			[]string{existing},
			nil,
		},
		{
			"one existing file",
			[]string{existingFile},
			nil,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			appFS = afero.NewMemMapFs()
			_ = appFS.MkdirAll(existing, os.ModeDir)
			_ = afero.WriteFile(appFS, existingFile, []byte("some contents"), 0644)

			result, err := ExistingDirsIn(test.input)
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

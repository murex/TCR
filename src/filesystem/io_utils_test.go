/*
Copyright (c) 2024 Murex

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

package filesystem

import (
	"github.com/murex/tcr/report"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"os"
	"testing"
)

const pathWithPermError = "/path/with/permission/err"

type mockFs struct {
	afero.MemMapFs
}

func (m *mockFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if name == pathWithPermError {
		return nil, fs.ErrPermission
	}
	return m.MemMapFs.OpenFile(name, flag, perm)
}

func Test_write_file_reports_fs_errors(t *testing.T) {
	appFs = &mockFs{MemMapFs: afero.MemMapFs{}}

	tests := []struct {
		description string
		path        string
		expected    int
	}{
		{
			description: "no error",
			path:        "my-file.txt",
			expected:    0,
		},
		{
			description: "permission error",
			path:        pathWithPermError,
			expected:    1,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			report.TestWithIsolatedReporterAndFilters(func(reporter *report.Reporter, sniffer *report.Sniffer) {
				WriteFile(test.path, []byte("some stuff"))
				sniffer.Stop()
				assert.Equal(t, test.expected, sniffer.GetMatchCount())
			}, func(msg report.Message) bool {
				return msg.Type.Category == report.Error
			})
		})
	}

}

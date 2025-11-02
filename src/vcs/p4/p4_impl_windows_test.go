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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_convert_to_p4_client_path_windows(t *testing.T) {
	testFlags := []struct {
		desc          string
		rootDir       string
		clientName    string
		dir           string
		expectedError error
		expected      string
	}{
		{
			"Dir is empty",
			"D:\\p4root",
			"test_client",
			"",
			errors.New("can not convert an empty path"),
			"",
		},
		{
			"Dir outside the root directory",
			"D:\\p4root",
			"test_client",
			"D:\\somewhere_else\\sub_dir",
			errors.New("path is outside p4 root directory"),
			"",
		},
		{
			"Root dir is at window drive level",
			"D:\\",
			"test_client",
			"D:\\sub_dir",
			nil,
			"//test_client/sub_dir/...",
		},
		{
			"Root dir and dir are at window drive level",
			"D:\\",
			"test_client",
			"D:\\",
			nil,
			"//test_client/...",
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			p, _ := newP4Impl(inMemoryDepotInit, tt.rootDir, true)
			p.clientName = tt.clientName
			clientPath, err := p.toP4ClientPath(tt.dir)
			assert.Equal(t, tt.expected, clientPath)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

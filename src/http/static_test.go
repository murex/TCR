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

package http

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed static
var testStaticFS embed.FS

func Test_mount_static_FS(t *testing.T) {
	fs := embedFolder(testStaticFS, "static")
	assert.True(t, fs.Exists("/", "index.html"))
	assert.False(t, fs.Exists("/", "dummy-file"))
}

func Test_mount_static_FS_with_invalid_path(t *testing.T) {
	assert.Panics(t, func() {
		// Panic is triggered by the trailing "/" on targetPath
		embedFolder(testStaticFS, "invalid-path/")
	})
}

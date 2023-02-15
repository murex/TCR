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

package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// IsSubPathOf indicates if aPath is a sub-path of refPath
func IsSubPathOf(aPath string, refPath string) bool {
	cleanDir := filepath.Clean(aPath)
	cleanRoot := filepath.Clean(refPath)

	// If refPath is empty, we consider it as being the root, thus aPath is a sub-path of refPath
	if cleanRoot == "" {
		return true
	}
	if cleanRoot == cleanDir || strings.HasPrefix(cleanDir, filepath.Clean(cleanRoot+string(os.PathSeparator))) {
		return true
	}
	return false
}

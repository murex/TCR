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
	"path/filepath"
	"strings"
)

// Cpp is the language implementation for C++
type Cpp struct {
}

// Name returns the language name. This name is used to detect if a directory contains Cpp files
func (lang Cpp) Name() string {
	return "cpp"
}

// SrcDirs returns the list of subdirectories that may contain Cpp source files
func (lang Cpp) SrcDirs() []string {
	return []string{
		filepath.Join("src"),
		filepath.Join("include"),
	}
}

// TestDirs returns the list of subdirectories that may contain Cpp test files
func (lang Cpp) TestDirs() []string {
	return []string{
		filepath.Join("test"),
	}
}

// IsSrcFile returns true if the provided filename is recognized as a Cpp source file
func (lang Cpp) IsSrcFile(filename string) bool {
	extension := filepath.Ext(filename)
	switch strings.ToLower(extension) {
	case ".cpp", ".hpp":
		return true
	case ".c", ".h":
		return true
	case ".cc", ".hh":
		return true
	default:
		return false
	}
}

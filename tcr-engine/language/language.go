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
	"errors"
	"fmt"
	"path/filepath"
)

// Language is the interface that any supported language implementation must comply with
// in order to be used by TCR engine
type Language interface {
	Name() string
	SrcDirs() []string
	TestDirs() []string
	IsSrcFile(filename string) bool
}

// DetectLanguage is used to identify the language used in the provided directory. The current implementation
// simply looks at the name of the directory and checks if it matches with one of the supported languages
func DetectLanguage(baseDir string) (Language, error) {
	dir := filepath.Base(baseDir)
	switch dir {
	case "java":
		return Java{}, nil
	case "cpp":
		return Cpp{}, nil
	default:
		return nil, errors.New(fmt.Sprint("Unrecognized language: ", dir))
	}
}

// DirsToWatch returns the list of directories that TCR engine needs to watch for the provided language
func DirsToWatch(baseDir string, lang Language) []string {
	dirList := append(lang.SrcDirs(), lang.TestDirs()...)
	for i := 0; i < len(dirList); i++ {
		dirList[i] = filepath.Join(baseDir, dirList[i])
	}
	//report.PostInfo(dirList)
	return dirList
}

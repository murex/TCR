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
	"github.com/stretchr/testify/assert"

	"path/filepath"
	"testing"
)

func Test_java_language_is_supported(t *testing.T) {
	assert.True(t, isSupported("java"))
	assert.True(t, isSupported("Java"))
	assert.True(t, isSupported("JAVA"))
}

func Test_get_java_language_instance(t *testing.T) {
	language, err := getLanguage("java")
	assert.Equal(t, Java{}, language)
	assert.Zero(t, err)
}

func Test_detect_java_language(t *testing.T) {
	dirPath := filepath.Join("dummy", "java")
	language, err := detectLanguage(dirPath)
	assert.Equal(t, Java{}, language)
	assert.Zero(t, err)
}

func Test_java_language_name(t *testing.T) {
	assert.Equal(t, "java", Java{}.Name())
}

func Test_list_of_dirs_to_watch_in_java(t *testing.T) {
	var expected = []string{
		filepath.Join("src", "main"),
		filepath.Join("src", "test"),
	}
	assert.Equal(t, expected, DirsToWatch("", Java{}))
}

func Test_filenames_recognized_as_java_src(t *testing.T) {
	expected := []filenameMatching{
		{"Dummy.java", true},
		{"Dummy.JAVA", true},
		{"/dummy/Dummy.java", true},
		{"Dummy.java~", false},
		{"Dummy.java.swp", false},

		{"", false},
		{"dummy", false},
		{"Dummy.cpp", false},
		{"Dummy.sh", false},
	}
	assertFilenames(t, expected, Java{})
}

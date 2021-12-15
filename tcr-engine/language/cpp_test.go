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
	"github.com/murex/tcr/tcr-engine/toolchain"
	"github.com/stretchr/testify/assert"

	"path/filepath"
	"testing"
)

func Test_cpp_language_is_supported(t *testing.T) {
	assert.True(t, isSupported("cpp"))
	assert.True(t, isSupported("Cpp"))
	assert.True(t, isSupported("CPP"))
}

func Test_get_cpp_language_instance(t *testing.T) {
	language, err := getLanguage("cpp")
	assert.Equal(t, Cpp{}, language)
	assert.Zero(t, err)
}

func Test_detect_cpp_language(t *testing.T) {
	dirPath := filepath.Join("dummy", "cpp")
	language, err := detectLanguage(dirPath)
	assert.Equal(t, Cpp{}, language)
	assert.Zero(t, err)
}

func Test_cpp_language_name(t *testing.T) {
	assert.Equal(t, "cpp", Cpp{}.Name())
}

func Test_list_of_dirs_to_watch_in_cpp(t *testing.T) {
	var expected = []string{
		filepath.Join("src"),
		filepath.Join("include"),
		filepath.Join("test"),
	}
	assert.Equal(t, expected, DirsToWatch("", Cpp{}))
}

func Test_filenames_recognized_as_cpp_src(t *testing.T) {
	expected := []filenameMatching{
		{"Dummy.cpp", true},
		{"Dummy.CPP", true},
		{"/dummy/Dummy.cpp", true},
		{"Dummy.cpp~", false},
		{"Dummy.cpp.swp", false},

		{"Dummy.hpp", true},
		{"Dummy.HPP", true},
		{"/dummy/Dummy.hpp", true},
		{"Dummy.hpp~", false},
		{"Dummy.hpp.swp", false},

		{"Dummy.cc", true},
		{"Dummy.CC", true},
		{"/dummy/Dummy.cc", true},
		{"Dummy.cc~", false},
		{"Dummy.cc.swp", false},

		{"Dummy.hh", true},
		{"Dummy.HH", true},
		{"/dummy/Dummy.hh", true},
		{"Dummy.hh~", false},
		{"Dummy.hh.swp", false},

		{"Dummy.c", true},
		{"Dummy.C", true},
		{"/dummy/Dummy.c", true},
		{"Dummy.c~", false},
		{"Dummy.c.swp", false},

		{"Dummy.h", true},
		{"Dummy.H", true},
		{"/dummy/Dummy.h", true},
		{"Dummy.h~", false},
		{"Dummy.h.swp", false},

		{"", false},
		{"dummy", false},
		{"Dummy.java", false},
		{"Dummy.sh", false},
		{"Dummy.swp", false},
	}
	assertFilenames(t, expected, Cpp{})
}

func Test_default_toolchain_for_cpp(t *testing.T) {
	expected, _ := toolchain.Get("cmake")
	assert.Equal(t, expected, Cpp{}.defaultToolchain())
}

func Test_cpp_works_with_cmake(t *testing.T) {
	cmake, _ := toolchain.Get("cmake")
	assert.True(t, Cpp{}.worksWithToolchain(cmake))
}

func Test_cpp_does_not_work_with_gradle(t *testing.T) {
	gradle, _ := toolchain.Get("gradle")
	assert.False(t, Cpp{}.worksWithToolchain(gradle))
}

func Test_cpp_does_not_work_with_maven(t *testing.T) {
	maven, _ := toolchain.Get("maven")
	assert.False(t, Cpp{}.worksWithToolchain(maven))
}

package language

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_detect_cpp_language(t *testing.T) {
	dirPath := filepath.Join("dummy", "cpp")
	assert.Equal(t, Cpp{}, DetectLanguage(dirPath))
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

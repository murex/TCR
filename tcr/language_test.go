package tcr

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

// Any language -----------------------------------------------------------------

type FakeLanguage struct {
}

func (language FakeLanguage) name() string {
	return "fake"
}

func (language FakeLanguage) toolchain() string {
	return "fake"
}

func (language FakeLanguage) srcDirs() []string {
	return []string{"src"}
}

func (language FakeLanguage) testDirs() []string {
	return []string{"test"}
}

func (language FakeLanguage) isSrcFile(_ string) bool {
	return true
}

func Test_dirs_to_watch_should_contain_both_src_and_test_dirs(t *testing.T) {
	var expected = append(FakeLanguage{}.srcDirs(), FakeLanguage{}.testDirs()...)
	assert.Equal(t, expected, dirsToWatch("", FakeLanguage{}))
}

func Test_dirs_to_watch_should_have_absolute_path(t *testing.T) {
	baseDir, _ := os.Getwd()
	var expected = []string{
		filepath.Join(baseDir, FakeLanguage{}.srcDirs()[0]),
		filepath.Join(baseDir, FakeLanguage{}.testDirs()[0]),
	}
	assert.Equal(t, expected, dirsToWatch(baseDir, FakeLanguage{}))
}

type filenameMatching struct {
	filename string
	match    bool
}

func assertFilenames(t *testing.T, params []filenameMatching, language Language) {
	for i := range params {
		assert.Equal(t, params[i].match, language.isSrcFile(params[i].filename),
			"Filename: %v", params[i].filename)
	}
}

// Java --------------------------------------------------------------------------

func Test_list_of_dirs_to_watch_in_java(t *testing.T) {
	var expected = []string{
		filepath.Join("src", "main"),
		filepath.Join("src", "test"),
	}
	assert.Equal(t, expected, dirsToWatch("", JavaLanguage{}))
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
	assertFilenames(t, expected, JavaLanguage{})
}

// C++ --------------------------------------------------------------------------

func Test_list_of_dirs_to_watch_in_cpp(t *testing.T) {
	var expected = []string{
		filepath.Join("src"),
		filepath.Join("include"),
		filepath.Join("test"),
	}
	assert.Equal(t, expected, dirsToWatch("", CppLanguage{}))
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
	assertFilenames(t, expected, CppLanguage{})
}

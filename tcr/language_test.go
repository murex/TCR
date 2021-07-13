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

func (Language FakeLanguage) name() string {
	return "fake"
}

func (Language FakeLanguage) toolchain() string {
	return "fake"
}

func (Language FakeLanguage) srcDirs() []string {
	return []string{"src"}
}

func (Language FakeLanguage) testDirs() []string {
	return []string{"test"}
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
	filename    string
	shouldMatch bool
}

func assertFilenames(t *testing.T,
	params []filenameMatching,
	matchingFunction func(filename string) bool) {

	for i := range params {
		assert.Equal(t, params[i].shouldMatch,
			matchingFunction(params[i].filename),
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

var (
	javaFilenames = []filenameMatching{
		{"Dummy.java", true},
		{"Dummy.JAVA", true},
		{"/dummy/Dummy.java", true},
		{"Dummy.java~", false},

		{"", false},
		{"dummy", false},
		{"Dummy.cpp", false},
		{"Dummy.sh", false},
	}
)

func Test_filenames_recognized_as_java_src(t *testing.T) {
	assertFilenames(t, javaFilenames, JavaLanguage{}.matchesSrcFile)
}

func Test_filenames_recognized_as_java_test_src(t *testing.T) {
	assertFilenames(t, javaFilenames, JavaLanguage{}.matchesTestFile)
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

var (
	cppFilenames = []filenameMatching{
		{"Dummy.cpp", true},
		{"Dummy.CPP", true},
		{"/dummy/Dummy.cpp", true},
		{"Dummy.cpp~", false},

		{"Dummy.hpp", true},
		{"Dummy.HPP", true},
		{"/dummy/Dummy.hpp", true},
		{"Dummy.hpp~", false},

		{"Dummy.cc", true},
		{"Dummy.CC", true},
		{"/dummy/Dummy.cc", true},
		{"Dummy.cc~", false},

		{"Dummy.hh", true},
		{"Dummy.HH", true},
		{"/dummy/Dummy.hh", true},
		{"Dummy.hh~", false},

		{"Dummy.c", true},
		{"Dummy.C", true},
		{"/dummy/Dummy.c", true},
		{"Dummy.c~", false},

		{"Dummy.h", true},
		{"Dummy.H", true},
		{"/dummy/Dummy.h", true},
		{"Dummy.h~", false},

		{"", false},
		{"dummy", false},
		{"Dummy.java", false},
		{"Dummy.sh", false},
	}
)

func Test_filenames_recognized_as_cpp_src(t *testing.T) {
	assertFilenames(t, cppFilenames, CppLanguage{}.matchesSrcFile)
}

func Test_filenames_recognized_as_cpp_test_src(t *testing.T) {
	assertFilenames(t, cppFilenames, CppLanguage{}.matchesTestFile)
}

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

func (language FakeLanguage) Name() string {
	return "fake"
}

func (language FakeLanguage) SrcDirs() []string {
	return []string{"src"}
}

func (language FakeLanguage) TestDirs() []string {
	return []string{"test"}
}

func (language FakeLanguage) IsSrcFile(_ string) bool {
	return true
}

func Test_does_not_detect_unknown_language(t *testing.T) {
	dirPath := filepath.Join("dummy", "dummy")
	assert.Zero(t, DetectLanguage(dirPath))
}

func Test_dirs_to_watch_should_contain_both_src_and_test_dirs(t *testing.T) {
	var expected = append(FakeLanguage{}.SrcDirs(), FakeLanguage{}.TestDirs()...)
	assert.Equal(t, expected, DirsToWatch("", FakeLanguage{}))
}

func Test_dirs_to_watch_should_have_absolute_path(t *testing.T) {
	baseDir, _ := os.Getwd()
	var expected = []string{
		filepath.Join(baseDir, FakeLanguage{}.SrcDirs()[0]),
		filepath.Join(baseDir, FakeLanguage{}.TestDirs()[0]),
	}
	assert.Equal(t, expected, DirsToWatch(baseDir, FakeLanguage{}))
}

type filenameMatching struct {
	filename string
	match    bool
}

func assertFilenames(t *testing.T, params []filenameMatching, language Language) {
	for i := range params {
		assert.Equal(t, params[i].match, language.IsSrcFile(params[i].filename),
			"Filename: %v", params[i].filename)
	}
}

// Java --------------------------------------------------------------------------

func Test_detect_java_language(t *testing.T) {
	dirPath := filepath.Join("dummy", "java")
	assert.Equal(t, JavaLanguage{}, DetectLanguage(dirPath))
}

func Test_java_language_name(t *testing.T) {
	assert.Equal(t, "java", JavaLanguage{}.Name())
}

func Test_list_of_dirs_to_watch_in_java(t *testing.T) {
	var expected = []string{
		filepath.Join("src", "main"),
		filepath.Join("src", "test"),
	}
	assert.Equal(t, expected, DirsToWatch("", JavaLanguage{}))
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

func Test_detect_cpp_language(t *testing.T) {
	dirPath := filepath.Join("dummy", "cpp")
	assert.Equal(t, CppLanguage{}, DetectLanguage(dirPath))
}

func Test_cpp_language_name(t *testing.T) {
	assert.Equal(t, "cpp", CppLanguage{}.Name())
}

func Test_list_of_dirs_to_watch_in_cpp(t *testing.T) {
	var expected = []string{
		filepath.Join("src"),
		filepath.Join("include"),
		filepath.Join("test"),
	}
	assert.Equal(t, expected, DirsToWatch("", CppLanguage{}))
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

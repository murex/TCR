package tcr

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

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

// --------------------------------------------------------------------------

func Test_list_of_dirs_to_watch_in_java(t *testing.T) {
	var expected = []string{
		filepath.Join("src", "main"),
		filepath.Join("src", "test"),
	}
	assert.Equal(t, expected, dirsToWatch("", JavaLanguage{}))
}

func Test_list_of_dirs_to_watch_in_cpp(t *testing.T) {
	var expected = []string{
		filepath.Join("src"),
		filepath.Join("include"),
		filepath.Join("test"),
	}
	assert.Equal(t, expected, dirsToWatch("", CppLanguage{}))
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

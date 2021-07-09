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
	return []string{ "src" }
}

func (Language FakeLanguage) testDirs() []string {
	return []string{ "test" }
}

// --------------------------------------------------------------------------

func Test_DirsToWatch_Java(t *testing.T) {
	var expected = []string{
		filepath.Join("src", "main"),
		filepath.Join("src", "test"),
	}
	assert.Equal(t, expected, dirsToWatch("", JavaLanguage{}))
}

func Test_DirsToWatch_Cpp(t *testing.T) {
	var expected = []string{
		filepath.Join("src"),
		filepath.Join("include"),
		filepath.Join("test"),
	}
	assert.Equal(t, expected, dirsToWatch("", CppLanguage{}))
}

func Test_DirsToWatch_PrependsBaseDirToAll(t *testing.T) {
	baseDir, _ := os.Getwd()
	var expected = []string{
		filepath.Join(baseDir, "src"),
		filepath.Join(baseDir, "test"),
	}
	assert.Equal(t, expected, dirsToWatch(baseDir, FakeLanguage{}))
}
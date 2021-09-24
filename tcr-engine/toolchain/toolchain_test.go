package toolchain

import (
	"github.com/mengdaming/tcr-engine/language"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

const (
	testDataRootDir = "../testdata"
)

func testLanguageRootDir(lang language.Language) string {
	return filepath.Join(testDataRootDir, lang.Name())
}

type FakeLanguage struct {
}

func (lang FakeLanguage) Name() string {
	return "fake"
}

func (lang FakeLanguage) SrcDirs() []string {
	return []string{"src"}
}

func (lang FakeLanguage) TestDirs() []string {
	return []string{"test"}
}

func (lang FakeLanguage) IsSrcFile(_ string) bool {
	return true
}

func runFromDir(t *testing.T, testDir string, testFunction func(t *testing.T)) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	initialDir, _ := os.Getwd()
	_ = os.Chdir(testDir)
	testFunction(t)
	_ = os.Chdir(initialDir)
}

func Test_unrecognized_toolchain_name(t *testing.T) {
	toolchain, err := New("dummy", nil)
	assert.Zero(t, toolchain)
	assert.NotZero(t, err)
}

func Test_language_with_no_toolchain(t *testing.T) {
	toolchain, err := New("", FakeLanguage{})
	assert.Zero(t, toolchain)
	assert.NotZero(t, err)
}

func Test_default_toolchain_for_java(t *testing.T) {
	toolchain, err := New("", language.Java{})
	assert.Equal(t, GradleToolchain{}, toolchain)
	assert.Zero(t, err)
}

func Test_default_toolchain_for_cpp(t *testing.T) {
	toolchain, err := New("", language.Cpp{})
	assert.Equal(t, CmakeToolchain{}, toolchain)
	assert.Zero(t, err)
}

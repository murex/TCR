package toolchain

import (
	"github.com/mengdaming/tcr/tcr/language"
	"github.com/mengdaming/tcr/trace"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	// Prevent trace.Error() from triggering os.Exit()
	trace.SetTestMode()
	os.Exit(m.Run())
}

const (
	testKataRootDir = "../../test/kata"
)

func testLanguageRootDir(lang language.Language) string {
	return filepath.Join(testKataRootDir, lang.Name())
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
	initialDir, _ := os.Getwd()
	_ = os.Chdir(testDir)
	workDir, _ := os.Getwd()
	trace.Info("Working directory: ", workDir)
	testFunction(t)
	_ = os.Chdir(initialDir)
}

func Test_unrecognized_toolchain_name(t *testing.T) {
	assert.Zero(t, NewToolchain("dummy", nil))
	assert.NotZero(t, trace.GetExitReturnCode())
}

func Test_language_with_no_toolchain(t *testing.T) {
	assert.Zero(t, NewToolchain("", FakeLanguage{}))
}

func Test_default_toolchain_for_java(t *testing.T) {
	assert.Equal(t, GradleToolchain{}, NewToolchain("", language.Java{}))
}

func Test_default_toolchain_for_cpp(t *testing.T) {
	assert.Equal(t, CmakeToolchain{}, NewToolchain("", language.Cpp{}))
}


package language

import (
	"github.com/mengdaming/tcr/tcr/trace"
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

func assertFilenames(t *testing.T, params []filenameMatching, lang Language) {
	for i := range params {
		assert.Equal(t, params[i].match, lang.IsSrcFile(params[i].filename),
			"Filename: %v", params[i].filename)
	}
}

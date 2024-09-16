/*
Copyright (c) 2022 Murex

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

package filesystem

import (
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	testDataRootDir = "../testdata"
)

var (
	testDataDirJava = filepath.Join(testDataRootDir, "java")
)

func Test_init_source_tree(t *testing.T) {
	testFlags := []struct {
		desc         string
		path         string
		expectError  bool
		expectedPath func() string
	}{
		{
			"with empty path",
			"",
			false,
			func() string { path, _ := os.Getwd(); return path },
		},
		{
			"with current directory",
			".",
			false,
			func() string { path, _ := os.Getwd(); return path },
		},
		{
			"with existing directory",
			testDataDirJava,
			false,
			func() string { path, _ := filepath.Abs(testDataDirJava); return path },
		},
		{
			"with non-existing path",
			filepath.Join(testDataDirJava, "dummy-dir"),
			true,
			nil,
		},
		{
			"with existing file",
			filepath.Join(testDataDirJava, "pom.xml"),
			true,
			nil,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			tree, err := New(tt.path)
			if tt.expectError {
				assert.Error(t, err)
				assert.Zero(t, tree)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tree)
				assert.True(t, tree.IsValid())
				assert.Equal(t, tt.expectedPath(), tree.GetBaseDir())
			}
		})
	}
}

func Test_watch_can_detect_changes_on_matching_files(t *testing.T) {
	// this test fails randomly on Windows for an unknown reason.
	utils.SkipOnWindows(t)

	baseDir, _ := os.MkdirTemp("", "tcr-test-watch")
	defer func(path string) { _ = os.RemoveAll(path) }(baseDir)

	srcDir := filepath.Join(baseDir, "src")
	_ = os.Mkdir(srcDir, 0700)
	file := filepath.Join(srcDir, "file.txt")
	_ = os.WriteFile(file, []byte("some contents\n"), 0600)

	tree, _ := New(baseDir)
	stopWatching := make(chan bool)
	changeDetected := make(chan bool)
	go func() {
		changeDetected <- tree.Watch([]string{srcDir}, func(_ string) bool { return true }, stopWatching)
	}()
	time.Sleep(1 * time.Millisecond)
	_ = os.WriteFile(file, []byte("some other contents\n"), 0600)
	assert.True(t, <-changeDetected)
}

func Test_watch_ignores_changes_on_non_matching_files(t *testing.T) {
	baseDir, _ := os.MkdirTemp("", "tcr-test-watch")
	defer func(path string) { _ = os.RemoveAll(path) }(baseDir)

	srcDir := filepath.Join(baseDir, "src")
	_ = os.Mkdir(srcDir, 0700)
	file := filepath.Join(srcDir, "file.txt")
	_ = os.WriteFile(file, []byte("some contents\n"), 0600)

	tree, _ := New(baseDir)
	stopWatching := make(chan bool)
	changeDetected := make(chan bool)
	go func() {
		changeDetected <- tree.Watch([]string{srcDir}, func(_ string) bool { return false }, stopWatching)
	}()
	time.Sleep(1 * time.Millisecond)
	_ = os.WriteFile(file, []byte("some other contents\n"), 0600)
	stopWatching <- true
	assert.False(t, <-changeDetected)
}

func Test_watch_can_be_stopped_on_request(t *testing.T) {
	baseDir, _ := os.MkdirTemp("", "tcr-test-watch")
	defer func(path string) { _ = os.RemoveAll(path) }(baseDir)

	srcDir := filepath.Join(baseDir, "src")
	_ = os.Mkdir(srcDir, 0700)
	file := filepath.Join(srcDir, "file.txt")
	_ = os.WriteFile(file, []byte("some contents\n"), 0600)

	tree, _ := New(baseDir)
	stopWatching := make(chan bool)
	changeDetected := make(chan bool)
	go func() {
		changeDetected <- tree.Watch([]string{srcDir}, func(_ string) bool { return true }, stopWatching)
	}()
	time.Sleep(1 * time.Millisecond)
	stopWatching <- true
	assert.False(t, <-changeDetected)
}

func Test_watch_reports_a_warning_on_missing_directory(t *testing.T) {
	baseDir, _ := os.MkdirTemp("", "tcr-test-watch")
	defer func(path string) { _ = os.RemoveAll(path) }(baseDir)

	srcDir := filepath.Join(baseDir, "src")

	tree, _ := New(baseDir)

	sniffer := report.NewSniffer(func(msg report.Message) bool {
		return msg.Type.Category == report.Warning
	})

	stopWatching := make(chan bool)
	changeDetected := make(chan bool)
	go func() {
		changeDetected <- tree.Watch([]string{srcDir}, func(_ string) bool { return false }, stopWatching)
	}()
	time.Sleep(1 * time.Millisecond)
	stopWatching <- true
	<-changeDetected
	sniffer.Stop()

	assert.Equal(t, 1, sniffer.GetMatchCount())
	assert.Equal(t, report.Warning, sniffer.GetAllMatches()[0].Type.Category)
	assert.NotEmpty(t, sniffer.GetAllMatches()[0].Payload.ToString())
}

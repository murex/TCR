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
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
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
			filepath.Join(testDataDirJava, "Makefile"),
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
	// TODO figure out why this is failing on Unix OS's
	if runtime.GOOS != "windows" {
		t.Skip("needs to be fixed on Unix OS's")
	}
	baseDir, _ := ioutil.TempDir("", "tcr-test-watch")
	defer func(path string) { _ = os.RemoveAll(path) }(baseDir)

	srcDir := filepath.Join(baseDir, "src")
	_ = os.Mkdir(srcDir, os.ModeDir)
	file := filepath.Join(srcDir, "file.txt")
	_ = os.WriteFile(file, []byte("some contents\n"), 0600)

	tree, _ := New(baseDir)
	stopWatching := make(chan bool)
	go func() {
		caughtFileUpdate := tree.Watch([]string{srcDir}, func(_ string) bool { return true }, stopWatching)
		time.Sleep(1 * time.Second)
		assert.True(t, caughtFileUpdate)
	}()

	_ = os.WriteFile(file, []byte("some other contents\n"), 0600)
}

func Test_watch_ignores_changes_on_non_matching_files(t *testing.T) {
	// TODO figure out why this is failing on Unix OS's
	if runtime.GOOS != "windows" {
		t.Skip("needs to be fixed on Unix OS's")
	}
	baseDir, _ := ioutil.TempDir("", "tcr-test-watch")
	defer func(path string) { _ = os.RemoveAll(path) }(baseDir)

	srcDir := filepath.Join(baseDir, "src")
	_ = os.Mkdir(srcDir, os.ModeDir)
	file := filepath.Join(srcDir, "file.txt")
	_ = os.WriteFile(file, []byte("some contents\n"), 0600)

	tree, _ := New(baseDir)
	stopWatching := make(chan bool)
	go func() {
		caughtFileUpdate := tree.Watch([]string{srcDir}, func(_ string) bool { return false }, stopWatching)
		time.Sleep(1 * time.Second)
		assert.False(t, caughtFileUpdate)
	}()

	_ = os.WriteFile(file, []byte("some other contents\n"), 0600)
	stopWatching <- true
}

func Test_watch_can_be_stopped_on_request(t *testing.T) {
	// TODO figure out why this is failing on Unix OS's
	if runtime.GOOS != "windows" {
		t.Skip("needs to be fixed on Unix OS's")
	}
	baseDir, _ := ioutil.TempDir("", "tcr-test-watch")
	defer func(path string) { _ = os.RemoveAll(path) }(baseDir)

	srcDir := filepath.Join(baseDir, "src")
	_ = os.Mkdir(srcDir, os.ModeDir)
	file := filepath.Join(srcDir, "file.txt")
	_ = os.WriteFile(file, []byte("some contents\n"), 0600)

	tree, _ := New(baseDir)
	stopWatching := make(chan bool)
	go func() {
		caughtFileUpdate := tree.Watch([]string{srcDir}, func(_ string) bool { return true }, stopWatching)
		time.Sleep(1 * time.Millisecond)
		assert.False(t, caughtFileUpdate)
	}()

	stopWatching <- true
}

func Test_watch_exits_if_missing_directory(t *testing.T) {
	// TODO figure out why this is failing on Unix OS's
	if runtime.GOOS != "windows" {
		t.Skip("needs to be fixed on Unix OS's")
	}
	baseDir, _ := ioutil.TempDir("", "tcr-test-watch")
	defer func(path string) { _ = os.RemoveAll(path) }(baseDir)

	srcDir := filepath.Join(baseDir, "src")

	tree, _ := New(baseDir)
	stopWatching := make(chan bool)
	go func() {
		caughtFileUpdate := tree.Watch([]string{srcDir}, func(_ string) bool { return true }, stopWatching)
		time.Sleep(1 * time.Millisecond)
		assert.False(t, caughtFileUpdate)
	}()
}

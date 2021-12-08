/*
Copyright (c) 2021 Murex

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

package language

import (
	"github.com/murex/tcr/tcr-engine/toolchain"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

type FakeLanguage struct {
}

func (lang FakeLanguage) GetToolchain(t string) (toolchain.Toolchain, error) {
	return nil, nil
}

func (lang FakeLanguage) defaultToolchain() toolchain.Toolchain {
	return nil
}

func (lang FakeLanguage) worksWithToolchain(_ toolchain.Toolchain) bool {
	return false
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

func Test_does_not_support_empty_language_name(t *testing.T) {
	assert.False(t, isSupported(""))
}

func Test_does_not_support_dummy_language_name(t *testing.T) {
	assert.False(t, isSupported("dummy"))
}

func Test_does_not_set_unknown_language_name(t *testing.T) {
	language, err := New("dummy", "")
	assert.Zero(t, language)
	assert.NotZero(t, err)
}

func Test_fallbacks_on_dir_name_if_language_is_not_specified(t *testing.T) {
	language, err := New("", Java{}.Name())
	assert.Equal(t, Java{}, language)
	assert.Zero(t, err)
}

func Test_does_not_detect_unknown_language(t *testing.T) {
	dirPath := filepath.Join("dummy", "dummy")
	language, err := detectLanguage(dirPath)
	assert.Zero(t, language)
	assert.NotZero(t, err)
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

// TODO
//func Test_language_with_no_toolchain(t *testing.T) {
//	toolchain, err := New("")
//	assert.Zero(t, toolchain)
//	assert.NotZero(t, err)
//}

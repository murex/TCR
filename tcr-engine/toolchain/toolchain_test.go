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

package toolchain

import (
	"github.com/murex/tcr-engine/language"
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

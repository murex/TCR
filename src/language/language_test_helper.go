//go:build test_helper

/*
Copyright (c) 2023 Murex

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
	"testing"

	"github.com/murex/tcr/language/built_in_test_data"
	"github.com/stretchr/testify/assert"
)

func assertDefaultToolchain(t *testing.T, languageName string, toolchainName string) {
	t.Helper()
	lang, err := Get(languageName)
	assert.NoError(t, err)
	if assert.NotNil(t, lang) {
		assert.Equal(t, toolchainName, lang.GetToolchains().Default)
	}
}

func assertListOfDirsToWatch(t *testing.T, languageName string, dirs ...string) {
	t.Helper()
	dirList := getBuiltIn(languageName).DirsToWatch("")
	for _, dir := range dirs {
		t.Run(dir, func(t *testing.T) {
			assert.Contains(t, dirList, toLocalPath(dir))
		})
	}
}

func assertCompatibleToolchains(t *testing.T, languageName string, toolchainNames ...string) {
	t.Helper()
	lang := getBuiltIn(languageName)
	for _, tchn := range toolchainNames {
		t.Run(tchn, func(t *testing.T) {
			assert.True(t, lang.worksWithToolchain(tchn))
		})
	}
}

func assertIncompatibleToolchains(t *testing.T, languageName string, toolchainNames ...string) {
	t.Helper()
	lang := getBuiltIn(languageName)
	for _, tchn := range toolchainNames {
		t.Run(tchn, func(t *testing.T) {
			assert.False(t, lang.worksWithToolchain(tchn))
		})
	}
}

func assertFilePathsMatching(t *testing.T, languageName string, matchers ...built_in_test_data.FilePathMatcher) {
	t.Helper()
	lang, err := GetLanguage(languageName, "")
	assert.NoError(t, err)
	if assert.NotNil(t, lang) {
		for _, matcher := range matchers {
			assert.Equal(t, matcher.IsSrcFile, lang.IsSrcFile(matcher.FilePath),
				"Should %v be a source file?", matcher.FilePath)
			assert.Equal(t, matcher.IsTestFile, lang.IsTestFile(matcher.FilePath),
				"Should %v be a test file?", matcher.FilePath)
		}
	}
}

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

package language

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

func Test_does_not_support_empty_language_name(t *testing.T) {
	assert.False(t, isSupported(""))
}

func Test_does_not_support_unregistered_language_name(t *testing.T) {
	assert.False(t, isSupported("unregistered-language"))
}

func Test_unrecognized_language_name(t *testing.T) {
	lang, err := Get("dummy-language")
	assert.Error(t, err)
	assert.Zero(t, lang)
}

func Test_can_add_a_built_in_language(t *testing.T) {
	const name = "new-built-in-language"
	assert.False(t, isBuiltIn(name))
	assert.NoError(t, addBuiltIn(ALanguage(WithName(name))))
	assert.True(t, isBuiltIn(name))
}

func Test_cannot_add_a_built_in_language_with_no_name(t *testing.T) {
	assert.Error(t, addBuiltIn(ALanguage(WithName(""))))
}

func Test_language_name_is_case_insensitive(t *testing.T) {
	const name = "miXeD-CasE"
	_ = Register(ALanguage(WithName(name)))
	assertNameIsNotCaseSensitive(t, name)
}

func Test_can_register_a_language(t *testing.T) {
	const name = "new-language"
	assert.False(t, isSupported(name))
	assert.NoError(t, Register(ALanguage(
		WithName(name),
		WithDefaultToolchain("default-toolchain"),
		WithCompatibleToolchain("default-toolchain"),
	)))
	assert.True(t, isSupported(name))
}

func Test_cannot_register_a_language_with_no_name(t *testing.T) {
	assert.Error(t, Register(ALanguage(WithName(""))))
}

func Test_cannot_register_a_language_with_no_compatible_toolchain(t *testing.T) {
	const name = "no-compatible-toolchain"
	assert.Error(t, Register(ALanguage(WithName(name), WithNoCompatibleToolchain())))
	assert.False(t, isSupported(name))
}

func Test_cannot_register_a_language_with_no_default_toolchain(t *testing.T) {
	const name = "no-default-toolchain"
	assert.Error(t, Register(ALanguage(WithName(name), WithNoDefaultToolchain())))
	assert.False(t, isSupported(name))
}

func Test_cannot_register_a_language_with_default_toolchain_not_compatible(t *testing.T) {
	const name = "default-toolchain-not-compatible"
	assert.Error(t, Register(ALanguage(
		WithName(name),
		WithDefaultToolchain("toolchain1"),
		WithCompatibleToolchain("toolchain2"),
	)))
	assert.False(t, isSupported(name))
}

func Test_does_not_detect_language_from_a_dir_name_not_matching_a_known_language(t *testing.T) {
	lang, err := detectLanguageFromDirName("dummy")
	assert.Error(t, err)
	assert.Zero(t, lang)
}

// Assert utility functions for language register

func assertIsABuiltInLanguage(t *testing.T, name string) {
	assert.True(t, isBuiltIn(name))
}

func assertIsSupported(t *testing.T, name string) {
	assert.True(t, isSupported(name))
}

func assertNameIsNotCaseSensitive(t *testing.T, name string) {
	assert.True(t, isSupported(name))
	assert.True(t, isSupported(strings.ToUpper(name)))
	assert.True(t, isSupported(strings.ToLower(name)))
	assert.True(t, isSupported(strings.Title(name)))
}

func assertLanguageInitialization(t *testing.T, name string) {
	lang, err := Get(name)
	assert.NoError(t, err)
	assert.Equal(t, name, lang.GetName())
}

func assertLanguageName(t *testing.T, name string) {
	lang, _ := Get(name)
	assert.Equal(t, name, lang.GetName())
}

func assertFallbacksOnDirNameIfLanguageIsNotSpecified(t *testing.T, dirName string) {
	lang, err := GetLanguage("", filepath.Join("some", "path", dirName))
	assert.NoError(t, err)
	assert.NotEmpty(t, lang)
}

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

import "testing"

const (
	kotlinLanguageName = "kotlin"
)

var (
	kotlinLanguageExtensions = []string{".kt"}
)

func init() {
	registerLanguageFileExtensionsForTests(kotlinLanguageExtensions...)
}

func Test_kotlin_is_a_built_in_language(t *testing.T) {
	assertIsABuiltInLanguage(t, kotlinLanguageName)
}

func Test_kotlin_language_is_supported(t *testing.T) {
	assertIsSupported(t, kotlinLanguageName)
}

func Test_kotlin_language_is_registered(t *testing.T) {
	assertIsRegistered(t, kotlinLanguageName)
}

func Test_kotlin_language_name(t *testing.T) {
	assertLanguageName(t, kotlinLanguageName)
}

func Test_kotlin_language_name_is_case_insensitive(t *testing.T) {
	assertNameIsNotCaseSensitive(t, kotlinLanguageName)
}

func Test_kotlin_language_initialization(t *testing.T) {
	assertLanguageInitialization(t, kotlinLanguageName)
}

func Test_fallbacks_on_kotlin_dir_name_if_language_is_not_specified(t *testing.T) {
	assertFallbacksOnDirNameIfLanguageIsNotSpecified(t, kotlinLanguageName)
}

func Test_list_of_dirs_to_watch_in_kotlin(t *testing.T) {
	assertListOfDirsToWatch(t, []string{"src/main", "src/test"}, kotlinLanguageName)
}

func Test_kotlin_default_toolchain(t *testing.T) {
	assertDefaultToolchain(t, "gradle-wrapper", kotlinLanguageName)
}

func Test_kotlin_compatible_toolchains(t *testing.T) {
	assertCompatibleToolchains(t, []string{"gradle", "gradle-wrapper", "maven", "maven-wrapper"}, kotlinLanguageName)
	assertCompatibleToolchains(t, []string{"make"}, kotlinLanguageName)
}

func Test_kotlin_incompatible_toolchains(t *testing.T) {
	assertIncompatibleToolchains(t, []string{"cmake"}, kotlinLanguageName)
	assertIncompatibleToolchains(t, []string{"go-tools", "gotestsum"}, kotlinLanguageName)
	assertIncompatibleToolchains(t, []string{"dotnet"}, kotlinLanguageName)
}

func Test_kotlin_valid_file_paths(t *testing.T) {
	languageName := kotlinLanguageName
	for _, ext := range kotlinLanguageExtensions {
		assertFilePathsMatching(t, buildFilePathMatchers(shouldMatchSrc, "src/main", "SomeSrcFile", ext), languageName)
		assertFilePathsMatching(t, buildFilePathMatchers(shouldMatchTest, "src/test", "SomeTestFile", ext), languageName)
	}
}

func Test_kotlin_invalid_file_paths(t *testing.T) {
	languageName := kotlinLanguageName
	for _, ext := range allLanguageFileExtensionsBut(kotlinLanguageExtensions...) {
		assertFilePathsMatching(t, buildFilePathMatchers(shouldNotMatch, "src/main", "SomeSrcFile", ext), languageName)
		assertFilePathsMatching(t, buildFilePathMatchers(shouldNotMatch, "src/test", "SomeTestFile", ext), languageName)
	}
}

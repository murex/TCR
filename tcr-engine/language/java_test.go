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

import "testing"

const (
	javaLanguageName = "java"
)

var (
	javaLanguageExtensions = []string{".java"}
)

func init() {
	registerLanguageFileExtensionsForTests(javaLanguageExtensions...)
}

func Test_java_is_a_built_in_language(t *testing.T) {
	assertIsABuiltInLanguage(t, javaLanguageName)
}

func Test_java_language_is_supported(t *testing.T) {
	assertIsSupported(t, javaLanguageName)
}

func Test_java_language_is_registered(t *testing.T) {
	assertIsRegistered(t, javaLanguageName)
}

func Test_java_language_name(t *testing.T) {
	assertLanguageName(t, javaLanguageName)
}

func Test_java_language_name_is_case_insensitive(t *testing.T) {
	assertNameIsNotCaseSensitive(t, javaLanguageName)
}

func Test_java_language_initialization(t *testing.T) {
	assertLanguageInitialization(t, javaLanguageName)
}

func Test_fallbacks_on_java_dir_name_if_language_is_not_specified(t *testing.T) {
	assertFallbacksOnDirNameIfLanguageIsNotSpecified(t, javaLanguageName)
}

func Test_list_of_dirs_to_watch_in_java(t *testing.T) {
	assertListOfDirsToWatch(t, []string{"src/main", "src/test"}, javaLanguageName)
}

func Test_java_default_toolchain(t *testing.T) {
	assertDefaultToolchain(t, "gradle-wrapper", javaLanguageName)
}

func Test_java_compatible_toolchains(t *testing.T) {
	assertCompatibleToolchains(t, []string{"gradle", "gradle-wrapper", "maven", "maven-wrapper"}, javaLanguageName)
	assertCompatibleToolchains(t, []string{"make"}, javaLanguageName)
}

func Test_java_incompatible_toolchains(t *testing.T) {
	assertIncompatibleToolchains(t, []string{"cmake"}, javaLanguageName)
	assertIncompatibleToolchains(t, []string{"go-tools", "gotestsum"}, javaLanguageName)
	assertIncompatibleToolchains(t, []string{"dotnet"}, javaLanguageName)
}

func Test_java_valid_file_paths(t *testing.T) {
	languageName := javaLanguageName
	for _, ext := range javaLanguageExtensions {
		assertFilePathsMatching(t, buildFilePathMatchers(shouldMatchSrc, "src/main", "SomeSrcFile", ext), languageName)
		assertFilePathsMatching(t, buildFilePathMatchers(shouldMatchTest, "src/test", "SomeTestFile", ext), languageName)
	}
}

func Test_java_invalid_file_paths(t *testing.T) {
	languageName := javaLanguageName
	for _, ext := range allLanguageFileExtensionsBut(javaLanguageExtensions...) {
		assertFilePathsMatching(t, buildFilePathMatchers(shouldNotMatch, "src/main", "SomeSrcFile", ext), languageName)
		assertFilePathsMatching(t, buildFilePathMatchers(shouldNotMatch, "src/test", "SomeTestFile", ext), languageName)
	}
}

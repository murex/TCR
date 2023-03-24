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
	goLanguageName = "go"
)

var (
	goLanguageExtensions   = []string{".go"}
	goCompatibleToolchains = []string{"gotestsum", "go-tools", "make"}
)

func init() {
	registerLanguageFileExtensionsForTests(goLanguageExtensions...)
	registerToolchainsForTests(goCompatibleToolchains...)
}

func Test_go_is_a_built_in_language(t *testing.T) {
	assertIsABuiltInLanguage(t, goLanguageName)
}

func Test_go_language_is_supported(t *testing.T) {
	assertIsSupported(t, goLanguageName)
}

func Test_go_language_is_registered(t *testing.T) {
	assertIsRegistered(t, goLanguageName)
}

func Test_go_language_name(t *testing.T) {
	assertLanguageName(t, goLanguageName)
}

func Test_go_language_name_is_case_insensitive(t *testing.T) {
	assertNameIsNotCaseSensitive(t, goLanguageName)
}

func Test_go_language_initialization(t *testing.T) {
	assertLanguageInitialization(t, goLanguageName)
}

func Test_fallbacks_on_go_dir_name_if_language_is_not_specified(t *testing.T) {
	assertFallbacksOnDirNameIfLanguageIsNotSpecified(t, goLanguageName)
}

func Test_list_of_dirs_to_watch_in_go(t *testing.T) {
	assertListOfDirsToWatch(t, goLanguageName, ".")
}

func Test_go_default_toolchain(t *testing.T) {
	assertDefaultToolchain(t, goLanguageName, "go-tools")
}

func Test_go_compatible_toolchains(t *testing.T) {
	assertCompatibleToolchains(t, goLanguageName, goCompatibleToolchains...)
}

func Test_go_incompatible_toolchains(t *testing.T) {
	assertIncompatibleToolchains(t, goLanguageName, allKnownToolchainsBut(goCompatibleToolchains...)...)
}

func Test_go_valid_file_paths(t *testing.T) {
	languageName := goLanguageName
	for _, ext := range goLanguageExtensions {
		assertFilePathsMatching(t, languageName,
			buildFilePathMatchers(shouldMatchSrc, ".", "some_src_file", ext)...)
		assertFilePathsMatching(t, languageName,
			buildFilePathMatchers(shouldMatchTest, ".", "some_src_file_test", ext)...)
	}
}

func Test_go_invalid_file_paths(t *testing.T) {
	languageName := goLanguageName
	for _, ext := range allLanguageFileExtensionsBut(goLanguageExtensions...) {
		assertFilePathsMatching(t, languageName,
			buildFilePathMatchers(shouldNotMatch, ".", "some_src_file", ext)...)
		assertFilePathsMatching(t, languageName,
			buildFilePathMatchers(shouldNotMatch, ".", "some_src_file_test", ext)...)
	}
}

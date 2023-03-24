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

import "testing"

const (
	pythonLanguageName = "python"
)

var (
	pythonLanguageExtensions   = []string{".py"}
	pythonCompatibleToolchains = []string{"pytest", "make"}
)

func init() {
	registerLanguageFileExtensionsForTests(pythonLanguageExtensions...)
	registerToolchainsForTests(pythonCompatibleToolchains...)
}

func Test_python_is_a_built_in_language(t *testing.T) {
	assertIsABuiltInLanguage(t, pythonLanguageName)
}

func Test_python_language_is_supported(t *testing.T) {
	assertIsSupported(t, pythonLanguageName)
}

func Test_python_language_is_registered(t *testing.T) {
	assertIsRegistered(t, pythonLanguageName)
}

func Test_python_language_name(t *testing.T) {
	assertLanguageName(t, pythonLanguageName)
}

func Test_python_language_name_is_case_insensitive(t *testing.T) {
	assertNameIsNotCaseSensitive(t, pythonLanguageName)
}

func Test_python_language_initialization(t *testing.T) {
	assertLanguageInitialization(t, pythonLanguageName)
}

func Test_fallbacks_on_python_dir_name_if_language_is_not_specified(t *testing.T) {
	assertFallbacksOnDirNameIfLanguageIsNotSpecified(t, pythonLanguageName)
}

func Test_list_of_dirs_to_watch_in_python(t *testing.T) {
	assertListOfDirsToWatch(t, pythonLanguageName, "src", "tests")
}

func Test_python_default_toolchain(t *testing.T) {
	assertDefaultToolchain(t, pythonLanguageName, "pytest")
}

func Test_python_compatible_toolchains(t *testing.T) {
	assertCompatibleToolchains(t, pythonLanguageName, pythonCompatibleToolchains...)
}

func Test_python_incompatible_toolchains(t *testing.T) {
	assertIncompatibleToolchains(t, pythonLanguageName, allKnownToolchainsBut(pythonCompatibleToolchains...)...)
}

func Test_python_valid_file_paths(t *testing.T) {
	languageName := pythonLanguageName
	for _, ext := range pythonLanguageExtensions {
		assertFilePathsMatching(t, languageName,
			buildFilePathMatchers(shouldMatchSrc, "src", "some_src_file", ext)...)
		assertFilePathsMatching(t, languageName,
			buildFilePathMatchers(shouldMatchTest, "tests", "some_src_file_test", ext)...)
	}
}

func Test_python_invalid_file_paths(t *testing.T) {
	languageName := pythonLanguageName
	for _, ext := range allLanguageFileExtensionsBut(pythonLanguageExtensions...) {
		assertFilePathsMatching(t, languageName,
			buildFilePathMatchers(shouldNotMatch, "src", "some_src_file", ext)...)
		assertFilePathsMatching(t, languageName,
			buildFilePathMatchers(shouldNotMatch, "tests", "some_src_file_test", ext)...)
	}
}

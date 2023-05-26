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
	csharpLanguageName = "csharp"
)

var (
	csharpLanguageExtensions   = []string{".cs", ".csx"}
	csharpCompatibleToolchains = []string{"bazel", "dotnet", "make"}
)

func init() {
	registerLanguageFileExtensionsForTests(csharpLanguageExtensions...)
	registerToolchainsForTests(csharpCompatibleToolchains...)
}

func Test_csharp_is_a_built_in_language(t *testing.T) {
	assertIsABuiltInLanguage(t, csharpLanguageName)
}

func Test_csharp_language_is_supported(t *testing.T) {
	assertIsSupported(t, csharpLanguageName)
}

func Test_csharp_language_is_registered(t *testing.T) {
	assertIsRegistered(t, csharpLanguageName)
}

func Test_csharp_language_name(t *testing.T) {
	assertLanguageName(t, csharpLanguageName)
}

func Test_csharp_language_name_is_case_insensitive(t *testing.T) {
	assertNameIsNotCaseSensitive(t, csharpLanguageName)
}

func Test_csharp_language_initialization(t *testing.T) {
	assertLanguageInitialization(t, csharpLanguageName)
}

func Test_fallbacks_on_csharp_dir_name_if_language_is_not_specified(t *testing.T) {
	assertFallbacksOnDirNameIfLanguageIsNotSpecified(t, csharpLanguageName)
}

func Test_list_of_dirs_to_watch_in_csharp(t *testing.T) {
	assertListOfDirsToWatch(t, csharpLanguageName, "src", "tests")
}

func Test_csharp_default_toolchain(t *testing.T) {
	assertDefaultToolchain(t, csharpLanguageName, "dotnet")
}

func Test_csharp_compatible_toolchains(t *testing.T) {
	assertCompatibleToolchains(t, csharpLanguageName, csharpCompatibleToolchains...)
}

func Test_csharp_incompatible_toolchains(t *testing.T) {
	assertIncompatibleToolchains(t, csharpLanguageName, allKnownToolchainsBut(csharpCompatibleToolchains...)...)
}

func Test_csharp_valid_file_paths(t *testing.T) {
	languageName := csharpLanguageName
	for _, ext := range csharpLanguageExtensions {
		assertFilePathsMatching(t, languageName,
			buildFilePathMatchers(shouldMatchSrc, "src", "SomeSrcFile", ext)...)
		assertFilePathsMatching(t, languageName,
			buildFilePathMatchers(shouldMatchTest, "tests", "SomeTestFile", ext)...)
	}
}

func Test_csharp_invalid_file_paths(t *testing.T) {
	languageName := csharpLanguageName
	for _, ext := range allLanguageFileExtensionsBut(csharpLanguageExtensions...) {
		assertFilePathsMatching(t, languageName,
			buildFilePathMatchers(shouldNotMatch, "src", "SomeSrcFile", ext)...)
		assertFilePathsMatching(t, languageName,
			buildFilePathMatchers(shouldNotMatch, "tests", "SomeTestFile", ext)...)
	}
}

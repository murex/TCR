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
	"testing"
)

const (
	cppLanguageName = "cpp"
)

func Test_cpp_is_a_built_in_language(t *testing.T) {
	assertIsABuiltInLanguage(t, cppLanguageName)
}

func Test_cpp_language_is_supported(t *testing.T) {
	assertIsSupported(t, cppLanguageName)
}

func Test_cpp_language_name(t *testing.T) {
	assertLanguageName(t, cppLanguageName)
}

func Test_cpp_language_name_is_case_insensitive(t *testing.T) {
	assertNameIsNotCaseSensitive(t, cppLanguageName)
}

func Test_cpp_language_initialization(t *testing.T) {
	assertLanguageInitialization(t, cppLanguageName)
}

func Test_fallbacks_on_cpp_dir_name_if_language_is_not_specified(t *testing.T) {
	assertFallbacksOnDirNameIfLanguageIsNotSpecified(t, cppLanguageName)
}

func Test_list_of_dirs_to_watch_in_cpp(t *testing.T) {
	assertListOfDirsToWatch(t, []string{"src", "include", "test"}, cppLanguageName)
}

func Test_cpp_default_toolchain(t *testing.T) {
	assertDefaultToolchain(t, "cmake", cppLanguageName)
}

func Test_cpp_compatible_toolchains(t *testing.T) {
	assertCompatibleToolchains(t, []string{"cmake", "cmake-local"}, cppLanguageName)
}

func Test_cpp_incompatible_toolchains(t *testing.T) {
	assertIncompatibleToolchains(t, []string{"gradle", "gradle-wrapper", "maven", "maven-wrapper"}, cppLanguageName)
	assertIncompatibleToolchains(t, []string{"go-tools"}, cppLanguageName)
}

func Test_file_paths_recognized_as_cpp_src(t *testing.T) {
	expected := []filePathMatcher{
		shouldMatchSrc("src/Dummy.cpp"),
		shouldMatchSrc("src/Dummy.CPP"),
		shouldMatchSrc("src/dummy/Dummy.cpp"),

		shouldNotMatch("src/Dummy.cpp~"),
		shouldNotMatch("src/Dummy.cpp.swp"),

		shouldMatchSrc("src/Dummy.hpp"),
		shouldMatchSrc("src/Dummy.HPP"),
		shouldMatchSrc("src/dummy/Dummy.hpp"),

		shouldNotMatch("src/Dummy.hpp~"),
		shouldNotMatch("src/Dummy.hpp.swp"),

		shouldMatchSrc("src/Dummy.cc"),
		shouldMatchSrc("src/Dummy.CC"),
		shouldMatchSrc("src/dummy/Dummy.cc"),

		shouldNotMatch("src/Dummy.cc~"),
		shouldNotMatch("src/Dummy.cc.swp"),

		shouldMatchSrc("src/Dummy.hh"),
		shouldMatchSrc("src/Dummy.HH"),
		shouldMatchSrc("src/dummy/Dummy.hh"),

		shouldNotMatch("src/Dummy.hh~"),
		shouldNotMatch("src/Dummy.hh.swp"),

		shouldMatchSrc("src/Dummy.c"),
		shouldMatchSrc("src/Dummy.C"),
		shouldMatchSrc("src//dummy/Dummy.c"),

		shouldNotMatch("src/Dummy.c~"),
		shouldNotMatch("src/Dummy.c.swp"),

		shouldMatchSrc("src/Dummy.h"),
		shouldMatchSrc("src/Dummy.H"),
		shouldMatchSrc("src/dummy/Dummy.h"),

		shouldNotMatch("src/Dummy.h~"),
		shouldNotMatch("src/Dummy.h.swp"),

		shouldNotMatch("src"),
		shouldNotMatch("src/dummy"),
		shouldNotMatch("src/dummy.java"),
		shouldNotMatch("src/dummy.go"),

		shouldNotMatch("src/Dummy.sh"),
		shouldNotMatch("src/Dummy.swp"),

		shouldMatchSrc("include/Dummy.hpp"),

		shouldMatchTest("test/Dummy.hpp"),
	}
	assertFilePathsMatching(t, expected, cppLanguageName)
}

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

	. "github.com/murex/tcr/language/built_in_test_data"
)

func Test_is_a_built_in_language(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertIsABuiltInLanguage(t, builtIn.Name)
		})
	}
}

func Test_built_in_language_is_supported(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertIsSupported(t, builtIn.Name)
		})
	}
}

func Test_built_in_language_is_registered(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertIsRegistered(t, builtIn.Name)
		})
	}
}

func Test_built_in_language_name(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertLanguageName(t, builtIn.Name)
		})
	}
}

func Test_built_in_language_name_is_case_insensitive(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertNameIsNotCaseSensitive(t, builtIn.Name)
		})
	}
}

func Test_built_in_language_initialization(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertLanguageInitialization(t, builtIn.Name)
		})
	}
}

func Test_fallbacks_on_built_in_language_dir_name_if_language_is_not_specified(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertFallbacksOnDirNameIfLanguageIsNotSpecified(t, builtIn.Name)
		})
	}
}

func Test_list_of_dirs_to_watch_in_built_in_language(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertListOfDirsToWatch(t, builtIn.Name, builtIn.DirsToWatch...)
		})
	}
}

func Test_built_in_language_default_toolchain(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertDefaultToolchain(t, builtIn.Name, builtIn.DefaultToolchain)
		})
	}
}

func Test_built_in_language_compatible_toolchains(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertCompatibleToolchains(t, builtIn.Name, builtIn.CompatibleToolchains...)
		})
	}
}

func Test_built_in_language_incompatible_toolchains(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			assertIncompatibleToolchains(t, builtIn.Name, AllKnownToolchainsBut(builtIn.CompatibleToolchains...)...)
		})
	}
}

func Test_built_in_language_valid_src_file_paths(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			for _, ext := range builtIn.Extensions {
				for _, def := range builtIn.SrcMatchersDefs {
					assertFilePathsMatching(t, builtIn.Name,
						BuildFilePathMatchers(ShouldMatchSrc, def, ext)...)
				}
			}
		})
	}
}

func Test_built_in_language_invalid_src_file_paths(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			for _, ext := range AllKnownFileExtensionsBut(builtIn.Extensions...) {
				for _, def := range builtIn.SrcMatchersDefs {
					assertFilePathsMatching(t, builtIn.Name,
						BuildFilePathMatchers(ShouldNotMatch, def, ext)...)
				}
			}
		})
	}
}

func Test_built_in_language_valid_test_file_paths(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			for _, ext := range builtIn.Extensions {
				for _, def := range builtIn.TestMatcherDefs {
					assertFilePathsMatching(t, builtIn.Name,
						BuildFilePathMatchers(ShouldMatchTest, def, ext)...)
				}
			}
		})
	}
}

func Test_built_in_language_invalid_test_file_paths(t *testing.T) {
	for _, builtIn := range BuiltInTests {
		t.Run(builtIn.Name, func(t *testing.T) {
			for _, ext := range AllKnownFileExtensionsBut(builtIn.Extensions...) {
				for _, def := range builtIn.TestMatcherDefs {
					assertFilePathsMatching(t, builtIn.Name,
						BuildFilePathMatchers(ShouldNotMatch, def, ext)...)
				}
			}
		})
	}
}

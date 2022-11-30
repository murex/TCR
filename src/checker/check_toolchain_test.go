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

package checker

import (
	"github.com/murex/tcr/language"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/toolchain"
	"testing"
)

func Test_check_toolchain_returns_ok_when_set_and_recognized(t *testing.T) {
	assertOk(t, checkToolchain,
		*params.AParamSet(
			params.WithToolchain("make"),
			params.WithLanguage("go"),
		),
	)
}

func Test_check_toolchain_returns_error_when_set_but_unknown(t *testing.T) {
	assertError(t, checkToolchain,
		*params.AParamSet(
			params.WithToolchain("unknown-toolchain"),
		),
	)
}

func Test_check_toolchain_returns_error_when_set_but_language_incompatible(t *testing.T) {
	assertError(t, checkToolchain,
		*params.AParamSet(
			params.WithToolchain("gradlew"),
			params.WithLanguage("go"),
		),
	)
}

func Test_check_toolchain_returns_ok_when_deduced_from_language(t *testing.T) {
	assertOk(t, checkToolchain,
		*params.AParamSet(
			params.WithToolchain(""),
			params.WithLanguage("go"),
		),
	)
}

func Test_check_toolchain_returns_error_when_cannot_be_deduced_from_language(t *testing.T) {
	assertError(t, checkToolchain,
		*params.AParamSet(
			params.WithToolchain(""),
			params.WithLanguage("unknown-language"),
		),
	)
}

func Test_check_toolchain_returns_ok_when_deduced_from_base_dir_name(t *testing.T) {
	assertOk(t, checkToolchain,
		*params.AParamSet(
			params.WithToolchain(""),
			params.WithLanguage(""),
			params.WithWorkDir(testDataDirJava),
			params.WithBaseDir(testDataDirJava),
		),
	)
}

func Test_check_toolchain_returns_error_when_cannot_be_deduced_from_language_nor_base_dir_name(t *testing.T) {
	assertError(t, checkToolchain,
		*params.AParamSet(
			params.WithToolchain(""),
			params.WithLanguage(""),
			params.WithBaseDir(""),
		),
	)
}

func Test_check_toolchain_return_value_for_test_result_dir(t *testing.T) {
	testFlags := []struct {
		desc     string
		dir      string
		expected CheckStatus
	}{
		{"when not set", "", CheckStatusWarning},
		{"when set", ".", CheckStatusOk},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			okCmd := toolchain.ACommand(toolchain.WithPath("true"))
			tchn := toolchain.AToolchain(
				toolchain.WithTestResultDir(tt.dir),
				toolchain.WithNoBuildCommand(), toolchain.WithBuildCommand(okCmd),
				toolchain.WithNoTestCommand(), toolchain.WithTestCommand(okCmd),
			)
			_ = toolchain.Register(tchn)
			lang := language.ALanguage(language.WithDefaultToolchain(tchn.GetName()))
			_ = language.Register(lang)

			assertStatus(t, tt.expected, checkToolchain,
				*params.AParamSet(
					params.WithToolchain(tchn.GetName()),
					params.WithLanguage(lang.GetName()),
				),
			)
		})
	}
}

// TODO test that the build and test commands can be found

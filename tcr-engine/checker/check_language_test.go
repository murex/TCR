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
	"github.com/murex/tcr/params"
	"testing"
)

func Test_check_language_returns_ok_when_set_and_recognized(t *testing.T) {
	assertOk(t, checkLanguage,
		*params.AParamSet(
			params.WithBaseDir(testDataDirJava),
			params.WithLanguage("java"),
		),
	)
}

func Test_check_language_returns_error_when_set_but_unknown(t *testing.T) {
	assertError(t, checkLanguage,
		*params.AParamSet(
			params.WithLanguage("unknown-language"),
		),
	)
}

func Test_check_language_returns_error_when_not_set_and_no_base_dir(t *testing.T) {
	assertError(t, checkLanguage,
		*params.AParamSet(
			params.WithLanguage(""),
		),
	)
}

func Test_check_language_returns_ok_when_not_set_but_with_valid_base_dir(t *testing.T) {
	assertOk(t, checkLanguage,
		*params.AParamSet(
			params.WithLanguage(""),
			params.WithBaseDir(testDataDirJava),
		),
	)
}

func Test_check_language_returns_error_when_not_set_and_with_invalid_base_dir(t *testing.T) {
	assertError(t, checkLanguage,
		*params.AParamSet(
			params.WithLanguage(""),
			params.WithBaseDir("invalid-base-dir"),
		),
	)
}

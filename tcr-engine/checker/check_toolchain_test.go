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
	"github.com/murex/tcr/tcr-engine/engine"
	"testing"
)

func Test_check_toolchain_returns_ok_when_set_and_recognized(t *testing.T) {
	assertOk(t, checkToolchain, *engine.AParamSet(engine.WithToolchain("make"), engine.WithLanguage("go")))
}

func Test_check_toolchain_returns_error_when_set_but_unknown(t *testing.T) {
	assertError(t, checkToolchain, *engine.AParamSet(engine.WithToolchain("unknown-toolchain")))
}

func Test_check_toolchain_returns_error_when_set_but_language_incompatible(t *testing.T) {
	assertError(t, checkToolchain, *engine.AParamSet(engine.WithToolchain("gradlew"), engine.WithLanguage("go")))
}

func Test_check_toolchain_returns_ok_when_deduced_from_language(t *testing.T) {
	assertOk(t, checkToolchain, *engine.AParamSet(engine.WithToolchain(""), engine.WithLanguage("go")))
}

func Test_check_toolchain_returns_error_when_cannot_be_deduced_from_language(t *testing.T) {
	assertError(t, checkToolchain, *engine.AParamSet(engine.WithToolchain(""), engine.WithLanguage("unknown-language")))
}

func Test_check_toolchain_returns_ok_when_deduced_from_base_dir_name(t *testing.T) {
	assertOk(t, checkToolchain, *engine.AParamSet(engine.WithToolchain(""), engine.WithLanguage(""), engine.WithBaseDir(testDataDirJava)))
}

func Test_check_toolchain_returns_error_when_cannot_be_deduced_from_language_nor_base_dir_name(t *testing.T) {
	assertError(t, checkToolchain, *engine.AParamSet(engine.WithToolchain(""), engine.WithLanguage(""), engine.WithBaseDir("")))
}

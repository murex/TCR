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

func Test_check_config_directory_returns_ok_when_set_and_exists(t *testing.T) {
	t.Skip("need to provide fake configuration settings for tests")
	assertOk(t, checkConfigDirectory, *engine.AParamSet(engine.WithConfigDir(".")))
}

func Test_check_config_directory_returns_ok_when_not_set(t *testing.T) {
	// When not set, config dir is automatically initialized to the current directory
	t.Skip("need to provide fake configuration settings for tests")
	assertOk(t, checkConfigDirectory, *engine.AParamSet(engine.WithConfigDir("")))
}

func Test_check_config_directory_returns_error_when_set_but_does_not_exist(t *testing.T) {
	t.Skip("need to provide fake configuration settings for tests")
	assertError(t, checkConfigDirectory, *engine.AParamSet(engine.WithConfigDir("missing-dir")))
}

func Test_check_config_directory_returns_error_when_set_but_insufficient_permissions(t *testing.T) {
	t.Skip("disabled until we plug an in-memory filesystem")
	assertError(t, checkConfigDirectory, *engine.AParamSet(engine.WithConfigDir("no-perm-dir")))
}

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

package factory

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_vcs_factory_returns_an_error_when_vcs_is_not_supported(t *testing.T) {
	v, err := InitVCS("unknown-vcs", "")
	assert.IsType(t, &UnsupportedVCSError{}, err)
	assert.Zero(t, v)
}

func Test_unsupported_vcs_message_format(t *testing.T) {
	err := UnsupportedVCSError{"some-vcs"}
	assert.Equal(t, "VCS not supported: \"some-vcs\"", err.Error())
}

func Test_vcs_factory_supports_git(t *testing.T) {
	_, err := InitVCS("git", "")
	assert.NotEqual(t, reflect.TypeOf(&UnsupportedVCSError{}), reflect.TypeOf(err))
}

func Test_vcs_factory_supports_p4(t *testing.T) {
	_, err := InitVCS("p4", "")
	assert.NotEqual(t, reflect.TypeOf(&UnsupportedVCSError{}), reflect.TypeOf(err))
}

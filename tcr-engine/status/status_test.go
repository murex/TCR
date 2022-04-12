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

package status

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_return_code_when_no_error(t *testing.T) {
	RecordState(Ok)
	assert.Equal(t, 0, GetReturnCode())
}

func Test_return_code_on_build_failure(t *testing.T) {
	RecordState(BuildFailed)
	assert.Equal(t, 1, GetReturnCode())
}

func Test_return_code_on_test_failure(t *testing.T) {
	RecordState(TestFailed)
	assert.Equal(t, 2, GetReturnCode())
}

func Test_return_code_on_config_error(t *testing.T) {
	RecordState(ConfigError)
	assert.Equal(t, 3, GetReturnCode())
}

func Test_return_code_on_git_error(t *testing.T) {
	RecordState(GitError)
	assert.Equal(t, 4, GetReturnCode())
}

func Test_return_code_on_miscellaneous_error(t *testing.T) {
	RecordState(OtherError)
	assert.Equal(t, 5, GetReturnCode())
}

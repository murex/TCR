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

package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_check_group_status(t *testing.T) {
	tests := []struct {
		desc        string
		checkpoints []CheckPoint
		expected    CheckStatus
	}{
		{
			"0 checkpoint",
			nil,
			CheckStatusOk,
		},
		{
			"1 ok",
			[]CheckPoint{OkCheckPoint("A")},
			CheckStatusOk,
		},
		{
			"2 ok",
			[]CheckPoint{OkCheckPoint("A"), OkCheckPoint("B")},
			CheckStatusOk,
		},
		{
			"1 warning",
			[]CheckPoint{WarningCheckPoint("A")},
			CheckStatusWarning,
		},
		{
			"1 error",
			[]CheckPoint{ErrorCheckPoint("A")},
			CheckStatusError,
		},
		{
			"1 ok 1 warning",
			[]CheckPoint{OkCheckPoint("A"), WarningCheckPoint("B")},
			CheckStatusWarning,
		},
		{
			"1 warning 1 error",
			[]CheckPoint{WarningCheckPoint("A"), ErrorCheckPoint("B")},
			CheckStatusError,
		},
		{
			"1 ok 1 warning 1 error",
			[]CheckPoint{OkCheckPoint("A"), WarningCheckPoint("B"), ErrorCheckPoint("C")},
			CheckStatusError,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cg := NewCheckGroup("check group")
			cg.Add(test.checkpoints...)
			assert.Equal(t, test.expected, cg.GetStatus())

		})
	}
}

func Test_add_ok_checkpoint(t *testing.T) {
	cg := NewCheckGroup("group")
	cg.Ok("ok checkpoint")
	assert.Equal(t, 1, len(cg.checkpoints))
	assert.Equal(t, CheckStatusOk, cg.GetStatus())
}

func Test_add_warning_checkpoint(t *testing.T) {
	cg := NewCheckGroup("group")
	cg.Warning("warning checkpoint")
	assert.Equal(t, 1, len(cg.checkpoints))
	assert.Equal(t, CheckStatusWarning, cg.GetStatus())
}

func Test_add_error_checkpoint(t *testing.T) {
	cg := NewCheckGroup("group")
	cg.Error("error checkpoint")
	assert.Equal(t, 1, len(cg.checkpoints))
	assert.Equal(t, CheckStatusError, cg.GetStatus())
}

func Test_check_group_print(t *testing.T) {
	// TODO
}

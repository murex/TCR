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
	"testing"

	"github.com/murex/tcr/status"
	"github.com/stretchr/testify/assert"
)

func Test_record_check_state(t *testing.T) {
	tests := []struct {
		desc        string
		checkStatus CheckStatus
		expectedRC  int
	}{
		{"status ok", CheckStatusOk, 0},
		{"status warning", CheckStatusWarning, 1},
		{"status error", CheckStatusError, 2},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			status.RecordState(status.Ok)
			RecordCheckState(test.checkStatus)
			assert.Equal(t, test.expectedRC, status.GetReturnCode())
		})
	}
}

func Test_update_return_state(t *testing.T) {
	tests := []struct {
		desc       string
		checkpoint CheckPoint
		expectedRC int
	}{
		{"status ok", OkCheckPoint("A"), 0},
		{"status warning", WarningCheckPoint("A"), 1},
		{"status error", ErrorCheckPoint("A"), 2},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cg := NewCheckGroup("check group")
			cg.Add(test.checkpoint)
			status.RecordState(status.Ok)
			UpdateReturnState(cg)
			assert.Equal(t, test.expectedRC, status.GetReturnCode())
		})
	}
}

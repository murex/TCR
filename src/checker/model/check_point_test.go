/*
Copyright (c) 2024 Murex

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
	"errors"
	"io/fs"
	"testing"

	"github.com/murex/tcr/report"
	"github.com/stretchr/testify/assert"
)

func Test_checkpoint_creation(t *testing.T) {
	tests := []struct {
		desc       string
		initFunc   func(a ...any) CheckPoint
		expectedRC CheckStatus
	}{
		{"ok", OkCheckPoint, CheckStatusOk},
		{"warning", WarningCheckPoint, CheckStatusWarning},
		{"error", ErrorCheckPoint, CheckStatusError},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cp := test.initFunc("some description")
			assert.Equal(t, test.expectedRC, cp.rc)
			assert.Equal(t, "some description", cp.description)
		})
	}
}

func Test_checkpoint_print(t *testing.T) {
	tests := []struct {
		desc             string
		checkpoint       CheckPoint
		expectedCategory report.Category
		expectedText     string
	}{
		{"ok", OkCheckPoint("A"), report.Info, "\t✔ A"},
		{"warning", WarningCheckPoint("A"), report.Warning, "\t● A"},
		{"error", ErrorCheckPoint("A"), report.Error, "\t▼ A"},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			report.TestWithIsolatedReporter(func(reporter *report.Reporter, sniffer *report.Sniffer) {
				test.checkpoint.Print()
				sniffer.Stop()
				assert.Equal(t, 1, sniffer.GetMatchCount())
				assert.Equal(t, test.expectedCategory, sniffer.GetAllMatches()[0].Type.Category)
				assert.Equal(t, test.expectedText, sniffer.GetAllMatches()[0].Payload.ToString())
			})
		})
	}
}

func Test_checkpoints_for_dir_access_error(t *testing.T) {
	tests := []struct {
		desc     string
		dir      string
		err      error
		expected CheckPoint
	}{
		{
			"missing directory",
			"some-dir",
			fs.ErrNotExist,
			ErrorCheckPoint("directory not found: some-dir"),
		},
		{
			"insufficient permissions",
			"some-dir",
			fs.ErrPermission,
			ErrorCheckPoint("cannot access directory some-dir"),
		},
		{
			"other error",
			"some-dir",
			errors.New("some other error"),
			ErrorCheckPoint("some other error"),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert.Equal(t, []CheckPoint{test.expected},
				CheckpointsForDirAccessError(test.dir, test.err))
		})
	}
}

func Test_checkpoints_for_lists(t *testing.T) {
	tests := []struct {
		desc      string
		headerMsg string
		emptyMsg  string
		values    []string
		expected  []CheckPoint
	}{
		{
			"0 item",
			"some header message",
			"message when empty",
			nil,
			[]CheckPoint{WarningCheckPoint("message when empty")},
		},
		{
			"1 item",
			"some header message",
			"message when empty",
			[]string{"A"},
			[]CheckPoint{
				OkCheckPoint("some header message"),
				OkCheckPoint("- A"),
			},
		},
		{
			"2 items",
			"some header message",
			"message when empty",
			[]string{"A", "B"},
			[]CheckPoint{
				OkCheckPoint("some header message"),
				OkCheckPoint("- A"),
				OkCheckPoint("- B"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert.Equal(t, test.expected,
				CheckpointsForList(test.headerMsg, test.emptyMsg, test.values...))
		})
	}
}

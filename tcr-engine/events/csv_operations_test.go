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

package events

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_append_tcr_event_to_csv_writer(t *testing.T) {
	testFlags := []struct {
		desc     string
		position int
		event    TcrEvent
		expected string
	}{
		{
			"modified source lines",
			0,
			*ATcrEvent(WithModifiedSrcLines(2)),
			"2",
		},
		{
			"modified test lines",
			1,
			*ATcrEvent(WithModifiedTestLines(25)),
			"25",
		},
		{
			"total test cases",
			2,
			*ATcrEvent(WithTotalTestsRun(3)),
			"3",
		},
		{
			"passed test cases",
			3,
			*ATcrEvent(WithTestsPassed(4)),
			"4",
		},
		{
			"failed test cases",
			4,
			*ATcrEvent(WithTestsFailed(10)),
			"10",
		},
		{
			"skipped test cases",
			5,
			*ATcrEvent(WithTestsSkipped(9)),
			"9",
		},
		{
			"test cases with errors",
			6,
			*ATcrEvent(WithTestsWithErrors(10)),
			"10",
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			csvRecord := toCsvRecord(tt.event)
			assert.Equal(t, tt.expected, csvRecord[tt.position])
		})
	}
}

func Test_append_event_to_writer(t *testing.T) {
	var b bytes.Buffer
	event := *ATcrEvent(
		WithModifiedSrcLines(12),
		WithModifiedTestLines(25),
		WithTotalTestsRun(14),
		WithTestsPassed(8),
		WithTestsFailed(2),
		WithTestsSkipped(1),
		WithTestsWithErrors(3),
	)
	_ = appendEvent(event, &b)
	assert.Equal(t, "12,25,14,8,2,1,3\n", b.String())
}

func Test_converts_a_csv_record_to_an_event(t *testing.T) {
	testFlags := []struct {
		desc     string
		eventStr string
		expected TcrEvent
	}{
		{
			"modified source lines",
			"2, 0, 0, 0, 0, 0, 0\n",
			*ATcrEvent(WithModifiedSrcLines(2)),
		},
		{
			"modified test lines",
			"0, 3, 0, 0, 0, 0, 0\n",
			*ATcrEvent(WithModifiedTestLines(3)),
		},
		{
			"total test cases run",
			"0, 0, 4, 0, 0, 0, 0\n",
			*ATcrEvent(WithTotalTestsRun(4)),
		},
		{
			"passed test cases",
			"0, 0, 0, 3, 0, 0, 0\n",
			*ATcrEvent(WithTestsPassed(3)),
		},
		{
			"failed test cases",
			"0, 0, 0, 0, 2, 0, 0\n",
			*ATcrEvent(WithTestsFailed(2)),
		},
		{
			"skipped test cases",
			"0, 0, 0, 0, 0, 5, 0\n",
			*ATcrEvent(WithTestsSkipped(5)),
		},
		{
			"test cases with errors",
			"0, 0, 0, 0, 0, 0, 4\n",
			*ATcrEvent(WithTestsWithErrors(4)),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			csvRecord := toTcrEvent(tt.eventStr)
			assert.Equal(t, tt.expected, csvRecord)
		})
	}
}

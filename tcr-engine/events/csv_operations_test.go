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
	"time"
)

func Test_append_tcr_event_to_csv_writer(t *testing.T) {
	testFlags := []struct {
		desc     string
		position int
		event    TcrEvent
		expected string
	}{
		{
			"timestamp in UTC",
			0,
			*ATcrEvent(WithTimestamp(time.Date(
				2022, 4, 11, 15, 52, 3, 0,
				time.UTC))),
			"2022-04-11 15:52:03",
		},
		{
			"timestamp not in UTC",
			0,
			*ATcrEvent(WithTimestamp(time.Date(
				2022, 4, 11, 15, 52, 3, 0,
				time.FixedZone("UTC-7", -7*60*60)))),
			"2022-04-11 22:52:03",
		},
		{
			"modified source lines",
			1,
			*ATcrEvent(WithModifiedSrcLines(2)),
			"2",
		},
		{
			"modified test lines",
			2,
			*ATcrEvent(WithModifiedTestLines(25)),
			"25",
		},
		{
			"added test cases",
			3,
			*ATcrEvent(WithTotalTestsRan(3)),
			"3",
		},
		{
			"build passing",
			4,
			*ATcrEvent(WithPassingBuild()),
			"1",
		},
		{
			"build failing",
			4,
			*ATcrEvent(WithFailingBuild()),
			"2",
		},
		{
			"tests passing",
			5,
			*ATcrEvent(WithPassingTests()),
			"1",
		},
		{
			"tests failing",
			5,
			*ATcrEvent(WithFailingTests()),
			"2",
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
		WithTimestamp(time.Date(2022, 4, 11, 15, 52, 3, 0, time.UTC)),
		WithModifiedSrcLines(12),
		WithModifiedTestLines(25),
		WithTotalTestsRan(3),
		WithPassingBuild(),
		WithFailingTests(),
	)
	_ = appendEvent(event, &b)
	assert.Equal(t, "2022-04-11 15:52:03,12,25,3,1,2\n", b.String())
}

func Test_converts_a_csv_record_to_an_event(t *testing.T) {
	testFlags := []struct {
		desc     string
		eventStr string
		expected TcrEvent
	}{
		{
			"timestamp in UTC",
			"2022-04-11 15:52:03, 0, 0, 0, 2, 2\n",
			*ATcrEvent(WithTimestamp(time.Date(
				2022, 4, 11, 15, 52, 3, 0,
				time.UTC)),
				WithFailingBuild(),
				WithFailingTests()),
		},
		{
			"modified source lines",
			"2022-04-11 15:52:03, 2, 0, 0, 2, 2\n",
			*ATcrEvent(WithTimestamp(time.Date(
				2022, 4, 11, 15, 52, 3, 0,
				time.UTC)),
				WithModifiedSrcLines(2),
				WithFailingBuild(),
				WithFailingTests()),
		},
		{
			"modified test lines",
			"2022-04-11 15:52:03, 0, 3, 0, 2, 2\n",
			*ATcrEvent(WithTimestamp(time.Date(
				2022, 4, 11, 15, 52, 3, 0,
				time.UTC)),
				WithModifiedTestLines(3),
				WithFailingBuild(),
				WithFailingTests()),
		},
		{
			"added test cases",
			"2022-04-11 15:52:03, 0, 0, 4, 2, 2\n",
			*ATcrEvent(WithTimestamp(time.Date(
				2022, 4, 11, 15, 52, 3, 0,
				time.UTC)),
				WithTotalTestsRan(4),
				WithFailingBuild(),
				WithFailingTests()),
		},
		{
			"with build passed",
			"2022-04-11 15:52:03, 0, 0, 0, 1, 2\n",
			*ATcrEvent(WithTimestamp(time.Date(
				2022, 4, 11, 15, 52, 3, 0,
				time.UTC)),
				WithPassingBuild(),
				WithFailingTests()),
		},
		{
			"with test passed",
			"2022-04-11 15:52:03, 0, 0, 0, 2, 1\n",
			*ATcrEvent(WithTimestamp(time.Date(
				2022, 4, 11, 15, 52, 3, 0,
				time.UTC)),
				WithFailingBuild(),
				WithPassingTests()),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			csvRecord := toTcrEvent(tt.eventStr)
			assert.Equal(t, tt.expected, csvRecord)
		})
	}
}

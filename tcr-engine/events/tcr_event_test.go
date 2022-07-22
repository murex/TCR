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
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_empty_tcr_event_conversion_to_yaml(t *testing.T) {
	event := *ATcrEvent()
	expected := buildYamlString("0", "0", "0", "0", "0", "0", "0", "0s")
	assert.Equal(t, expected, event.ToYaml())
}

func Test_sample_tcr_event_conversion_to_yaml(t *testing.T) {
	event := ATcrEvent(
		WithModifiedSrcLines(1),
		WithModifiedTestLines(2),
		WithTotalTestsRun(12),
		WithTestsPassed(3),
		WithTestsFailed(4),
		WithTestsSkipped(5),
		WithTestsWithErrors(6),
		WithTestsDuration(10*time.Second),
	)
	expected := buildYamlString("1", "2", "12", "3", "4", "5", "6", "10s")
	assert.Equal(t, expected, event.ToYaml())
}

func Test_events_nb_records(t *testing.T) {
	now := time.Now().UTC()
	testFlags := []struct {
		desc     string
		events   TcrEvents
		expected int
	}{
		{
			"nil",
			nil,
			0,
		},
		{
			"no record",
			TcrEvents{},
			0,
		},
		{
			"2 records",
			TcrEvents{now: *ATcrEvent(), now.Add(1 * time.Second): *ATcrEvent()},
			2,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.events.NbRecords())
		})
	}
}

func Test_events_records_counters(t *testing.T) {
	now := time.Now().UTC()
	testFlags := []struct {
		desc            string
		events          TcrEvents
		expectedPassing int
		expectedFailing int
	}{
		{
			desc:            "nil",
			events:          nil,
			expectedPassing: 0,
			expectedFailing: 0,
		},
		{
			"no record",
			TcrEvents{},
			0,
			0,
		},
		{
			"1 passing record",
			TcrEvents{now: *ATcrEvent(WithCommandStatus(StatusPass))},
			1,
			0,
		},
		// TODO add more test cases
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expectedPassing, tt.events.NbPassingRecords())
			assert.Equal(t, tt.expectedFailing, tt.events.NbFailingRecords())
		})
	}
}

func Test_events_time_span(t *testing.T) {
	now := time.Now().UTC()
	testFlags := []struct {
		desc     string
		events   TcrEvents
		expected time.Duration
	}{
		{
			"nil",
			nil,
			0,
		},
		{
			"no record",
			TcrEvents{},
			0,
		},
		{
			"1 records",
			TcrEvents{now: *ATcrEvent()},
			0,
		},
		{
			"2 records",
			TcrEvents{now: *ATcrEvent(), now.Add(1 * time.Second): *ATcrEvent()},
			1 * time.Second,
		},
		{
			"3 records unsorted",
			TcrEvents{now.Add(2 * time.Second): *ATcrEvent(), now: *ATcrEvent(), now.Add(1 * time.Second): *ATcrEvent()},
			2 * time.Second,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.events.TimeSpan())
		})
	}
}

func Test_events_time_boundaries(t *testing.T) {
	now := time.Now().UTC()
	zeroTime := time.Unix(0, 0).UTC()
	testFlags := []struct {
		desc          string
		events        TcrEvents
		expectedStart time.Time
		expectedEnd   time.Time
	}{
		{
			"nil",
			nil,
			zeroTime,
			zeroTime,
		},
		{
			"no record",
			TcrEvents{},
			zeroTime,
			zeroTime,
		},
		{
			"1 records",
			TcrEvents{now: *ATcrEvent()},
			now,
			now,
		},
		{
			"2 records",
			TcrEvents{now: *ATcrEvent(), now.Add(1 * time.Second): *ATcrEvent()},
			now,
			now.Add(1 * time.Second),
		},
		{
			"3 records unsorted",
			TcrEvents{now.Add(2 * time.Second): *ATcrEvent(), now: *ATcrEvent(), now.Add(1 * time.Second): *ATcrEvent()},
			now,
			now.Add(2 * time.Second),
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expectedStart, tt.events.StartingTime())
			assert.Equal(t, tt.expectedEnd, tt.events.EndingTime())
		})
	}
}

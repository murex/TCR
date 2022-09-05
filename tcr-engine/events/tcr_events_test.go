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

func Test_events_nb_records(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(1 * time.Second)

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
			TcrEvents{
				now:   *ADatedTcrEvent(WithTimestamp(now)),
				later: *ADatedTcrEvent(WithTimestamp(later)),
			},
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
	later := now.Add(1 * time.Second)

	testFlags := []struct {
		desc                      string
		events                    TcrEvents
		expectedPassingCount      int
		expectedPassingPercentage int
		expectedFailingCount      int
		expectedFailingPercentage int
	}{
		{
			"nil",
			nil,
			0,
			0,
			0,
			0,
		},
		{
			"no record",
			TcrEvents{},
			0,
			0,
			0,
			0,
		},
		{
			"1 passing",
			TcrEvents{
				now: *ADatedTcrEvent(
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusPass))),
				),
			},
			1,
			100,
			0,
			0,
		},
		{
			"1 failing",
			TcrEvents{
				now: *ADatedTcrEvent(
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusFail))),
				),
			},
			0,
			0,
			1,
			100,
		},
		{
			"1 passing 1 failing",
			TcrEvents{
				now: *ADatedTcrEvent(
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusPass))),
				),
				later: *ADatedTcrEvent(
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusFail))),
				),
			},
			1,
			50,
			1,
			50,
		},
		{
			"2 passing",
			TcrEvents{
				now: *ADatedTcrEvent(
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusPass))),
				),
				later: *ADatedTcrEvent(
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusPass))),
				),
			},
			2,
			100,
			0,
			0,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expectedPassingCount, tt.events.NbPassingRecords(), "passing count")
			assert.Equal(t, tt.expectedPassingPercentage, tt.events.PercentPassing(), "passing percentage")
			assert.Equal(t, tt.expectedFailingCount, tt.events.NbFailingRecords(), " failing count")
			assert.Equal(t, tt.expectedFailingPercentage, tt.events.PercentFailing(), "failing percentage")
		})
	}
}

func Test_events_time_span_and_boundaries(t *testing.T) {
	now := time.Now().UTC()
	oneSecLater := now.Add(1 * time.Second)
	twoSecLater := now.Add(2 * time.Second)
	zeroTime := time.Unix(0, 0).UTC()

	testFlags := []struct {
		desc             string
		events           TcrEvents
		expectedStart    time.Time
		expectedEnd      time.Time
		expectedTimespan time.Duration
	}{
		{
			"nil",
			nil,
			zeroTime,
			zeroTime,
			0,
		},
		{
			"no record",
			*NewTcrEvents(0),
			zeroTime,
			zeroTime,
			0,
		},
		{
			"1 record",
			TcrEvents{
				now: *ADatedTcrEvent(WithTimestamp(now)),
			},
			now,
			now,
			0,
		},
		{
			"2 records",
			TcrEvents{
				now:         *ADatedTcrEvent(WithTimestamp(now)),
				oneSecLater: *ADatedTcrEvent(WithTimestamp(oneSecLater)),
			},
			now,
			oneSecLater,
			1 * time.Second,
		},
		{
			"3 records unsorted",
			TcrEvents{
				twoSecLater: *ADatedTcrEvent(WithTimestamp(twoSecLater)),
				now:         *ADatedTcrEvent(WithTimestamp(now)),
				oneSecLater: *ADatedTcrEvent(WithTimestamp(oneSecLater)),
			},
			now,
			twoSecLater,
			2 * time.Second,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expectedStart, tt.events.StartingTime(), "starting time")
			assert.Equal(t, tt.expectedEnd, tt.events.EndingTime(), "ending time")
			assert.Equal(t, tt.expectedTimespan, tt.events.TimeSpan(), "time span")
		})
	}
}

func Test_events_duration_in_green_and_red(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(1 * time.Second)

	testFlags := []struct {
		desc                    string
		events                  TcrEvents
		expectedDurationInGreen time.Duration
		expectedDurationInRed   time.Duration
	}{
		{
			"nil",
			nil,
			0,
			0,
		},
		{
			"1 record",
			TcrEvents{now: *ADatedTcrEvent(WithTimestamp(now))},
			0,
			0,
		},
		{
			"2 records starting green",
			TcrEvents{
				now: *ADatedTcrEvent(
					WithTimestamp(now),
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusPass))),
				),
				later: *ADatedTcrEvent(
					WithTimestamp(later),
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusPass))),
				),
			},
			1 * time.Second,
			0,
		},
		// TODO (WIP)
		//{
		//	"2 records starting red",
		//	TcrEvents{now: *ATcrEvent(WithCommandStatus(StatusFail)), now.Add(1 * time.Second): *ATcrEvent()},
		//	0,
		//	1 * time.Second,
		//},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expectedDurationInGreen, tt.events.DurationInGreen(), "duration in green")
			assert.Equal(t, tt.expectedDurationInRed, tt.events.DurationInRed(), "duration in red")
		})
	}
}

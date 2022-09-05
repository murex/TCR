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
	"sort"
	"testing"
	"time"
)

func Test_events_nb_records(t *testing.T) {
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
				*ADatedTcrEvent(),
				*ADatedTcrEvent(),
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
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusPass)))),
			},
			1,
			100,
			0,
			0,
		},
		{
			"1 failing",
			TcrEvents{
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusFail)))),
			},
			0,
			0,
			1,
			100,
		},
		{
			"1 passing 1 failing",
			TcrEvents{
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusPass)))),
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusFail)))),
			},
			1,
			50,
			1,
			50,
		},
		{
			"2 passing",
			TcrEvents{
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusPass)))),
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusPass)))),
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

func Test_events_adding(t *testing.T) {
	now := time.Now().UTC()

	testFlags := []struct {
		desc         string
		input        TcrEvents
		expectedSize int
	}{
		{
			"adding to nil",
			nil,
			1,
		},
		{
			"adding to empty slice",
			TcrEvents{},
			1,
		},
		{
			"adding to slice with 1 event",
			TcrEvents{*ADatedTcrEvent()},
			2,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			tt.input.Add(now, *ATcrEvent())
			assert.Equal(t, tt.expectedSize, tt.input.Len())
		})
	}
}

func Test_events_sorting(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(1 * time.Second)
	earlier := now.Add(-1 * time.Second)

	testFlags := []struct {
		desc     string
		input    TcrEvents
		expected TcrEvents
	}{
		{
			"nil",
			nil,
			nil,
		},
		{
			"empty slice",
			TcrEvents{},
			TcrEvents{},
		},
		{
			"1 event",
			TcrEvents{*ADatedTcrEvent(WithTimestamp(now))},
			TcrEvents{*ADatedTcrEvent(WithTimestamp(now))},
		},
		{
			"2 time-ordered events",
			TcrEvents{
				*ADatedTcrEvent(WithTimestamp(now)),
				*ADatedTcrEvent(WithTimestamp(later)),
			},
			TcrEvents{
				*ADatedTcrEvent(WithTimestamp(now)),
				*ADatedTcrEvent(WithTimestamp(later)),
			},
		},
		{
			"2 time-inverted events",
			TcrEvents{
				*ADatedTcrEvent(WithTimestamp(later)),
				*ADatedTcrEvent(WithTimestamp(now)),
			},
			TcrEvents{
				*ADatedTcrEvent(WithTimestamp(now)),
				*ADatedTcrEvent(WithTimestamp(later)),
			},
		},
		{
			"3 unsorted events",
			TcrEvents{
				*ADatedTcrEvent(WithTimestamp(later)),
				*ADatedTcrEvent(WithTimestamp(earlier)),
				*ADatedTcrEvent(WithTimestamp(now)),
			},
			TcrEvents{
				*ADatedTcrEvent(WithTimestamp(earlier)),
				*ADatedTcrEvent(WithTimestamp(now)),
				*ADatedTcrEvent(WithTimestamp(later)),
			},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			sort.Sort(&tt.input)
			assert.Equal(t, tt.expected, tt.input)
		})
	}
}

func Test_events_time_span_and_boundaries(t *testing.T) {
	now := time.Now().UTC()
	oneSecLater := now.Add(1 * time.Second)
	twoSecLater := now.Add(2 * time.Second)

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
			ZeroTime,
			ZeroTime,
			0,
		},
		{
			"no record",
			*NewTcrEvents(),
			ZeroTime,
			ZeroTime,
			0,
		},
		{
			"1 record",
			TcrEvents{
				*ADatedTcrEvent(WithTimestamp(now)),
			},
			now,
			now,
			0,
		},
		{
			"2 records",
			TcrEvents{
				*ADatedTcrEvent(WithTimestamp(now)),
				*ADatedTcrEvent(WithTimestamp(oneSecLater)),
			},
			now,
			oneSecLater,
			1 * time.Second,
		},
		{
			"3 records unsorted",
			TcrEvents{
				*ADatedTcrEvent(WithTimestamp(twoSecLater)),
				*ADatedTcrEvent(WithTimestamp(now)),
				*ADatedTcrEvent(WithTimestamp(oneSecLater)),
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
	oneSecLater := now.Add(1 * time.Second)

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
			TcrEvents{*ADatedTcrEvent(WithTimestamp(now))},
			0,
			0,
		},
		{
			"2 records starting green",
			TcrEvents{
				*ADatedTcrEvent(
					WithTimestamp(now),
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusPass))),
				),
				*ADatedTcrEvent(
					WithTimestamp(oneSecLater),
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusPass))),
				),
			},
			1 * time.Second,
			0,
		},
		{
			"2 records starting red",
			TcrEvents{
				*ADatedTcrEvent(
					WithTimestamp(now),
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusFail))),
				),
				*ADatedTcrEvent(
					WithTimestamp(oneSecLater),
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusPass))),
				),
			},
			0,
			1 * time.Second,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expectedDurationInGreen, tt.events.DurationInGreen(), "duration in green")
			assert.Equal(t, tt.expectedDurationInRed, tt.events.DurationInRed(), "duration in red")
		})
	}
}

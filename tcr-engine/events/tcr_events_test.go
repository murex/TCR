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
	twoSecLater := now.Add(2 * time.Second)

	testFlags := []struct {
		desc                      string
		events                    TcrEvents
		expectedDurationInGreen   time.Duration
		expectedPercentageInGreen int
		expectedDurationInRed     time.Duration
		expectedPercentageInRed   int
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
			"1 record",
			TcrEvents{*ADatedTcrEvent(WithTimestamp(now))},
			0,
			0,
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
			100,
			0,
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
			0,
			1 * time.Second,
			100,
		},
		{
			"3 records",
			TcrEvents{
				*ADatedTcrEvent(
					WithTimestamp(now),
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusFail))),
				),
				*ADatedTcrEvent(
					WithTimestamp(oneSecLater),
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusPass))),
				),
				*ADatedTcrEvent(
					WithTimestamp(twoSecLater),
					WithTcrEvent(*ATcrEvent(WithCommandStatus(StatusFail))),
				),
			},
			1 * time.Second,
			50,
			1 * time.Second,
			50,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expectedDurationInGreen, tt.events.DurationInGreen(), "duration in green")
			assert.Equal(t, tt.expectedPercentageInGreen, tt.events.PercentDurationInGreen(), "percentage in green")
			assert.Equal(t, tt.expectedDurationInRed, tt.events.DurationInRed(), "duration in red")
			assert.Equal(t, tt.expectedPercentageInRed, tt.events.PercentDurationInRed(), "percentage in red")
		})
	}
}

func Test_events_time_between_commits(t *testing.T) {
	now := time.Now().UTC()
	oneSecLater := now.Add(1 * time.Second)
	twoSecLater := now.Add(2 * time.Second)
	oneMinLater := now.Add(1 * time.Minute)

	testFlags := []struct {
		desc     string
		events   TcrEvents
		expected DurationAggregates
	}{
		{
			"nil",
			nil,
			DurationAggregates{0, 0, 0},
		},
		{
			"no record",
			*NewTcrEvents(),
			DurationAggregates{0, 0, 0},
		},
		{
			"1 record",
			TcrEvents{
				*ADatedTcrEvent(WithTimestamp(now)),
			},
			DurationAggregates{0, 0, 0},
		},
		{
			"2 records",
			TcrEvents{
				*ADatedTcrEvent(WithTimestamp(now)),
				*ADatedTcrEvent(WithTimestamp(oneSecLater)),
			},
			DurationAggregates{1 * time.Second, 1 * time.Second, 1 * time.Second},
		},
		{
			"3 records unsorted",
			TcrEvents{
				*ADatedTcrEvent(WithTimestamp(twoSecLater)),
				*ADatedTcrEvent(WithTimestamp(now)),
				*ADatedTcrEvent(WithTimestamp(oneSecLater)),
			},
			DurationAggregates{1 * time.Second, 1 * time.Second, 1 * time.Second},
		},
		{
			"4 records",
			TcrEvents{
				*ADatedTcrEvent(WithTimestamp(now)),
				*ADatedTcrEvent(WithTimestamp(oneSecLater)),
				*ADatedTcrEvent(WithTimestamp(twoSecLater)),
				*ADatedTcrEvent(WithTimestamp(oneMinLater)),
			},
			DurationAggregates{1 * time.Second, 20 * time.Second, 58 * time.Second},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.events.TimeBetweenCommits())
		})
	}
}

func Test_events_line_changes_per_commit(t *testing.T) {
	testFlags := []struct {
		desc         string
		events       TcrEvents
		expectedSrc  IntAggregates
		expectedTest IntAggregates
		expectedAll  IntAggregates
	}{
		{
			"nil",
			nil,
			IntAggregates{0, 0, 0},
			IntAggregates{min: 0, avg: 0, max: 0},
			IntAggregates{min: 0, avg: 0, max: 0},
		},
		{
			"no record",
			*NewTcrEvents(),
			IntAggregates{0, 0, 0},
			IntAggregates{min: 0, avg: 0, max: 0},
			IntAggregates{min: 0, avg: 0, max: 0},
		},
		{
			"1 record",
			TcrEvents{
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(
					WithModifiedSrcLines(1),
					WithModifiedTestLines(2),
				))),
			},
			IntAggregates{1, 1, 1},
			IntAggregates{min: 2, avg: 2, max: 2},
			IntAggregates{min: 3, avg: 3, max: 3},
		},
		{
			"2 records",
			TcrEvents{
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(
					WithModifiedSrcLines(1),
					WithModifiedTestLines(2),
				))),
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(
					WithModifiedSrcLines(2),
					WithModifiedTestLines(3),
				))),
			},
			IntAggregates{1, 1.5, 2},
			IntAggregates{min: 2, avg: 2.5, max: 3},
			IntAggregates{min: 3, avg: 4, max: 5},
		},
		{
			"3 records with rounded avg",
			TcrEvents{
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(
					WithModifiedSrcLines(1),
					WithModifiedTestLines(2),
				))),
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(
					WithModifiedSrcLines(2),
					WithModifiedTestLines(3),
				))),
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(
					WithModifiedSrcLines(4),
					WithModifiedTestLines(5),
				))),
			},
			IntAggregates{1, 2.3, 4},
			IntAggregates{min: 2, avg: 3.3, max: 5},
			IntAggregates{min: 3, avg: 5.6, max: 9},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expectedSrc, tt.events.SrcLineChangesPerCommit(), "src lines aggregates")
			assert.Equal(t, tt.expectedTest, tt.events.TestLineChangesPerCommit(), "test lines aggregates")
			assert.Equal(t, tt.expectedAll, tt.events.AllLineChangesPerCommit(), "all lines aggregates")
		})
	}
}

func Test_test_counters_evolution(t *testing.T) {
	testFlags := []struct {
		desc             string
		events           TcrEvents
		expectedPassing  IntValueEvolution
		expectedFailing  IntValueEvolution
		expectedSkipped  IntValueEvolution
		expectedDuration DurationValueEvolution
	}{
		{
			"nil",
			nil,
			IntValueEvolution{0, 0},
			IntValueEvolution{0, 0},
			IntValueEvolution{0, 0},
			DurationValueEvolution{0, 0},
		},
		{
			"no record",
			*NewTcrEvents(),
			IntValueEvolution{0, 0},
			IntValueEvolution{0, 0},
			IntValueEvolution{0, 0},
			DurationValueEvolution{0, 0},
		},
		{
			"1 record",
			TcrEvents{
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(
					WithTestsPassed(1),
					WithTestsFailed(2),
					WithTestsSkipped(3),
					WithTestsDuration(1*time.Second),
				))),
			},
			IntValueEvolution{1, 1},
			IntValueEvolution{2, 2},
			IntValueEvolution{3, 3},
			DurationValueEvolution{1 * time.Second, 1 * time.Second},
		},
		{
			"2 records",
			TcrEvents{
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(
					WithTestsPassed(1),
					WithTestsFailed(2),
					WithTestsSkipped(3),
					WithTestsDuration(500*time.Millisecond),
				))),
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(
					WithTestsPassed(5),
					WithTestsFailed(1),
					WithTestsSkipped(2),
					WithTestsDuration(1*time.Second),
				))),
			},
			IntValueEvolution{1, 5},
			IntValueEvolution{2, 1},
			IntValueEvolution{3, 2},
			DurationValueEvolution{500 * time.Millisecond, 1 * time.Second},
		},
		{
			"3 records",
			TcrEvents{
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(
					WithTestsPassed(1),
					WithTestsFailed(2),
					WithTestsSkipped(3),
					WithTestsDuration(500*time.Millisecond),
				))),
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(
					WithTestsPassed(10),
					WithTestsFailed(2),
					WithTestsSkipped(4),
					WithTestsDuration(1*time.Second),
				))),
				*ADatedTcrEvent(WithTcrEvent(*ATcrEvent(
					WithTestsPassed(5),
					WithTestsFailed(3),
					WithTestsSkipped(6),
					WithTestsDuration(400*time.Millisecond),
				))),
			},
			IntValueEvolution{1, 5},
			IntValueEvolution{2, 3},
			IntValueEvolution{3, 6},
			DurationValueEvolution{500 * time.Millisecond, 400 * time.Millisecond},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expectedPassing, tt.events.PassingTestsEvolution(), "passing tests")
			assert.Equal(t, tt.expectedFailing, tt.events.FailingTestsEvolution(), "failing tests")
			assert.Equal(t, tt.expectedSkipped, tt.events.SkippedTestsEvolution(), "skipped tests")
			assert.Equal(t, tt.expectedDuration, tt.events.TestDurationEvolution(), "test duration")
		})
	}
}

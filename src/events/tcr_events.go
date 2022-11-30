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
	"sort"
	"time"
)

// TcrEvents is a slice of DatedTcrEvent instances
type TcrEvents []DatedTcrEvent

// NewTcrEvents initializes a TcrEvents instance
func NewTcrEvents() *TcrEvents {
	events := make(TcrEvents, 0)
	return &events
}

// Add adds a new event to the TcrEvents instance
func (events *TcrEvents) Add(t time.Time, e TcrEvent) {
	*events = append(*events, NewDatedTcrEvent(t, e))
}

// NbRecords provides the number of records in TcrEvents
func (events *TcrEvents) NbRecords() int {
	return events.Len()
}

// TimeSpan returns the total duration of TcrEvents (from first record to the last).
// Returns 0 if TcrEvents contains 0 or 1 record
func (events *TcrEvents) TimeSpan() time.Duration {
	if events.NbRecords() < 2 {
		return 0
	}
	start, end := events.getBoundaryTimes()
	return end.Sub(start)
}

// sortByTime sorts items by timestamp.
// This is a prerequisite for all time-related stats operations
// to return an accurate value
func (events *TcrEvents) sortByTime() {
	sort.Sort(events)
}

// Less is the function indicating ordering between events.
func (events *TcrEvents) Less(i, j int) bool {
	return (*events)[i].Timestamp.Before((*events)[j].Timestamp)
}

// Swap allows to swap 2 elements in a TcrEvents slice
func (events *TcrEvents) Swap(i, j int) {
	(*events)[i], (*events)[j] = (*events)[j], (*events)[i]
}

// Len returns the number of elements in a TcrEvents slice
func (events *TcrEvents) Len() int {
	return len(*events)
}

// StartingTime provides the starting times for all stored events.
// Returns Unix 0-time if there is no record
func (events *TcrEvents) StartingTime() time.Time {
	if len(*events) == 0 {
		return ZeroTime
	}
	start, _ := events.getBoundaryTimes()
	return start
}

// EndingTime provides the ending times for all stored events
// Returns Unix 0-time if there is no record
func (events *TcrEvents) EndingTime() time.Time {
	if len(*events) == 0 {
		return ZeroTime
	}
	_, end := events.getBoundaryTimes()
	return end
}

func (events *TcrEvents) getBoundaryTimes() (start, end time.Time) {
	events.sortByTime()
	start = (*events)[0].Timestamp
	end = (*events)[events.Len()-1].Timestamp
	return start, end
}

// DurationInGreen returns the duration spent in green, e.g. with no failing tests,
// and its percentage vs the total duration.
func (events *TcrEvents) DurationInGreen() DurationValueAndRatio {
	return DurationValueAndRatio{
		value:      events.durationInState(StatusPass),
		percentage: events.percentDurationInState(StatusPass),
	}
}

// DurationInRed returns the total duration spent in red, e.g. with at least 1 failing test,
// and its percentage vs the total duration.
func (events *TcrEvents) DurationInRed() DurationValueAndRatio {
	return DurationValueAndRatio{
		value:      events.durationInState(StatusFail),
		percentage: events.percentDurationInState(StatusFail),
	}
}

func (events *TcrEvents) durationInState(status CommandStatus) (t time.Duration) {
	if len(*events) < 2 {
		return 0
	}
	events.sortByTime()
	for i := range (*events)[:(events.Len() - 1)] {
		t += (*events)[i].timeInState(status, &(*events)[i+1])
	}
	return t
}

// percentDurationInState provides the percentage (rounded) of time spent in
// the provided status.
// Returns 0 if there are less than 2 records.
func (events *TcrEvents) percentDurationInState(status CommandStatus) int {
	return asPercentage(
		inSeconds(events.durationInState(status)),
		inSeconds(events.TimeSpan()),
	)
}

// TimeBetweenCommits returns the minimum, average and maximum time between commits
func (events *TcrEvents) TimeBetweenCommits() DurationAggregates {
	if len(*events) < 2 {
		return DurationAggregates{0 * time.Second, 0 * time.Second, 0 * time.Second}
	}

	var result DurationAggregates
	events.sortByTime()
	for i := range (*events)[:(events.Len() - 1)] {
		t := (*events)[i].timeSpanUntil(&(*events)[i+1])
		if i == 0 || t < result.min {
			result.min = t
		}
		if t > result.max {
			result.max = t
		}
	}
	result.avg = events.averageDuration()
	return result
}

func (events *TcrEvents) averageDuration() time.Duration {
	if len(*events) < 2 {
		return 0 * time.Second
	}
	return time.Duration(inSeconds(events.TimeSpan())/(events.Len()-1)) * time.Second
}

// SrcLineChangesPerCommit returns the minimum, average and maximum number of changed
// source lines per commit
func (events *TcrEvents) SrcLineChangesPerCommit() IntAggregates {
	return events.lineChangesPerCommit(
		func(e DatedTcrEvent) int {
			return e.Event.Changes.Src
		})
}

// TestLineChangesPerCommit returns the minimum, average and maximum number of changed
// test lines per commit
func (events *TcrEvents) TestLineChangesPerCommit() IntAggregates {
	return events.lineChangesPerCommit(
		func(e DatedTcrEvent) int {
			return e.Event.Changes.Test
		})
}

// AllLineChangesPerCommit returns the minimum, average and maximum number of changed
// lines per commit (source and tests)
func (events *TcrEvents) AllLineChangesPerCommit() IntAggregates {
	return events.lineChangesPerCommit(
		func(e DatedTcrEvent) int {
			return e.Event.Changes.Src + e.Event.Changes.Test
		})
}

func (events *TcrEvents) lineChangesPerCommit(metricFunc func(e DatedTcrEvent) int) IntAggregates {
	if len(*events) == 0 {
		return IntAggregates{0, 0, 0}
	}

	var result IntAggregates
	var totalChangedLines int

	for i, e := range *events {
		changedLines := metricFunc(e)
		if i == 0 || changedLines < result.min {
			result.min = changedLines
		}
		if changedLines > result.max {
			result.max = changedLines
		}
		totalChangedLines += changedLines
	}

	// We keep only 1 decimal for average value
	// (higher precision would just add meaningless noise)
	result.avg = float64(10*totalChangedLines/len(*events)) / 10 //nolint:revive

	return result
}

// PassingTestsEvolution returns the evolution of the number of passing tests
// from the first to the last TCR event
func (events *TcrEvents) PassingTestsEvolution() IntValueEvolution {
	return events.testCounterEvolution(
		func(events *TcrEvents, index int) int {
			return (*events)[index].Event.Tests.Passed
		})
}

// FailingTestsEvolution returns the evolution of the number of failing tests
// from the first to the last TCR event
func (events *TcrEvents) FailingTestsEvolution() IntValueEvolution {
	return events.testCounterEvolution(
		func(events *TcrEvents, index int) int {
			return (*events)[index].Event.Tests.Failed
		})
}

// SkippedTestsEvolution returns the evolution of the number of skipped tests
// from the first to the last TCR event
func (events *TcrEvents) SkippedTestsEvolution() IntValueEvolution {
	return events.testCounterEvolution(
		func(events *TcrEvents, index int) int {
			return (*events)[index].Event.Tests.Skipped
		})
}

func (events *TcrEvents) testCounterEvolution(counterFunc func(events *TcrEvents, index int) int) IntValueEvolution {
	if len(*events) == 0 {
		return IntValueEvolution{from: 0, to: 0}
	}
	return IntValueEvolution{
		from: counterFunc(events, 0),
		to:   counterFunc(events, len(*events)-1),
	}
}

// TestDurationEvolution returns the evolution of the test execution duration
// from the first to the last TCR event
func (events *TcrEvents) TestDurationEvolution() DurationValueEvolution {
	if len(*events) == 0 {
		return DurationValueEvolution{from: 0, to: 0}
	}
	return DurationValueEvolution{
		from: (*events)[0].Event.Tests.Duration,
		to:   (*events)[len(*events)-1].Event.Tests.Duration,
	}
}

// PassingRecords provides the total number of passing records and their percentage
// vs the total number of records
func (events *TcrEvents) PassingRecords() IntValueAndRatio {
	return events.recordsWithState(StatusPass)
}

// FailingRecords provides the total number of failing records and their percentage
// vs the total number of records
func (events *TcrEvents) FailingRecords() IntValueAndRatio {
	return events.recordsWithState(StatusFail)
}

func (events *TcrEvents) recordsWithState(status CommandStatus) IntValueAndRatio {
	if len(*events) == 0 {
		return IntValueAndRatio{0, 0}
	}
	count := events.nbRecordsWithState(status)
	return IntValueAndRatio{
		value:      count,
		percentage: asPercentage(count, events.NbRecords()),
	}
}

func (events *TcrEvents) nbRecordsWithState(status CommandStatus) (count int) {
	for _, te := range *events {
		if te.Event.Status == status {
			count++
		}
	}
	return count
}

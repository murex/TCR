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

// ZeroTime is Unix starting time, e.g. Jan 1, 1970
var ZeroTime = time.Unix(0, 0).UTC()

// DatedTcrEvent is a TcrEvent with a timestamp
type DatedTcrEvent struct {
	Timestamp time.Time
	Event     TcrEvent
}

// TcrEvents is a slice of DatedTcrEvent instances
type TcrEvents []DatedTcrEvent

// NewDatedTcrEvent creates a new DatedTcrEvent instance
func NewDatedTcrEvent(t time.Time, e TcrEvent) DatedTcrEvent {
	return DatedTcrEvent{
		Timestamp: t,
		Event:     e,
	}
}

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
	return len(*events)
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

func (events *TcrEvents) getBoundaryTimes() (start, end time.Time) {
	sort.Sort(events)
	start = (*events)[0].Timestamp
	end = (*events)[events.Len()-1].Timestamp
	return
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

// NbPassingRecords provides the total number of passing records
func (events *TcrEvents) NbPassingRecords() int {
	return events.nbRecordsWithStatus(StatusPass)
}

// NbFailingRecords provides the total number of failing records
func (events *TcrEvents) NbFailingRecords() (count int) {
	return events.nbRecordsWithStatus(StatusFail)
}

func (events *TcrEvents) nbRecordsWithStatus(status CommandStatus) (count int) {
	for _, te := range *events {
		if te.Event.Status == status {
			count++
		}
	}
	return count
}

// PercentPassing provides the percentage (rounded) of passing records out of the total number of records.
// Returns 0 if there is no record.
func (events *TcrEvents) PercentPassing() int {
	if len(*events) == 0 {
		return 0
	}
	return asPercentage(events.NbPassingRecords(), events.NbRecords())
}

// PercentFailing provides the percentage (rounded) of failing records out of the total number of records.
// Returns 0 if there is no record.
func (events *TcrEvents) PercentFailing() int {
	if len(*events) == 0 {
		return 0
	}
	return asPercentage(events.NbFailingRecords(), events.NbRecords())
}

// DurationInGreen returns the total duration spent in green, e.g. with no failing tests.
func (events *TcrEvents) DurationInGreen() (t time.Duration) {
	if len(*events) < 2 {
		return 0
	}
	sort.Sort(events)
	for i := range (*events)[:(events.Len() - 1)] {
		startEvent := (*events)[i]
		endEvent := (*events)[i+1]
		t += timeInState(StatusPass, startEvent, endEvent)
	}
	return
}

// DurationInRed returns the total duration spent in red, e.g. with at least 1 failing test.
func (events *TcrEvents) DurationInRed() (t time.Duration) {
	if len(*events) < 2 {
		return 0
	}
	sort.Sort(events)
	for i := range (*events)[:(events.Len() - 1)] {
		startEvent := (*events)[i]
		endEvent := (*events)[i+1]
		t += timeInState(StatusFail, startEvent, endEvent)
	}
	return
}

func timeInState(status CommandStatus, startEvent DatedTcrEvent, endEvent DatedTcrEvent) time.Duration {
	if startEvent.Event.Status == status {
		return timeSpan(startEvent, endEvent)
	}
	return 0
}

func timeSpan(startEvent DatedTcrEvent, endEvent DatedTcrEvent) time.Duration {
	return endEvent.Timestamp.Sub(startEvent.Timestamp)
}

func asPercentage(dividend, divisor int) int {
	return roundClosestInt(100*dividend, divisor) //nolint:revive
}

func roundClosestInt(dividend, divisor int) int {
	return (dividend + (divisor / 2)) / divisor
}

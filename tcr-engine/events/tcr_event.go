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
	"time"
)

type (
	// CommandStatus is the status of a test command
	CommandStatus string

	// ChangedLines is the structure containing info related to the lines changes in src and test
	ChangedLines struct {
		Src  int
		Test int
	}

	// TestStats is the structure containing info related to the tests execution
	TestStats struct {
		Run      int
		Passed   int
		Failed   int
		Skipped  int
		Error    int
		Duration time.Duration
	}

	// TcrEvent is the structure containing information related to a TCR event
	TcrEvent struct {
		Status  CommandStatus
		Changes ChangedLines
		Tests   TestStats
	}

	// TcrEvents is a map of TcrEvent instances. Key is the commit timestamp
	TcrEvents map[time.Time]TcrEvent
)

// Possible values for CommandStatus
const (
	StatusPass    CommandStatus = "pass"
	StatusFail    CommandStatus = "fail"
	StatusUnknown CommandStatus = "unknown"
)

// NewTcrEvent creates a new TcrEvent instance
func NewTcrEvent(status CommandStatus, changes ChangedLines, stats TestStats) TcrEvent {
	return TcrEvent{
		Status:  status,
		Changes: changes,
		Tests:   stats,
	}
}

// NewTestStats creates a new TestStats instance
func NewTestStats(run int, passed int, failed int, skipped int, errors int, duration time.Duration) TestStats {
	return TestStats{
		Run:      run,
		Passed:   passed,
		Failed:   failed,
		Skipped:  skipped,
		Error:    errors,
		Duration: duration,
	}
}

// NewChangedLines creates a new ChangedLines instance
func NewChangedLines(srcLines, testLines int) ChangedLines {
	return ChangedLines{
		Src:  srcLines,
		Test: testLines,
	}
}

// ToYaml converts a TcrEvent to a yaml string
func (event TcrEvent) ToYaml() string {
	return tcrEventToYaml(event)
}

// FromYaml converts a yaml string to a TcrEvent
func FromYaml(yaml string) TcrEvent {
	return yamlToTcrEvent(yaml)
}

// NewTcrEvents initializes a TcrEvents instance
func NewTcrEvents(n int) *TcrEvents {
	events := make(TcrEvents, n)
	return &events
}

// NbRecords provides the number of records in TcrEvents
func (events TcrEvents) NbRecords() int {
	return len(events)
}

// TimeSpan returns the total duration of TcrEvents (from first record to the last).
// Returns 0 if TcrEvents contains 0 or 1 record
func (events TcrEvents) TimeSpan() time.Duration {
	if events.NbRecords() < 2 {
		return 0
	}
	start, end := events.getBoundaryTimes()
	return end.Sub(start)
}

func (events TcrEvents) getBoundaryTimes() (start, end time.Time) {
	var MinTime = time.Unix(0, 0) // Jan 1, 1970
	var MaxTime = MinTime.Add(1<<63 - 1)

	start, end = MaxTime, MinTime
	for key := range events {
		if key.Before(start) {
			start = key
		}
		if key.After(end) {
			end = key
		}
	}
	return
}

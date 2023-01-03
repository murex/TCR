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

	// TCREvent is the structure containing information related to a TCR event
	TCREvent struct {
		Status  CommandStatus
		Changes ChangedLines
		Tests   TestStats
	}
)

// Possible values for CommandStatus
const (
	StatusPass    CommandStatus = "pass"
	StatusFail    CommandStatus = "fail"
	StatusUnknown CommandStatus = "unknown"
)

// NewTCREvent creates a new TCREvent instance
func NewTCREvent(status CommandStatus, changes ChangedLines, stats TestStats) TCREvent {
	return TCREvent{
		Status:  status,
		Changes: changes,
		Tests:   stats,
	}
}

// NewTestStats creates a new TestStats instance
func NewTestStats(run int, passed int, failed int, skipped int, errors int, duration time.Duration) TestStats { //nolint:revive
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

// ToYAML converts a TCREvent to a YAML string
func (event TCREvent) ToYAML() string {
	return tcrEventToYAML(event)
}

// FromYAML converts a YAML string to a TCREvent
func FromYAML(yaml string) TCREvent {
	return yamlToTCREvent(yaml)
}

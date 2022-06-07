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

import "time"

type (
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
		Changes ChangedLines
		Tests   TestStats
	}
)

// NewTcrEvent creates a new TcrEvent instance
func NewTcrEvent(changes ChangedLines, stats TestStats) TcrEvent {
	return TcrEvent{
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

//go:build test_helper

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

var (
	// snapshotTime is used as a base time for all tests using a timestamp
	snapshotTime = time.Now()
)

// ATcrEvent is a test data builder for a TCR event
func ATcrEvent(builders ...func(tcrEvent *TcrEvent)) *TcrEvent {
	tcrEvent := NewTcrEvent(0, 0, 0, StatusUnknown, StatusUnknown)
	tcrEvent.Timestamp = snapshotTime

	for _, build := range builders {
		build(&tcrEvent)
	}
	return &tcrEvent
}

// WithTimestamp adds timestamp to TcrEvent test data builder
func WithTimestamp(timestamp time.Time) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.Timestamp = timestamp
	}
}

// WithDelay adds delayed timestamp to TcrEvent test data builder
func WithDelay(delay time.Duration) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.Timestamp = tcrEvent.Timestamp.Add(delay)
	}
}

// TodayAt sets TcrEvent test data builder timestamp to today and provided time
func TodayAt(hour int, min int, sec int) time.Time {
	return time.Date(
		snapshotTime.Year(), snapshotTime.Month(), snapshotTime.Day(),
		hour, min, sec,
		0, snapshotTime.Location(),
	)
}

// WithModifiedSrcLines sets modified source lines to TCR event test data builder
func WithModifiedSrcLines(count int) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.ModifiedSrcLines = count
	}
}

// WithModifiedTestLines sets modified test lines to TCR event test data builder
func WithModifiedTestLines(count int) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.ModifiedTestLines = count
	}
}

// WithAddedTestCases sets added test cases to TCR event test data builder
func WithAddedTestCases(count int) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.AddedTestCases = count
	}
}

// WithFailingBuild sets TCR event test data builder with failing build
func WithFailingBuild() func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.BuildStatus = StatusFailed
	}
}

// WithPassingBuild sets TCR event test data builder with passing build
func WithPassingBuild() func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.BuildStatus = StatusPassed
	}
}

// WithFailingTests sets TCR event test data builder with failing tests
func WithFailingTests() func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.TestsStatus = StatusFailed
	}
}

// WithPassingTests sets TCR event test data builder with passing tests
func WithPassingTests() func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.TestsStatus = StatusPassed
	}
}

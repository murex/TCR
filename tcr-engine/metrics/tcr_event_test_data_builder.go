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

package metrics

import "time"

var (
	// snapshotTime is used as a base time for all tests using a timestamp
	snapshotTime = time.Now()
)

// ATcrEvent is a test data builder for a TCR event
func ATcrEvent(builders ...func(tcrEvent *TcrEvent)) *TcrEvent {
	tcrEvent := &TcrEvent{
		timestamp:         snapshotTime,
		modifiedSrcLines:  0,
		modifiedTestLines: 0,
		addedTestCases:    0,
		buildPassed:       true,
		testsPassed:       true,
	}

	for _, build := range builders {
		build(tcrEvent)
	}
	return tcrEvent
}

// WithTimestamp adds timestamp to TcrEvent test data builder
func WithTimestamp(timestamp time.Time) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.timestamp = timestamp
	}
}

// WithDelay adds delayed timestamp to TcrEvent test data builder
func WithDelay(delay time.Duration) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.timestamp = tcrEvent.timestamp.Add(delay)
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
		tcrEvent.modifiedSrcLines = count
	}
}

// WithModifiedTestLines sets modified test lines to TCR event test data builder
func WithModifiedTestLines(count int) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.modifiedTestLines = count
	}
}

// WithAddedTestCases sets added test cases to TCR event test data builder
func WithAddedTestCases(count int) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.addedTestCases = count
	}
}

// WithFailingBuild sets TCR event test data builder with failing build
func WithFailingBuild() func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.buildPassed = false
	}
}

// WithPassingBuild sets TCR event test data builder with passing build
func WithPassingBuild() func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.buildPassed = true
	}
}

// WithFailingTests sets TCR event test data builder with failing tests
func WithFailingTests() func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.testsPassed = false
	}
}

// WithPassingTests sets TCR event test data builder with passing tests
func WithPassingTests() func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.testsPassed = true
	}
}

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

// ATcrEvent is a test data builder for a TCR event
func ATcrEvent(builders ...func(tcrEvent *TcrEvent)) *TcrEvent {
	tcrEvent := NewTcrEvent(
		StatusUnknown,
		NewChangedLines(0, 0),
		NewTestStats(0, 0, 0, 0, 0, 0),
	)

	for _, build := range builders {
		build(&tcrEvent)
	}
	return &tcrEvent
}

// WithModifiedSrcLines sets modified source lines to TCR event test data builder
func WithModifiedSrcLines(count int) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.Changes.Src = count
	}
}

// WithModifiedTestLines sets modified test lines to TCR event test data builder
func WithModifiedTestLines(count int) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.Changes.Test = count
	}
}

// WithTotalTestsRun sets the total number of tests run to TCR event test data builder
func WithTotalTestsRun(count int) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.Tests.Run = count
	}
}

// WithTestsPassed sets the number of passed test cases to TCR event test data builder
func WithTestsPassed(count int) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.Tests.Passed = count
	}
}

// WithTestsFailed sets the number of failed test cases to TCR event test data builder
func WithTestsFailed(count int) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.Tests.Failed = count
	}
}

// WithTestsSkipped sets the number of skipped test cases to TCR event test data builder
func WithTestsSkipped(count int) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.Tests.Skipped = count
	}
}

// WithTestsWithErrors sets the number of test cases with errors to TCR event test data builder
func WithTestsWithErrors(count int) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.Tests.Error = count
	}
}

// WithTestsDuration sets the test duration to TCR event test data builder
func WithTestsDuration(duration time.Duration) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.Tests.Duration = duration
	}
}

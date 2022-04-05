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

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	// snapshotTime is used as a base time for all tests using a timestamp
	snapshotTime = time.Now()
)

// TCR Event test data builder

func aTcrEvent(builders ...func(tcrEvent *TcrEvent)) *TcrEvent {
	tcrEvent := &TcrEvent{
		timestamp:         snapshotTime,
		modifiedSrcCount:  0,
		modifiedTestCount: 0,
		addedTestCount:    0,
		buildPassed:       true,
		testPassed:        true,
	}

	for _, build := range builders {
		build(tcrEvent)
	}
	return tcrEvent
}

func withTimestamp(timestamp time.Time) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.timestamp = timestamp
	}
}

func withDelay(delay time.Duration) func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.timestamp = tcrEvent.timestamp.Add(delay)
	}
}

func withFailingTests() func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.testPassed = false
	}
}

func withNoFailingTests() func(filter *TcrEvent) {
	return func(tcrEvent *TcrEvent) {
		tcrEvent.testPassed = true
	}
}

func todayAt(hour int, min int, sec int) time.Time {
	return time.Date(
		snapshotTime.Year(), snapshotTime.Month(), snapshotTime.Day(),
		hour, min, sec,
		0, snapshotTime.Location(),
	)
}

// -------------------------------------------------------------------------

func Test_compute_score(t *testing.T) {
	var timeInGreenRatio = .5
	var savingRate float64 = 60
	var changesPerCommit float64 = 3
	var expected Score = 10
	assert.Equal(t, expected, computeScore(timeInGreenRatio, savingRate, changesPerCommit))
}

// TODO: Case where changesPerCommit = 0

func Test_compute_duration_between_2_records(t *testing.T) {
	startEvent := aTcrEvent(withTimestamp(todayAt(4, 39, 31)))
	endEvent := aTcrEvent(withTimestamp(todayAt(4, 41, 02)))
	assert.Equal(t, 1*time.Minute+31*time.Second, computeDuration(*startEvent, *endEvent))
}

// TODO: Case where timestamps are inverted
// TODO: Case where a timestamp is not available

func Test_compute_durations_with_no_failing_tests_between_2_records(t *testing.T) {
	startEvent := aTcrEvent(withNoFailingTests())
	endEvent := aTcrEvent(withDelay(1 * time.Second))
	assert.Equal(t, 1*time.Second, computeDurationInGreen(*startEvent, *endEvent))
	assert.Equal(t, 0*time.Second, computeDurationInRed(*startEvent, *endEvent))
}

func Test_compute_durations_with_failing_tests_between_2_records(t *testing.T) {
	startEvent := aTcrEvent(withFailingTests())
	endEvent := aTcrEvent(withDelay(1 * time.Second))
	assert.Equal(t, 0*time.Second, computeDurationInGreen(*startEvent, *endEvent))
	assert.Equal(t, 1*time.Second, computeDurationInRed(*startEvent, *endEvent))
}

func Test_compute_time_ratios_with_no_failing_tests_between_2_records(t *testing.T) {
	startEvent := aTcrEvent()
	endEvent := aTcrEvent(withDelay(1 * time.Second))
	assert.Equal(t, float64(1), computeTimeInGreenRatio(*startEvent, *endEvent))
	assert.Equal(t, float64(0), computeTimeInRedRatio(*startEvent, *endEvent))
}

// TODO: case where timestamps are equal (division by zero)

func Test_compute_time_ratios_with_failing_tests_between_2_records(t *testing.T) {
	startEvent := aTcrEvent(withFailingTests())
	endEvent := aTcrEvent(withDelay(1 * time.Second))
	assert.Equal(t, float64(0), computeTimeInGreenRatio(*startEvent, *endEvent))
	assert.Equal(t, float64(1), computeTimeInRedRatio(*startEvent, *endEvent))
}

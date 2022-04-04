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

func Test_compute_score(t *testing.T) {
	var timeInGreenRatio = .5
	var savingRate float64 = 60
	var changesPerCommit float64 = 3
	var expected Score = 10
	assert.Equal(t, expected, computeScore(timeInGreenRatio, savingRate, changesPerCommit))
}

func Test_compute_duration_between_2_records(t *testing.T) {
	startEvent := TcrEvent{
		timestamp:         getTimestamp(4, 39, 31),
		modifiedSrcCount:  0,
		modifiedTestCount: 0,
		addedTestCount:    0,
		buildPassed:       true,
		testPassed:        true,
	}
	endEvent := TcrEvent{
		timestamp:         getTimestamp(4, 41, 02),
		modifiedSrcCount:  6,
		modifiedTestCount: 0,
		addedTestCount:    0,
		buildPassed:       true,
		testPassed:        true,
	}
	assert.Equal(t, 1*time.Minute+31*time.Second, computeDuration(startEvent, endEvent))
}

func Test_compute_duration_in_green_between_2_records(t *testing.T) {
	startEvent := TcrEvent{
		timestamp:         getTimestamp(4, 39, 31),
		modifiedSrcCount:  0,
		modifiedTestCount: 0,
		addedTestCount:    0,
		buildPassed:       true,
		testPassed:        true,
	}
	endEvent := TcrEvent{
		timestamp:         getTimestamp(4, 41, 02),
		modifiedSrcCount:  6,
		modifiedTestCount: 0,
		addedTestCount:    0,
		buildPassed:       true,
		testPassed:        true,
	}
	assert.Equal(t, 1*time.Minute+31*time.Second, computeDurationInGreen(startEvent, endEvent))
}

func Test_compute_duration_in_red_between_2_records(t *testing.T) {
	startEvent := TcrEvent{
		timestamp:         getTimestamp(4, 39, 31),
		modifiedSrcCount:  0,
		modifiedTestCount: 0,
		addedTestCount:    0,
		buildPassed:       true,
		testPassed:        false,
	}
	endEvent := TcrEvent{
		timestamp:         getTimestamp(4, 41, 02),
		modifiedSrcCount:  6,
		modifiedTestCount: 0,
		addedTestCount:    0,
		buildPassed:       true,
		testPassed:        true,
	}
	assert.Equal(t, 1*time.Minute+31*time.Second, computeDurationInRed(startEvent, endEvent))
}

func Test_compute_time_in_green_ratio_between_2_records(t *testing.T) {
	startEvent := TcrEvent{
		timestamp:         getTimestamp(4, 39, 31),
		modifiedSrcCount:  0,
		modifiedTestCount: 0,
		addedTestCount:    0,
		buildPassed:       true,
		testPassed:        true,
	}
	endEvent := TcrEvent{
		timestamp:         getTimestamp(4, 41, 02),
		modifiedSrcCount:  6,
		modifiedTestCount: 0,
		addedTestCount:    0,
		buildPassed:       true,
		testPassed:        true,
	}
	assert.Equal(t, float64(1), computeTimeInGreenRatio(startEvent, endEvent))
}

func Test_compute_time_in_red_ratio_between_2_records(t *testing.T) {
	startEvent := TcrEvent{
		timestamp:         getTimestamp(4, 39, 31),
		modifiedSrcCount:  0,
		modifiedTestCount: 0,
		addedTestCount:    0,
		buildPassed:       true,
		testPassed:        true,
	}
	endEvent := TcrEvent{
		timestamp:         getTimestamp(4, 41, 02),
		modifiedSrcCount:  6,
		modifiedTestCount: 0,
		addedTestCount:    0,
		buildPassed:       true,
		testPassed:        true,
	}
	assert.Equal(t, float64(0), computeTimeInRedRatio(startEvent, endEvent))
}

func getTimestamp(hour int, min int, sec int) time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, now.Location())
}

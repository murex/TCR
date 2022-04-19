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
	"time"
)

// Score is the computed TCR score
type Score float64

// TcrEvent is the structure containing information related to a TCR event
type TcrEvent struct {
	timestamp         time.Time
	modifiedSrcLines  int
	modifiedTestLines int
	addedTestCases    int
	buildPassed       bool
	testsPassed       bool
}

// NewTcrEvent create a new TCREvent instance
func NewTcrEvent(
	modifiedSrcLines, modifiedTestLines, addedTestCases int,
	buildPassed, testPassed bool) TcrEvent {
	return TcrEvent{
		timestamp:         time.Now(),
		modifiedSrcLines:  modifiedSrcLines,
		modifiedTestLines: modifiedTestLines,
		addedTestCases:    addedTestCases,
		buildPassed:       buildPassed,
		testsPassed:       testPassed,
	}
}

// GetBuildPassed indicates if the build passed in this event
func (event TcrEvent) GetBuildPassed() bool {
	return event.buildPassed
}

// GetTestsPassed indicates if the tests passed in this event
func (event TcrEvent) GetTestsPassed() bool {
	return event.testsPassed
}

func computeScore(timeInGreenRatio float64, savingRate float64, changesPerCommit float64) Score {
	if changesPerCommit == 0 {
		return 0
	}
	return Score(timeInGreenRatio * savingRate / changesPerCommit)
}

func computeDuration(from TcrEvent, to TcrEvent) time.Duration {
	if from.timestamp.After(to.timestamp) {
		return computeDuration(to, from)
	}
	return to.timestamp.Sub(from.timestamp)
}

func computeDurationInGreen(from TcrEvent, to TcrEvent) time.Duration {
	if from.testsPassed {
		return computeDuration(from, to)
	}
	return 0
}

func computeDurationInRed(from TcrEvent, to TcrEvent) time.Duration {
	return computeDuration(from, to) - computeDurationInGreen(from, to)
}

func computeTimeInGreenRatio(from TcrEvent, _ TcrEvent) float64 {
	if from.testsPassed {
		return 1
	}
	return 0
}

func computeTimeInRedRatio(from TcrEvent, to TcrEvent) float64 {
	return 1 - computeTimeInGreenRatio(from, to)
}

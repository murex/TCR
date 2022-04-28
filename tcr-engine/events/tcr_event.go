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

// TcrEventStatus is the status of the build and/or tests in a TCR Event record
type TcrEventStatus int

// List of possible values of TcrEventStatus
const (
	StatusUnknown TcrEventStatus = iota
	StatusPassed
	StatusFailed
)

// TcrEvent is the structure containing information related to a TCR event
type TcrEvent struct {
	Timestamp         time.Time
	ModifiedSrcLines  int
	ModifiedTestLines int
	BuildStatus       TcrEventStatus
	TestsStatus       TcrEventStatus
	TotalTestsRan     int
	TestsFailed       int
	TestsPassed       int
	TestsSkipped      int
	TestsWithErrors   int
}

// NewTcrEvent create a new TCREvent instance
func NewTcrEvent(modifiedSrcLines, modifiedTestLines int,
	buildStatus, testStatus TcrEventStatus,
	totalTestsRan, testsPassed, testsFailed, testsSkipped, testsWithErrors int) TcrEvent {
	return TcrEvent{
		Timestamp:         time.Now(),
		ModifiedSrcLines:  modifiedSrcLines,
		ModifiedTestLines: modifiedTestLines,
		BuildStatus:       buildStatus,
		TestsStatus:       testStatus,
		TotalTestsRan:     totalTestsRan,
		TestsFailed:       testsFailed,
		TestsPassed:       testsPassed,
		TestsSkipped:      testsSkipped,
		TestsWithErrors:   testsWithErrors,
	}
}

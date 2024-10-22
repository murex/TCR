/*
Copyright (c) 2024 Murex

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

package retro

import (
	"github.com/murex/tcr/events"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_generate_retrospective_md_for_empty_tcr_events(t *testing.T) {
	tcrEvents := events.NewTcrEvents()
	md := GenerateMarkdown(tcrEvents)
	assert.Contains(t, md, "# Quick Retrospective")
	assert.Contains(t, md, "Average passed commit size: 0")
	assert.Contains(t, md, "Average failed commit size: 0")
}

func Test_generate_retrospective_md_with_one_passing_and_one_failing_commit(t *testing.T) {
	now := time.Now().UTC()
	tcrEvents := events.NewTcrEvents()
	passedStatus := events.WithCommandStatus(events.StatusPass)
	passedLines := events.WithModifiedSrcLines(20)
	passedTestLines := events.WithModifiedTestLines(10)

	passedEvent := events.ATcrEvent(passedStatus, passedLines, passedTestLines)

	failedStatus := events.WithCommandStatus(events.StatusFail)
	failedLines := events.WithModifiedSrcLines(10)
	failedTestLines := events.WithModifiedTestLines(5)
	failedEvent := events.ATcrEvent(failedStatus, failedLines, failedTestLines)

	tcrEvents.Add(now, *passedEvent)
	tcrEvents.Add(now, *failedEvent)

	md := GenerateMarkdown(tcrEvents)
	assert.Contains(t, md, "# Quick Retrospective")
	assert.Contains(t, md, "Average passed commit size: 30")
	assert.Contains(t, md, "Average failed commit size: 15")
}

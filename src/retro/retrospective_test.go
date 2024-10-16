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
	e "github.com/murex/tcr/events"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_generate_retrospective_md_for_empty_tcr_events(t *testing.T) {
	tcrEvents := e.NewTcrEvents()
	md := GenerateMarkdown(tcrEvents)
	assert.Contains(t, md, "# Quick Retrospective")
	assert.Contains(t, md, "Average passed commit size: 0")
	assert.Contains(t, md, "Average failed commit size: 0")
}

func Test_generate_retrospective_md_with_one_passing_and_one_failing_commit(t *testing.T) {
	now := time.Now().UTC()
	tcrEvents := e.NewTcrEvents()
	tcrEvents.Add(now, *e.ATcrEvent(
		e.WithCommandStatus(e.StatusPass),
		e.WithModifiedSrcLines(20),
		e.WithModifiedTestLines(10)))
	tcrEvents.Add(now, *e.ATcrEvent(
		e.WithCommandStatus(e.StatusFail),
		e.WithModifiedSrcLines(10),
		e.WithModifiedTestLines(5)))

	md := GenerateMarkdown(tcrEvents)
	assert.Contains(t, md, "# Quick Retrospective")
	assert.Contains(t, md, "Average passed commit size: 30")
	assert.Contains(t, md, "Average failed commit size: 15")
}

func Test_include_last_commit_date_in_generated_md(t *testing.T) {
	d := time.Date(2022, 9, 22, 11, 0, 0, 0, time.UTC)
	tcrEvents := e.TcrEvents{
		*e.ADatedTcrEvent(e.WithTimestamp(d.Add(-24 * time.Hour))),
		*e.ADatedTcrEvent(e.WithTimestamp(d)),
	}
	md := GenerateMarkdown(&tcrEvents)
	assert.Contains(t, md, "2022/09/22")
}

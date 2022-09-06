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

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_dated_event_timespan_until(t *testing.T) {
	now := time.Now().UTC()
	oneSecLater := now.Add(1 * time.Second)
	oneSecEarlier := now.Add(-1 * time.Second)

	testFlags := []struct {
		desc      string
		event     *DatedTcrEvent
		nextEvent *DatedTcrEvent
		expected  time.Duration
	}{
		{
			"next event is nil",
			ADatedTcrEvent(),
			nil,
			0,
		},
		{
			"next event is synchronous",
			ADatedTcrEvent(WithTimestamp(now)),
			ADatedTcrEvent(WithTimestamp(now)),
			0,
		},
		{
			"next event happens later",
			ADatedTcrEvent(WithTimestamp(now)),
			ADatedTcrEvent(WithTimestamp(oneSecLater)),
			1 * time.Second,
		},
		{
			"next event happens earlier",
			ADatedTcrEvent(WithTimestamp(now)),
			ADatedTcrEvent(WithTimestamp(oneSecEarlier)),
			0,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			result := tt.event.timeSpanUntil(tt.nextEvent)
			assert.Equal(t, tt.expected, result, "timespan")
		})
	}
}

// TODO add test case for timeInState method

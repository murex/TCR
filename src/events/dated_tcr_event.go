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

// DatedTcrEvent is a TCREvent with a timestamp
type DatedTcrEvent struct {
	Timestamp time.Time
	Event     TCREvent
}

// NewDatedTcrEvent creates a new DatedTcrEvent instance
func NewDatedTcrEvent(t time.Time, e TCREvent) DatedTcrEvent {
	return DatedTcrEvent{Timestamp: t, Event: e}
}

func (datedEvent DatedTcrEvent) timeInState(status CommandStatus, nextEvent *DatedTcrEvent) time.Duration {
	if datedEvent.Event.Status == status {
		return datedEvent.timeSpanUntil(nextEvent)
	}
	return 0
}

func (datedEvent DatedTcrEvent) timeSpanUntil(nextEvent *DatedTcrEvent) time.Duration {
	if nextEvent == nil || nextEvent.Timestamp.Before(datedEvent.Timestamp) {
		return 0
	}
	return nextEvent.Timestamp.Sub(datedEvent.Timestamp)
}

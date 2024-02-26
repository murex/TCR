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

package timer

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// EventTrigger represents what triggers a timer event
type EventTrigger string

// List of possible EventTrigger values
const (
	TriggerStart     EventTrigger = "start"
	TriggerCountdown EventTrigger = "countdown"
	TriggerStop      EventTrigger = "stop"
	TriggerTimeout   EventTrigger = "timeout"
)

const separator = ":"

// EventMessage contains a timer event information
type EventMessage struct {
	Trigger   EventTrigger
	Timeout   time.Duration
	Elapsed   time.Duration
	Remaining time.Duration
}

// WithEmphasis indicates whether the event message should be reported with emphasis flag
func (em EventMessage) WithEmphasis() bool {
	return em.Trigger != TriggerTimeout || em.Remaining >= 0
}

// WrapEventMessage wraps a TimerEventMessage into a string
func WrapEventMessage(em EventMessage) string {
	return fmt.Sprint(em.Trigger,
		separator, int(em.Timeout.Seconds()),
		separator, int(em.Elapsed.Seconds()),
		separator, int(em.Remaining.Seconds()))
}

// UnwrapEventMessage unwraps a timer event message string into a TimerEventMessage
func UnwrapEventMessage(message string) EventMessage {
	parts := strings.Split(message, separator)
	timeout, _ := strconv.Atoi(parts[1])
	elapsed, _ := strconv.Atoi(parts[2])
	remaining, _ := strconv.Atoi(parts[3])
	return EventMessage{
		Trigger:   EventTrigger(parts[0]),
		Timeout:   time.Duration(timeout) * time.Second,
		Elapsed:   time.Duration(elapsed) * time.Second,
		Remaining: time.Duration(remaining) * time.Second,
	}
}

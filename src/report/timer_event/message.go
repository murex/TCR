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

package timer_event //nolint:revive

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Trigger represents what triggers a timer event
type Trigger string

// List of possible Trigger values
const (
	TriggerStart     Trigger = "start"
	TriggerCountdown Trigger = "countdown"
	TriggerStop      Trigger = "stop"
	TriggerTimeout   Trigger = "timeout"
)

const separator = ":"

// Message contains a timer event information
type Message struct {
	Trigger   Trigger
	Timeout   time.Duration
	Elapsed   time.Duration
	Remaining time.Duration
}

// ToString returns the string representation of the message
func (m Message) ToString() string {
	return fmt.Sprint(
		m.Trigger,
		separator, int(m.Timeout.Seconds()),
		separator, int(m.Elapsed.Seconds()),
		separator, int(m.Remaining.Seconds()))
}

// UnwrapMessage unwraps a timer event message string into a Message instance
func UnwrapMessage(str string) Message {
	parts := strings.Split(str, separator)
	timeout, _ := strconv.Atoi(parts[1])
	elapsed, _ := strconv.Atoi(parts[2])
	remaining, _ := strconv.Atoi(parts[3])
	return Message{
		Trigger:   Trigger(parts[0]),
		Timeout:   time.Duration(timeout) * time.Second,
		Elapsed:   time.Duration(elapsed) * time.Second,
		Remaining: time.Duration(remaining) * time.Second,
	}
}

// WithEmphasis indicates whether the message should be reported with emphasis flag
func (m Message) WithEmphasis() bool {
	return m.Trigger != TriggerTimeout || m.Remaining >= 0
}

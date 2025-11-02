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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_get_current_state(t *testing.T) {
	tests := []struct {
		desc     string
		state    reminderState
		expected string
	}{
		{
			desc:     "not started",
			state:    notStarted,
			expected: StatePending,
		},
		{
			desc:     "running",
			state:    running,
			expected: StateRunning,
		},
		{
			desc:     "after timeout",
			state:    afterTimeOut,
			expected: StateTimeout,
		},
		{
			desc:     "stopped after interruption",
			state:    stoppedAfterInterruption,
			expected: StateStopped,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			reminder := NewPeriodicReminder(0, 0, func(ctx ReminderContext) {})
			reminder.state = test.state
			reminder.startTime = time.Now()
			// We don't test the remaining and elapsed values are they
			// are time-sensitive
			assert.Equal(t, test.expected, GetCurrentState(reminder).State)
		})
	}
}

func Test_get_current_state_for_timer_not_initialized(t *testing.T) {
	assert.Equal(t, StateOff, GetCurrentState(nil).State)
}

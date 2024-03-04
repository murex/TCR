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
	"time"
)

// CurrentState provides the current state of the PeriodicReminder
type CurrentState struct {
	State     string
	Timeout   time.Duration
	Elapsed   time.Duration
	Remaining time.Duration
}

// Possible values for CurrentState.State
const (
	StateOff     = "off"
	StatePending = "pending"
	StateRunning = "running"
	StateStopped = "stopped"
	StateTimeout = "timeout"
)

// GetCurrentState returns the current state of the PeriodicReminder
func GetCurrentState(r *PeriodicReminder) CurrentState {
	if r == nil {
		return CurrentState{State: StateOff}
	}
	var state string
	switch r.state {
	case notStarted:
		state = StatePending
	case running:
		if r.GetRemainingTime() > 0 {
			state = StateRunning
		} else {
			state = StateTimeout
		}
	case afterTimeOut:
		state = StateTimeout
	case stoppedAfterInterruption:
		state = StateStopped
	default:
		state = StateOff
	}
	return CurrentState{
		State:     state,
		Timeout:   r.timeout,
		Elapsed:   r.GetElapsedTime(),
		Remaining: r.GetRemainingTime(),
	}
}

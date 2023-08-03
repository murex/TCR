/*
Copyright (c) 2021 Murex

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

const defaultTimeout = 5 * time.Minute
const defaultTickPeriod = 1 * time.Minute

// ReminderState type used for managing reminder state
type ReminderState int

// List of possible values for ReminderState
const (
	NotStarted ReminderState = iota
	Running
	AfterTimeOut
	StoppedAfterInterruption
)

// ReminderEventType type used for managing ticker state
type ReminderEventType int

// List of possible values for ReminderState
const (
	StartEvent ReminderEventType = iota
	PeriodicEvent
	InterruptEvent
	TimeoutEvent
)

// PeriodicReminder provides a mechanism allowing to trigger an action every tickPeriod, until timeout expires.
type PeriodicReminder struct {
	timeout       time.Duration
	tickPeriod    time.Duration
	onEventAction func(ctx ReminderContext)
	state         ReminderState
	startTime     time.Time
	stopTime      time.Time
	tickCounter   int
	lastTickIndex int
	ticker        *time.Ticker
	done          chan bool
}

// ReminderContext provides the context related to a specific reminder event
type ReminderContext struct {
	eventType ReminderEventType
	index     int
	indexMax  int
	timestamp time.Time
	elapsed   time.Duration
	remaining time.Duration
}

// NewPeriodicReminder returns a new PeriodicReminder that will trigger action onEventAction() every tickPeriod,
// until timeout expires.
// The returned PeriodicReminder is ready to start, but is not counting yet.
func NewPeriodicReminder(
	timeout time.Duration,
	tickPeriod time.Duration,
	onEventAction func(ctx ReminderContext),
) *PeriodicReminder {
	r := PeriodicReminder{
		timeout:       defaultTimeout,
		tickPeriod:    defaultTickPeriod,
		tickCounter:   0,
		onEventAction: onEventAction,
		state:         NotStarted,
	}
	if timeout > 0 {
		r.timeout = timeout
	}
	if tickPeriod > 0 {
		r.tickPeriod = tickPeriod
	}
	r.lastTickIndex = int(r.timeout/r.tickPeriod) - 1
	return &r
}

// Start triggers the PeriodicReminder's beginning of counting.
func (r *PeriodicReminder) Start() {
	// Create the ticker and stopTicking it for now
	r.ticker = time.NewTicker(r.tickPeriod)
	r.state = Running
	r.startTime = time.Now()
	r.done = make(chan bool)

	r.onEventAction(r.buildEventContext(StartEvent, r.startTime))

	go func() {
		for {
			select {
			case <-r.done:
				if r.state == StoppedAfterInterruption {
					r.onEventAction(r.buildEventContext(InterruptEvent, time.Now()))
				}
				return
			case timestamp := <-r.ticker.C:
				if r.state == AfterTimeOut {
					r.onEventAction(r.buildEventContext(TimeoutEvent, timestamp))
				} else {
					r.onEventAction(r.buildEventContext(PeriodicEvent, timestamp))
				}
				r.tickCounter++
			}
		}
	}()

	go func() {
		time.Sleep(r.timeout)
		if r.state == Running {
			r.state = AfterTimeOut
		}
	}()
}

func (r *PeriodicReminder) buildEventContext(eventType ReminderEventType, timestamp time.Time) ReminderContext {
	var ctx ReminderContext
	switch eventType {
	case StartEvent:
		ctx = ReminderContext{
			eventType: eventType,
			index:     -1,
			indexMax:  r.lastTickIndex,
			timestamp: timestamp,
			elapsed:   0,
			remaining: r.timeout,
		}
	case PeriodicEvent:
		elapsed := time.Duration(r.tickCounter+1) * r.tickPeriod
		ctx = ReminderContext{
			eventType: eventType,
			index:     r.tickCounter,
			indexMax:  r.lastTickIndex,
			timestamp: timestamp,
			elapsed:   elapsed,
			remaining: r.timeout - elapsed,
		}
	case InterruptEvent:
		ctx = ReminderContext{
			eventType: eventType,
			index:     -1,
			indexMax:  r.lastTickIndex,
			timestamp: timestamp,
			elapsed:   time.Since(r.startTime),
			remaining: 0,
		}
	case TimeoutEvent:
		elapsed := time.Duration(r.tickCounter+1) * r.tickPeriod
		ctx = ReminderContext{
			eventType: eventType,
			index:     r.tickCounter,
			indexMax:  r.lastTickIndex,
			timestamp: timestamp,
			elapsed:   elapsed,
			remaining: r.timeout - elapsed,
		}
	}
	return ctx
}

func (r *PeriodicReminder) stopTicking(s ReminderState) {
	if r.state == Running || r.state == AfterTimeOut {
		r.ticker.Stop()
		r.state = s
		r.stopTime = time.Now()
		r.done <- true
	}
}

// Stop stops the PeriodicReminder, even if it has not yet timed out.
func (r *PeriodicReminder) Stop() {
	r.stopTicking(StoppedAfterInterruption)
}

// GetElapsedTime returns the time elapsed since the timer was started
func (r *PeriodicReminder) GetElapsedTime() time.Duration {
	switch r.state {
	case NotStarted:
		return 0
	case Running, AfterTimeOut:
		return time.Since(r.startTime)
	default:
		return r.stopTime.Sub(r.startTime)
	}
}

// GetRemainingTime returns the time remaining until the timer ends
func (r *PeriodicReminder) GetRemainingTime() time.Duration {
	switch r.state {
	case NotStarted:
		return r.timeout
	case Running, AfterTimeOut:
		return r.timeout - time.Since(r.startTime)
	default:
		return 0
	}
}

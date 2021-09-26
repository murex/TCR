package timer

import (
	"time"
)

const defaultTimeout = 5 * time.Minute
const defaultTickPeriod = 1 * time.Minute

// ReminderState type used for managing ticker state
type ReminderState int

// List of possible values for ReminderState
const (
	NotStarted ReminderState = iota
	Running
	StoppedAfterTimeOut
	StoppedAfterInterruption
)

// PeriodicReminder provides a mechanism allowing to trigger an action every tickPeriod, until timeout expires.
type PeriodicReminder struct {
	timeout     time.Duration
	tickPeriod  time.Duration
	tickCounter int
	onTick      func(tickIndex int, timestamp time.Time)
	ticker      *time.Ticker
	state       ReminderState
	timedOut    chan bool
	interrupted chan bool
	startTime   time.Time
	stopTime    time.Time
}

// New returns a new PeriodicReminder that will trigger action onTick() every tickPeriod, until timeout expires.
// The returned PeriodicReminder is ready to start, but is not counting yet.
func New(
	timeout time.Duration,
	tickPeriod time.Duration,
	onTick func(tickIndex int, timestamp time.Time),
) *PeriodicReminder {
	r := PeriodicReminder{
		timeout:     defaultTimeout,
		tickPeriod:  defaultTickPeriod,
		tickCounter: 0,
		onTick:      onTick,
		state:       NotStarted,
	}
	if timeout > 0 {
		r.timeout = timeout
	}
	if tickPeriod > 0 {
		r.tickPeriod = tickPeriod
	}
	return &r
}

// Start triggers the PeriodicReminder's beginning of counting.
func (r *PeriodicReminder) Start() {
	// Create the ticker and stop it for now
	r.ticker = time.NewTicker(r.tickPeriod)
	r.startTime = time.Now()
	r.timedOut = make(chan bool)
	r.interrupted = make(chan bool)
	r.state = Running

	go func() {
		for {
			select {
			case <-r.timedOut:
				r.state = StoppedAfterTimeOut
				r.stopTime = time.Now()
				return
			case <-r.interrupted:
				r.state = StoppedAfterInterruption
				r.stopTime = time.Now()
				return
			case timestamp := <-r.ticker.C:
				r.onTick(r.tickCounter, timestamp)
				r.tickCounter++
			}
		}
	}()

	go func() {
		time.Sleep(r.timeout)

		if r.state == Running {
			r.ticker.Stop()
			r.timedOut <- true
		}
	}()
}

// Stop stops the PeriodicReminder, even if it has not yet timed out.
func (r *PeriodicReminder) Stop() {
	if r.state == Running {
		r.ticker.Stop()
		r.interrupted <- true
	}
}

// GetElapsedTime returns the time elapsed since the timer was started
func (r *PeriodicReminder) GetElapsedTime() time.Duration {
	switch r.state {
	case NotStarted:
		return 0
	case Running:
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
	case Running:
		return r.timeout - time.Since(r.startTime)
	default:
		return 0
	}
}

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
	StoppedAfterTimeOut
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

// NewPeriodicReminder returns a new PeriodicReminder that will trigger action onEventAction() every tickPeriod, until timeout expires.
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
				if r.state == StoppedAfterTimeOut {
					r.onEventAction(r.buildEventContext(TimeoutEvent, time.Now()))
				}
				if r.state == StoppedAfterInterruption {
					r.onEventAction(r.buildEventContext(InterruptEvent, time.Now()))
				}
				return
			case timestamp := <-r.ticker.C:
				r.onEventAction(r.buildEventContext(PeriodicEvent, timestamp))
				r.tickCounter++
			}
		}
	}()

	go func() {
		time.Sleep(r.timeout)
		r.stopTicking(StoppedAfterTimeOut)
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
		ctx = ReminderContext{
			eventType: eventType,
			index:     -1,
			indexMax:  r.lastTickIndex,
			timestamp: timestamp,
			elapsed:   r.timeout,
			remaining: 0,
		}
	}
	return ctx
}

func (r *PeriodicReminder) stopTicking(s ReminderState) {
	if r.state == Running {
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

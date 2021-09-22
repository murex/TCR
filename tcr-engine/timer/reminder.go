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

type Reminder struct {
	timeout     time.Duration
	tickPeriod  time.Duration
	tickCounter int
	onTick      func(tickIndex int, timestamp time.Time)
	ticker      *time.Ticker
	state       ReminderState
	timedOut    chan bool
	interrupted chan bool
}

func New(
	timeout time.Duration,
	tickPeriod time.Duration,
	onTick func(tickIndex int, timestamp time.Time),
) *Reminder {
	r := Reminder{
		timeout: defaultTimeout,
		tickPeriod: defaultTickPeriod,
		tickCounter: 0,
		onTick: onTick,
		state: NotStarted,
	}
	if timeout > 0 {
		r.timeout = timeout
	}
	if tickPeriod > 0 {
		r.tickPeriod = tickPeriod
	}
	return &r
}

func (r *Reminder) Start() {
	// Create the ticker and stop it for now
	r.ticker = time.NewTicker(r.tickPeriod)
	r.timedOut = make(chan bool)
	r.interrupted = make(chan bool)
	r.state = Running

	go func() {
		for {
			select {
			case <-r.timedOut:
				r.state = StoppedAfterTimeOut
				return
			case <-r.interrupted:
				r.state = StoppedAfterInterruption
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

func (r *Reminder) Stop() {
	if r.state == Running {
		r.ticker.Stop()
		r.interrupted <- true
	}
}

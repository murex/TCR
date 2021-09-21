package timer

import (
	"fmt"
	"time"
)

const defaultTimeout = 5 * time.Minute
const defaultTickPeriod = 1 * time.Minute

type Reminder struct {
	timeout    time.Duration
	tickPeriod time.Duration
	onTick     func(t time.Time)
	ticker     *time.Ticker
	done       chan bool
}

func New(
	timeout time.Duration,
	tickPeriod time.Duration,
	onTick func(t time.Time),
) Reminder {
	r := Reminder{timeout: defaultTimeout, tickPeriod: defaultTickPeriod, onTick: onTick}
	if timeout > 0 {
		r.timeout = timeout
	}
	if tickPeriod > 0 {
		r.tickPeriod = tickPeriod
	}

	// Create the ticker and stop it for now
	r.ticker = time.NewTicker(r.tickPeriod)
	r.ticker.Stop()
	r.done = make(chan bool)

	return r
}

func (r Reminder) Start() {
	r.ticker.Reset(r.tickPeriod)

	go func() {
		for {
			select {
			case <-r.done:
				return
			case t := <-r.ticker.C:
				fmt.Println("Tick at", t)
				r.onTick(t)
			}
		}
	}()

	go func() {
		time.Sleep(r.timeout)
		r.Stop()
	}()
}

func (r Reminder) Stop() {
	r.ticker.Stop()
	r.done <- true
}

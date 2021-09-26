package timer

import (
	"github.com/mengdaming/tcr-engine/report"
	"time"
)

const tickPeriod = 1 * time.Minute

// NewDriverCountdownTimer creates a PeriodicReminder that sends a countdown message every minute
// until timeout expires, after which it sends a message notifying the end of driver's turn
func NewDriverCountdownTimer(timeout time.Duration) *PeriodicReminder {
	return New(timeout, tickPeriod,
		func(tickIndex int, timestamp time.Time) {
			// TODO wire in the PeriodicReminder instance
			report.PostWarning("Time spent so far...")
		})
}

package timer

import (
	"github.com/mengdaming/tcr-engine/report"
	"strings"
	"time"
)

const tickPeriod = 1 * time.Minute

const messagePrefix = "(Mob Timer) "

// NewMobTurnCountdown creates a PeriodicReminder that starts when entering driver mode, and
// then sends a countdown message every minute until the driver turn expires, after which it
// sends a message notifying the end of driver's turn
func NewMobTurnCountdown(timeout time.Duration) *PeriodicReminder {
	return New(timeout, tickPeriod,
		func(ctx ReminderContext) {
			switch ctx.eventType {
			case StartEvent:
				report.PostWarning(messagePrefix, "Starting ", fmtDuration(timeout), " countdown")
			case PeriodicEvent:
				if ctx.index < ctx.indexMax {
					report.PostWarning(messagePrefix, "Your turn ends in ", fmtDuration(ctx.remaining))
				}
			case InterruptEvent:
				report.PostWarning(messagePrefix, "Stopping countdown after ", fmtDuration(ctx.elapsed))
			case TimeoutEvent:
				report.PostWarning(messagePrefix, "Time's up. Time to rotate!")
			}
		},
	)
}

func fmtDuration(d time.Duration) string {
	s := d.Round(time.Second).String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	return s
}

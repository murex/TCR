package timer

import (
	"github.com/mengdaming/tcr-engine/report"
	"github.com/mengdaming/tcr-engine/runmode"
	"strings"
	"time"
)

const messagePrefix = "(Mob Timer) "

// NewMobTurnCountdown creates a PeriodicReminder that starts when entering driver mode, and
// then sends a countdown message every minute until the driver turn expires, after which it
// sends a message notifying the end of driver's turn
func NewMobTurnCountdown(mode runmode.RunMode, timeout time.Duration) *PeriodicReminder {
	if mode.NeedsCountdownTimer() {
		tickPeriod := findBestTickPeriodFor(timeout)
		return New(timeout, tickPeriod,
			func(ctx ReminderContext) {
				switch ctx.eventType {
				case StartEvent:
					report.PostEvent(messagePrefix, "Starting ", fmtDuration(timeout), " countdown")
				case PeriodicEvent:
					if ctx.remaining > 0 {
						report.PostEvent(messagePrefix, "Your turn ends in ", fmtDuration(ctx.remaining))
					}
				case InterruptEvent:
					report.PostEvent(messagePrefix, "Stopping countdown after ", fmtDuration(ctx.elapsed))
				case TimeoutEvent:
					report.PostEvent(messagePrefix, "Time's up. Time to rotate!")
				}
			},
		)
	}
	return New(0, 0, func(ctx ReminderContext) {})
}

func findBestTickPeriodFor(timeout time.Duration) time.Duration {
	// TODO adjust tick period based on timeout ranges
	return defaultTickPeriod
}

func fmtDuration(d time.Duration) string {
	s := d.Round(time.Second).String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	return s
}

package timer

import (
	"fmt"
	"github.com/mengdaming/tcr-engine/report"
	"strings"
	"time"
)

const tickPeriod = 1 * time.Minute

// NewMobTurnCountdown creates a PeriodicReminder that starts when entering driver mode, and
// then sends a countdown message every minute until the driver turn expires, after which it
// sends a message notifying the end of driver's turn
func NewMobTurnCountdown(timeout time.Duration) *PeriodicReminder {
	return New(timeout, tickPeriod,
		func(tc TickContext) {
			bar := buildProgressBar(tc.index, tc.indexMax)
			report.PostWarning(bar, " Your turn ends in ", fmtDuration(tc.remaining))
		},
		func(tc TickContext) {
			bar := buildProgressBar(tc.indexMax, tc.indexMax)
			report.PostWarning(bar, " Time to rotate!")
		})
}

func buildProgressBar(current int, max int) string {
	done := strings.Repeat("X", current+1)
	remaining := strings.Repeat(".", max-current)
	return "[" + done + remaining + "]"
}

func fmtDuration(d time.Duration) string {
	m := d.Round(time.Minute) / time.Minute
	return fmt.Sprintf("%d min", m)
}

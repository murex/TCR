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
	"github.com/murex/tcr/tcr-engine/report"
	"github.com/murex/tcr/tcr-engine/runmode"
	"strings"
	"time"
)

const messagePrefix = "(Mob Timer) "

// NewMobTurnCountdown creates a PeriodicReminder that starts when entering driver mode, and
// then sends a countdown message periodically until the driver turn expires, after which it
// sends a message notifying the end of driver's turn
func NewMobTurnCountdown(mode runmode.RunMode, timeout time.Duration) *PeriodicReminder {
	if mode.NeedsCountdownTimer() {
		tickPeriod := findBestTickPeriodFor(timeout)
		return NewPeriodicReminder(timeout, tickPeriod,
			func(ctx ReminderContext) {
				switch ctx.eventType {
				case StartEvent:
					report.PostNotification(messagePrefix, "Starting ", fmtDuration(timeout), " countdown")
				case PeriodicEvent:
					if ctx.remaining > 0 {
						report.PostNotification(messagePrefix, "Your turn ends in ", fmtDuration(ctx.remaining))
					}
				case InterruptEvent:
					report.PostNotification(messagePrefix, "Stopping countdown after ", fmtDuration(ctx.elapsed))
				case TimeoutEvent:
					report.PostNotification(messagePrefix, "Time's up. Time to rotate!")
				}
			},
		)
	}
	return NewPeriodicReminder(0, 0, func(ctx ReminderContext) {})
}

func findBestTickPeriodFor(timeout time.Duration) time.Duration {
	if timeout <= 10*time.Second {
		return 1 * time.Second
	}
	if timeout <= 1*time.Minute {
		return 10 * time.Second
	}
	return defaultTickPeriod
}

func fmtDuration(d time.Duration) string {
	s := d.Round(time.Second).String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	return s
}

// ReportCountDownStatus Reports the status for the provided PeriodicReminder,
// If the PeriodicReminder is in running state, indicates time spent and time remaining.
func ReportCountDownStatus(t *PeriodicReminder) {
	if t == nil {
		report.PostInfo("Mob Timer is off")
	} else {
		switch t.state {
		case NotStarted:
			report.PostInfo("Mob Timer is not started")
		case Running:
			report.PostInfo("Mob Timer: ",
				fmtDuration(t.GetElapsedTime()), " done, ",
				fmtDuration(t.GetRemainingTime()), " to go")
		case StoppedAfterTimeOut:
			report.PostInfo("Mob Timer has timed out")
		case StoppedAfterInterruption:
			report.PostInfo("Mob Timer was interrupted")
		}
	}
}

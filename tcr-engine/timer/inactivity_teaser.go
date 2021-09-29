package timer

import (
	"github.com/mengdaming/tcr-engine/report"
	"time"
)

var message = map[int]string{
	0: "It's been a while since you last saved your work. Is there anything wrong?",
	1: "Still nothing worth saving? Shall we start to worry?",
	2: "Remember: the more you wait, the more reverting will hurt!",
	3: "Come on, let's at least give it a try...",
	4: "Are you still there?!",
}

// NewInactivityTeaser creates a PeriodicReminder that sends a message every teasingPeriod
// until timeout expires.
func NewInactivityTeaser(timeout time.Duration, teasingPeriod time.Duration) *PeriodicReminder {
	return New(timeout, teasingPeriod,
		func(ctx ReminderContext) {
			msg, ok := message[ctx.index]
			if ok {
				report.PostWarning(msg)
			}
		},
	)
}

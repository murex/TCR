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

func NewInactivityReminder(timeout time.Duration, tickPeriod time.Duration) *Reminder {
	return New(timeout, tickPeriod, func(tickIndex int, timestamp time.Time) {

		msg, ok := message[tickIndex]
		if ok {
			report.PostWarning(msg)
		}
	})
}

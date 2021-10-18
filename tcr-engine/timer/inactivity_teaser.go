package timer

import (
	"github.com/murex/tcr-engine/report"
	"github.com/murex/tcr-engine/settings"
	"sync"
	"time"
)

var teasingMessage = map[int]string{
	0: "It's been a while since you last saved your work. Is there anything wrong?",
	1: "Still nothing worth saving? Shall we start to worry?",
	2: "Remember: the more you wait, the more reverting will hurt!",
	3: "Come on, let's at least give it a try...",
	4: "Are you still there?!",
}

var once sync.Once

// InactivityTeaser sends a teasing message every period until it times out
type InactivityTeaser struct {
	timeout  time.Duration
	period   time.Duration
	reminder *PeriodicReminder
}

var teaserInstance *InactivityTeaser

// GetInactivityTeaserInstance returns the InactivityTeaser instance (singleton)
func GetInactivityTeaserInstance() *InactivityTeaser {
	if teaserInstance == nil {
		once.Do(createTeaser)
	}
	return teaserInstance
}

func createTeaser() {
	teaserInstance = &InactivityTeaser{
		timeout: settings.DefaultInactivityTimeout,
		period:  settings.DefaultInactivityPeriod,
	}
	teaserInstance.reminder = createReminder(teaserInstance.timeout, teaserInstance.period)
}

// createReminder creates a PeriodicReminder that sends a teasing message every teasingPeriod
// until timeout expires.
// The message is taken from teasingMessage map defined at the top of this file
func createReminder(timeout time.Duration, teasingPeriod time.Duration) *PeriodicReminder {
	if settings.EnableTcrInactivityTeaser {
		return NewPeriodicReminder(
			timeout,
			teasingPeriod,
			func(ctx ReminderContext) {
				if ctx.eventType == PeriodicEvent {
					msg, ok := teasingMessage[ctx.index]
					if ok {
						report.PostWarning(msg)
					}
				}
			},
		)
	}
	return nil
}

// Start starts sending periodic teasing messages
func (teaser *InactivityTeaser) Start() {
	if teaser.reminder != nil {
		teaser.reminder.Start()
	}
}

// Stop stops sending periodic teasing messages
func (teaser *InactivityTeaser) Stop() {
	if teaser.reminder != nil {
		teaser.reminder.Stop()
	}
}

// Reset resets the InactivityTeaser. Next call to Start will run as if the InactivityTeaser was
// just created
func (teaser *InactivityTeaser) Reset() {
	teaser.Stop()
	teaser.reminder = createReminder(teaser.timeout, teaser.period)
}

package engine

import "time"

// Default values

const (
	// DefaultPollingPeriod is the waiting time between 2 consecutive calls to git pull when running as Navigator
	DefaultPollingPeriod = 2 * time.Second
	// DefaultMobTurnDuration is the default duration for a mob turn
	DefaultMobTurnDuration = 5 * time.Minute
	// DefaultInactivityPeriod is the default inactivity period until TCR sends an inactivity teaser message
	DefaultInactivityPeriod = 1 * time.Minute
	// DefaultInactivityTimeout is the default timeout after which TCR stops sending inactivity teaser messages
	DefaultInactivityTimeout = 5 * time.Minute
)


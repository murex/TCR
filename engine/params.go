package engine

import (
	"github.com/mengdaming/tcr/runmode"
	"time"
)

// Params contains the main parameter values that TCR engine is using
type Params struct {
	CfgFile       string
	Toolchain     string
	AutoPush      bool
	BaseDir       string
	Mode          runmode.RunMode
	PollingPeriod time.Duration
}

const (
	// DefaultPollingPeriod is the waiting time between 2 consecutive calls to git pull
	DefaultPollingPeriod = 2 * time.Second
)

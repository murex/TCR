package engine

import (
	"github.com/murex/tcr-engine/runmode"
	"time"
)

// Params contains the main parameter values that TCR engine is using
type Params struct {
	CfgFile         string
	Toolchain       string
	AutoPush        bool
	BaseDir         string
	Mode            runmode.RunMode
	PollingPeriod   time.Duration
	MobTurnDuration time.Duration
}


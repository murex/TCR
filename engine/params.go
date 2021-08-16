package engine

import (
	"github.com/mengdaming/tcr/runmode"
	"time"
)

type Params struct {
	CfgFile       string
	Toolchain     string
	AutoPush      bool
	BaseDir       string
	Mode          runmode.RunMode
	PollingPeriod time.Duration
}

const (
	DefaultPollingPeriod = 2 * time.Second
)

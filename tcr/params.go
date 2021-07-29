package tcr

import (
	"github.com/mengdaming/tcr/tcr/role"
	"time"
)

type Params struct {
	CfgFile       string
	Toolchain     string
	AutoPush      bool
	BaseDir       string
	Mode          WorkMode
	PollingPeriod time.Duration
}

type WorkMode string

type Role role.Role

const (
	Solo = "solo"
	Mob  = "mob"

	DefaultPollingPeriod = 2 * time.Second
)

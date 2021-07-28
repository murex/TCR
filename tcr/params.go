package tcr

import "time"

type Params struct {
	CfgFile       string
	Toolchain     string
	AutoPush      bool
	BaseDir       string
	Mode          WorkMode
	PollingPeriod time.Duration
}

type WorkMode string

type Role string

const (
	Solo = "solo"
	Mob  = "mob"

	DriverRole = "driver"
	NavigatorRole = "navigator"

	DefaultPollingPeriod = 2 * time.Second
)

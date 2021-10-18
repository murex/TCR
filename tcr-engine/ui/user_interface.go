package ui

import (
	"github.com/murex/tcr-engine/role"
	"github.com/murex/tcr-engine/runmode"
)

// UserInterface provides the interface to be satisfied by a UI implementation to be able to interact
// with TCR engine
type UserInterface interface {
	Start()
	ShowRunningMode(mode runmode.RunMode)
	NotifyRoleStarting(r role.Role)
	NotifyRoleEnding(r role.Role)
	ShowSessionInfo()
	Confirm(message string, def bool) bool
	StartReporting()
	StopReporting()
}

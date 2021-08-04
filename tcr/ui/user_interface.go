package ui

import (
	"github.com/mengdaming/tcr/tcr/role"
	"github.com/mengdaming/tcr/tcr/runmode"
)

type UserInterface interface {
	RunInMode(mode runmode.RunMode)
	ShowRunningMode(mode runmode.RunMode)
	NotifyRoleStarting(r role.Role)
	NotifyRoleEnding(r role.Role)
	ShowSessionInfo()
	Info(a ...interface{})
	Warning(a ...interface{})
	Error(a ...interface{})
	Trace(a ...interface{})
	Confirm(message string, def bool) bool
}

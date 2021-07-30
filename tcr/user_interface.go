package tcr

type UserInterface interface {
	WaitForAction()
	ShowRunningMode(mode WorkMode)
	NotifyRoleStarting(role Role)
	NotifyRoleEnding(role Role)
	ShowSessionInfo()
	Info(a ...interface{})
	Warning(a ...interface{})
	Error(a ...interface{})
	Trace(a ...interface{})
	Confirm(message string, def bool) bool
}

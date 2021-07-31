package tcr

type UserInterface interface {
	RunInMode(mode WorkMode)
	ShowRunningMode(mode WorkMode)
	NotifyRoleStarting(r Role)
	NotifyRoleEnding(r Role)
	ShowSessionInfo()
	Info(a ...interface{})
	Warning(a ...interface{})
	Error(a ...interface{})
	Trace(a ...interface{})
	Confirm(message string, def bool) bool
}

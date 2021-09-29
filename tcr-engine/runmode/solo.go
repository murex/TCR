package runmode

// Solo is a type of run mode useful when the application is used by a single user
type Solo struct {
}

// Name returns the name of this run mode
func (mode Solo) Name() string {
	return "solo"
}

// AutoPushDefault returns the default value of git auto-push option with this run mode
func (mode Solo) AutoPushDefault() bool {
	return false
}

// NeedsCountdownTimer indicates if a countdown timer is needed with this run mode
func (mode Solo) NeedsCountdownTimer() bool {
	return false
}

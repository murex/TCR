package runmode

// Mob is a type of run mode useful when the application is used by a mob of users
type Mob struct {
}

// Name returns the name of this run mode
func (mode Mob) Name() string {
	return "mob"
}

// AutoPushDefault returns the default value of git auto-push option with this run mode
func (mode Mob) AutoPushDefault() bool {
	return true
}

// NeedsCountdownTimer indicates if a countdown timer is needed with this run mode
func (mode Mob) NeedsCountdownTimer() bool {
	return true
}

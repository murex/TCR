package role

// Driver is a role where the user is actively adding changes to the project
type Driver struct {
}

// Name returns the name of the current role
func (role Driver) Name() string {
	return "driver"
}

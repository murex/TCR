package role

// Navigator is a role where the user does not apply any change, but regularly retrieves updates from git
type Navigator struct {
}

// Name returns the name of the current role
func (role Navigator) Name() string {
	return "navigator"
}

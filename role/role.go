package role

// Role provides the interface that a role must implement in order to be used by TCR engine
type Role interface {
	Name() string
}

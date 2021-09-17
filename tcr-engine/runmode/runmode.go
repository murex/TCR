package runmode

// RunMode is the interface that any run mode needs to satisfy to bee used by the TCR engine
type RunMode interface {
	Name() string
	AutoPushDefault() bool
}

var (
	allModes = []RunMode{Mob{}, Solo{}}
)

// Names returns the list of available run mode names
func Names() []string {
	var names []string
	for _, mode := range allModes {
		names = append(names, mode.Name())
	}
	return names
}

// Map returns the list of available run modes as a map of strings
func Map() map[string]RunMode {
	var m = make(map[string]RunMode)
	for _, mode := range allModes {
		m[mode.Name()] = mode
	}
	return m
}

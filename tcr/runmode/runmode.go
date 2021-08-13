package runmode

type RunMode interface {
	Name() string
	AutoPushDefault() bool
}

var (
	allModes = []RunMode{Mob{}, Solo{}}
)

func Names() []string {
	var names []string
	for _, mode := range allModes {
		names = append(names, mode.Name())
	}
	return names
}

func Map() map[string]RunMode {
	var m = make(map[string]RunMode)
	for _, mode := range allModes {
		m[mode.Name()] = mode
	}
	return m
}

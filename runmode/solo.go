package runmode

type Solo struct {
}

func (mode Solo) Name() string {
	return "solo"
}

func (mode Solo) AutoPushDefault() bool {
	return false
}

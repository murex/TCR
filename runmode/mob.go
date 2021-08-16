package runmode

type Mob struct {
}

func (mode Mob) Name() string {
	return "mob"
}

func (mode Mob) AutoPushDefault() bool {
	return true
}

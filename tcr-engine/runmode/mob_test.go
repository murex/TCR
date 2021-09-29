package runmode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mob_mode_name(t *testing.T) {
	assert.Equal(t, "mob", Mob{}.Name())
}

func Test_mob_mode_default_auto_push_if_true(t *testing.T) {
	assert.True(t, Mob{}.AutoPushDefault())
}

func Test_mob_mode_requires_a_countdown_timer(t *testing.T) {
	assert.True(t, Mob{}.NeedsCountdownTimer())
}

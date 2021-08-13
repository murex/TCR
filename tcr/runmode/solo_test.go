package runmode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_solo_mode_name(t *testing.T) {
	assert.Equal(t, "solo", Solo{}.Name())
}

func Test_solo_mode_default_auto_push_if_false(t *testing.T) {
	assert.False(t, Solo{}.AutoPushDefault())
}

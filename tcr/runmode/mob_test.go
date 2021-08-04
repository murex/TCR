package runmode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mob_mode_name(t *testing.T) {
	assert.Equal(t, "mob", Mob{}.Name())
}
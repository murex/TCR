package runmode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_solo_mode_name(t *testing.T) {
	assert.Equal(t, "solo", Solo{}.Name())
}

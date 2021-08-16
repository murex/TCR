package role

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_driver_role_name(t *testing.T) {
	assert.Equal(t, "driver", Driver{}.Name())
}

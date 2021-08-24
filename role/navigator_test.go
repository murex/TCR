package role

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_navigator_role_name(t *testing.T) {
	assert.Equal(t, "navigator", Navigator{}.Name())
}

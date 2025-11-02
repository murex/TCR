package helloworld

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_say_hello(t *testing.T) {
	assert.Equal(t, "Hello Sue!", sayHello("Sue"))
}

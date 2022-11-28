package helloworld

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_say_hello(t *testing.T) {
	assert.Equal(t, "Hello Sue!", sayHello("Sue"))
}

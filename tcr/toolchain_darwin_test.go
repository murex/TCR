package tcr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_cmake_toolchain_build_command_name(t *testing.T) {
	// TODO -- OS-specific cmake path
	assert.Equal(t, "cmake", CmakeToolchain{}.buildCommandName())
}


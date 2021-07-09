package tcr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_gradle_toolchain_name(t *testing.T) {
	assert.Equal(t, "gradle", GradleToolchain{}.name())
}

func Test_gradle_toolchain_build_command_name(t *testing.T) {
	assert.Equal(t, "./gradlew", GradleToolchain{}.buildCommandName())
}

func Test_gradle_toolchain_build_command_args(t *testing.T) {
	assert.Equal(t, []string{"-x", "test"}, GradleToolchain{}.buildCommandArgs())
}

// --------------------------------------------------------------------------

func Test_cmake_toolchain_name(t *testing.T) {
	assert.Equal(t, "cmake", CmakeToolchain{}.name())
}

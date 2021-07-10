package tcr

import (
	"github.com/mengdaming/tcr/trace"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	trace.SetTestMode()
	os.Exit(m.Run())
}

func Test_unrecognized_toolchain_name(t *testing.T) {
	assert.Zero(t, NewToolchain("dummy"))
	assert.NotZero(t, trace.GetExitReturnCode())
}

func Test_gradle_toolchain_creation(t *testing.T) {
	assert.Equal(t, GradleToolchain{}, NewToolchain("gradle"))
}

func Test_gradle_toolchain_name(t *testing.T) {
	assert.Equal(t, "gradle", GradleToolchain{}.name())
}

func Test_gradle_toolchain_build_command_name(t *testing.T) {
	assert.Equal(t, "./gradlew", GradleToolchain{}.buildCommandName())
}

func Test_gradle_toolchain_build_command_args(t *testing.T) {
	assert.Equal(t, []string{"-x", "test"}, GradleToolchain{}.buildCommandArgs())
}

func Test_gradle_toolchain_returns_non_0_when_build_fails(t *testing.T) {
	// TODO How to setup failing/successful test environments?
	assert.NotZero(t, GradleToolchain{}.runBuild())
}

// --------------------------------------------------------------------------

func Test_cmake_toolchain_creation(t *testing.T) {
	assert.Equal(t, CmakeToolchain{}, NewToolchain("cmake"))
}

func Test_cmake_toolchain_name(t *testing.T) {
	assert.Equal(t, "cmake", CmakeToolchain{}.name())
}

func Test_cmake_toolchain_build_command_args(t *testing.T) {
	assert.Equal(t, []string{
		"--build", ".",
		"--config", "Debug",
	}, CmakeToolchain{}.buildCommandArgs())
}

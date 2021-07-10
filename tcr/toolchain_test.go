package tcr

import (
	"github.com/mengdaming/tcr/trace"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Prevent trace.Error() from triggering os.Exit()
	trace.SetTestMode()
	os.Exit(m.Run())
}

func runFromDir(t *testing.T, testDir string, testFunction func(t *testing.T)) {
	initialDir, _ := os.Getwd()
	_ = os.Chdir(testDir)
	workDir, _ := os.Getwd()
	trace.Info("Working directory: ", workDir)
	testFunction(t)
	_ = os.Chdir(initialDir)
}

// Gradle --------------------------------------------------------------------------

func Test_unrecognized_toolchain_name(t *testing.T) {
	assert.Zero(t, NewToolchain("dummy"))
	assert.NotZero(t, trace.GetExitReturnCode())
}

func Test_gradle_toolchain_initialization(t *testing.T) {
	assert.Equal(t, GradleToolchain{}, NewToolchain("gradle"))
}

func Test_gradle_toolchain_name(t *testing.T) {
	assert.Equal(t, "gradle", GradleToolchain{}.name())
}

func Test_gradle_toolchain_build_command_name(t *testing.T) {
	assert.Equal(t, "gradlew", GradleToolchain{}.buildCommandName())
}

func Test_gradle_toolchain_build_command_args(t *testing.T) {
	assert.Equal(t, []string{"build", "-x", "test"}, GradleToolchain{}.buildCommandArgs())
}

func Test_gradle_toolchain_returns_error_when_build_fails(t *testing.T) {
	runFromDir(t,"../test/kata",
		func(t *testing.T) {
			assert.NotZero(t, GradleToolchain{}.runBuild())
		})
}

func Test_gradle_toolchain_returns_ok_when_build_passes(t *testing.T) {
	runFromDir(t,"../test/kata/java",
		func(t *testing.T) {
			assert.Zero(t, GradleToolchain{}.runBuild())
		})
}

// CMake -------------------------------------------------------------------------

func Test_cmake_toolchain_initialization(t *testing.T) {
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

func Test_cmake_toolchain_returns_error_when_build_fails(t *testing.T) {
	runFromDir(t,"../test/kata",
		func(t *testing.T) {
			assert.NotZero(t, CmakeToolchain{}.runBuild())
		})
}

// TODO Figure out a way to provide a cmake wrapper
func test_cmake_toolchain_returns_ok_when_build_passes(t *testing.T) {
	runFromDir(t,"../test/kata/cpp",
		func(t *testing.T) {
			assert.Zero(t, CmakeToolchain{}.runBuild())
		})
}
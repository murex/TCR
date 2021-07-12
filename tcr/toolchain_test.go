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

func Test_unrecognized_toolchain_name(t *testing.T) {
	assert.Zero(t, NewToolchain("dummy", nil))
	assert.NotZero(t, trace.GetExitReturnCode())
}

// Default toolchain for each language ---------------------------------------------

func Test_default_toolchain_for_java(t *testing.T) {
	assert.Equal(t, GradleToolchain{}, defaultToolchain(JavaLanguage{}))
}

func Test_default_toolchain_for_cpp(t *testing.T) {
	assert.Equal(t, CmakeToolchain{}, defaultToolchain(CppLanguage{}))
}

// Gradle --------------------------------------------------------------------------

func Test_gradle_toolchain_initialization(t *testing.T) {
	assert.Equal(t, GradleToolchain{}, NewToolchain("gradle", JavaLanguage{}))
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
	runFromDir(t, "../test/kata",
		func(t *testing.T) {
			assert.NotZero(t, GradleToolchain{}.runBuild())
		})
}

func Test_gradle_toolchain_returns_ok_when_build_passes(t *testing.T) {
	runFromDir(t, "../test/kata/java",
		func(t *testing.T) {
			assert.Zero(t, GradleToolchain{}.runBuild())
		})
}

func Test_gradle_toolchain_test_command_name(t *testing.T) {
	assert.Equal(t, "gradlew", GradleToolchain{}.testCommandName())
}

func Test_gradle_toolchain_test_command_args(t *testing.T) {
	assert.Equal(t, []string{"test"}, GradleToolchain{}.testCommandArgs())
}

func Test_gradle_toolchain_returns_error_when_tests_fail(t *testing.T) {
	runFromDir(t, "../test/kata",
		func(t *testing.T) {
			assert.NotZero(t, GradleToolchain{}.runTests())
		})
}

func Test_gradle_toolchain_returns_ok_when_tests_pass(t *testing.T) {
	runFromDir(t, "../test/kata/java",
		func(t *testing.T) {
			assert.Zero(t, GradleToolchain{}.runTests())
		})
}

func Test_gradle_toolchain_supports_java(t *testing.T) {
	assert.True(t, GradleToolchain{}.supports(JavaLanguage{}))
}

func Test_gradle_toolchain_does_not_support_cpp(t *testing.T) {
	assert.False(t, GradleToolchain{}.supports(CppLanguage{}))
}

// Maven --------------------------------------------------------------------------

func Test_maven_toolchain_initialization(t *testing.T) {
	assert.Equal(t, MavenToolchain{}, NewToolchain("maven", JavaLanguage{}))
}

func Test_maven_toolchain_name(t *testing.T) {
	assert.Equal(t, "maven", MavenToolchain{}.name())
}

func Test_maven_toolchain_build_command_name(t *testing.T) {
	assert.Equal(t, "mvnw", MavenToolchain{}.buildCommandName())
}

func Test_maven_toolchain_build_command_args(t *testing.T) {
	assert.Equal(t, []string{"test-compile"}, MavenToolchain{}.buildCommandArgs())
}

func Test_maven_toolchain_returns_error_when_build_fails(t *testing.T) {
	runFromDir(t, "../test/kata",
		func(t *testing.T) {
			assert.NotZero(t, MavenToolchain{}.runBuild())
		})
}

func Test_maven_toolchain_returns_ok_when_build_passes(t *testing.T) {
	runFromDir(t, "../test/kata/java",
		func(t *testing.T) {
			assert.Zero(t, MavenToolchain{}.runBuild())
		})
}

func Test_maven_toolchain_test_command_name(t *testing.T) {
	assert.Equal(t, "mvnw", MavenToolchain{}.testCommandName())
}

func Test_maven_toolchain_test_command_args(t *testing.T) {
	assert.Equal(t, []string{"test"}, MavenToolchain{}.testCommandArgs())
}

func Test_maven_toolchain_returns_error_when_tests_fail(t *testing.T) {
	runFromDir(t, "../test/kata",
		func(t *testing.T) {
			assert.NotZero(t, MavenToolchain{}.runTests())
		})
}

func Test_maven_toolchain_returns_ok_when_tests_pass(t *testing.T) {
	runFromDir(t, "../test/kata/java",
		func(t *testing.T) {
			assert.Zero(t, MavenToolchain{}.runTests())
		})
}

func Test_maven_toolchain_supports_java(t *testing.T) {
	assert.True(t, MavenToolchain{}.supports(JavaLanguage{}))
}

func Test_maven_toolchain_does_not_support_cpp(t *testing.T) {
	assert.False(t, MavenToolchain{}.supports(CppLanguage{}))
}

// CMake -------------------------------------------------------------------------

func Test_cmake_toolchain_initialization(t *testing.T) {
	assert.Equal(t, CmakeToolchain{}, NewToolchain("cmake", CppLanguage{}))
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
	runFromDir(t, "../test/kata",
		func(t *testing.T) {
			assert.NotZero(t, CmakeToolchain{}.runBuild())
		})
}

// TODO Figure out a way to provide a cmake wrapper
func test_cmake_toolchain_returns_ok_when_build_passes(t *testing.T) {
	runFromDir(t, "../test/kata/cpp",
		func(t *testing.T) {
			assert.Zero(t, CmakeToolchain{}.runBuild())
		})
}

func Test_cmake_toolchain_test_command_args(t *testing.T) {
	assert.Equal(t, []string{
		"--output-on-failure",
		"-C", "Debug",
	}, CmakeToolchain{}.testCommandArgs())
}

func Test_cmake_toolchain_returns_error_when_tests_fail(t *testing.T) {
	runFromDir(t, "../test/kata",
		func(t *testing.T) {
			assert.NotZero(t, CmakeToolchain{}.runTests())
		})
}

// TODO Figure out a way to provide a cmake wrapper
func test_cmake_toolchain_returns_ok_when_tests_pass(t *testing.T) {
	runFromDir(t, "../test/kata/cpp",
		func(t *testing.T) {
			assert.Zero(t, CmakeToolchain{}.runTests())
		})
}

func Test_cmake_toolchain_supports_cpp(t *testing.T) {
	assert.True(t, CmakeToolchain{}.supports(CppLanguage{}))
}

func Test_cmake_toolchain_does_not_support_java(t *testing.T) {
	assert.False(t, CmakeToolchain{}.supports(JavaLanguage{}))
}

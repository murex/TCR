/*
Copyright (c) 2021 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package toolchain

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

const (
	testDataRootDir = "../testdata"
)

var (
	testDataDirJava = filepath.Join(testDataRootDir, "java")
	//testDataDirCpp  = filepath.Join(testDataRootDir, "cpp")
)

func assertToolchainName(t *testing.T, name string) {
	toolchain, _ := GetToolchain(name)
	assert.Equal(t, name, toolchain.GetName())
}

func assertBuildCommandPath(t *testing.T, expected string, name string) {
	toolchain, _ := GetToolchain(name)
	assert.Equal(t, expected, toolchain.BuildCommandPath())
}

func assertBuildCommandArgs(t *testing.T, expected []string, name string) {
	toolchain, _ := GetToolchain(name)
	assert.Equal(t, expected, toolchain.BuildCommandArgs())
}

func assertTestCommandPath(t *testing.T, expected string, name string) {
	toolchain, _ := GetToolchain(name)
	assert.Equal(t, expected, toolchain.TestCommandPath())
}

func assertTestCommandArgs(t *testing.T, expected []string, name string) {
	toolchain, _ := GetToolchain(name)
	assert.Equal(t, expected, toolchain.TestCommandArgs())
}

func assertErrorWhenBuildFails(t *testing.T, name string, workDir string) {
	toolchain, _ := GetToolchain(name)
	runFromDir(t, workDir,
		func(t *testing.T) {
			assert.Error(t, toolchain.RunBuild())
		})
}

func assertNoErrorWhenBuildPasses(t *testing.T, name string, workDir string) {
	toolchain, _ := GetToolchain(name)
	runFromDir(t, workDir,
		func(t *testing.T) {
			assert.NoError(t, toolchain.RunBuild())
		})
}

func assertErrorWhenTestFails(t *testing.T, name string, workDir string) {
	toolchain, _ := GetToolchain(name)
	runFromDir(t, workDir,
		func(t *testing.T) {
			assert.Error(t, toolchain.RunTests())
		})
}

func assertNoErrorWhenTestPasses(t *testing.T, name string, workDir string) {
	toolchain, _ := GetToolchain(name)
	runFromDir(t, workDir,
		func(t *testing.T) {
			assert.NoError(t, toolchain.RunTests())
		})
}

func runFromDir(t *testing.T, workDir string, testFunction func(t *testing.T)) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	initialDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	testFunction(t)
	_ = os.Chdir(initialDir)
}

func assertRunsOnAllOsWithAmd64(t *testing.T, name string) {
	// We don't check all platforms, just verifying that at least all OS's with amd64 are supported
	for _, osName := range GetAllOsNames() {
		assertRunsOnPlatform(t, name, osName, ArchAmd64)
	}
}

func assertRunsOnPlatform(t *testing.T, name string, osName OsName, archName ArchName) {
	toolchain, _ := GetToolchain(name)
	assert.True(t, toolchain.runsOnPlatform(osName, archName))
}

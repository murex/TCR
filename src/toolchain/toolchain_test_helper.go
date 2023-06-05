//go:build test_helper

/*
Copyright (c) 2023 Murex

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
)

func assertToolchainName(t *testing.T, toolchainName string) {
	t.Helper()
	toolchain, _ := Get(toolchainName)
	assert.Equal(t, toolchainName, toolchain.GetName())
}

func assertBuildCommandPath(t *testing.T, toolchainName string, expected string) {
	t.Helper()
	toolchain, _ := Get(toolchainName)
	assert.Equal(t, expected, toolchain.BuildCommandPath())
}

func assertBuildCommandArgs(t *testing.T, toolchainName string, expected []string) {
	t.Helper()
	toolchain, _ := Get(toolchainName)
	assert.Equal(t, expected, toolchain.BuildCommandArgs())
}

func assertTestCommandPath(t *testing.T, toolchainName string, expected string) {
	t.Helper()
	toolchain, _ := Get(toolchainName)
	assert.Equal(t, expected, toolchain.TestCommandPath())
}

func assertTestCommandArgs(t *testing.T, toolchainName string, expected []string) {
	t.Helper()
	toolchain, _ := Get(toolchainName)
	assert.Equal(t, expected, toolchain.TestCommandArgs())
}

func assertTestResultDir(t *testing.T, toolchainName string, expected string) {
	t.Helper()
	toolchain, _ := Get(toolchainName)
	assert.Equal(t, expected, toolchain.GetTestResultDir())
}

func assertErrorWhenBuildFails(t *testing.T, toolchainName string, workDir string) {
	t.Helper()
	toolchain, _ := Get(toolchainName)
	runFromDir(t, workDir,
		func(t *testing.T) {
			assert.Equal(t, toolchain.RunBuild().Status, CommandStatusFail)
		})
}

func assertNoErrorWhenBuildPasses(t *testing.T, toolchainName string, workDir string) {
	t.Helper()
	toolchain, _ := Get(toolchainName)
	runFromDir(t, workDir,
		func(t *testing.T) {
			assert.Equal(t, toolchain.RunBuild().Status, CommandStatusPass)
		})
}

func assertErrorWhenTestFails(t *testing.T, toolchainName string, workDir string) {
	t.Helper()
	toolchain, _ := Get(toolchainName)
	runFromDir(t, workDir,
		func(t *testing.T) {
			assert.Equal(t, toolchain.RunTests().Status, CommandStatusFail)
		})
}

func assertNoErrorWhenTestPasses(t *testing.T, toolchainName string, workDir string) {
	t.Helper()
	toolchain, _ := Get(toolchainName)
	runFromDir(t, workDir,
		func(t *testing.T) {
			assert.Equal(t, toolchain.RunTests().Status, CommandStatusPass)
		})
}

func runFromDir(t *testing.T, workDir string, testFunction func(t *testing.T)) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	initialDir, _ := os.Getwd()
	err := os.Chdir(workDir)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_ = os.Chdir(initialDir)
	}()
	testFunction(t)
}

func assertRunsOnAllOsWithAmd64(t *testing.T, name string) {
	// We don't check all platforms, just verifying that at least all OS's with amd64 are supported
	t.Helper()
	for _, osName := range GetAllOsNames() {
		t.Run(string(osName), func(t *testing.T) {
			assertRunsOnPlatform(t, name, osName, ArchAmd64)
		})
	}
}

func assertRunsOnPlatform(t *testing.T, toolchainName string, osName OsName, archName ArchName) {
	t.Helper()
	toolchain, _ := Get(toolchainName)
	assert.True(t, toolchain.runsOnPlatform(osName, archName))
}

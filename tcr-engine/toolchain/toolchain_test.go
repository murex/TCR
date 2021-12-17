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
	"strings"
	"testing"
)

const (
	testDataRootDir = "../testdata"
)

var (
	testDataDirJava = filepath.Join(testDataRootDir, "java")
	//testDataDirCpp  = filepath.Join(testDataRootDir, "cpp")
)

// aToolchain is a test data builder for type Toolchain
func aToolchain(toolchainBuilders ...func(tchn *Toolchain)) *Toolchain {
	tchn := &Toolchain{
		Name:          "default-toolchain",
		BuildCommands: []Command{*aCommand(withPath("build-path"))},
		TestCommands:  []Command{*aCommand(withPath("test-path"))},
	}

	for _, build := range toolchainBuilders {
		build(tchn)
	}
	return tchn
}

func withName(name string) func(tchn *Toolchain) {
	return func(tchn *Toolchain) { tchn.Name = name }
}

func withNoBuildCommand() func(tchn *Toolchain) {
	return func(tchn *Toolchain) { tchn.BuildCommands = nil }
}

func withNoTestCommand() func(tchn *Toolchain) {
	return func(tchn *Toolchain) { tchn.TestCommands = nil }
}

func runFromDir(t *testing.T, testDir string, testFunction func(t *testing.T)) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	initialDir, _ := os.Getwd()
	_ = os.Chdir(testDir)
	testFunction(t)
	_ = os.Chdir(initialDir)
}

func Test_does_not_support_empty_toolchain_name(t *testing.T) {
	assert.False(t, isSupported(""))
}

func Test_does_not_support_unregistered_toolchain_name(t *testing.T) {
	assert.False(t, isSupported("unregistered-toolchain"))
}

func Test_unrecognized_toolchain_name(t *testing.T) {
	toolchain, err := Get("dummy")
	assert.Error(t, err)
	assert.Zero(t, toolchain)
}

func Test_can_add_a_built_in_toolchain(t *testing.T) {
	const name = "new-built-in-toolchain"
	assert.False(t, isBuiltIn(name))
	assert.NoError(t, addBuiltIn(*aToolchain(withName(name))))
	assert.True(t, isBuiltIn(name))
}

func Test_cannot_add_a_built_in_toolchain_with_no_name(t *testing.T) {
	assert.Error(t, addBuiltIn(*aToolchain(withName(""))))
}

func Test_toolchain_name_is_case_insensitive(t *testing.T) {
	const name = "miXeD-CasE"
	_ = Register(*aToolchain(withName(name)))
	assert.True(t, isSupported(name))
	assert.True(t, isSupported(strings.ToUpper(name)))
	assert.True(t, isSupported(strings.ToLower(name)))
	assert.True(t, isSupported(strings.Title(name)))
}

func Test_can_register_a_toolchain(t *testing.T) {
	const name = "new-toolchain"
	assert.False(t, isSupported(name))
	assert.NoError(t, Register(*aToolchain(withName(name))))
	assert.True(t, isSupported(name))
}

func Test_cannot_register_a_toolchain_with_no_name(t *testing.T) {
	assert.Error(t, Register(*aToolchain(withName(""))))
}

func Test_cannot_register_a_toolchain_with_no_build_command(t *testing.T) {
	const name = "no-build-command"
	assert.Error(t, Register(*aToolchain(withName(name), withNoBuildCommand())))
	assert.False(t, isSupported(name))
}

func Test_cannot_register_a_toolchain_with_no_test_command(t *testing.T) {
	const name = "no-test-command"
	assert.Error(t, Register(*aToolchain(withName(name), withNoTestCommand())))
	assert.False(t, isSupported(name))
}

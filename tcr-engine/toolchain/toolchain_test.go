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

func Test_does_not_support_dummy_toolchain_name(t *testing.T) {
	assert.False(t, isSupported("dummy"))
}

func Test_unrecognized_toolchain_name(t *testing.T) {
	toolchain, err := GetToolchain("dummy")
	assert.Zero(t, toolchain)
	assert.NotZero(t, err)
}

func Test_can_add_a_built_in_toolchain(t *testing.T) {
	tchn := Toolchain{
		Name:          "dummy",
		BuildCommands: nil,
		TestCommands:  nil,
	}
	assert.False(t, isBuiltIn("dummy"))
	err := addBuiltInToolchain(tchn)
	assert.Zero(t, err)
	assert.True(t, isBuiltIn("dummy"))
}

func Test_cannot_add_a_built_in_toolchain_with_no_name(t *testing.T) {
	var tchn = Toolchain{}
	err := addBuiltInToolchain(tchn)
	assert.NotZero(t, err)
}

func Test_toolchain_name_is_case_insensitive(t *testing.T) {
	tchn := Toolchain{
		Name:          "dummy",
		BuildCommands: nil,
		TestCommands:  nil,
	}
	_ = addBuiltInToolchain(tchn)
	assert.True(t, isSupported("dummy"))
	assert.True(t, isSupported("Dummy"))
	assert.True(t, isSupported("DUMMY"))
}

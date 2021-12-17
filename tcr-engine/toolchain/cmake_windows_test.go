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
	"runtime"
	"testing"
)

func Test_cmake_toolchain_build_command_path_on_windows(t *testing.T) {
	toolchain, _ := Get(cmakeToolchainName)
	var expected string
	switch runtime.GOARCH {
	case ArchAmd64:
		expected = "build\\cmake\\cmake-windows-x86_64\\bin\\cmake.exe"
	case Arch386:
		expected = "build\\cmake\\cmake-windows-i386\\bin\\cmake.exe"
	default:
		t.Error("Architecture ", runtime.GOARCH, " is not supported by cmake on ", runtime.GOOS)
	}
	assert.Equal(t, expected, toolchain.buildCommandPath())
}

func Test_cmake_toolchain_test_command_path_on_windows(t *testing.T) {
	toolchain, _ := Get(cmakeToolchainName)
	var expected string
	switch runtime.GOARCH {
	case ArchAmd64:
		expected = "build\\cmake\\cmake-windows-x86_64\\bin\\ctest.exe"
	case Arch386:
		expected = "build\\cmake\\cmake-windows-i386\\bin\\ctest.exe"
	default:
		t.Error("Architecture ", runtime.GOARCH, " is not supported by cmake on ", runtime.GOOS)
	}
	assert.Equal(t, expected, toolchain.testCommandPath())
}

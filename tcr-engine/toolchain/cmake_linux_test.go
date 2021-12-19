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
	"runtime"
	"testing"
)

func Test_cmake_toolchain_build_command_path_on_linux(t *testing.T) {
	switch runtime.GOARCH {
	case ArchAmd64:
		assertBuildCommandPath(t, "build/cmake/cmake-linux-x86_64/bin/cmake", cmakeToolchainName)
	case ArchArm64:
		assertBuildCommandPath(t, "build/cmake/cmake-linux-aarch64/bin/cmake", cmakeToolchainName)
	default:
		t.Error("Architecture ", runtime.GOARCH, " is not supported by cmake on ", runtime.GOOS)
	}
}

func Test_cmake_toolchain_test_command_path_on_linux(t *testing.T) {
	switch runtime.GOARCH {
	case ArchAmd64:
		assertTestCommandPath(t, "build/cmake/cmake-linux-x86_64/bin/ctest", cmakeToolchainName)
	case ArchArm64:
		assertTestCommandPath(t, "build/cmake/cmake-linux-aarch64/bin/ctest", cmakeToolchainName)
	default:
		t.Error("Architecture ", runtime.GOARCH, " is not supported by ctest on ", runtime.GOOS)
	}
}

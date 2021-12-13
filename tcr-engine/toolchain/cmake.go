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

import "path/filepath"

func init() {
	// TODO refactor initialization to remove redundancies
	// TODO add other architectures than amd64

	var buildWindows = Command{
		Os:        []OsName{OsWindows},
		Arch:      []ArchName{ArchAmd64},
		Path:      filepath.Join("build", "cmake", "cmake-windows-x86_64", "bin", "cmake.exe"),
		Arguments: []string{"--build", "build", "--config", "Debug"},
	}

	var testWindows = Command{
		Os:        []OsName{OsWindows},
		Arch:      []ArchName{ArchAmd64},
		Path:      filepath.Join("build", "cmake", "cmake-windows-x86_64", "bin", "ctest.exe"),
		Arguments: []string{"--output-on-failure", "--test-dir", "build", "--build-config", "Debug"},
	}

	var buildLinux = Command{
		Os:        []OsName{OsLinux},
		Arch:      []ArchName{ArchAmd64},
		Path:      "build/cmake/cmake-linux-x86_64/bin/cmake",
		Arguments: []string{"--build", "build", "--config", "Debug"},
	}

	var testLinux = Command{
		Os:        []OsName{OsLinux},
		Arch:      []ArchName{ArchAmd64},
		Path:      "build/cmake/cmake-linux-x86_64/bin/ctest",
		Arguments: []string{"--output-on-failure", "--test-dir", "build", "--build-config", "Debug"},
	}

	var buildDarwin = Command{
		Os:        []OsName{OsDarwin},
		Arch:      []ArchName{ArchAmd64},
		Path:      "build/cmake/cmake-macos-universal/CMake.app/Contents/bin/cmake",
		Arguments: []string{"--build", "build", "--config", "Debug"},
	}

	var testDarwin = Command{
		Os:        []OsName{OsDarwin},
		Arch:      []ArchName{ArchAmd64},
		Path:      "build/cmake/cmake-macos-universal/CMake.app/Contents/bin/ctest",
		Arguments: []string{"--output-on-failure", "--test-dir", "build", "--build-config", "Debug"},
	}

	_ = addBuiltInToolchain(
		Toolchain{
			Name:          "cmake",
			BuildCommands: []Command{buildDarwin, buildLinux, buildWindows},
			TestCommands:  []Command{testDarwin, testLinux, testWindows},
		},
	)
}

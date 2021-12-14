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

func init() {
	// TODO add other architectures than amd64

	_ = addBuiltInToolchain(
		Toolchain{
			Name: "cmake",
			BuildCommands: []Command{
				initCmakeCommand(OsDarwin, ArchAmd64),
				initCmakeCommand(OsLinux, ArchAmd64),
				initCmakeCommand(OsWindows, ArchAmd64),
			},
			TestCommands: []Command{
				initCtestCommand(OsDarwin, ArchAmd64),
				initCtestCommand(OsLinux, ArchAmd64),
				initCtestCommand(OsWindows, ArchAmd64),
			},
		},
	)
}

func initCmakeCommand(osName OsName, archName ArchName) Command {
	return Command{
		Os:        []OsName{osName},
		Arch:      []ArchName{archName},
		Path:      initCommandPath("cmake", osName, archName),
		Arguments: cmakeArguments(),
	}
}

func initCtestCommand(osName OsName, archName ArchName) Command {
	return Command{
		Os:        []OsName{osName},
		Arch:      []ArchName{archName},
		Path:      initCommandPath("ctest", osName, archName),
		Arguments: ctestArguments(),
	}
}

func initCommandPath(cmd string, osName OsName, archName ArchName) string {
	var cmakeOsArchDirName = "cmake-" + osNameInPath(osName) + "-" + archNameInPath(osName, archName)

	switch osName {
	case OsDarwin:
		return "build/cmake/" + cmakeOsArchDirName + "/CMake.app/Contents/bin/" + cmd
	case OsLinux:
		return "build/cmake/" + cmakeOsArchDirName + "/bin/" + cmd
	case OsWindows:
		return "build\\cmake\\" + cmakeOsArchDirName + "\\bin\\" + cmd + ".exe"
	default:
		return ""
	}
}

func archNameInPath(osName OsName, archName ArchName) string {
	switch {
	case osName == OsDarwin:
		return "universal"
	case archName == ArchAmd64:
		return "x86_64"
	case archName == Arch386:
		return "i386"
	case archName == ArchArm64:
		return "aarch64"
	default:
		return string(archName)
	}
}

func osNameInPath(osName OsName) string {
	if osName == OsDarwin {
		return "macos"
	}
	return string(osName)
}

func cmakeArguments() []string {
	return []string{"--build", "build", "--config", "Debug"}
}

func ctestArguments() []string {
	return []string{"--output-on-failure", "--test-dir", "build", "--build-config", "Debug"}
}

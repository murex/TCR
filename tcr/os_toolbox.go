package tcr

import "path"

type OsToolbox interface {
	fsWatchCommand() string
	cmakeBinPath() string
	cmakeCommand() string
	ctestCommand() string
}

// ========================================================================

type MacOSToolbox struct {
}

func (OsToolbox MacOSToolbox) fsWatchCommand() string {
	return "fswatch -1 -r"
}

func (OsToolbox MacOSToolbox) cmakeBinPath() string {
	return "./cmake/cmake-macos-universal/CMake.app/Contents/bin"
}

func (OsToolbox MacOSToolbox) cmakeCommand() string {
	return path.Join(OsToolbox.cmakeBinPath(), "cmake")
}

func (OsToolbox MacOSToolbox) ctestCommand() string {
	return path.Join(OsToolbox.cmakeBinPath(), "ctest")
}

// ========================================================================

type LinuxToolbox struct {
}

func (OsToolbox LinuxToolbox) fsWatchCommand() string {
	return "inotifywait -r -e modify"
}

func (OsToolbox LinuxToolbox) cmakeBinPath() string {
	return "./cmake/cmake-LinuxToolbox-x86_64/bin"
}

func (OsToolbox LinuxToolbox) cmakeCommand() string {
	return path.Join(OsToolbox.cmakeBinPath(), "cmake")
}

func (OsToolbox LinuxToolbox) ctestCommand() string {
	return path.Join(OsToolbox.cmakeBinPath(), "ctest")
}

// ========================================================================

type WindowsToolbox struct {
}

func (OsToolbox WindowsToolbox) fsWatchCommand() string {
	return path.Join(ScriptDir(), "inotify-win.exe") + " -r -e modify"
}

func (OsToolbox WindowsToolbox) cmakeBinPath() string {
	return "./cmake/cmake-win64-x64/bin"
}

func (OsToolbox WindowsToolbox) cmakeCommand() string {
	return path.Join(OsToolbox.cmakeBinPath(), "cmake.exe")
}

func (OsToolbox WindowsToolbox) ctestCommand() string {
	return path.Join(OsToolbox.cmakeBinPath(), "ctest.exe")
}
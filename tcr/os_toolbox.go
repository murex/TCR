package tcr

import (
	"github.com/mengdaming/tcr/trace"
	"path/filepath"
	"runtime"
)

type OSToolbox interface {
	fsWatchCommand() string
	cmakeBinPath() string
	cmakeCommand() string
	ctestCommand() string
}

// ========================================================================

type MacOSToolbox struct {
}

func (osToolbox MacOSToolbox) fsWatchCommand() string {
	return "fswatch -1 -r"
}

func (osToolbox MacOSToolbox) cmakeBinPath() string {
	return "./cmake/cmake-macos-universal/CMake.app/Contents/bin"
}

func (osToolbox MacOSToolbox) cmakeCommand() string {
	return filepath.Join(osToolbox.cmakeBinPath(), "cmake")
}

func (osToolbox MacOSToolbox) ctestCommand() string {
	return filepath.Join(osToolbox.cmakeBinPath(), "ctest")
}

// ========================================================================

type LinuxToolbox struct {
}

func (osToolbox LinuxToolbox) fsWatchCommand() string {
	return "inotifywait -r -e modify"
}

func (osToolbox LinuxToolbox) cmakeBinPath() string {
	return "./cmake/cmake-LinuxToolbox-x86_64/bin"
}

func (osToolbox LinuxToolbox) cmakeCommand() string {
	return filepath.Join(osToolbox.cmakeBinPath(), "cmake")
}

func (osToolbox LinuxToolbox) ctestCommand() string {
	return filepath.Join(osToolbox.cmakeBinPath(), "ctest")
}

// ========================================================================

type WindowsToolbox struct {
}

func (osToolbox WindowsToolbox) fsWatchCommand() string {
	return filepath.Join(ScriptDir(), "inotify-win.exe") + " -r -e modify"
}

func (osToolbox WindowsToolbox) cmakeBinPath() string {
	return "./cmake/cmake-win64-x64/bin"
}

func (osToolbox WindowsToolbox) cmakeCommand() string {
	return filepath.Join(osToolbox.cmakeBinPath(), "cmake.exe")
}

func (osToolbox WindowsToolbox) ctestCommand() string {
	return filepath.Join(osToolbox.cmakeBinPath(), "ctest.exe")
}

func initOSToolbox() OSToolbox {
	toolbox := osToolbox
	switch runtime.GOOS {
	case "darwin":
		toolbox = MacOSToolbox{}
	case "linux":
		toolbox = LinuxToolbox{}
	case "windows":
		toolbox = WindowsToolbox{}
	default:
		trace.Error("OS ", runtime.GOOS, " is currently not supported")
	}
	return toolbox
}

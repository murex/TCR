package tcr

import (
	"fmt"
	"github.com/mengdaming/tcr/trace"
	"runtime"
)

var osToolbox OsToolbox
var language Language

func Start() {
	initOsToolbox()
	detectKataLanguage()
	detectGitWorkingBranch()
	whatShallWeDo()
}

func whatShallWeDo() {

	trace.HorizontalLine()
	h1 := fmt.Sprintf("Language=%v, Toolchain=%v", language.name(), language.toolchain())
	trace.Info(h1)
	// TODO
}

func detectGitWorkingBranch() {
	// TODO
}

func detectKataLanguage() {
	// TODO Add C++ case
	language = JavaLanguage{}
}

func initOsToolbox() {
	switch runtime.GOOS {
	case "darwin":
		osToolbox = MacOSToolbox{}
	case "linux":
		osToolbox = LinuxToolbox{}
	case "windows":
		osToolbox = WindowsToolbox{}
	default:
		message := fmt.Sprintf("OS %v is currently not supported", runtime.GOOS)
		trace.Error(message)
	}
}

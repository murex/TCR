package tcr

import (
	"fmt"
	"github.com/mengdaming/tcr/trace"
	"runtime"
)

var osToolbox OSToolbox
var language Language
var toolchain string
var autoPush bool
var gitWorkingBranch string

func Start(t string, ap bool) {
	toolchain = t
	autoPush = ap
	initOSToolbox()
	detectKataLanguage()
	detectGitWorkingBranch()
	whatShallWeDo()
}

func whatShallWeDo() {
	printTraceHeader()
	// TODO
}

func printTraceHeader() {
	trace.HorizontalLine()

	th1 := fmt.Sprintf("Language=%v, Toolchain=%v", language.name(), language.toolchain())
	trace.Info(th1)

	var autoPushStr string
	if autoPush {
		autoPushStr = "enabled"
	} else {
		autoPushStr = "disabled"
	}
	th2 := fmt.Sprintf("Running on git branch \"%v\" with auto-push %v", gitWorkingBranch, autoPushStr)
	trace.Info(th2)

}

func detectGitWorkingBranch() {
	// TODO Hardcoded for now. Replace with branch detection
	gitWorkingBranch = "main"
}

func detectKataLanguage() {
	// TODO Add language detection. Hardcoding java for now
	language = JavaLanguage{}
}

func initOSToolbox() {
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

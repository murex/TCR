package tcr

import (
	"fmt"
	"github.com/mengdaming/tcr/trace"
	"runtime"
)

var (
	osToolbox        OSToolbox
	language         Language
	toolchain        string
	autoPush         bool
	gitWorkingBranch string
)

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
	trace.HorizontalLine()
	trace.Info("What shall we do?")
	trace.Info("\tD -> Driver mode")
	trace.Info("\tN -> Navigator mode")
	trace.Info("\tQ -> Quit")

}

func printTraceHeader() {
	trace.HorizontalLine()

	th1 := fmt.Sprintf("Language=%v, Toolchain=%v",
		language.name(), language.toolchain())
	trace.Info(th1)

	autoPushStr := "disabled"
	if autoPush {
		autoPushStr = "enabled"
	}
	th2 := fmt.Sprintf("Running on git branch \"%v\" with auto-push %v",
		gitWorkingBranch, autoPushStr)
	trace.Info(th2)
}

func detectGitWorkingBranch() {
	// TODO Hardcoded for now. Replace with branch detection
	gitWorkingBranch = "main"
}

func detectKataLanguage() {
	// TODO Add language detection. Hard-coding java for now
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

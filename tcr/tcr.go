package tcr

import (
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

	trace.Info(
		"Language:", language.name(),
		"Toolchain:", language.toolchain())

	autoPushStr := "disabled"
	if autoPush {
		autoPushStr = "enabled"
	}
	trace.Info(
		"Running on git branch", gitWorkingBranch,
		"with auto-push", autoPushStr)
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
		trace.Error("OS", runtime.GOOS, "is currently not supported")
	}
}

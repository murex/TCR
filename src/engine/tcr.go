/*
Copyright (c) 2023 Murex

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

package engine

import (
	"github.com/murex/tcr/checker"
	"github.com/murex/tcr/events"
	"github.com/murex/tcr/filesystem"
	"github.com/murex/tcr/language"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/role"
	"github.com/murex/tcr/runmode"
	"github.com/murex/tcr/settings"
	"github.com/murex/tcr/stats"
	"github.com/murex/tcr/status"
	"github.com/murex/tcr/timer"
	"github.com/murex/tcr/toolchain"
	"github.com/murex/tcr/ui"
	"github.com/murex/tcr/vcs"
	"github.com/murex/tcr/vcs/factory"
	"gopkg.in/tomb.v2"
	"os"
	"strings"
	"time"
)

type (
	// TCRInterface provides the API for interacting with TCR engine
	TCRInterface interface {
		Init(u ui.UserInterface, p params.Params)
		setVCS(vcsInterface vcs.Interface)
		ToggleAutoPush()
		SetAutoPush(flag bool)
		SetCommitOnFail(flag bool)
		GetCurrentRole() role.Role
		RunAsDriver()
		RunAsNavigator()
		Stop()
		RunTCRCycle()
		GetSessionInfo() SessionInfo
		ReportMobTimerStatus()
		SetRunMode(m runmode.RunMode)
		RunCheck(p params.Params)
		PrintLog(p params.Params)
		PrintStats(p params.Params)
		VCSPull()
		VCSPush()
		Quit()
	}

	// TCREngine is the engine running all TCR operations
	TCREngine struct {
		mode            runmode.RunMode
		ui              ui.UserInterface
		vcs             vcs.Interface
		language        language.LangInterface
		toolchain       toolchain.TchnInterface
		sourceTree      filesystem.SourceTree
		pollingPeriod   time.Duration
		mobTurnDuration time.Duration
		mobTimer        *timer.PeriodicReminder
		currentRole     role.Role
		commitOnFail    bool
		messageSuffix   string
		// shoot channel is used for handling interruptions coming from the UI
		shoot chan bool
		// fsWatchRearmDelay is the waiting time until TCR starts watching the filesystem again
		// after a filesystem event was detected. The default value should not be changed except
		// when running tests
		fsWatchRearmDelay time.Duration
	}
)

const traceReporterWaitingTime = 100 * time.Millisecond

const fsWatchRearmDelay = 2 * time.Second

const (
	commitMessageOk     = "✅ TCR - tests passing"
	commitMessageFail   = "❌ TCR - tests failing"
	commitMessageRevert = "⏪ TCR - revert changes"
	buildFailureMessage = "There are build errors! I can't go any further"
	testFailureMessage  = "Some tests are failing! That's unfortunate"
	testSuccessMessage  = "Tests passed!"
)

var (
	// TCR is TCR Engine singleton instance
	TCR TCRInterface
)

// NewTCREngine instantiates TCR engine instance
func NewTCREngine() (engine *TCREngine) {
	engine = &TCREngine{fsWatchRearmDelay: fsWatchRearmDelay}
	TCR = engine
	return engine
}

// Init initializes the TCR engine with the provided parameters, and wires it to the user interface.
// This function should be called only once during the lifespan of the application
func (tcr *TCREngine) Init(u ui.UserInterface, p params.Params) {
	var err error
	status.RecordState(status.Ok)
	tcr.ui = u

	report.PostInfo("Starting ", settings.ApplicationName, " version ", settings.BuildVersion, "...")

	tcr.SetRunMode(p.Mode)
	if !tcr.mode.IsActive() {
		tcr.ui.ShowRunningMode(tcr.mode)
		return
	}

	tcr.pollingPeriod = p.PollingPeriod

	tcr.initSourceTree(p)

	tcr.language, err = language.GetLanguage(p.Language, tcr.sourceTree.GetBaseDir())
	tcr.handleError(err, true, status.ConfigError)

	tcr.toolchain, err = tcr.language.GetToolchain(p.Toolchain)
	tcr.handleError(err, true, status.ConfigError)

	err = toolchain.SetWorkDir(p.WorkDir)
	tcr.handleError(err, true, status.ConfigError)
	report.PostInfo("Work directory is ", toolchain.GetWorkDir())

	tcr.initVCS(p.VCS, p.Trace)
	tcr.setMessageSuffix(p.MessageSuffix)
	tcr.vcs.EnablePush(p.AutoPush)

	tcr.SetCommitOnFail(p.CommitFailures)

	tcr.setMobTimerDuration(p.MobTurnDuration)

	tcr.ui.ShowRunningMode(tcr.mode)
	tcr.ui.ShowSessionInfo()
	tcr.warnIfOnRootBranch(tcr.mode.IsInteractive())
}

// SetCommitOnFail sets VCS commit-on-fail option to the provided value
func (tcr *TCREngine) SetCommitOnFail(flag bool) {
	tcr.commitOnFail = flag
	if tcr.commitOnFail {
		report.PostInfo("Test-breaking changes will be committed")
	} else {
		report.PostInfo("Test-breaking changes will not be committed")
	}
}

func (tcr *TCREngine) setMobTimerDuration(duration time.Duration) {
	if settings.EnableMobTimer && tcr.mode.NeedsCountdownTimer() {
		tcr.mobTurnDuration = duration
		report.PostInfo("Timer duration is ", tcr.mobTurnDuration)
	}
}

// RunCheck checks the provided parameters and prints out corresponding report
func (*TCREngine) RunCheck(p params.Params) {
	checker.Run(p)
}

// PrintLog prints the TCR VCS commit history
func (tcr *TCREngine) PrintLog(p params.Params) {
	tcrLogs := tcr.queryVCSLogs(p)
	report.PostInfo("Printing TCR log for ", tcr.vcs.SessionSummary())
	for _, log := range tcrLogs {
		report.PostTitle("commit:    ", log.Hash)
		report.PostInfo("timestamp: ", log.Timestamp)
		report.PostInfo("message:   ", log.Message)
		// Giving trace reporter some time to flush its contents
		time.Sleep(traceReporterWaitingTime)
	}
}

// PrintStats prints the TCR execution stats
func (tcr *TCREngine) PrintStats(p params.Params) {
	tcrLogs := tcr.queryVCSLogs(p)
	stats.Print(tcr.vcs.SessionSummary(), tcrLogsToEvents(tcrLogs))
}

func tcrLogsToEvents(tcrLogs vcs.LogItems) (tcrEvents events.TcrEvents) {
	tcrEvents = *events.NewTcrEvents()
	for _, log := range tcrLogs {
		tcrEvents.Add(log.Timestamp, parseCommitMessage(log.Message))
	}
	return tcrEvents
}

func (tcr *TCREngine) queryVCSLogs(p params.Params) vcs.LogItems {
	tcr.initSourceTree(p)
	tcr.initVCS(p.VCS, p.Trace)

	logs, err := tcr.vcs.Log(isTCRCommitMessage)
	if err != nil {
		report.PostError(err)
	}
	if len(logs) == 0 {
		report.PostWarning("no TCR commit found in ", tcr.vcs.SessionSummary(), "'s history")
	}
	return logs
}

func isTCRCommitMessage(msg string) bool {
	return strings.Index(msg, commitMessageOk) == 0 || strings.Index(msg, commitMessageFail) == 0
}

func parseCommitMessage(message string) (event events.TCREvent) {
	// First line is the main commit message
	// Second line is a blank line
	// The YAML-structured data starts on the third line until we reach a blank line
	// The user-specified message prefix, if any, is after the blank line

	var header string
	var statsYAML strings.Builder
	var section = 1
	for _, line := range strings.Split(message, "\n") {
		switch section {
		case 1: // main commit message
			header = line
			section++
		case 2: // blank line between header and TCR event stats
			section++
		case 3: // YAML-structured data containing TCR event stats
			if line == "" {
				// First empty line or end of message should mark the end of YAML data
				section++
			} else {
				_, _ = statsYAML.WriteString(line)
				_, _ = statsYAML.WriteRune('\n')
			}
		case 4: // commit message suffix, if any
			// Ignoring commit message suffix for now. May be useful in the future if
			// we want to filter commit history based on its contents
			continue
		}
	}

	event = events.FromYAML(statsYAML.String())
	switch header {
	case commitMessageOk:
		event.Status = events.StatusPass
	case commitMessageFail:
		event.Status = events.StatusFail
	default:
		event.Status = events.StatusUnknown
	}
	return event
}

func (tcr *TCREngine) setMessageSuffix(suffix string) {
	tcr.messageSuffix = suffix
}

func (tcr *TCREngine) wrapCommitMessages(statusMessage string, event *events.TCREvent) []string {
	messages := []string{statusMessage}
	if event != nil {
		messages = append(messages, event.ToYAML())
	}
	if tcr.messageSuffix != "" {
		messages = append(messages, "\n"+tcr.messageSuffix)
	}
	return messages
}

func (tcr *TCREngine) initVCS(vcsName string, trace string) {
	if tcr.vcs != nil {
		return // VCS should be initialized only once
	}
	// Set VCS trace trace flag
	vcs.SetTrace(trace == "vcs")
	var err error
	tcr.vcs, err = factory.InitVCS(vcsName, tcr.sourceTree.GetBaseDir())
	switch err.(type) {
	case *factory.UnsupportedVCSError:
		tcr.handleError(err, true, status.ConfigError)
	default:
		tcr.handleError(err, true, status.VCSError)
	}
}

func (tcr *TCREngine) initSourceTree(p params.Params) {
	var err error
	tcr.sourceTree, err = filesystem.New(p.BaseDir)
	tcr.handleError(err, true, status.ConfigError)
	report.PostInfo("Base directory is ", tcr.sourceTree.GetBaseDir())
}

func (tcr *TCREngine) setVCS(vcsInterface vcs.Interface) {
	tcr.vcs = vcsInterface
}

func (tcr *TCREngine) warnIfOnRootBranch(interactive bool) {
	if tcr.vcs.IsOnRootBranch() {
		message := "Running " + settings.ApplicationName +
			" on " + tcr.vcs.SessionSummary() + " is not recommended"
		if interactive {
			if !tcr.ui.Confirm(message, false) {
				tcr.Quit()
			}
		} else {
			report.PostWarning(message)
		}
	}
}

// ToggleAutoPush toggles VCS auto-push state
func (tcr *TCREngine) ToggleAutoPush() {
	tcr.vcs.EnablePush(!tcr.vcs.IsPushEnabled())
}

// SetAutoPush sets VCS auto-push to the provided value
func (tcr *TCREngine) SetAutoPush(ap bool) {
	tcr.vcs.EnablePush(ap)
}

// GetCurrentRole returns the role currently used for running TCR.
// Returns nil when TCR engine is in standby
func (tcr *TCREngine) GetCurrentRole() role.Role {
	return tcr.currentRole
}

// RunAsDriver tells TCR engine to start running with driver role
func (tcr *TCREngine) RunAsDriver() {
	tcr.initTimer()

	go tcr.fromBirthTillDeath(
		func() {
			tcr.currentRole = role.Driver{}
			tcr.ui.NotifyRoleStarting(tcr.currentRole)
			tcr.handleError(tcr.vcs.Pull(), false, status.VCSError)
			tcr.startTimer()
		},
		func(interrupt <-chan bool) bool {
			if tcr.waitForChange(interrupt) {
				// Some file changes were detected
				tcr.RunTCRCycle()
				return true
			}
			// If we arrive here this means that the end of waitForChange
			// was triggered by the user
			return false
		},
		func() {
			tcr.stopTimer()
			tcr.ui.NotifyRoleEnding(tcr.currentRole)
			tcr.currentRole = nil
		},
	)
}

// RunAsNavigator tells TCR engine to start running with navigator role
func (tcr *TCREngine) RunAsNavigator() {
	go tcr.fromBirthTillDeath(
		func() {
			tcr.currentRole = role.Navigator{}
			tcr.ui.NotifyRoleStarting(tcr.currentRole)
		},
		func(interrupt <-chan bool) bool {
			select {
			case <-interrupt:
				return false
			default:
				tcr.handleError(tcr.vcs.Pull(), false, status.VCSError)
				time.Sleep(tcr.pollingPeriod)
				return true
			}
		},
		func() {
			tcr.ui.NotifyRoleEnding(tcr.currentRole)
			tcr.currentRole = nil
		},
	)
}

// Stop is the entry point for telling TCR engine to stop its current operations
func (tcr *TCREngine) Stop() {
	tcr.shoot <- true
}

func (tcr *TCREngine) fromBirthTillDeath(
	birth func(),
	dailyLife func(interrupt <-chan bool) bool,
	death func(),
) {
	var tmb tomb.Tomb
	tcr.shoot = make(chan bool)

	// The goroutine doing the work
	tmb.Go(func() error {
		birth()
		for oneMoreDay := true; oneMoreDay; {
			oneMoreDay = dailyLife(tcr.shoot)
		}
		death()
		return nil
	})
	tcr.handleError(tmb.Wait(), true, status.OtherError)
}

func (tcr *TCREngine) waitForChange(interrupt <-chan bool) bool {
	report.PostInfo("Going to sleep until something interesting happens")
	// We need to wait a bit to make sure the file watcher
	// does not get triggered again following a revert operation
	time.Sleep(tcr.fsWatchRearmDelay)
	return tcr.sourceTree.Watch(
		tcr.language.DirsToWatch(tcr.sourceTree.GetBaseDir()),
		tcr.language.IsLanguageFile,
		interrupt)
}

// RunTCRCycle is the core of TCR engine: e.g. it runs one test && commit || revert cycle
func (tcr *TCREngine) RunTCRCycle() {
	status.RecordState(status.Ok)
	if tcr.build().Failed() {
		return
	}
	result := tcr.test()
	event := tcr.createTCREvent(result)
	if result.Passed() {
		tcr.commit(event)
	} else {
		tcr.revert(event)
	}
}

func (tcr *TCREngine) createTCREvent(testResult toolchain.TestCommandResult) (event events.TCREvent) {
	diffs, err := tcr.vcs.Diff()
	if err != nil {
		report.PostWarning(err)
	}
	commandStatus := events.StatusFail
	if testResult.Passed() {
		commandStatus = events.StatusPass
	}
	return events.NewTCREvent(
		commandStatus,
		events.NewChangedLines(
			diffs.ChangedLines(tcr.language.IsSrcFile),
			diffs.ChangedLines(tcr.language.IsTestFile),
		),
		events.NewTestStats(
			testResult.Stats.TotalRun,
			testResult.Stats.Passed,
			testResult.Stats.Failed,
			testResult.Stats.Skipped,
			testResult.Stats.WithErrors,
			testResult.Stats.Duration,
		),
	)
}

func (tcr *TCREngine) build() (result toolchain.CommandResult) {
	report.PostInfo("Launching Build")
	result = tcr.toolchain.RunBuild()
	if result.Failed() {
		status.RecordState(status.BuildFailed)
		report.PostWarningWithEmphasis(buildFailureMessage)
	}
	return result
}

func (tcr *TCREngine) test() (result toolchain.TestCommandResult) {
	report.PostInfo("Running Tests")
	result = tcr.toolchain.RunTests()
	if result.Failed() {
		status.RecordState(status.TestFailed)
		report.PostErrorWithEmphasis(testFailureMessage)
	} else {
		report.PostSuccessWithEmphasis(testSuccessMessage)
	}
	return result
}

func (tcr *TCREngine) commit(event events.TCREvent) {
	report.PostInfo("Committing changes on ", tcr.vcs.SessionSummary())
	var err error
	err = tcr.vcs.Add()
	tcr.handleError(err, false, status.VCSError)
	if err != nil {
		return
	}
	err = tcr.vcs.Commit(false, tcr.wrapCommitMessages(commitMessageOk, &event)...)
	tcr.handleError(err, false, status.VCSError)
	if err != nil {
		return
	}
	tcr.handleError(tcr.vcs.Push(), false, status.VCSError)
}

func (tcr *TCREngine) revert(event events.TCREvent) {
	if tcr.commitOnFail {
		err := tcr.commitTestBreakingChanges(event)
		tcr.handleError(err, false, status.VCSError)
		if err != nil {
			return
		}
		tcr.handleError(tcr.vcs.Push(), false, status.VCSError)
	}
	tcr.revertSrcFiles()
}

func (tcr *TCREngine) commitTestBreakingChanges(event events.TCREvent) (err error) {
	// Create stash with the changes
	err = tcr.vcs.Stash(commitMessageFail)
	if err != nil {
		return err
	}
	// Apply changes back in the working tree
	err = tcr.vcs.UnStash(true)
	if err != nil {
		return err
	}
	// Commit changes with failure message into VCS index
	err = tcr.vcs.Add()
	if err != nil {
		return err
	}
	err = tcr.vcs.Commit(false, tcr.wrapCommitMessages(commitMessageFail, &event)...)
	if err != nil {
		return err
	}
	// Revert changes (both in VCS index and working tree)
	err = tcr.vcs.Revert()
	if err != nil {
		return err
	}
	// Amend commit message on revert operation in VCS index
	err = tcr.vcs.Commit(false, tcr.wrapCommitMessages(commitMessageRevert, nil)...)
	if err != nil {
		return err
	}
	// Re-apply changes in the working tree and get rid of stash
	err = tcr.vcs.UnStash(false)
	return err
}

func (tcr *TCREngine) revertSrcFiles() {
	diffs, err := tcr.vcs.Diff()
	tcr.handleError(err, false, status.VCSError)
	if err != nil {
		return
	}
	var reverted int
	for _, diff := range diffs {
		if tcr.language.IsSrcFile(diff.Path) {
			err := tcr.revertFile(diff.Path)
			tcr.handleError(err, false, status.VCSError)
			if err == nil {
				reverted++
			}
		}
	}
	if reverted > 0 {
		report.PostWarning(reverted, " file(s) reverted")
	} else {
		report.PostInfo("No file reverted (only test files were updated since last commit)")
	}
}

func (tcr *TCREngine) revertFile(file string) error {
	return tcr.vcs.Restore(file)
}

// GetSessionInfo provides the information related to the current TCR session.
// Used mainly by the user interface packages to retrieve and display this information
func (tcr *TCREngine) GetSessionInfo() SessionInfo {
	return SessionInfo{
		BaseDir:           tcr.sourceTree.GetBaseDir(),
		WorkDir:           toolchain.GetWorkDir(),
		LanguageName:      tcr.language.GetName(),
		ToolchainName:     tcr.toolchain.GetName(),
		VCSName:           tcr.vcs.Name(),
		VCSSessionSummary: tcr.vcs.SessionSummary(),
		GitAutoPush:       tcr.vcs.IsPushEnabled(),
		CommitOnFail:      tcr.commitOnFail,
		MessageSuffix:     tcr.messageSuffix,
	}
}

func (tcr *TCREngine) initTimer() {
	if settings.EnableMobTimer {
		tcr.mobTimer = timer.NewMobTurnCountdown(tcr.mode, tcr.mobTurnDuration)
	}
}

func (tcr *TCREngine) startTimer() {
	if settings.EnableMobTimer && tcr.mobTimer != nil {
		tcr.mobTimer.Start()
	}
}

func (tcr *TCREngine) stopTimer() {
	if settings.EnableMobTimer && tcr.mobTimer != nil {
		tcr.mobTimer.Stop()
		tcr.mobTimer = nil
	}
}

// ReportMobTimerStatus reports the status of the mob timer
func (tcr *TCREngine) ReportMobTimerStatus() {
	if settings.EnableMobTimer {
		timer.ReportCountDownStatus(tcr.mobTimer)
	}
}

// SetRunMode sets the run mode for TCR engine
func (tcr *TCREngine) SetRunMode(m runmode.RunMode) {
	tcr.mode = m
}

// Quit is the exit point for TCR application
func (*TCREngine) Quit() {
	report.PostInfo("That's All Folks!")
	// Give trace reporter some time to flush whatever has not been posted yet
	time.Sleep(traceReporterWaitingTime)
	rc := status.GetReturnCode()
	os.Exit(rc) //nolint:revive
}

func (*TCREngine) handleError(err error, fatal bool, s status.Status) {
	if err != nil {
		status.RecordState(s)
		if fatal {
			report.PostError(err)
			time.Sleep(traceReporterWaitingTime)
			os.Exit(status.GetReturnCode()) //nolint:revive
		}
		report.PostWarning(err)
	} else {
		status.RecordState(status.Ok)
	}
}

// VCSPull runs a VCS pull command on demand
func (tcr *TCREngine) VCSPull() {
	if tcr.vcs.Pull() != nil {
		report.PostError("VCS pull command failed!")
	}
}

// VCSPush runs a VCS push command on demand
func (tcr *TCREngine) VCSPush() {
	if tcr.vcs.Push() != nil {
		report.PostError("VCS push command failed!")
	}
}

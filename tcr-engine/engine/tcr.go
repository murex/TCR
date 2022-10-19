/*
Copyright (c) 2022 Murex

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
	"github.com/murex/tcr/tcr-engine/checker"
	"github.com/murex/tcr/tcr-engine/events"
	"github.com/murex/tcr/tcr-engine/filesystem"
	"github.com/murex/tcr/tcr-engine/language"
	"github.com/murex/tcr/tcr-engine/params"
	"github.com/murex/tcr/tcr-engine/report"
	"github.com/murex/tcr/tcr-engine/role"
	"github.com/murex/tcr/tcr-engine/runmode"
	"github.com/murex/tcr/tcr-engine/settings"
	"github.com/murex/tcr/tcr-engine/stats"
	"github.com/murex/tcr/tcr-engine/status"
	"github.com/murex/tcr/tcr-engine/timer"
	"github.com/murex/tcr/tcr-engine/toolchain"
	"github.com/murex/tcr/tcr-engine/ui"
	"github.com/murex/tcr/tcr-engine/vcs"
	"gopkg.in/tomb.v2"
	"os"
	"strings"
	"time"
)

type (
	// TcrInterface provides the API for interacting with TCR engine
	TcrInterface interface {
		Init(u ui.UserInterface, p params.Params)
		setVcs(gitInterface vcs.GitInterface)
		ToggleAutoPush()
		SetAutoPush(flag bool)
		SetCommitOnFail(flag bool)
		GetCurrentRole() role.Role
		RunAsDriver()
		RunAsNavigator()
		Stop()
		RunTCRCycle()
		build() (result toolchain.CommandResult)
		test() (result toolchain.TestCommandResult)
		commit(event events.TcrEvent)
		revert(events.TcrEvent)
		GetSessionInfo() SessionInfo
		ReportMobTimerStatus()
		SetRunMode(m runmode.RunMode)
		RunCheck(p params.Params)
		PrintLog(p params.Params)
		PrintStats(p params.Params)
		Quit()
	}

	// TcrEngine is the engine running all TCR operations
	TcrEngine struct {
		mode            runmode.RunMode
		ui              ui.UserInterface
		vcs             vcs.GitInterface
		language        language.LangInterface
		toolchain       toolchain.TchnInterface
		sourceTree      filesystem.SourceTree
		pollingPeriod   time.Duration
		mobTurnDuration time.Duration
		mobTimer        *timer.PeriodicReminder
		currentRole     role.Role
		commitOnFail    bool
		// shoot channel is used for handling interruptions coming from the UI
		shoot chan bool
	}
)

const traceReporterWaitingTime = 100 * time.Millisecond

const (
	commitMessageOk     = "✅ TCR - tests passing"
	commitMessageFail   = "❌ TCR - tests failing"
	commitMessageRevert = "⏪ TCR - revert changes"
)

var (
	// Tcr is TCR Engine singleton instance
	Tcr TcrInterface
)

// NewTcrEngine instantiates TCR engine instance
func NewTcrEngine() TcrInterface {
	Tcr = &TcrEngine{}
	return Tcr
}

// Init initializes the TCR engine with the provided parameters, and wires it to the user interface.
// This function should be called only once during the lifespan of the application
func (tcr *TcrEngine) Init(u ui.UserInterface, p params.Params) {
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

	tcr.initVcs()
	tcr.vcs.EnablePush(p.AutoPush)

	tcr.SetCommitOnFail(p.CommitFailures)

	tcr.setMobTimerDuration(p.MobTurnDuration)

	tcr.ui.ShowRunningMode(tcr.mode)
	tcr.ui.ShowSessionInfo()
	tcr.warnIfOnRootBranch(tcr.vcs.GetWorkingBranch(), tcr.mode.IsInteractive())
}

// SetCommitOnFail sets git commit-on-fail option to the provided value
func (tcr *TcrEngine) SetCommitOnFail(flag bool) {
	tcr.commitOnFail = flag
	if tcr.commitOnFail {
		report.PostInfo("Test-breaking changes will be committed")
	} else {
		report.PostInfo("Test-breaking changes will not be committed")
	}
}

func (tcr *TcrEngine) setMobTimerDuration(duration time.Duration) {
	if settings.EnableMobTimer && tcr.mode.NeedsCountdownTimer() {
		tcr.mobTurnDuration = duration
		report.PostInfo("Timer duration is ", tcr.mobTurnDuration)
	}
}

// RunCheck checks the provided parameters and prints out corresponding report
func (*TcrEngine) RunCheck(p params.Params) {
	checker.Run(p)
}

// PrintLog prints the TCR git commit history
func (tcr *TcrEngine) PrintLog(p params.Params) {
	tcrLogs := tcr.queryGitLogs(p)
	report.PostInfo("Printing TCR log for branch ", tcr.vcs.GetWorkingBranch())
	for _, log := range tcrLogs {
		report.PostTitle("commit:    ", log.Hash)
		report.PostInfo("timestamp: ", log.Timestamp)
		report.PostInfo("message:   ", log.Message)
		// Giving trace reporter some time to flush its contents
		time.Sleep(traceReporterWaitingTime)
	}
}

// PrintStats prints the TCR execution stats
func (tcr *TcrEngine) PrintStats(p params.Params) {
	tcrLogs := tcr.queryGitLogs(p)
	stats.Print(tcr.vcs.GetWorkingBranch(), tcrLogsToEvents(tcrLogs))
}

func tcrLogsToEvents(tcrLogs vcs.GitLogItems) (tcrEvents events.TcrEvents) {
	tcrEvents = *events.NewTcrEvents()
	for _, log := range tcrLogs {
		tcrEvents.Add(log.Timestamp, parseCommitMessage(log.Message))
	}
	return tcrEvents
}

func (tcr *TcrEngine) queryGitLogs(p params.Params) vcs.GitLogItems {
	tcr.initSourceTree(p)
	tcr.initVcs()

	logs, err := tcr.vcs.Log(isTcrCommitMessage)
	if err != nil {
		report.PostError(err)
	}
	if len(logs) == 0 {
		report.PostWarning("no TCR commit found in branch ", tcr.vcs.GetWorkingBranch(), "'s history")
	}
	return logs
}

func isTcrCommitMessage(msg string) bool {
	return strings.Index(msg, commitMessageOk) == 0 || strings.Index(msg, commitMessageFail) == 0
}

func parseCommitMessage(message string) (event events.TcrEvent) {
	var header string
	// First line is the main commit message
	// Second line is a blank line
	// The yaml-structured data starts on the third line
	const nbParts = 3
	parts := strings.SplitN(message, "\n", nbParts)
	if len(parts) == nbParts {
		header = parts[0]
		event = events.FromYaml(parts[nbParts-1])
	}
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

func (tcr *TcrEngine) initVcs() {
	if tcr.vcs == nil {
		var err error
		tcr.vcs, err = vcs.New(tcr.sourceTree.GetBaseDir())
		tcr.handleError(err, true, status.GitError)
	}
}

func (tcr *TcrEngine) initSourceTree(p params.Params) {
	var err error
	tcr.sourceTree, err = filesystem.New(p.BaseDir)
	tcr.handleError(err, true, status.ConfigError)
	report.PostInfo("Base directory is ", tcr.sourceTree.GetBaseDir())
}

func (tcr *TcrEngine) setVcs(gitInterface vcs.GitInterface) {
	tcr.vcs = gitInterface
}

func (tcr *TcrEngine) warnIfOnRootBranch(branch string, interactive bool) {
	if vcs.IsRootBranch(branch) {
		message := "Running " + settings.ApplicationName + " on branch \"" + branch + "\" is not recommended"
		if interactive {
			if !tcr.ui.Confirm(message, false) {
				tcr.Quit()
			}
		} else {
			report.PostWarning(message)
		}
	}
}

// ToggleAutoPush toggles git auto-push state
func (tcr *TcrEngine) ToggleAutoPush() {
	tcr.vcs.EnablePush(!tcr.vcs.IsPushEnabled())
}

// SetAutoPush sets git auto-push to the provided value
func (tcr *TcrEngine) SetAutoPush(ap bool) {
	tcr.vcs.EnablePush(ap)
}

// GetCurrentRole returns the role currently used for running TCR.
// Returns nil when TCR engine is in standby
func (tcr *TcrEngine) GetCurrentRole() role.Role {
	return tcr.currentRole
}

// RunAsDriver tells TCR engine to start running with driver role
func (tcr *TcrEngine) RunAsDriver() {
	tcr.initTimer()

	go tcr.fromBirthTillDeath(
		func() {
			tcr.currentRole = role.Driver{}
			tcr.ui.NotifyRoleStarting(tcr.currentRole)
			tcr.handleError(tcr.vcs.Pull(), false, status.GitError)
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
func (tcr *TcrEngine) RunAsNavigator() {
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
				tcr.handleError(tcr.vcs.Pull(), false, status.GitError)
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
func (tcr *TcrEngine) Stop() {
	tcr.shoot <- true
}

func (tcr *TcrEngine) fromBirthTillDeath(
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

func (tcr *TcrEngine) waitForChange(interrupt <-chan bool) bool {
	report.PostInfo("Going to sleep until something interesting happens")
	// We need to wait a bit to make sure the file watcher
	// does not get triggered again following a revert operation
	time.Sleep(2 * time.Second)
	return tcr.sourceTree.Watch(
		tcr.language.DirsToWatch(tcr.sourceTree.GetBaseDir()),
		tcr.language.IsLanguageFile,
		interrupt)
}

// RunTCRCycle is the core of TCR engine: e.g. it runs one test && commit || revert cycle
func (tcr *TcrEngine) RunTCRCycle() {
	status.RecordState(status.Ok)
	if tcr.build().Failed() {
		return
	}
	result := tcr.test()
	event := tcr.createTcrEvent(result)
	if result.Passed() {
		tcr.commit(event)
	} else {
		tcr.revert(event)
	}
}

func (tcr *TcrEngine) createTcrEvent(testResult toolchain.TestCommandResult) (event events.TcrEvent) {
	diffs, err := tcr.vcs.Diff()
	if err != nil {
		report.PostWarning(err)
	}
	commandStatus := events.StatusFail
	if testResult.Passed() {
		commandStatus = events.StatusPass
	}
	return events.NewTcrEvent(
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

func (tcr *TcrEngine) build() (result toolchain.CommandResult) {
	report.PostInfo("Launching Build")
	result = tcr.toolchain.RunBuild()
	if result.Failed() {
		status.RecordState(status.BuildFailed)
		report.PostWarning("There are build errors! I can't go any further")
	}
	return result
}

func (tcr *TcrEngine) test() (result toolchain.TestCommandResult) {
	report.PostInfo("Running Tests")
	result = tcr.toolchain.RunTests()
	if result.Failed() {
		status.RecordState(status.TestFailed)
		report.PostWarning("Some tests are failing! That's unfortunate")
		report.PostNotification("Some tests are failing! That's unfortunate")
	}
	return result
}

func (tcr *TcrEngine) commit(event events.TcrEvent) {
	report.PostInfo("Committing changes on branch ", tcr.vcs.GetWorkingBranch())
	var err error
	err = tcr.vcs.Add()
	tcr.handleError(err, false, status.GitError)
	if err != nil {
		return
	}
	err = tcr.vcs.Commit(false, commitMessageOk, event.ToYaml())
	tcr.handleError(err, false, status.GitError)
	if err != nil {
		return
	}
	tcr.handleError(tcr.vcs.Push(), false, status.GitError)
}

func (tcr *TcrEngine) revert(event events.TcrEvent) {
	if tcr.commitOnFail {
		err := tcr.commitTestBreakingChanges(event)
		tcr.handleError(err, false, status.GitError)
		if err != nil {
			return
		}
		tcr.handleError(tcr.vcs.Push(), false, status.GitError)
	}
	tcr.revertSrcFiles()
}

func (tcr *TcrEngine) commitTestBreakingChanges(event events.TcrEvent) (err error) {
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
	// Commit changes with failure message into git index
	err = tcr.vcs.Add()
	if err != nil {
		return err
	}
	err = tcr.vcs.Commit(false, commitMessageFail, event.ToYaml())
	if err != nil {
		return err
	}
	// Revert changes (both in git index and working tree)
	err = tcr.vcs.Revert()
	if err != nil {
		return err
	}
	// Amend commit message on revert operation in git index
	err = tcr.vcs.Commit(true, commitMessageRevert)
	if err != nil {
		return err
	}
	// Re-apply changes in the working tree and get rid of stash
	err = tcr.vcs.UnStash(false)
	return err
}

func (tcr *TcrEngine) revertSrcFiles() {
	diffs, err := tcr.vcs.Diff()
	tcr.handleError(err, false, status.GitError)
	if err != nil {
		return
	}
	var reverted int
	for _, diff := range diffs {
		if tcr.language.IsSrcFile(diff.Path) {
			err := tcr.revertFile(diff.Path)
			tcr.handleError(err, false, status.GitError)
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

func (tcr *TcrEngine) revertFile(file string) error {
	return tcr.vcs.Restore(file)
}

// GetSessionInfo provides the information related to the current TCR session.
// Used mainly by the user interface packages to retrieve and display this information
func (tcr *TcrEngine) GetSessionInfo() SessionInfo {
	return SessionInfo{
		BaseDir:       tcr.sourceTree.GetBaseDir(),
		WorkDir:       toolchain.GetWorkDir(),
		LanguageName:  tcr.language.GetName(),
		ToolchainName: tcr.toolchain.GetName(),
		AutoPush:      tcr.vcs.IsPushEnabled(),
		CommitOnFail:  tcr.commitOnFail,
		BranchName:    tcr.vcs.GetWorkingBranch(),
	}
}

func (tcr *TcrEngine) initTimer() {
	if settings.EnableMobTimer {
		tcr.mobTimer = timer.NewMobTurnCountdown(tcr.mode, tcr.mobTurnDuration)
	}
}

func (tcr *TcrEngine) startTimer() {
	if settings.EnableMobTimer && tcr.mobTimer != nil {
		tcr.mobTimer.Start()
	}
}

func (tcr *TcrEngine) stopTimer() {
	if settings.EnableMobTimer && tcr.mobTimer != nil {
		tcr.mobTimer.Stop()
		tcr.mobTimer = nil
	}
}

// ReportMobTimerStatus reports the status of the mob timer
func (tcr *TcrEngine) ReportMobTimerStatus() {
	if settings.EnableMobTimer {
		timer.ReportCountDownStatus(tcr.mobTimer)
	}
}

// SetRunMode sets the run mode for TCR engine
func (tcr *TcrEngine) SetRunMode(m runmode.RunMode) {
	tcr.mode = m
}

// Quit is the exit point for TCR application
func (*TcrEngine) Quit() {
	report.PostInfo("That's All Folks!")
	// Give trace reporter some time to flush whatever has not been posted yet
	time.Sleep(traceReporterWaitingTime)
	rc := status.GetReturnCode()
	os.Exit(rc) //nolint:revive
}

func (*TcrEngine) handleError(err error, fatal bool, s status.Status) {
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

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
		Init(u ui.UserInterface, params params.Params)
		setVcs(vcs vcs.GitInterface)
		ToggleAutoPush()
		SetAutoPush(flag bool)
		SetCommitOnFail(flag bool)
		GetCurrentRole() role.Role
		RunAsDriver()
		RunAsNavigator()
		Stop()
		RunTCRCycle()
		build() (result toolchain.CommandResult)
		test() (testStats toolchain.TestStats, result toolchain.CommandResult)
		commit(event events.TcrEvent)
		revert(events.TcrEvent)
		GetSessionInfo() SessionInfo
		ReportMobTimerStatus()
		SetRunMode(m runmode.RunMode)
		RunCheck(params params.Params)
		PrintLog(params params.Params)
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
func (tcr *TcrEngine) Init(u ui.UserInterface, params params.Params) {
	var err error
	status.RecordState(status.Ok)
	tcr.ui = u

	report.PostInfo("Starting ", settings.ApplicationName, " version ", settings.BuildVersion, "...")

	tcr.SetRunMode(params.Mode)
	if !tcr.mode.IsActive() {
		tcr.ui.ShowRunningMode(tcr.mode)
		return
	}

	tcr.pollingPeriod = params.PollingPeriod

	tcr.initSourceTree(params)

	tcr.language, err = language.GetLanguage(params.Language, tcr.sourceTree.GetBaseDir())
	tcr.handleError(err, true, status.ConfigError)

	tcr.toolchain, err = tcr.language.GetToolchain(params.Toolchain)
	tcr.handleError(err, true, status.ConfigError)

	err = toolchain.SetWorkDir(params.WorkDir)
	tcr.handleError(err, true, status.ConfigError)
	report.PostInfo("Work directory is ", toolchain.GetWorkDir())

	tcr.initVcs()
	tcr.vcs.EnablePush(params.AutoPush)

	tcr.SetCommitOnFail(params.CommitFailures)

	tcr.setMobTimerDuration(params.MobTurnDuration)

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
func (tcr *TcrEngine) RunCheck(params params.Params) {
	checker.Run(params)
}

// PrintLog prints the TCR git commit history
func (tcr *TcrEngine) PrintLog(params params.Params) {
	tcr.initSourceTree(params)
	tcr.initVcs()

	logs, _ := tcr.vcs.Log(func(msg string) bool {
		return strings.Index(msg, commitMessageOk) == 0 || strings.Index(msg, commitMessageFail) == 0
	})

	for _, log := range logs {
		report.PostInfo("commit: " + log.Hash)
		report.PostInfo("timestamp: " + log.Timestamp.String())
		report.PostInfo("message: " + log.Message)
	}
}

func parseCommitMessage(message string) (string, events.TcrEvent) {
	// First line is the main commit message
	// Second line is a blank line
	// The yaml-structured data starts on the third line
	parts := strings.SplitN(message, "\n", 3)
	if len(parts) == 3 {
		return parts[0], events.FromYaml(parts[2])
	}
	return "", events.TcrEvent{}
}

func (tcr *TcrEngine) initVcs() {
	if tcr.vcs == nil {
		var err error
		tcr.vcs, err = vcs.New(tcr.sourceTree.GetBaseDir())
		tcr.handleError(err, true, status.GitError)
	}
}

func (tcr *TcrEngine) initSourceTree(params params.Params) {
	var err error
	tcr.sourceTree, err = filesystem.New(params.BaseDir)
	tcr.handleError(err, true, status.ConfigError)
	report.PostInfo("Base directory is ", tcr.sourceTree.GetBaseDir())
}

func (tcr *TcrEngine) setVcs(vcs vcs.GitInterface) {
	tcr.vcs = vcs
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
	time.Sleep(1 * time.Second)
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
	stats, result := tcr.test()
	event := tcr.createTcrEvent(stats)
	if result.Passed() {
		tcr.commit(event)
	} else {
		tcr.revert(event)
	}
}

func (tcr *TcrEngine) createTcrEvent(testStats toolchain.TestStats) (event events.TcrEvent) {
	diffs, err := tcr.vcs.Diff()
	if err != nil {
		report.PostWarning(err)
	}
	return events.NewTcrEvent(
		events.NewChangedLines(
			diffs.ChangedLines(tcr.language.IsSrcFile),
			diffs.ChangedLines(tcr.language.IsTestFile),
		),
		events.NewTestStats(
			testStats.TotalRun,
			testStats.Passed,
			testStats.Failed,
			testStats.Skipped,
			testStats.WithErrors,
			testStats.Duration,
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
	return
}

func (tcr *TcrEngine) test() (testStats toolchain.TestStats, result toolchain.CommandResult) {
	report.PostInfo("Running Tests")
	result, testStats = tcr.toolchain.RunTests()
	if result.Failed() {
		status.RecordState(status.TestFailed)
		report.PostWarning("Some tests are failing! That's unfortunate")
	}
	return
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
		return
	}
	// Apply changes back in the working tree
	err = tcr.vcs.UnStash(true)
	if err != nil {
		return
	}
	// Commit changes with failure message into git index
	err = tcr.vcs.Add()
	if err != nil {
		return
	}
	err = tcr.vcs.Commit(false, commitMessageFail, event.ToYaml())
	if err != nil {
		return
	}
	// Revert changes (both in git index and working tree)
	err = tcr.vcs.Revert()
	if err != nil {
		return
	}
	// Amend commit message on revert operation in git index
	err = tcr.vcs.Commit(true, commitMessageRevert)
	if err != nil {
		return
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
func (tcr *TcrEngine) Quit() {
	report.PostInfo("That's All Folks!")
	time.Sleep(1 * time.Millisecond)
	os.Exit(status.GetReturnCode())
}

func (tcr *TcrEngine) handleError(err error, fatal bool, s status.Status) {
	if err != nil {
		status.RecordState(s)
		if fatal {
			report.PostError(err)
			time.Sleep(1 * time.Millisecond)
			os.Exit(status.GetReturnCode())
		} else {
			report.PostWarning(err)
		}
	} else {
		status.RecordState(status.Ok)
	}
}

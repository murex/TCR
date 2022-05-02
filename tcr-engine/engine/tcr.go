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
	"time"
)

type (
	// TcrInterface provides the API for interacting with TCR engine
	TcrInterface interface {
		Init(u ui.UserInterface, params params.Params)
		setVcs(vcs vcs.GitInterface)
		ToggleAutoPush()
		SetAutoPush(ap bool)
		GetCurrentRole() role.Role
		RunAsDriver()
		RunAsNavigator()
		Stop()
		RunTCRCycle()
		build() error
		test() error
		commit()
		revert()
		GetSessionInfo() SessionInfo
		ReportMobTimerStatus()
		SetRunMode(m runmode.RunMode)
		RunCheck(params params.Params)
		Quit()
	}

	// TcrEngine is the engine running all TCR operations
	TcrEngine struct {
		mode            runmode.RunMode
		uitf            ui.UserInterface
		vcs             vcs.GitInterface
		lang            language.LangInterface
		tchn            toolchain.TchnInterface
		sourceTree      filesystem.SourceTree
		pollingPeriod   time.Duration
		mobTurnDuration time.Duration
		mobTimer        *timer.PeriodicReminder
		currentRole     role.Role
		// shoot channel is used for handling interruptions coming from the UI
		shoot chan bool
	}
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
	tcr.uitf = u

	report.PostInfo("Starting ", settings.ApplicationName, " version ", settings.BuildVersion, "...")

	tcr.SetRunMode(params.Mode)
	if !tcr.mode.IsActive() {
		tcr.uitf.ShowRunningMode(tcr.mode)
		return
	}

	tcr.pollingPeriod = params.PollingPeriod

	tcr.sourceTree, err = filesystem.New(params.BaseDir)
	tcr.handleError(err, true, status.ConfigError)
	report.PostInfo("Base directory is ", tcr.sourceTree.GetBaseDir())

	tcr.lang, err = language.GetLanguage(params.Language, tcr.sourceTree.GetBaseDir())
	tcr.handleError(err, true, status.ConfigError)

	tcr.tchn, err = tcr.lang.GetToolchain(params.Toolchain)
	tcr.handleError(err, true, status.ConfigError)

	err = toolchain.SetWorkDir(params.WorkDir)
	tcr.handleError(err, true, status.ConfigError)
	report.PostInfo("Work directory is ", toolchain.GetWorkDir())

	git, err := vcs.New(tcr.sourceTree.GetBaseDir())
	tcr.handleError(err, true, status.GitError)
	tcr.setVcs(git)
	tcr.vcs.EnablePush(params.AutoPush)

	if settings.EnableMobTimer && tcr.mode.NeedsCountdownTimer() {
		tcr.mobTurnDuration = params.MobTurnDuration
		report.PostInfo("Timer duration is ", tcr.mobTurnDuration)
	}

	tcr.uitf.ShowRunningMode(tcr.mode)
	tcr.uitf.ShowSessionInfo()
	tcr.warnIfOnRootBranch(tcr.vcs.GetWorkingBranch(), tcr.mode.IsInteractive())
}

// RunCheck checks the provided parameters and prints out corresponding report
func (tcr *TcrEngine) RunCheck(params params.Params) {
	checker.Run(params)
}

func (tcr *TcrEngine) setVcs(vcs vcs.GitInterface) {
	tcr.vcs = vcs
}

func (tcr *TcrEngine) warnIfOnRootBranch(branch string, interactive bool) {
	if vcs.IsRootBranch(branch) {
		message := "Running " + settings.ApplicationName + " on branch \"" + branch + "\" is not recommended"
		if interactive {
			if !tcr.uitf.Confirm(message, false) {
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
	if settings.EnableMobTimer {
		tcr.mobTimer = timer.NewMobTurnCountdown(tcr.mode, tcr.mobTurnDuration)
	}

	go tcr.fromBirthTillDeath(
		func() {
			tcr.currentRole = role.Driver{}
			tcr.uitf.NotifyRoleStarting(tcr.currentRole)
			tcr.handleError(tcr.vcs.Pull(), false, status.GitError)
			if settings.EnableMobTimer {
				tcr.mobTimer.Start()
			}
		},
		func(interrupt <-chan bool) bool {
			inactivityTeaser := timer.GetInactivityTeaserInstance()
			inactivityTeaser.Start()
			if tcr.waitForChange(interrupt) {
				// Some file changes were detected
				inactivityTeaser.Reset()
				tcr.RunTCRCycle()
				inactivityTeaser.Start()
				return true
			}
			// If we arrive here this means that the end of waitForChange
			// was triggered by the user
			inactivityTeaser.Reset()
			return false
		},
		func() {
			if settings.EnableMobTimer {
				tcr.mobTimer.Stop()
				tcr.mobTimer = nil
			}
			tcr.uitf.NotifyRoleEnding(tcr.currentRole)
			tcr.currentRole = nil
		},
	)
}

// RunAsNavigator tells TCR engine to start running with navigator role
func (tcr *TcrEngine) RunAsNavigator() {
	go tcr.fromBirthTillDeath(
		func() {
			tcr.currentRole = role.Navigator{}
			tcr.uitf.NotifyRoleStarting(tcr.currentRole)
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
			tcr.uitf.NotifyRoleEnding(tcr.currentRole)
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
		tcr.lang.DirsToWatch(tcr.sourceTree.GetBaseDir()),
		tcr.lang.IsLanguageFile,
		interrupt)
}

// RunTCRCycle is the core of TCR engine: e.g. it runs one test && commit || revert cycle
func (tcr *TcrEngine) RunTCRCycle() {
	status.RecordState(status.Ok)
	if tcr.build() != nil {
		tcr.logEvent(events.StatusFailed, events.StatusUnknown)
		return
	}
	if tcr.test() == nil {
		tcr.logEvent(events.StatusPassed, events.StatusPassed)
		tcr.commit()
	} else {
		tcr.logEvent(events.StatusPassed, events.StatusFailed)
		tcr.revert()
	}
}

func (tcr *TcrEngine) logEvent(buildStatus, testsStatus events.TcrEventStatus) {
	changedFiles, err := tcr.vcs.Diff()
	if err != nil {
		report.PostWarning(err)
	}

	events.EventRepository.Add(
		events.NewTcrEvent(
			computeSrcLinesChanged(tcr.lang, changedFiles),
			computeTestLinesChanged(tcr.lang, changedFiles),
			buildStatus,
			testsStatus,
			0,
			0,
			0,
			0,
			0),
	)
}

func (tcr *TcrEngine) build() error {
	report.PostInfo("Launching Build")
	err := tcr.tchn.RunBuild()
	if err != nil {
		status.RecordState(status.BuildFailed)
		report.PostWarning("There are build errors! I can't go any further")
	}
	return err
}

func (tcr *TcrEngine) test() error {
	report.PostInfo("Running Tests")
	_, err := tcr.tchn.RunTests()
	if err != nil {
		status.RecordState(status.TestFailed)
		report.PostWarning("Some tests are failing! That's unfortunate")
	}
	return err
}

func (tcr *TcrEngine) commit() {
	report.PostInfo("Committing changes on branch ", tcr.vcs.GetWorkingBranch())
	err := tcr.vcs.Commit()
	tcr.handleError(err, false, status.GitError)
	if err == nil {
		tcr.handleError(tcr.vcs.Push(), false, status.GitError)
	}
}

func (tcr *TcrEngine) revert() {
	changedFiles, err := tcr.vcs.ListChanges()
	tcr.handleError(err, false, status.GitError)
	if err != nil {
		return
	}

	var reverted int
	for _, file := range changedFiles {
		if tcr.lang.IsSrcFile(file) {
			err := tcr.revertFile(file)
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
		LanguageName:  tcr.lang.GetName(),
		ToolchainName: tcr.tchn.GetName(),
		AutoPush:      tcr.vcs.IsPushEnabled(),
		BranchName:    tcr.vcs.GetWorkingBranch(),
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

/*
Copyright (c) 2024 Murex

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
	"errors"
	"github.com/murex/tcr/checker"
	"github.com/murex/tcr/events"
	"github.com/murex/tcr/filesystem"
	"github.com/murex/tcr/language"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/report/role_event"
	"github.com/murex/tcr/retro"
	"github.com/murex/tcr/role"
	"github.com/murex/tcr/runmode"
	"github.com/murex/tcr/settings"
	"github.com/murex/tcr/stats"
	"github.com/murex/tcr/status"
	"github.com/murex/tcr/timer"
	"github.com/murex/tcr/toolchain"
	"github.com/murex/tcr/toolchain/command"
	"github.com/murex/tcr/ui"
	"github.com/murex/tcr/variant"
	"github.com/murex/tcr/vcs"
	"github.com/murex/tcr/vcs/factory"
	"gopkg.in/tomb.v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type (
	// TCRInterface provides the API for interacting with TCR engine
	TCRInterface interface {
		AttachUI(u ui.UserInterface, primary bool)
		Init(p params.Params)
		setVCS(vcsInterface vcs.Interface)
		ToggleAutoPush()
		SetAutoPush(flag bool)
		SetVariant(name string)
		GetCurrentRole() role.Role
		RunAsDriver()
		RunAsNavigator()
		Stop()
		RunTCRCycle()
		AbortCommand()
		GetSessionInfo() SessionInfo
		GetMobTimerStatus() timer.CurrentState
		SetRunMode(m runmode.RunMode)
		RunCheck(p params.Params)
		PrintLog(p params.Params)
		PrintStats(p params.Params)
		VCSPull()
		VCSPush()
		Quit()
		GenerateRetro(p params.Params)
	}

	// TCREngine is the engine running all TCR operations
	TCREngine struct {
		mode            runmode.RunMode
		ui              ui.Multicaster
		vcs             vcs.Interface
		language        language.LangInterface
		toolchain       toolchain.TchnInterface
		sourceTree      filesystem.SourceTree
		pollingPeriod   time.Duration
		mobTurnDuration time.Duration
		mobTimer        *timer.PeriodicReminder
		currentRole     role.Role
		// roleMutex is used to prevent the engine from starting 2 different
		// roles simultaneously: we wait for it to leave the previous role
		// before starting a new one
		roleMutex     sync.Mutex
		variant       *variant.Variant
		messageSuffix string
		// shoot channel is used for handling interruptions coming from the UI
		shoot chan bool
		// traceReporterWaitingTime is used to prevent trace reporter overflow when
		// due to slowness of terminal output when there is a large quantity
		// of information to report (such as when printing VCS log outcome)
		traceReporterWaitingTime time.Duration
		// fsWatchRearmDelay is the waiting time until TCR starts watching the filesystem again
		// after a filesystem event was detected. The default value should not be changed except
		// when running tests
		fsWatchRearmDelay time.Duration
	}
)

const retroFileName = "tcr-retro.md"

const traceReporterWaitingTime = 100 * time.Millisecond

const fsWatchRearmDelay = 100 * time.Millisecond

const (
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
	engine = &TCREngine{
		ui:                       *ui.NewMulticaster(),
		fsWatchRearmDelay:        fsWatchRearmDelay,
		traceReporterWaitingTime: traceReporterWaitingTime,
	}
	TCR = engine
	return engine
}

// AttachUI plugs a user interface to TCR
func (tcr *TCREngine) AttachUI(u ui.UserInterface, primary bool) {
	tcr.ui.Register(u, primary)
}

// Init initializes the TCR engine with the provided parameters, and wires it to the user interface.
// This function should be called only once during the lifespan of the application
// nolint:revive
func (tcr *TCREngine) Init(p params.Params) {
	var err error
	status.RecordState(status.Ok)

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
	tcr.reportFileStats()

	tcr.toolchain, err = tcr.language.GetToolchain(p.Toolchain)
	tcr.handleError(err, true, status.ConfigError)

	err = toolchain.SetWorkDir(p.WorkDir)
	tcr.handleError(err, true, status.ConfigError)
	report.PostInfo("Work directory is ", toolchain.GetWorkDir())

	tcr.initVCS(p.VCS, p.GitRemote, p.Trace)
	tcr.setMessageSuffix(p.MessageSuffix)
	tcr.vcs.EnableAutoPush(p.AutoPush)

	tcr.SetVariant(p.Variant)
	tcr.setMobTimerDuration(p.MobTurnDuration)

	tcr.ui.ShowRunningMode(tcr.mode)
	tcr.ui.ShowSessionInfo()
	tcr.warnIfOnRootBranch(tcr.mode.IsInteractive())
}

// SetVariant sets the TCR variant that will be used by TCR engine
func (tcr *TCREngine) SetVariant(name string) {
	var err error
	tcr.variant, err = variant.Select(name)
	if err != nil {
		var unsupportedVariantError *variant.UnsupportedVariantError
		if errors.As(err, &unsupportedVariantError) {
			tcr.handleError(err, true, status.ConfigError)
		}
	}
}

func (tcr *TCREngine) setMobTimerDuration(duration time.Duration) {
	if settings.EnableMobTimer {
		if tcr.mode.IsMultiRole() {
			tcr.mobTurnDuration = duration
			report.PostInfo("Timer duration is ", tcr.mobTurnDuration)
		} else {
			report.PostInfo("Timer is not used in " + tcr.mode.Name() + " mode")
		}
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
		time.Sleep(tcr.traceReporterWaitingTime)
	}
}

// PrintStats prints the TCR execution stats
func (tcr *TCREngine) PrintStats(p params.Params) {
	tcrLogs := tcr.queryVCSLogs(p)
	stats.Print(tcr.vcs.SessionSummary(), tcrLogsToEvents(tcrLogs))
}

// GenerateRetro generates a retrospective markdown file template using stats
func (tcr *TCREngine) GenerateRetro(p params.Params) {
	tcrEvents := tcrLogsToEvents(tcr.queryVCSLogs(p))
	markdown := retro.GenerateMarkdown(filepath.Base(tcr.vcs.GetRootDir()), &tcrEvents)
	retroPath := filepath.Join(tcr.sourceTree.GetBaseDir(), retroFileName)
	filesystem.WriteFile(retroPath, []byte(markdown))
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
	tcr.initVCS(p.VCS, "", p.Trace)

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
	return parseCommitStatus(msg) != events.StatusUnknown
}

func parseCommitMessage(message string) events.TCREvent {
	header, event := parseCommitHeaderAndEvents(message)
	event.Status = parseCommitStatus(header)
	return event
}

func parseCommitHeaderAndEvents(message string) (string, events.TCREvent) {
	// First line is the main commit message
	// Second line is a blank line
	// The YAML-structured data starts on the third line until we reach a blank line
	// The user-specified message prefix, if any, is after the blank line

	var header string
	var statsYAML strings.Builder
	var section = 1
	for line := range strings.SplitSeq(message, "\n") {
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
	event := events.FromYAML(statsYAML.String())
	return header, event
}

func parseCommitStatus(header string) events.CommandStatus {
	if strings.Contains(header, messagePassed.Tag) {
		return events.StatusPass
	}
	if strings.Contains(header, messageFailed.Tag) {
		return events.StatusFail
	}
	return events.StatusUnknown
}

func (tcr *TCREngine) setMessageSuffix(suffix string) {
	tcr.messageSuffix = suffix
}

func (tcr *TCREngine) wrapCommitMessages(header CommitMessage, event *events.TCREvent) []string {
	messages := []string{header.toString(tcr.vcs.SupportsEmojis())}
	if event != nil {
		messages = append(messages, event.ToYAML())
	}
	if tcr.messageSuffix != "" {
		messages = append(messages, "\n"+tcr.messageSuffix)
	}
	return messages
}

func (tcr *TCREngine) initVCS(vcsName string, remoteName string, trace string) {
	if tcr.vcs != nil {
		return // VCS should be initialized only once
	}
	// Set VCS trace flag
	vcs.SetTrace(trace == "vcs")
	var err error
	tcr.vcs, err = factory.InitVCS(vcsName, tcr.sourceTree.GetBaseDir(), remoteName)
	var unsupportedVCSError *factory.UnsupportedVCSError
	switch {
	case errors.As(err, &unsupportedVCSError):
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
	tcr.vcs.EnableAutoPush(!tcr.vcs.IsAutoPushEnabled())
}

// SetAutoPush sets VCS auto-push to the provided value
func (tcr *TCREngine) SetAutoPush(ap bool) {
	tcr.vcs.EnableAutoPush(ap)
}

// resetCurrentRole sets the current role to nil (TCR engine in standby).
// This is a mandatory step prior to starting a new role.
func (tcr *TCREngine) resetCurrentRole() {
	if tcr.currentRole != nil {
		report.PostRoleEvent(role_event.TriggerEnd, tcr.currentRole)
		tcr.currentRole = nil
	}
	tcr.roleMutex.Unlock()
}

// setCurrentRole sets the current role.
// Setting it to nil is the same as calling resetCurrentRole().
// This operation is mutex-protected, e.g. it will wait until
// the currentRole is reset before setting the new role.
func (tcr *TCREngine) setCurrentRole(r role.Role) {
	if r == nil {
		tcr.resetCurrentRole()
		return
	}
	tcr.roleMutex.Lock()
	if r != tcr.currentRole {
		tcr.currentRole = r
		report.PostRoleEvent(role_event.TriggerStart, tcr.currentRole)
	}
}

// GetCurrentRole returns the role currently used for running TCR.
// Returns nil when TCR engine is in standby
func (tcr *TCREngine) GetCurrentRole() role.Role {
	return tcr.currentRole
}

// RunAsDriver tells TCR engine to start running with driver role
func (tcr *TCREngine) RunAsDriver() {
	// Force previous role to quit if needed
	if tcr.GetCurrentRole() != nil {
		tcr.Stop()
	}
	// Prepare the timer
	tcr.initTimer()

	go tcr.fromBirthTillDeath(
		func() {
			// the goroutine waits until currenRole is reset
			// prior to starting driver role
			tcr.setCurrentRole(role.Driver{})
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
			tcr.resetCurrentRole()
		},
	)
}

// RunAsNavigator tells TCR engine to start running with navigator role
func (tcr *TCREngine) RunAsNavigator() {
	// Force previous role to quit if needed
	if tcr.GetCurrentRole() != nil {
		tcr.Stop()
	}
	// preemptionDuration is the max time to wait before navigator role can be interrupted
	const preemptionDuration = 100 * time.Millisecond
	countdown := 0 * time.Millisecond

	go tcr.fromBirthTillDeath(
		func() {
			// the goroutine waits until currenRole is reset
			// prior to starting navigator role
			tcr.setCurrentRole(role.Navigator{})
		},
		func(interrupt <-chan bool) bool {
			select {
			case <-interrupt:
				return false
			default:
				if countdown <= 0 {
					tcr.handleError(tcr.vcs.Pull(), false, status.VCSError)
					countdown = tcr.pollingPeriod
				} else {
					time.Sleep(preemptionDuration)
					countdown -= preemptionDuration
				}
				return true
			}
		},
		func() {
			tcr.resetCurrentRole()
		},
	)
}

// Stop is the entry point for telling TCR engine to stop its current operations
func (tcr *TCREngine) Stop() {
	if tcr.shoot != nil {
		tcr.shoot <- true
	}
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
	languageDirs := tcr.language.DirsToWatch(tcr.sourceTree.GetBaseDir())
	existingDirs, err := language.ExistingDirsIn(languageDirs)
	if err != nil {
		tcr.handleError(err, true, status.OtherError)
	}
	report.PostInfo("Going to sleep until something interesting happens")
	// We need to wait a bit to make sure the file watcher
	// does not get triggered again following a revert operation
	time.Sleep(tcr.fsWatchRearmDelay)

	return tcr.sourceTree.Watch(
		existingDirs,
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

// AbortCommand triggers interruption of an ongoing TCR cycle operation
func (tcr *TCREngine) AbortCommand() {
	_ = tcr.toolchain.AbortExecution()
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

func (tcr *TCREngine) build() (result command.Result) {
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
	err = tcr.vcs.Commit(tcr.wrapCommitMessages(messagePassed, &event)...)
	tcr.handleError(err, false, status.VCSError)
	if err != nil {
		return
	}
	tcr.handleError(tcr.vcsPushAuto(), false, status.VCSError)
}

func (tcr *TCREngine) revert(e events.TCREvent) {
	var err error

	switch *tcr.variant {
	case variant.Introspective:
		err = tcr.introspectiveRevert(e)
	default:
		err = tcr.simpleRevert()
	}

	tcr.handleError(err, false, status.VCSError)
}

func (tcr *TCREngine) simpleRevert() error {
	diffs, err := tcr.vcs.Diff()
	if err != nil {
		return err
	}
	var reverted int
	for _, diff := range diffs {
		if tcr.shouldRevertFile(diff.Path) {
			err := tcr.revertFile(diff.Path)
			if err != nil {
				return err
			}
			reverted++
		}
	}
	if reverted > 0 {
		report.PostWarning(reverted, " file(s) reverted")
	} else {
		report.PostInfo(tcr.noFilesRevertedMessage())
	}
	return nil
}

func (tcr *TCREngine) introspectiveRevert(event events.TCREvent) (err error) {
	err = tcr.vcs.Add()
	if err != nil {
		return err
	}
	err = tcr.vcs.Commit(tcr.wrapCommitMessages(messageFailed, &event)...)
	if err != nil {
		return err
	}
	err = tcr.vcs.RollbackLastCommit()
	if err != nil {
		return err
	}
	err = tcr.vcs.Commit(tcr.wrapCommitMessages(messageReverted, nil)...)
	return err
}

func (tcr *TCREngine) noFilesRevertedMessage() string {
	if *tcr.variant == variant.Relaxed {
		return "No file reverted (only test files were updated since last commit)"
	}
	return "No file reverted"
}

func (tcr *TCREngine) shouldRevertFile(path string) bool {
	return *tcr.variant == variant.BTCR || tcr.language.IsSrcFile(path)
}

func (tcr *TCREngine) revertFile(file string) error {
	return tcr.vcs.RevertLocal(file)
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
		GitAutoPush:       tcr.vcs.IsAutoPushEnabled(),
		Variant:           tcr.variant.Name(),
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

// GetMobTimerStatus returns the status of the mob timer
func (tcr *TCREngine) GetMobTimerStatus() timer.CurrentState {
	return timer.GetCurrentState(tcr.mobTimer)
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

// vcsPushAuto runs a VCS push command if the auto-push option is enabled
func (tcr *TCREngine) vcsPushAuto() error {
	if tcr.vcs.IsAutoPushEnabled() {
		return tcr.vcs.Push()
	}
	return nil
}

// reportFileStats traces summary information about the source and test files and directories
func (tcr *TCREngine) reportFileStats() {
	srcFileCount := countFiles("source", tcr.language.AllSrcFiles)
	testFileCount := countFiles("test", tcr.language.AllTestFiles)
	if srcFileCount+testFileCount == 0 {
		report.PostWarning("No matching ", tcr.language.GetName(), " file found")
	} else {
		report.PostInfo("Found ", srcFileCount, " source and ",
			testFileCount, " test file(s) for ", tcr.language.GetName(), " language")
	}
}

func countFiles(desc string, matcher func() ([]string, error)) int {
	matches, err := matcher()
	switch err := err.(type) {
	case nil:
		// do nothing
	case *language.UnreachableDirectoryError:
		// unreachable directories: we display a warning for each, then continue
		for _, dir := range err.DirList() {
			report.PostWarning("Unreachable ", desc, " directory: ", dir)
		}
	default:
		report.PostError(err)
	}
	return len(matches)
}

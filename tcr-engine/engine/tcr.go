package engine

import (
	"github.com/mengdaming/tcr-engine/filesystem"
	"github.com/mengdaming/tcr-engine/language"
	"github.com/mengdaming/tcr-engine/report"
	"github.com/mengdaming/tcr-engine/role"
	"github.com/mengdaming/tcr-engine/runmode"
	"github.com/mengdaming/tcr-engine/timer"
	"github.com/mengdaming/tcr-engine/toolchain"
	"github.com/mengdaming/tcr-engine/ui"
	"github.com/mengdaming/tcr-engine/vcs"
	"gopkg.in/tomb.v2"
	"os"
	"path/filepath"
	"time"
)

var (
	mode          runmode.RunMode
	uitf          ui.UserInterface
	git           vcs.GitInterface
	lang          language.Language
	tchn          toolchain.Toolchain
	sourceTree    filesystem.SourceTree
	pollingPeriod time.Duration
)

// Init initializes the TCR engine with the provided parameters, and wires it to the user interface.
// This function should be called only once during the lifespan of the application
func Init(u ui.UserInterface, params Params) {
	var err error

	uitf = u

	report.PostInfo("Starting TCR version ", Version, "...")

	mode = params.Mode
	pollingPeriod = params.PollingPeriod
	sourceTree, err = filesystem.New(params.BaseDir)
	handleError(err)
	report.PostInfo("Working directory is ", sourceTree.GetBaseDir())
	lang, err = language.DetectLanguage(sourceTree.GetBaseDir())
	handleError(err)
	tchn, err = toolchain.New(params.Toolchain, lang)
	handleError(err)
	git, err = vcs.New(sourceTree.GetBaseDir())
	handleError(err)
	git.EnablePush(params.AutoPush)

	uitf.ShowRunningMode(mode)
	uitf.ShowSessionInfo()
	warnIfOnRootBranch(git.WorkingBranch())
}

func warnIfOnRootBranch(branch string) {
	for _, b := range []string{"main", "master"} {
		if b == branch {
			if !uitf.Confirm("Running TCR on branch \""+branch+"\" is not recommended", false) {
				Quit()
			}
			break
		}
	}
}

// ToggleAutoPush toggles git auto-push state
func ToggleAutoPush() {
	git.EnablePush(!git.IsPushEnabled())
}

// SetAutoPush sets git auto-push to the provided value
func SetAutoPush(ap bool) {
	git.EnablePush(ap)
}

// RunAsDriver tells TCR engine to start running with driver role
func RunAsDriver() {
	go fromBirthTillDeath(
		func() {
			uitf.NotifyRoleStarting(role.Driver{})
			_ = git.Pull()
		},
		func(interrupt <-chan bool) bool {
			// TODO pass default values as parameters
			r := timer.NewInactivityTeaser(DefaultInactivityTimeout, DefaultInactivityPeriod)
			r.Start()
			if waitForChange(interrupt) {
				// Some file changes were detected
				r.Stop()
				runTCR()
				r.Start()
				return true
			}
			// If we arrive here this means that the end of waitForChange
			// was triggered by the user
			r.Stop()
			return false

		},
		func() {
			uitf.NotifyRoleEnding(role.Driver{})
		},
	)
}

// RunAsNavigator tells TCR engine to start running with navigator role
func RunAsNavigator() {
	go fromBirthTillDeath(
		func() {
			uitf.NotifyRoleStarting(role.Navigator{})
		},
		func(interrupt <-chan bool) bool {
			select {
			case <-interrupt:
				return false
			default:
				_ = git.Pull()
				time.Sleep(pollingPeriod)
				return true
			}
		},
		func() {
			uitf.NotifyRoleEnding(role.Navigator{})
		},
	)
}

// shoot channel is used to handle interruptions coming from the UI
var shoot chan bool

// Stop is the entry point for telling TCR engine to stop its current operations
func Stop() {
	shoot <- true
}

func fromBirthTillDeath(
	birth func(),
	dailyLife func(interrupt <-chan bool) bool,
	death func(),
) {
	var tmb tomb.Tomb
	shoot = make(chan bool)

	// The goroutine doing the work
	tmb.Go(func() error {
		birth()
		for oneMoreDay := true; oneMoreDay; {
			oneMoreDay = dailyLife(shoot)
		}
		death()
		return nil
	})

	err := tmb.Wait()
	if err != nil {
		report.PostError("tmb.Wait(): ", err)
	}
}

func waitForChange(interrupt <-chan bool) bool {
	report.PostInfo("Going to sleep until something interesting happens")
	return sourceTree.Watch(
		language.DirsToWatch(sourceTree.GetBaseDir(), lang),
		lang.IsSrcFile,
		interrupt)
}

func runTCR() {
	if build() != nil {
		return
	}
	if test() == nil {
		commit()
	} else {
		revert()
	}
}

func build() error {
	report.PostInfo("Launching Build")
	err := tchn.RunBuild()
	if err != nil {
		report.PostWarning("There are build errors! I can't go any further")
	}
	return err
}

func test() error {
	report.PostInfo("Running Tests")
	err := tchn.RunTests()
	if err != nil {
		report.PostWarning("Some tests are failing! That's unfortunate")
	}
	return err
}

func commit() {
	report.PostInfo("Committing changes on branch ", git.WorkingBranch())
	_ = git.Commit()
	_ = git.Push()
}

func revert() {
	report.PostWarning("Reverting changes")
	for _, dir := range lang.SrcDirs() {
		_ = git.Restore(filepath.Join(sourceTree.GetBaseDir(), dir))
	}
}

// GetSessionInfo provides the information (as strings) related to the current TCR session.
// Used mainly by the user interface packages to retrieve and display this information
func GetSessionInfo() (d string, l string, t string, ap bool, b string) {
	d = sourceTree.GetBaseDir()
	l = lang.Name()
	t = tchn.Name()
	ap = git.IsPushEnabled()
	b = git.WorkingBranch()

	return d, l, t, ap, b
}

// Quit is the exit point for TCR application
func Quit() {
	report.PostInfo("That's All Folks!")
	time.Sleep(1 * time.Millisecond)
	os.Exit(0)
}

func handleError(err error) {
	if err != nil {
		report.PostError(err)
		time.Sleep(1 * time.Millisecond)
		os.Exit(1)
	}
}

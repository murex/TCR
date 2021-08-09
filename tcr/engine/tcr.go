package engine

import (
	"github.com/mengdaming/tcr/tcr"
	"github.com/mengdaming/tcr/tcr/filesystem"
	"github.com/mengdaming/tcr/tcr/language"
	"github.com/mengdaming/tcr/tcr/report"
	"github.com/mengdaming/tcr/tcr/role"
	"github.com/mengdaming/tcr/tcr/runmode"
	"github.com/mengdaming/tcr/tcr/toolchain"
	"github.com/mengdaming/tcr/tcr/ui"
	"github.com/mengdaming/tcr/tcr/vcs"
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

func Init(u ui.UserInterface, params tcr.Params) {
	var err error

	uitf = u

	report.PostInfo("Starting TCR version ", tcr.Version, "...")

	mode = params.Mode
	pollingPeriod = params.PollingPeriod
	sourceTree, err = filesystem.New(params.BaseDir); handleError(err)
	lang, err = language.DetectLanguage(sourceTree.GetBaseDir()); handleError(err)
	tchn, err = toolchain.New(params.Toolchain, lang); handleError(err)
	git, err = vcs.New(sourceTree.GetBaseDir()); handleError(err)
	git.EnablePush(params.AutoPush)

	uitf.ShowRunningMode(mode)
	uitf.ShowSessionInfo()
	warnIfOnRootBranch(git.WorkingBranch())

	uitf.RunInMode(mode)
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

func ToggleAutoPush() {
	git.EnablePush(!git.IsPushEnabled())
}

func RunAsDriver() {
	go fromBirthTillDeath(
		func() {
			uitf.NotifyRoleStarting(role.Driver{})
			git.Pull()
		},
		func(interrupt <-chan bool) bool {
			if waitForChange(interrupt) {
				// Some file changes were detected
				runTCR()
				return true
			} else {
				// If we enter here this means that the end of waitForChange
				// was triggered by the user
				return false
			}
		},
		func() {
			uitf.NotifyRoleEnding(role.Driver{})
		},
	)
}

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
				git.Pull()
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
	git.Commit()
	git.Push()
}

func revert() {
	report.PostWarning("Reverting changes")
	for _, dir := range lang.SrcDirs() {
		git.Restore(filepath.Join(sourceTree.GetBaseDir(), dir))
	}
}

func GetSessionInfo() (d string, l string, t string, ap bool, b string) {
	d = sourceTree.GetBaseDir()
	l = lang.Name()
	t = tchn.Name()
	ap = git.IsPushEnabled()
	b = git.WorkingBranch()

	return d, l, t, ap, b
}

func Quit() {
	report.PostInfo("That's All Folks!")
	os.Exit(0)
}

func handleError(err error) {
	if err != nil {
		report.PostError(err)
		time.Sleep(1 * time.Millisecond)
		os.Exit(1)
	}
}

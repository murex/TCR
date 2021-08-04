package engine

import (
	"github.com/mengdaming/tcr/tcr"
	"github.com/mengdaming/tcr/tcr/filesystem"
	"github.com/mengdaming/tcr/tcr/language"
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

func Start(u ui.UserInterface, params tcr.Params) {
	uitf = u

	mode = params.Mode
	pollingPeriod = params.PollingPeriod
	sourceTree = filesystem.NewSourceTreeImpl(params.BaseDir)
	lang = language.DetectLanguage(sourceTree.GetBaseDir())
	tchn = toolchain.NewToolchain(params.Toolchain, lang)
	git = vcs.NewGitImpl(sourceTree.GetBaseDir())
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

func RunAsDriver(stopRequest <-chan bool) {
	fromBirthTillDeath(
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
		stopRequest)
}

func RunAsNavigator(stopRequest <-chan bool) {
	fromBirthTillDeath(
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
		stopRequest)
}

func fromBirthTillDeath(
	birth func(),
	dailyLife func(interrupt <-chan bool) bool,
	death func(),
	shoot <-chan bool) {

	var tmb tomb.Tomb

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
		uitf.Error("tmb.Wait(): ", err)
	}
}

func waitForChange(interrupt <-chan bool) bool {
	uitf.Info("Going to sleep until something interesting happens")
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
	uitf.Info("Launching Build")
	err := tchn.RunBuild()
	if err != nil {
		uitf.Warning("There are build errors! I can't go any further")
	}
	return err
}

func test() error {
	uitf.Info("Running Tests")
	err := tchn.RunTests()
	if err != nil {
		uitf.Warning("Some tests are failing! That's unfortunate")
	}
	return err
}

func commit() {
	uitf.Info("Committing changes on branch ", git.WorkingBranch())
	git.Commit()
	git.Push()
}

func revert() {
	uitf.Warning("Reverting changes")
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
	uitf.Info("That's All Folks!")
	os.Exit(0)
}

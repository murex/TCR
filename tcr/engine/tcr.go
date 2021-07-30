package engine

import (
	"github.com/mengdaming/tcr/tcr"
	"github.com/mengdaming/tcr/tcr/filesystem"
	"github.com/mengdaming/tcr/tcr/language"
	"github.com/mengdaming/tcr/tcr/role"
	"github.com/mengdaming/tcr/tcr/toolchain"
	"github.com/mengdaming/tcr/tcr/vcs"

	"gopkg.in/tomb.v2"

	"os"
	"path/filepath"
	"time"
)

var (
	mode          tcr.WorkMode
	ui            tcr.UserInterface
	git           vcs.GitInterface
	lang          language.Language
	tchn          toolchain.Toolchain
	sourceTree    filesystem.SourceTree
	pollingPeriod time.Duration
)

func Start(u tcr.UserInterface, params tcr.Params) {
	ui = u

	mode = params.Mode
	pollingPeriod = params.PollingPeriod
	sourceTree = filesystem.NewSourceTreeImpl(params.BaseDir)
	lang = language.DetectLanguage(sourceTree.GetBaseDir())
	tchn = toolchain.NewToolchain(params.Toolchain, lang)
	git = vcs.NewGitImpl(sourceTree.GetBaseDir())
	git.EnablePush(params.AutoPush)

	ui.ShowRunningMode(mode)
	ui.ShowSessionInfo()
	warnIfOnRootBranch(git.WorkingBranch())

	switch mode {
	case tcr.Solo:
		// When running TCR in solo mode, there's no
		// selection menu: we directly enter driver mode
		// TODO Put back -- should rely on UI which will handle interruption
		//stopEngine := make(chan bool)
		//RunAsDriver(stopEngine)
	case tcr.Mob:
		// When running TCR in mob mode, every participant
		// is given the possibility to switch between
		// driver and navigator modes
		ui.WaitForAction()
	}
}

func warnIfOnRootBranch(branch string) {
	for _, b := range []string{"main", "master"} {
		if b == branch {
			if !ui.Confirm("Running TCR on branch \""+branch+"\" is not recommended", false) {
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
			ui.NotifyRoleStarting(role.Driver{})
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
			ui.NotifyRoleEnding(role.Driver{})
		},
		stopRequest)
}

func RunAsNavigator(stopRequest <-chan bool) {
	fromBirthTillDeath(
		func() {
			ui.NotifyRoleStarting(role.Navigator{})
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
			ui.NotifyRoleEnding(role.Navigator{})
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
		ui.Error("tmb.Wait(): ", err)
	}
}

func waitForChange(interrupt <-chan bool) bool {
	ui.Info("Going to sleep until something interesting happens")
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
	ui.Info("Launching Build")
	err := tchn.RunBuild()
	if err != nil {
		ui.Warning("There are build errors! I can't go any further")
	}
	return err
}

func test() error {
	ui.Info("Running Tests")
	err := tchn.RunTests()
	if err != nil {
		ui.Warning("Some tests are failing! That's unfortunate")
	}
	return err
}

func commit() {
	ui.Info("Committing changes on branch ", git.WorkingBranch())
	git.Commit()
	git.Push()
}

func revert() {
	ui.Warning("Reverting changes")
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
	ui.Info("That's All Folks!")
	os.Exit(0)
}

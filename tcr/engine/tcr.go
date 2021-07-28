package engine

import (
	"github.com/mengdaming/tcr/tcr"
	"github.com/mengdaming/tcr/tcr/filesystem"
	"github.com/mengdaming/tcr/tcr/language"
	"github.com/mengdaming/tcr/tcr/toolchain"
	"github.com/mengdaming/tcr/tcr/vcs"

	"gopkg.in/tomb.v2"

	"os"
	"os/signal"
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

	switch mode {
	case tcr.Solo:
		// When running TCR in solo mode, there's no
		// selection menu: we directly enter driver mode
		RunAsDriver()
	case tcr.Mob:
		// When running TCR in mob mode, every participant
		// is given the possibility to switch between
		// driver and navigator modes
		ui.WaitForAction()
	}
}

func ToggleAutoPush() {
	git.EnablePush(!git.IsPushEnabled())
}

func RunAsDriver() {
	runInLoop(
		func() {
			ui.NotifyRoleStarting(tcr.DriverRole)
			git.Pull()
		},
		func(interrupt <-chan bool) {
			if waitForChange(interrupt) {
				runTCR()
			}
		},
		func() {
			ui.NotifyRoleEnding(tcr.DriverRole)
		},
	)
}

func RunAsNavigator() {
	runInLoop(
		func() {
			ui.NotifyRoleStarting(tcr.NavigatorRole)
		},
		func(interrupt <-chan bool) {
			git.Pull()
			time.Sleep(pollingPeriod)
		},
		func() {
			ui.NotifyRoleEnding(tcr.NavigatorRole)
		},
	)
}

func runInLoop(
	preLoopAction func(),
	inLoopAction func(interrupt <-chan bool),
	afterLoopAction func()) {

	// watch for interruption requests
	interrupt := make(chan bool)

	var tmb tomb.Tomb

	// The goroutine doing the work
	tmb.Go(func() error {
		preLoopAction()
		for {
			select {
			case <-tmb.Dying():
				afterLoopAction()
				return nil
			case <-interrupt:
				afterLoopAction()
				return nil
			default:
				inLoopAction(interrupt)
			}
		}
	})

	// The goroutine watching for Ctrl-C
	tmb.Go(func() error {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		ui.Warning("OK, let's stop here")
		interrupt <- true
		tmb.Kill(nil)
		return nil
	})

	err := tmb.Wait()
	if err != nil {
		ui.Error("tmb.Wait(): ", err)
	}
}

func Quit() {
	ui.Info("That's All Folks!")
	os.Exit(0)
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

func waitForChange(interrupt <-chan bool) bool {
	ui.Info("Going to sleep until something interesting happens")
	return sourceTree.Watch(
		language.DirsToWatch(sourceTree.GetBaseDir(), lang),
		lang.IsSrcFile,
		interrupt)
}

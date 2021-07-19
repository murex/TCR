package tcr

import (
	"github.com/mengdaming/tcr/tcr/language"
	"github.com/mengdaming/tcr/tcr/toolchain"
	"github.com/mengdaming/tcr/trace"
	"gopkg.in/tomb.v2"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

type WorkMode string

const (
	Solo = "solo"
	Mob  = "mob"

	GitPollingPeriod = 1 * time.Second
)

var (
	baseDir  string
	mode     WorkMode
	gitItf   GitInterface
	lang     language.Language
	tchn     toolchain.Toolchain
	autoPush bool
)

func Start(b string, m WorkMode, t string, ap bool) {
	mode = m
	autoPush = ap

	baseDir = changeDir(b)
	lang = language.DetectLanguage(baseDir)
	tchn = toolchain.NewToolchain(t, lang)

	// TODO For C++ special case (build subdirectory)
	//mkdir -p "${WORK_DIR}"
	//cd "${WORK_DIR}" || exit 1
	gitItf = NewGoGit(baseDir)

	printRunningMode(mode)
	printTCRHeader()

	switch mode {
	case Solo:
		// When running TCR in solo mode, there's no
		// selection menu: we directly enter driver mode
		runAsDriver()
	case Mob:
		// When running TCR in mob mode, every participant
		// is given the possibility to switch between
		// driver and navigator modes
		mobMainMenu()
	}
}

func printRunningMode(mode WorkMode) {
	trace.HorizontalLine()
	trace.Info("Running in ", mode, " mode")
}

func toggleAutoPush() {
	if autoPush {
		autoPush = false
	} else {
		autoPush = true
	}
}

func runAsDriver() {
	runInLoop(
		func() {
			trace.HorizontalLine()
			trace.Info("Entering Driver mode. Press CTRL-C to go back to the main menu")
			gitItf.Pull()
		},
		func(interrupt <-chan bool) {
			if waitForChanges(interrupt) {
				tcr()
			}
		},
		func() {
			trace.Info("Exiting Driver mode")
		},
	)
}

func runAsNavigator() {
	runInLoop(
		func() {
			trace.HorizontalLine()
			trace.Info("Entering Navigator mode. Press CTRL-C to go back to the main menu")
		},
		func(interrupt <-chan bool) {
			gitItf.Pull()
			time.Sleep(GitPollingPeriod)
		},
		func() {
			trace.Info("Exiting Navigator mode")
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
		trace.Warning("OK, let's stop here")
		interrupt <- true
		tmb.Kill(nil)
		return nil
	})

	err := tmb.Wait()
	if err != nil {
		trace.Error("tmb.Wait(): ", err)
	}
}

func quit() {
	trace.Info("That's All Folks!")
	os.Exit(0)
}

func tcr() {
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
	trace.Info("Launching Build")
	err := tchn.RunBuild()
	if err != nil {
		trace.Warning("There are build errors! I can't go any further")
	}
	return err
}

func test() error {
	trace.Info("Running Tests")
	err := tchn.RunTests()
	if err != nil {
		trace.Warning("Some tests are failing! That's unfortunate")
	}
	return err
}

func commit() {
	trace.Info("Committing changes on branch ", gitItf.WorkingBranch())
	gitItf.Commit()
	if autoPush {
		gitItf.Push()
	}
}

func revert() {
	trace.Warning("Reverting changes")
	for _, dir := range lang.SrcDirs() {
		gitItf.Restore(filepath.Join(baseDir, dir))
	}
}

func printTCRHeader() {
	trace.HorizontalLine()
	trace.Info("Language=", lang.Name(), ", Toolchain=", tchn.Name())

	autoPushStr := "disabled"
	if autoPush {
		autoPushStr = "enabled"
	}
	trace.Info(
		"Running on git branch \"", gitItf.WorkingBranch(),
		"\" with auto-push ", autoPushStr)
}

func changeDir(baseDir string) string {
	_, err := os.Stat(baseDir)
	switch {
	case os.IsNotExist(err):
		trace.Error("Directory ", baseDir, " does not exist")
	case os.IsPermission(err):
		trace.Error("Can't access directory ", baseDir)
	}

	err = os.Chdir(baseDir)
	if err != nil {
		trace.Error("Cannot change directory to ", baseDir)
	}

	getwd, _ := os.Getwd()
	trace.Info("Current Working Directory: ", getwd)
	return getwd
}

func waitForChanges(interrupt <-chan bool) bool {
	trace.Info("Going to sleep until something interesting happens")
	return WatchRecursive(
		language.DirsToWatch(baseDir, lang),
		lang.IsSrcFile,
		interrupt)
}

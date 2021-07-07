package tcr

import (
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
)

var (
	baseDir   string
	mode      WorkMode
	osToolbox OSToolbox
	language  Language
	toolchain string
	autoPush  bool
)

func Start(b string, m WorkMode, t string, ap bool) {
	baseDir = b
	mode = m
	toolchain = t
	autoPush = ap

	osToolbox = initOSToolbox()

	baseDir = changeDir(baseDir)
	language = detectKataLanguage(baseDir)
	// TODO For C++ special case (build subdirectory)
	//mkdir -p "${WORK_DIR}"
	//cd "${WORK_DIR}" || exit 1
	detectGitWorkingBranch(baseDir)

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
	loopWithTomb(
		func() {
			trace.HorizontalLine()
			trace.Info("Entering Driver mode. Press CTRL-C to go back to the main menu")
			pull()
		},
		func(interrupt <-chan bool) {
			if watchFileSystem(dirsToWatch(baseDir, language), interrupt) {
				tcr()
			}
		},
		func() {
			trace.Info("Exiting Driver mode")
		},
	)
}

func runAsNavigator() {
	loopWithTomb(
		func() {
			trace.HorizontalLine()
			trace.Info("Entering Navigator mode. Press CTRL-C to go back to the main menu")
		},
		func(interrupt <-chan bool) {
			pull()
		},
		func() {
			trace.Info("Exiting Navigator mode")
		},
	)
}

func loopWithTomb(
	preLoopAction func(),
	inLoopAction func(interrupt <-chan bool),
	afterLoopAction func()) {

	// watch for interruption requests
	interrupt := make(chan bool)

	var t tomb.Tomb

	// The goroutine doing the work
	t.Go(func() error {
		preLoopAction()
		for {
			select {
			case <-t.Dying():
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
	t.Go(func() error {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		//signal.Notify(sig, syscall.SIGTERM)
		<-sig
		trace.Warning("OK, let's stop here")
		interrupt <- true
		t.Kill(nil)
		return nil
	})

	err := t.Wait()
	if err != nil {
		trace.Error("t.Wait(): ", err)
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

	// TODO Rewrite in Go
	time.Sleep(1 * time.Second)
	//build_rc=0
	//case "${TOOLCHAIN}" in
	//gradle)
	//./gradlew build -x test
	//build_rc=$?
	//;;
	//maven)
	//./mvnw test-compile
	//build_rc=$?
	//;;
	//cmake)
	//${CMAKE_CMD} --build . --config Debug
	//build_rc=$?
	//;;
	//*)
	//tcr_error "Toolchain ${TOOLCHAIN} is not supported"
	//;;
	//esac
	//
	//[ $build_rc -ne 0 ] && tcr_warning "There are build errors! I can't go any further"
	//return $build_rc

	return nil
}

func test() error {
	trace.Info("Running Tests")

	// TODO Rewrite in Go
	time.Sleep(1 * time.Second)
	//test_rc=0
	//case ${TOOLCHAIN} in
	//gradle)
	//./gradlew test
	//test_rc=$?
	//;;
	//maven)
	//./mvnw test
	//test_rc=$?
	//;;
	//cmake)
	//${CTEST_CMD} --output-on-failure -C Debug
	//test_rc=$?
	//;;
	//*)
	//tcr_error "Toolchain ${TOOLCHAIN} is not supported"
	//;;
	//esac
	//
	//[ $test_rc -ne 0 ] && tcr_warning "Some tests are failing! That's unfortunate"
	//return $test_rc

	return nil
}

func commit() {
	trace.Info("Committing changes on branch ", gitWorkingBranch)
	time.Sleep(1 * time.Second)
	// TODO Call to git commit -am TCR
	if autoPush {
		push()
	}
}

func revert() {
	trace.Warning("Reverting changes")
	time.Sleep(1 * time.Second)
	// TODO Call to git checkout HEAD -- ${SRC_DIRS}
}

func printTCRHeader() {
	trace.HorizontalLine()

	trace.Info(
		"Language=", language.name(),
		", Toolchain=", language.toolchain())

	autoPushStr := "disabled"
	if autoPush {
		autoPushStr = "enabled"
	}
	trace.Info(
		"Running on git branch \"", gitWorkingBranch,
		"\" with auto-push ", autoPushStr)
}

func detectKataLanguage(baseDir string) Language {
	dir := filepath.Base(baseDir)
	switch dir {
	case "java":
		return JavaLanguage{}
	case "cpp":
		return CppLanguage{}
	default:
		trace.Error("Unrecognized language: ", dir)
	}
	return nil
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

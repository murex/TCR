package tcr

import (
	"github.com/mengdaming/tcr/trace"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"time"
)

type WorkMode string

const (
	Solo = "solo"
	Mob  = "mob"
)

var (
	mode                           WorkMode
	osToolbox                      OSToolbox
	language                       Language
	toolchain                      string
	autoPush                       bool
)

func Start(m WorkMode, t string, ap bool) {
	toolchain = t
	autoPush = ap
	mode = m

	osToolbox = initOSToolbox()

	detectKataLanguage()
	// TODO For C++ special case (build subdirectory)
	//mkdir -p "${WORK_DIR}"
	//cd "${WORK_DIR}" || exit 1
	detectGitWorkingBranch()

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
		autoPush  = true
	}
}

func runAsDriver() {
	loopWithInterrupt(
		func() {
			trace.HorizontalLine()
			trace.Info("Entering Driver mode. Press CTRL-C to go back to the main menu")
			pull()
		},
		func() {
			watchFileSystem()
			tcr()
		},
		func() {
			trace.Info("Leaving Driver mode")
		},
	)
}

func runAsNavigator() {
	loopWithInterrupt(
		func() {
			trace.HorizontalLine()
			trace.Info("Entering Navigator mode. Press CTRL-C to go back to the main menu")
		},
		func() {
			pull()
		},
		func() {
			trace.Info("Leaving Navigator mode")
		},
	)
}

func loopWithInterrupt(
	preLoopAction func(),
	inLoopAction func(),
	afterLoopAction func()) {

	c1, cancel := context.WithCancel(context.Background())

	exitCh := make(chan struct{})
	go func(ctx context.Context) {
		preLoopAction()
		for {
			inLoopAction()

			select {
			case <-ctx.Done():
				afterLoopAction()
				exitCh <- struct{}{}
				return
			default:
			}
		}
	}(c1)

	// For handling Ctrl-C interruption
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		select {
		case <-signalCh:
			cancel()
			return
		}
	}()
	<-exitCh
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

func watchFileSystem() {
	trace.Info("Going to sleep until something interesting happens")
	time.Sleep(1 * time.Second)
	// TODO ${FS_WATCH_CMD} ${SRC_DIRS} ${TEST_DIRS}
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

func detectKataLanguage() {
	// TODO Add language detection. Hard-coding java for now
	language = JavaLanguage{}
}


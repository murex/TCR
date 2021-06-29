package tcr

import (
	"github.com/mengdaming/tcr/trace"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"runtime"
	"time"
)

var (
	osToolbox                      OSToolbox
	language                       Language
	toolchain                      string
	autoPush                       bool
	gitWorkingBranch               string
	gitWorkingBranchExistsOnOrigin bool
)

func Start(t string, ap bool) {
	toolchain = t
	autoPush = ap
	initOSToolbox()
	detectKataLanguage()
	// TODO For C++ special case (build subdirectory)
	//mkdir -p "${WORK_DIR}"
	//cd "${WORK_DIR}" || exit 1
	detectGitWorkingBranch()
	whatShallWeDo()
}

func whatShallWeDo() {
	printTCRHeader()
	mainMenu()
}

func runAsDriver() {
	loopWithInterrupt(
		func() {
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

	// TODO
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
	trace.Info( "Running Tests")

	// TODO
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

func push() {
	trace.Info("Pushing changes to origin/", gitWorkingBranch)
	time.Sleep(1 * time.Second)
	// TODO Call to git push --no-recurse-submodules origin "${GIT_WORKING_BRANCH}"
	// TODO [ ${git_rc} -eq 0 ] && GIT_WORKING_BRANCH_EXISTS_ON_ORIGIN=1
	// TODO	return ${git_rc}
}

func pull() {
	if gitWorkingBranchExistsOnOrigin {
		trace.Info("Pulling latest changes from origin/", gitWorkingBranch)
		time.Sleep(1 * time.Second)
		// TODO Call to git pull --no-recurse-submodules origin "${GIT_WORKING_BRANCH}"
	} else {
		trace.Info("Working locally on branch ", gitWorkingBranch)
		time.Sleep(1 * time.Second)
	}
}

func watchFileSystem() {
	trace.Info("Going to sleep until something interesting happens")
	time.Sleep(1 * time.Second)
	// TODO ${FS_WATCH_CMD} ${SRC_DIRS} ${TEST_DIRS}
}

func printOptionsMenu() {
	trace.HorizontalLine()
	trace.Info("What shall we do?")
	trace.Info("\tD -> Driver mode")
	trace.Info("\tN -> Navigator mode")
	trace.Info("\tQ -> Quit")
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

func detectGitWorkingBranch() {
	// TODO Hardcoded for now

	// TODO GIT_WORKING_BRANCH=$(git rev-parse --abbrev-ref HEAD)
	// TODO GIT_WORKING_BRANCH_EXISTS_ON_ORIGIN=$(git branch -r | grep -c "origin/${GIT_WORKING_BRANCH}" || [ $? = 1 ])
	gitWorkingBranch = "main"
	gitWorkingBranchExistsOnOrigin = true
}

func detectKataLanguage() {
	// TODO Add language detection. Hard-coding java for now
	language = JavaLanguage{}
}

func initOSToolbox() {
	switch runtime.GOOS {
	case "darwin":
		osToolbox = MacOSToolbox{}
	case "linux":
		osToolbox = LinuxToolbox{}
	case "windows":
		osToolbox = WindowsToolbox{}
	default:
		trace.Error("OS ", runtime.GOOS, " is currently not supported")
	}
}

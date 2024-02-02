/*
Copyright (c) 2022 Murex

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

package cli

import (
	"fmt"
	"github.com/murex/tcr/desktop"
	"github.com/murex/tcr/engine"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/role"
	"github.com/murex/tcr/runmode"
	"github.com/murex/tcr/settings"
	"github.com/murex/tcr/vcs/git"
	"github.com/murex/tcr/vcs/p4"
	"os"
)

// TerminalUI is the user interface implementation when using the Command Line Interface
type TerminalUI struct {
	reportingChannel chan bool
	tcr              engine.TCRInterface
	params           params.Params
	desktop          *desktop.Desktop
	mainMenu         *menu
	roleMenu         *menu
}

const (
	enterKey = 0x0a
)

const (
	pullMenuHelper               = "Pull from remote"
	pushMenuHelper               = "Push to remote"
	syncMenuHelper               = "Synchronize with depot"
	enterDriverRoleMenuHelper    = "Driver role"
	enterNavigatorRoleMenuHelper = "Navigator role"
	openBrowserMenuHelper        = "Open in browser"
	gitAutoPushMenuHelper        = "Turn on/off git auto-push"
	quitMenuHelper               = "Quit"
	optionsMenuHelper            = "List available options"
	timerStatusMenuHelper        = "Timer status"
	quitDriverRoleMenuHelper     = "Quit Driver role"
	quitNavigatorRoleMenuHelper  = "Quit Navigator role"
)

// New creates a new instance of terminal
func New(p params.Params, tcr engine.TCRInterface) *TerminalUI {
	setLinePrefix("[" + settings.ApplicationName + "]")
	var term = TerminalUI{params: p, tcr: tcr, desktop: desktop.NewDesktop(nil)}
	tcr.AttachUI(&term, true)
	term.mainMenu = term.initMainMenu()
	term.roleMenu = term.initRoleMenu()
	term.MuteDesktopNotifications(false)
	term.StartReporting()
	StartInterruptHandler()
	return &term
}

// StartReporting tells the terminal to start reporting information
func (term *TerminalUI) StartReporting() {
	term.reportingChannel = report.Subscribe(term)
}

// MuteDesktopNotifications allows preventing desktop Notification popups from being displayed.
// Used for test automation at the moment. Could be turned into a feature later if there is need for it.
func (term *TerminalUI) MuteDesktopNotifications(muted bool) {
	if muted {
		term.desktop.MuteNotifications()
	} else {
		term.desktop.UnmuteNotifications()
	}
}

// StopReporting tells the terminal to stop reporting information
func (term *TerminalUI) StopReporting() {
	if term.reportingChannel != nil {
		report.Unsubscribe(term.reportingChannel)
	}
}

// NotifyRoleStarting tells the user that TCR engine is starting with the provided role
func (term *TerminalUI) NotifyRoleStarting(r role.Role) {
	term.ReportTitle(false, "Starting with ", r.LongName(), ". Press ? for options")
}

// NotifyRoleEnding tells the user that TCR engine is ending the provided role
func (term *TerminalUI) NotifyRoleEnding(r role.Role) {
	term.ReportInfo(false, "Ending ", r.LongName())
}

// ReportSimple reports simple messages
func (*TerminalUI) ReportSimple(_ bool, a ...any) {
	printUntouched(a...)
}

// ReportInfo reports info messages
func (*TerminalUI) ReportInfo(_ bool, a ...any) {
	printInCyan(a...)
}

// ReportTitle reports title messages
func (*TerminalUI) ReportTitle(_ bool, a ...any) {
	printHorizontalLine()
	printInCyan(a...)
}

// ReportTimer reports timer messages
func (term *TerminalUI) ReportTimer(emphasis bool, a ...any) {
	printInGreen(a...)
	term.notifyOnEmphasis(emphasis, "⏳", a...)
}

// ReportSuccess reports success messages
func (term *TerminalUI) ReportSuccess(emphasis bool, a ...any) {
	printInGreen(a...)
	term.notifyOnEmphasis(emphasis, "🟢", a...)
}

// ReportWarning reports warning messages
func (term *TerminalUI) ReportWarning(emphasis bool, a ...any) {
	printInYellow(a...)
	term.notifyOnEmphasis(emphasis, "🔶", a...)
}

// ReportError reports error messages
func (term *TerminalUI) ReportError(emphasis bool, a ...any) {
	printInRed(a...)
	term.notifyOnEmphasis(emphasis, "🟥", a...)
}

func (term *TerminalUI) notifyOnEmphasis(emphasis bool, emoji string, a ...any) {
	if emphasis {
		err := term.desktop.ShowNotification(desktop.NormalLevel, emoji+" "+settings.ApplicationName, fmt.Sprint(a...))
		if err != nil {
			term.ReportWarning(false, "Failed to show desktop notification: ", err.Error())
		}
	}
}

func (term *TerminalUI) enterMainMenu() {
	term.whatShallWeDo()
	term.runMenuLoop(term.mainMenu)
}

func (term *TerminalUI) enterRole(r role.Role) {
	// We ask first TCR engine to start...
	if err := term.runTCR(r); err != nil {
		term.ReportError(false, err)
		return
	}
	// Then we enter the role menu loop, waiting for user input
	term.runMenuLoop(term.roleMenu)
}

func (term *TerminalUI) runTCR(r role.Role) error {
	switch r {
	case role.Navigator{}:
		term.tcr.RunAsNavigator()
	case role.Driver{}:
		term.tcr.RunAsDriver()
	default:
		return fmt.Errorf("no action defined for %s", r.LongName())
	}
	return nil
}

func (term *TerminalUI) runMenuLoop(m *menu) {
	for {
		input := term.readKeyboardInput()
		matched, quit := m.matchAndRun(input)
		if quit {
			return
		}
		if !matched && input != enterKey {
			term.keyNotRecognizedMessage()
			term.listMenuOptions(m, "Please choose one of the following:")
		}
	}
}

func (term *TerminalUI) readKeyboardInput() byte {
	keyboardInput := make([]byte, 1)
	_, err := os.Stdin.Read(keyboardInput)
	if err != nil {
		term.ReportWarning(false, "Something went wrong while reading from stdin: ", err)
	}
	return keyboardInput[0]
}

func (term *TerminalUI) vcsPull() {
	term.tcr.VCSPull()
}

func (term *TerminalUI) vcsPush() {
	term.tcr.VCSPush()
}

func (term *TerminalUI) whatShallWeDo() {
	term.listMenuOptions(term.mainMenu, "What shall we do?")
}

func (term *TerminalUI) keyNotRecognizedMessage() {
	term.ReportWarning(false, "Key not recognized. Press ? for available options")
}

func (term *TerminalUI) showTimerStatus() {
	if settings.EnableMobTimer {
		if r := term.tcr.GetCurrentRole(); r != nil && r.RunsWithTimer() {
			term.tcr.ReportMobTimerStatus()
		} else {
			term.keyNotRecognizedMessage()
		}
	}
}

// ShowRunningMode shows the current running mode
func (term *TerminalUI) ShowRunningMode(mode runmode.RunMode) {
	term.ReportTitle(false, "Running in ", mode.Name(), " mode")
}

// ShowSessionInfo shows main information related to the current TCR session
func (term *TerminalUI) ShowSessionInfo() {
	info := term.tcr.GetSessionInfo()
	term.ReportTitle(false, "Base Directory: ", info.BaseDir)
	term.ReportInfo(false, "Work Directory: ", info.WorkDir)
	term.ReportInfo(false, "Language=", info.LanguageName, ", Toolchain=", info.ToolchainName)
	term.reportVCSInfo(info)
	term.reportMessageSuffix(info.MessageSuffix)
}

func (term *TerminalUI) reportVCSInfo(info engine.SessionInfo) {
	switch info.VCSName {
	case git.Name:
		autoPush := "disabled"
		if info.GitAutoPush {
			autoPush = "enabled"
		}
		term.ReportInfo(false, "Running on ", info.VCSSessionSummary, " with auto-push ", autoPush)
	case p4.Name:
		term.ReportInfo(false, "Running with ", info.VCSSessionSummary)
	default:
		term.ReportWarning(false, "VCS \"", info.VCSName, "\" is unknown")
	}
}

func (term *TerminalUI) reportMessageSuffix(suffix string) {
	if suffix == "" {
		return
	}
	term.ReportInfo(false, "Commit message suffix: \"", suffix, "\"")
}

// Confirm asks the user for confirmation
func (term *TerminalUI) Confirm(message string, defaultAnswer bool) bool {
	_ = SetRaw()
	defer Restore()

	term.ReportWarning(false, message)
	term.ReportWarning(false, "Do you want to proceed? ", yesOrNoAdvice(defaultAnswer))

	keyboardInput := make([]byte, 1)
	for {
		_, _ = os.Stdin.Read(keyboardInput)
		switch keyboardInput[0] {
		case 'y', 'Y':
			return true
		case 'n', 'N':
			return false
		case enterKey:
			return defaultAnswer
		}
	}
}

func yesOrNoAdvice(defaultAnswer bool) string {
	if defaultAnswer {
		return "[Y/n]"
	}
	return "[y/N]"
}

// Start runs the terminal session
func (term *TerminalUI) Start() {
	if term.params.Mode.IsInteractive() {
		_ = SetRaw()
		defer Restore()
	}

	switch term.params.Mode {
	case runmode.Solo{}:
		// When running TCR in solo mode, there's no selection menu:
		// we directly enter driver mode, and quit when done
		term.enterRole(role.Driver{})
		Restore()
		term.tcr.Quit()
	case runmode.Mob{}:
		// When running TCR in mob mode, every participant
		// is given the possibility to switch between
		// driver and navigator modes
		term.enterMainMenu()
	case runmode.OneShot{}:
		// When running TCR in one-shot mode, there's no selection menu:
		// we directly ask TCR engine to run one cycle and quit when done
		term.tcr.RunTCRCycle()
		term.tcr.Quit()
	case runmode.Check{}:
		// When running TCR in check mode, there's no selection menu:
		// we directly ask TCR engine to run a check and quit when done
		term.tcr.RunCheck(term.params)
		term.tcr.Quit()
	case runmode.Log{}:
		// When running TCR in log mode, there's no selection menu:
		// we directly ask TCR engine to print the commit history and quit when done
		term.tcr.PrintLog(term.params)
		term.tcr.Quit()
	case runmode.Stats{}:
		// When running TCR in stats mode, there's no selection menu:
		// we directly ask TCR engine to print the stats and quit when done
		term.tcr.PrintStats(term.params)
		term.tcr.Quit()
	default:
		term.ReportError(false, "Unknown run mode: ", term.params.Mode)
	}
}

func (term *TerminalUI) listMenuOptions(m *menu, title string) {
	term.ReportTitle(false, title)
	for _, option := range m.getOptions() {
		term.ReportInfo(false, (*option).toString())
	}
}

func (term *TerminalUI) initMainMenu() *menu {
	m := newMenu("Main menu")
	m.addOptions(
		newMenuOption('O', openBrowserMenuHelper,
			term.webMenuEnabler(),
			term.openBrowserMenuAction(), false),
		newMenuOption('D', enterDriverRoleMenuHelper, nil,
			term.enterRoleMenuAction(role.Driver{}), false),
		newMenuOption('N', enterNavigatorRoleMenuHelper, nil,
			term.enterRoleMenuAction(role.Navigator{}), false),
		newMenuOption('P', gitAutoPushMenuHelper,
			term.gitMenuEnabler(),
			term.autoPushMenuAction(), false),
		newMenuOption('L', pullMenuHelper,
			term.gitMenuEnabler(),
			term.vcsPullMenuAction(), false),
		newMenuOption('S', pushMenuHelper,
			term.gitMenuEnabler(),
			term.vcsPushMenuAction(), false),
		newMenuOption('Y', syncMenuHelper,
			term.p4MenuEnabler(),
			term.vcsPullMenuAction(), false),
		newMenuOption('Q', quitMenuHelper, nil,
			term.quitMenuAction(), true),
		newMenuOption('?', optionsMenuHelper, nil,
			term.optionsMenuAction(m), false),
	)
	return m
}

func (term *TerminalUI) initRoleMenu() *menu {
	m := newMenu("Role menu")
	m.addOptions(
		newMenuOption('T', timerStatusMenuHelper,
			term.timerStatusMenuEnabler(),
			term.timerStatusMenuAction(), false),
		newMenuOption('Q', quitDriverRoleMenuHelper,
			term.quitRoleMenuEnabler(role.Driver{}),
			term.quitRoleMenuAction(), true),
		newMenuOption('Q', quitNavigatorRoleMenuHelper,
			term.quitRoleMenuEnabler(role.Navigator{}),
			term.quitRoleMenuAction(), true),
		newMenuOption('?', optionsMenuHelper, nil,
			term.optionsMenuAction(m), false),
	)
	return m
}

func (term *TerminalUI) enterRoleMenuAction(r role.Role) menuAction {
	return func() {
		term.enterRole(r)
		term.whatShallWeDo()
	}
}

func (term *TerminalUI) openBrowserMenuAction() menuAction {
	return func() {
		desktop.OpenBrowser(term.params.PortNumber)
	}
}

func (term *TerminalUI) gitMenuEnabler() menuEnabler {
	return func() bool {
		return term.params.VCS == git.Name
	}
}

func (term *TerminalUI) p4MenuEnabler() menuEnabler {
	return func() bool {
		return term.params.VCS == p4.Name
	}
}

func (term *TerminalUI) webMenuEnabler() menuEnabler {
	return func() bool {
		// When port number is set to 0, HTTP server is not running
		return term.params.PortNumber != 0
	}
}

func (term *TerminalUI) autoPushMenuAction() menuAction {
	return func() {
		term.tcr.ToggleAutoPush()
		term.ShowSessionInfo()
		term.whatShallWeDo()
	}
}

func (term *TerminalUI) vcsPullMenuAction() menuAction {
	return func() {
		term.vcsPull()
		term.whatShallWeDo()
	}
}

func (term *TerminalUI) vcsPushMenuAction() menuAction {
	return func() {
		term.vcsPush()
		term.whatShallWeDo()
	}
}

func (term *TerminalUI) optionsMenuAction(m *menu) menuAction {
	return func() {
		term.listMenuOptions(m, "Available Options:")
	}
}

func (term *TerminalUI) quitMenuAction() menuAction {
	return func() {
		Restore()
		term.tcr.Quit()
	}
}

func (term *TerminalUI) timerStatusMenuEnabler() menuEnabler {
	return func() bool {
		r := term.tcr.GetCurrentRole()
		return settings.EnableMobTimer && r != nil && r.RunsWithTimer()
	}
}

func (term *TerminalUI) timerStatusMenuAction() menuAction {
	return func() {
		term.showTimerStatus()
	}
}

func (term *TerminalUI) quitRoleMenuEnabler(r role.Role) menuEnabler {
	return func() bool {
		return term.tcr.GetCurrentRole() == r
	}
}

func (term *TerminalUI) quitRoleMenuAction() menuAction {
	return func() {
		term.ReportWarning(false, "OK, I heard you")
		term.tcr.Stop()
	}
}

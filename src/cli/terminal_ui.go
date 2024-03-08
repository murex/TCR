/*
Copyright (c) 2024 Murex

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
	"github.com/murex/tcr/report/role_event"
	"github.com/murex/tcr/report/text"
	"github.com/murex/tcr/report/timer_event"
	"github.com/murex/tcr/role"
	"github.com/murex/tcr/runmode"
	"github.com/murex/tcr/settings"
	"github.com/murex/tcr/timer"
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
	soloMenu         *menu
	mobMenu          *menu
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
	quitTCRMenuHelper            = "Quit TCR"
)

const timerMessagePrefix = "(Mob Timer) "

// New creates a new instance of terminal
func New(p params.Params, tcr engine.TCRInterface) *TerminalUI {
	setLinePrefix("[" + settings.ApplicationName + "]")
	var term = TerminalUI{params: p, tcr: tcr, desktop: desktop.NewDesktop(nil)}
	tcr.AttachUI(&term, true)
	term.soloMenu = term.initSoloMenu()
	term.mobMenu = term.initMobMenu()
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

// ReportSimple reports simple messages
func (*TerminalUI) ReportSimple(_ bool, payload text.Message) {
	printUntouched(payload.ToString())
}

// ReportInfo reports info messages
func (term *TerminalUI) ReportInfo(_ bool, payload text.Message) {
	term.printInfo(payload.ToString())
}

// ReportTitle reports title messages
func (term *TerminalUI) ReportTitle(_ bool, payload text.Message) {
	term.printTitle(payload.ToString())
}

// ReportSuccess reports success messages
func (term *TerminalUI) ReportSuccess(emphasis bool, payload text.Message) {
	txt := payload.ToString()
	term.printSuccess(txt)
	term.notifyOnEmphasis(emphasis, "🟢", txt)
}

// ReportWarning reports warning messages
func (term *TerminalUI) ReportWarning(emphasis bool, payload text.Message) {
	txt := payload.ToString()
	term.printWarning(txt)
	term.notifyOnEmphasis(emphasis, "🔶", txt)
}

// ReportError reports error messages
func (term *TerminalUI) ReportError(emphasis bool, payload text.Message) {
	txt := payload.ToString()
	term.printError(txt)
	term.notifyOnEmphasis(emphasis, "🟥", txt)
}

// ReportRoleEvent reports role event messages
func (term *TerminalUI) ReportRoleEvent(_ bool, payload role_event.Message) {
	switch payload.Trigger {
	case role_event.TriggerStart:
		term.printRoleEvent("Starting with ", payload.Role.LongName(), ". Press ? for options")
	case role_event.TriggerEnd:
		term.printRoleEvent("Ending ", payload.Role.LongName())
	}
}

// ReportTimerEvent reports timer event messages
func (term *TerminalUI) ReportTimerEvent(emphasis bool, payload timer_event.Message) {
	var txt string
	switch payload.Trigger {
	case timer_event.TriggerStart:
		txt = fmt.Sprint(timerMessagePrefix, "Starting ",
			timer_event.FormatDuration(payload.Timeout), " countdown")
	case timer_event.TriggerCountdown:
		txt = fmt.Sprint(timerMessagePrefix, "Your turn ends in ",
			timer_event.FormatDuration(payload.Remaining))
	case timer_event.TriggerStop:
		txt = fmt.Sprint(timerMessagePrefix, "Stopping countdown after ",
			timer_event.FormatDuration(payload.Elapsed))
	case timer_event.TriggerTimeout:
		txt = fmt.Sprint(timerMessagePrefix, "Time's up. Time to rotate! You are ",
			timer_event.FormatDuration(payload.Remaining.Abs()), " over!")
	}
	term.printTimerEvent(payload.Trigger == timer_event.TriggerTimeout, txt)
	term.notifyOnEmphasis(emphasis, "⏳", txt)
}

// printTitle prints title messages locally in the terminal
func (*TerminalUI) printTitle(a ...any) {
	printHorizontalLine()
	printInCyan(a...)
}

// printInfo prints info messages locally in the terminal
func (*TerminalUI) printInfo(a ...any) {
	printInCyan(a...)
}

// printSuccess prints info messages locally in the terminal
func (*TerminalUI) printSuccess(a ...any) {
	printInGreen(a...)
}

// printWarning prints warning messages locally in the terminal
func (*TerminalUI) printWarning(a ...any) {
	printInYellow(a...)
}

// printError prints error messages locally in the terminal
func (*TerminalUI) printError(a ...any) {
	printInRed(a...)
}

// printRoleEvent prints role event messages locally in the terminal
func (*TerminalUI) printRoleEvent(a ...any) {
	printInYellow(a...)
}

// printTimerEvent prints timer event messages locally in the terminal
func (*TerminalUI) printTimerEvent(timeout bool, a ...any) {
	if timeout {
		printInYellow(a...)
	} else {
		printInGreen(a...)
	}
}

func (term *TerminalUI) notifyOnEmphasis(emphasis bool, emoji string, a ...any) {
	if emphasis {
		err := term.desktop.ShowNotification(desktop.NormalLevel, emoji+" "+settings.ApplicationName, fmt.Sprint(a...))
		if err != nil {
			term.printWarning("Failed to show desktop notification: ", err.Error())
		}
	}
}

func (term *TerminalUI) enterSoloMenu() {
	term.enterRole(role.Driver{})
	// Then we enter the solo menu loop, waiting for user input
	term.runMenuLoop(term.soloMenu)
}

func (term *TerminalUI) enterMobMenu() {
	term.whatShallWeDo()
	term.runMenuLoop(term.mobMenu)
}

func (term *TerminalUI) enterRole(r role.Role) {
	if err := term.runTCR(r); err != nil {
		term.printError(err.Error())
	}
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
		term.printWarning("Something went wrong while reading from stdin: ", err)
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
	if term.params.Mode.IsMultiRole() {
		term.listMenuOptions(term.mobMenu, "What shall we do?")
	}
}

func (term *TerminalUI) keyNotRecognizedMessage() {
	term.ReportWarning(false, "Key not recognized. Press ? for available options")
}

func (term *TerminalUI) showTimerStatus() {
	mts := term.tcr.GetMobTimerStatus()
	switch mts.State {
	case timer.StateOff:
		term.printInfo("Mob Timer is off")
	case timer.StateRunning:
		term.printInfo("Mob Timer: ",
			timer_event.FormatDuration(mts.Elapsed), " done, ",
			timer_event.FormatDuration(mts.Remaining), " to go")
	case timer.StateTimeout:
		term.printWarning("Mob Timer has timed out: ",
			timer_event.FormatDuration(mts.Remaining.Abs()), " over!")
	case timer.StateStopped:
		term.printInfo("Mob Timer was interrupted")
	}
}

// ShowRunningMode shows the current running mode
func (term *TerminalUI) ShowRunningMode(mode runmode.RunMode) {
	term.printTitle("Running in ", mode.Name(), " mode")
}

// ShowSessionInfo shows main information related to the current TCR session
func (term *TerminalUI) ShowSessionInfo() {
	info := term.tcr.GetSessionInfo()
	term.printTitle("Base Directory: ", info.BaseDir)
	term.printInfo("Work Directory: ", info.WorkDir)
	term.printInfo("Language=", info.LanguageName, ", Toolchain=", info.ToolchainName)
	term.printVCSInfo(info)
	term.printMessageSuffix(info.MessageSuffix)
}

func (term *TerminalUI) printVCSInfo(info engine.SessionInfo) {
	switch info.VCSName {
	case git.Name:
		autoPush := "disabled"
		if info.GitAutoPush {
			autoPush = "enabled"
		}
		term.printInfo("Running on ", info.VCSSessionSummary, " with auto-push ", autoPush)
	case p4.Name:
		term.printInfo("Running with ", info.VCSSessionSummary)
	default:
		term.printWarning("VCS \"", info.VCSName, "\" is unknown")
	}
}

func (term *TerminalUI) printMessageSuffix(suffix string) {
	if suffix == "" {
		return
	}
	term.printInfo("Commit message suffix: \"", suffix, "\"")
}

// Confirm asks the user for confirmation
func (term *TerminalUI) Confirm(message string, defaultAnswer bool) bool {
	_ = SetRaw()
	defer Restore()

	term.printWarning(message)
	term.printWarning("Do you want to proceed? ", yesOrNoAdvice(defaultAnswer))

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
		term.enterSoloMenu()
		term.tcr.Quit()
	case runmode.Mob{}:
		// When running TCR in mob mode, every participant
		// is given the possibility to switch between
		// driver and navigator modes
		term.enterMobMenu()
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
		term.printError("Unknown run mode: ", term.params.Mode)
	}
}

func (term *TerminalUI) listMenuOptions(m *menu, title string) {
	term.printTitle(title)
	for _, option := range m.getOptions() {
		term.printInfo((*option).toString())
	}
}

func (term *TerminalUI) initSoloMenu() *menu {
	m := newMenu("Solo menu")
	m.addOptions(
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
		newMenuOption('Q', quitTCRMenuHelper,
			term.quitRoleMenuEnabler(role.Driver{}),
			term.quitRoleMenuAction(), true),
		newMenuOption('?', optionsMenuHelper, nil,
			term.optionsMenuAction(m), false),
	)
	return m
}

func (term *TerminalUI) initMobMenu() *menu {
	m := newMenu("Mob main menu")
	m.addOptions(
		newMenuOption('O', openBrowserMenuHelper,
			term.webMenuEnabler(),
			term.openBrowserMenuAction(), false),
		newMenuOption('D', enterDriverRoleMenuHelper,
			term.enterRoleMenuEnabler(role.Driver{}),
			term.enterRoleMenuAction(role.Driver{}), false),
		newMenuOption('N', enterNavigatorRoleMenuHelper,
			term.enterRoleMenuEnabler(role.Navigator{}),
			term.enterRoleMenuAction(role.Navigator{}), false),
		newMenuOption('T', timerStatusMenuHelper,
			term.timerStatusMenuEnabler(),
			term.timerStatusMenuAction(), false),
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
		newMenuOption('Q', quitMenuHelper,
			term.quitRoleMenuEnabler(nil),
			term.quitMenuAction(), true),
		newMenuOption('Q', quitDriverRoleMenuHelper,
			term.quitRoleMenuEnabler(role.Driver{}),
			term.quitRoleMenuAction(), false),
		newMenuOption('Q', quitNavigatorRoleMenuHelper,
			term.quitRoleMenuEnabler(role.Navigator{}),
			term.quitRoleMenuAction(), false),
		newMenuOption('?', optionsMenuHelper, nil,
			term.optionsMenuAction(m), false),
	)
	return m
}

func (term *TerminalUI) enterRoleMenuEnabler(r role.Role) menuEnabler {
	return func() bool {
		return term.tcr.GetCurrentRole() != r
	}
}

func (term *TerminalUI) enterRoleMenuAction(r role.Role) menuAction {
	return func() {
		term.enterRole(r)
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

func (*TerminalUI) timerStatusMenuEnabler() menuEnabler {
	return func() bool {
		return settings.EnableMobTimer
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
		term.whatShallWeDo()
	}
}

//go:build test_helper

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

package engine

import (
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/role"
	"github.com/murex/tcr/status"
	"github.com/murex/tcr/timer"
	"github.com/murex/tcr/ui"
	"github.com/murex/tcr/variant"
)

// TCRCall is used to track calls to TCR operations
type TCRCall string

// Possible values for TCRCall
const (
	TCRCallQuit              TCRCall = "quit"
	TCRCallToggleAutoPush    TCRCall = "toggle-auto-push"
	TCRCallGetSessionInfo    TCRCall = "get-session-info"
	TCRCallRunAsDriver       TCRCall = "run-as-driver"
	TCRCallRunAsNavigator    TCRCall = "run-as-navigator"
	TCRCallStop              TCRCall = "stop"
	TCRCallAbortCommand      TCRCall = "abort-command"
	TCRCallGetMobTimerStatus TCRCall = "get-mob-timer-status"
	TCRCallRunTcrCycle       TCRCall = "run-tcr-cycle"
	TCRCallRunCheck          TCRCall = "run-check"
	TCRCallPrintLog          TCRCall = "print-log"
	TCRCallPrintStats        TCRCall = "print-stats"
	TCRCallVCSPull           TCRCall = "vcs-pull"
	TCRCallVCSPush           TCRCall = "vcs-push"
	TCRCallGenerateRetro     TCRCall = "generate-retro"
)

var NoTCRCall []TCRCall

// FakeTCREngine is a TCR engine fake. Used mainly for testing peripheral packages
// such as cli.
type FakeTCREngine struct {
	TCREngine
	timerStatus timer.CurrentState
	callRecord  []TCRCall
	returnCode  int
	info        *SessionInfo
}

// NewFakeTCREngine creates a FakeToolchain instance
func NewFakeTCREngine() *FakeTCREngine {
	return &FakeTCREngine{
		returnCode:  0,
		timerStatus: timer.CurrentState{State: timer.StateOff, Timeout: 0, Elapsed: 0, Remaining: 0},
		info: &SessionInfo{
			BaseDir:           "fake",
			WorkDir:           "fake",
			LanguageName:      "fake",
			ToolchainName:     "fake",
			VCSName:           "fake",
			VCSSessionSummary: "VCS session \"fake\"",
			GitAutoPush:       false,
			Variant:           variant.Relaxed.Name(),
		},
	}
}

func (fake *FakeTCREngine) recordCall(call TCRCall) {
	fake.callRecord = append(fake.callRecord, call)
}

// GetCallHistory returns the list of TCRCall events tracked by FakeTCREngine
func (fake *FakeTCREngine) GetCallHistory() []TCRCall {
	return fake.callRecord
}

// AttachUI plugs a user interface to TCR
func (fake *FakeTCREngine) AttachUI(_ ui.UserInterface, _ bool) {
	// nothing to do in fake context
}

// Init initializes the TCR engine with the provided parameters
func (fake *FakeTCREngine) Init(_ params.Params) {
	// nothing to do in fake context
}

// GetSessionInfo returns a SessionInfo struct filled with "fake" values
func (fake *FakeTCREngine) GetSessionInfo() SessionInfo {
	fake.recordCall(TCRCallGetSessionInfo)
	return *fake.info
}

// Quit is the exit point for TCR application. The FakeTCREngine implementation
// overrides the default behaviour in order to bypass the call to os.Exit().
// Instead, the return code is stored in returnCode attribute
func (fake *FakeTCREngine) Quit() {
	fake.recordCall(TCRCallQuit)
	fake.returnCode = status.GetReturnCode()
}

// ToggleAutoPush toggles VCS auto-push state
func (fake *FakeTCREngine) ToggleAutoPush() {
	fake.recordCall(TCRCallToggleAutoPush)
}

// RunAsDriver tells TCR engine to start running with driver role
func (fake *FakeTCREngine) RunAsDriver() {
	fake.currentRole = role.Driver{}
	fake.recordCall(TCRCallRunAsDriver)
}

// RunAsNavigator tells TCR engine to start running with navigator role
func (fake *FakeTCREngine) RunAsNavigator() {
	fake.currentRole = role.Navigator{}
	fake.recordCall(TCRCallRunAsNavigator)
}

// Stop is the entry point for telling TCR engine to stop its current operations
func (fake *FakeTCREngine) Stop() {
	fake.currentRole = nil
	fake.recordCall(TCRCallStop)
}

// GetMobTimerStatus returns the status of the mob timer
func (fake *FakeTCREngine) GetMobTimerStatus() timer.CurrentState {
	fake.recordCall(TCRCallGetMobTimerStatus)
	return fake.timerStatus
}

// SetMobTimerStatus sets the status of the mob timer
func (fake *FakeTCREngine) SetMobTimerStatus(state timer.CurrentState) {
	fake.timerStatus = state
}

// AbortCommand triggers interruption of an ongoing TCR cycle operation
func (fake *FakeTCREngine) AbortCommand() {
	fake.recordCall(TCRCallAbortCommand)
}

// RunTCRCycle is the core of TCR engine: e.g. it runs one test && commit || revert cycle
func (fake *FakeTCREngine) RunTCRCycle() {
	fake.recordCall(TCRCallRunTcrCycle)
}

// RunCheck checks the provided parameters and prints out corresponding report
func (fake *FakeTCREngine) RunCheck(_ params.Params) {
	fake.recordCall(TCRCallRunCheck)
}

// PrintLog prints the TCR VCS commit history
func (fake *FakeTCREngine) PrintLog(_ params.Params) {
	fake.recordCall(TCRCallPrintLog)
}

// PrintStats prints the TCR execution stats
func (fake *FakeTCREngine) PrintStats(_ params.Params) {
	fake.recordCall(TCRCallPrintStats)
}

// VCSPull runs a fake VCS pull command
func (fake *FakeTCREngine) VCSPull() {
	fake.recordCall(TCRCallVCSPull)
}

// VCSPush runs a fake VCS push command
func (fake *FakeTCREngine) VCSPush() {
	fake.recordCall(TCRCallVCSPush)
}

// GenerateRetro Generates a retrospective file using stats
func (fake *FakeTCREngine) GenerateRetro(_ params.Params) {
	fake.recordCall(TCRCallGenerateRetro)
}

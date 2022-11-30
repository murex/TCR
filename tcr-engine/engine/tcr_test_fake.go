//go:build test_helper

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

package engine

import (
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/role"
	"github.com/murex/tcr/status"
	"github.com/murex/tcr/ui"
)

// TcrCall is used to track calls to TCR operations
type TcrCall string

// Possible values for TcrCall
const (
	TcrCallQuit                 TcrCall = "quit"
	TcrCallToggleAutoPush       TcrCall = "toggle-auto-push"
	TcrCallGetSessionInfo       TcrCall = "get-session-info"
	TcrCallRunAsDriver          TcrCall = "run-as-driver"
	TcrCallRunAsNavigator       TcrCall = "run-as-navigator"
	TcrCallStop                 TcrCall = "stop"
	TcrCallReportMobTimerStatus TcrCall = "report-mob-timer-status"
	TcrCallRunTcrCycle          TcrCall = "run-tcr-cycle"
	TcrCallRunCheck             TcrCall = "run-check"
	TcrCallPrintLog             TcrCall = "print-log"
	TcrCallPrintStats           TcrCall = "print-stats"
)

// FakeTcrEngine is a TCR engine fake. Used mainly for testing peripheral packages
// such as cli.
type FakeTcrEngine struct {
	TcrEngine
	callRecord []TcrCall
	returnCode int
	info       *SessionInfo
}

// NewFakeTcrEngine creates a FakeToolchain instance
func NewFakeTcrEngine() *FakeTcrEngine {
	return &FakeTcrEngine{
		returnCode: 0,
		info: &SessionInfo{
			BaseDir:       "fake",
			WorkDir:       "fake",
			LanguageName:  "fake",
			ToolchainName: "fake",
			AutoPush:      false,
			CommitOnFail:  false,
			BranchName:    "fake",
		},
	}
}

func (fake *FakeTcrEngine) recordCall(call TcrCall) {
	fake.callRecord = append(fake.callRecord, call)
}

// GetCallHistory returns the list of TcrCall events tracked by FakeTcrEngine
func (fake *FakeTcrEngine) GetCallHistory() []TcrCall {
	return fake.callRecord
}

// Init initializes the TCR engine with the provided parameters, and wires it to the user interface.
func (fake *FakeTcrEngine) Init(_ ui.UserInterface, _ params.Params) {}

// GetSessionInfo returns a SessionInfo struct filled with "fake" values
func (fake *FakeTcrEngine) GetSessionInfo() SessionInfo {
	fake.recordCall(TcrCallGetSessionInfo)
	return *fake.info
}

// Quit is the exit point for TCR application. The FakeTcrEngine implementation
// overrides the default behaviour in order to bypass the call to os.Exit().
// Instead, the return code is stored in returnCode attribute
func (fake *FakeTcrEngine) Quit() {
	fake.recordCall(TcrCallQuit)
	fake.returnCode = status.GetReturnCode()
}

// ToggleAutoPush toggles git auto-push state
func (fake *FakeTcrEngine) ToggleAutoPush() {
	fake.recordCall(TcrCallToggleAutoPush)
}

// RunAsDriver tells TCR engine to start running with driver role
func (fake *FakeTcrEngine) RunAsDriver() {
	fake.currentRole = role.Driver{}
	fake.recordCall(TcrCallRunAsDriver)
}

// RunAsNavigator tells TCR engine to start running with navigator role
func (fake *FakeTcrEngine) RunAsNavigator() {
	fake.currentRole = role.Navigator{}
	fake.recordCall(TcrCallRunAsNavigator)
}

// Stop is the entry point for telling TCR engine to stop its current operations
func (fake *FakeTcrEngine) Stop() {
	fake.recordCall(TcrCallStop)
}

// ReportMobTimerStatus reports the status of the mob timer
func (fake *FakeTcrEngine) ReportMobTimerStatus() {
	fake.recordCall(TcrCallReportMobTimerStatus)
}

// RunTCRCycle is the core of TCR engine: e.g. it runs one test && commit || revert cycle
func (fake *FakeTcrEngine) RunTCRCycle() {
	fake.recordCall(TcrCallRunTcrCycle)
}

// RunCheck checks the provided parameters and prints out corresponding report
func (fake *FakeTcrEngine) RunCheck(_ params.Params) {
	fake.recordCall(TcrCallRunCheck)
}

// PrintLog prints the TCR git commit history
func (fake *FakeTcrEngine) PrintLog(_ params.Params) {
	fake.recordCall(TcrCallPrintLog)
}

// PrintStats prints the TCR execution stats
func (fake *FakeTcrEngine) PrintStats(params params.Params) {
	fake.recordCall(TcrCallPrintStats)
}

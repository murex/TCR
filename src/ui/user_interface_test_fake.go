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

package ui

import (
	"github.com/murex/tcr/runmode"
)

// Call is used to track calls to UI operations
type Call string

// Possible values for Call
const (
	CallStart                    Call = "start"
	CallShowRunningMode          Call = "show-running-mode"
	CallShowSessionInfo          Call = "show-session-info"
	CallConfirm                  Call = "confirm"
	CallStartReporting           Call = "start-reporting"
	CallStopReporting            Call = "stop-reporting"
	CallMuteDesktopNotifications Call = "mute-desktop-notifications"
)

var NoCall []Call

// FakeUI provides a fake UI for running tests. It basically does nothing
// apart from stubbing out all UI behaviors
type FakeUI struct {
	callRecord []Call
}

// NewFakeUI creates a new instance of a fake UI
func NewFakeUI() *FakeUI {
	return &FakeUI{}
}

func (ui *FakeUI) recordCall(call Call) {
	ui.callRecord = append(ui.callRecord, call)
}

// GetCallHistory returns the list of Call events tracked by FakeUI
func (ui *FakeUI) GetCallHistory() []Call {
	return ui.callRecord
}

// Start does nothing in FakeUI
func (ui *FakeUI) Start() {
	ui.recordCall(CallStart)
}

// ShowRunningMode does nothing in FakeUI
func (ui *FakeUI) ShowRunningMode(_ runmode.RunMode) {
	ui.recordCall(CallShowRunningMode)
}

// ShowSessionInfo does nothing in FakeUI
func (ui *FakeUI) ShowSessionInfo() {
	ui.recordCall(CallShowSessionInfo)
}

// Confirm always returns true in FakeUI
func (ui *FakeUI) Confirm(_ string, _ bool) bool {
	ui.recordCall(CallConfirm)
	return true
}

// StartReporting does nothing in FakeUI
func (ui *FakeUI) StartReporting() {
	ui.recordCall(CallStartReporting)
}

// StopReporting does nothing in FakeUI
func (ui *FakeUI) StopReporting() {
	ui.recordCall(CallStopReporting)
}

// MuteDesktopNotifications does nothing in FakeUI
func (ui *FakeUI) MuteDesktopNotifications(_ bool) {
	ui.recordCall(CallMuteDesktopNotifications)
}

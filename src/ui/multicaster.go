/*
Copyright (c) 2023 Murex

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
	"github.com/murex/tcr/role"
	"github.com/murex/tcr/runmode"
)

// Multicaster provides the mechanisms to distribute UI messages to several user interfaces.
// It implements UserInterface so that it can be easily plugged between TCR engine and
// user interface instance(s)
type Multicaster struct {
	uiList    []UserInterface
	uiPrimary UserInterface
}

// NewMulticaster creates a new Multicaster instance
func NewMulticaster() *Multicaster {
	return &Multicaster{uiList: []UserInterface{}}
}

// Register allows registering of a new user interface.
// If primary is true, the new user interface becomes the primary UI
func (m *Multicaster) Register(u UserInterface, primary bool) {
	m.uiList = append(m.uiList, u)
	if primary {
		m.uiPrimary = u
	}
}

// Start sends Start message to all registered user interfaces
func (m *Multicaster) Start() {
	for _, u := range m.uiList {
		u.Start()
	}
}

// ShowRunningMode sends ShowRunningMode message to all registered user interfaces
func (m *Multicaster) ShowRunningMode(mode runmode.RunMode) {
	for _, u := range m.uiList {
		u.ShowRunningMode(mode)
	}
}

// NotifyRoleStarting sends NotifyRoleStarting message to all registered user interfaces
func (m *Multicaster) NotifyRoleStarting(r role.Role) {
	for _, u := range m.uiList {
		u.NotifyRoleStarting(r)
	}
}

// NotifyRoleEnding sends NotifyRoleEnding message to all registered user interfaces
func (m *Multicaster) NotifyRoleEnding(r role.Role) {
	for _, u := range m.uiList {
		u.NotifyRoleEnding(r)
	}
}

// ShowSessionInfo sends ShowSessionInfo message to all registered user interfaces
func (m *Multicaster) ShowSessionInfo() {
	for _, u := range m.uiList {
		u.ShowSessionInfo()
	}
}

// Confirm sends Confirm request to the primary user interface.
// If no primary user interface is set, returns automatic confirmation
func (m *Multicaster) Confirm(message string, def bool) bool {
	// If no main UI is defined, do automatic confirmation
	return m.uiPrimary == nil || m.uiPrimary.Confirm(message, def)
}

// StartReporting sends StartReporting message to all registered user interfaces
func (m *Multicaster) StartReporting() {
	for _, u := range m.uiList {
		u.StartReporting()
	}
}

// StopReporting sends StopReporting message to all registered user interfaces
func (m *Multicaster) StopReporting() {
	for _, u := range m.uiList {
		u.StopReporting()
	}
}

// MuteDesktopNotifications sends MuteDesktopNotifications message to all registered user interfaces
func (m *Multicaster) MuteDesktopNotifications(muted bool) {
	for _, u := range m.uiList {
		u.MuteDesktopNotifications(muted)
	}
}

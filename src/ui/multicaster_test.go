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
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_primary_user_interface_registration(t *testing.T) {
	multicaster := NewMulticaster()
	ui := NewFakeUI()
	multicaster.Register(ui, true)

	assert.Equal(t, []UserInterface{ui}, multicaster.uiList)
	assert.Equal(t, ui, multicaster.uiPrimary)
}

func Test_secondary_user_interface_registration(t *testing.T) {
	multicaster := NewMulticaster()
	ui := NewFakeUI()
	multicaster.Register(ui, false)

	assert.Equal(t, []UserInterface{ui}, multicaster.uiList)
	assert.NotEqual(t, ui, multicaster.uiPrimary)
}

func Test_multiple_user_interfaces_registration(t *testing.T) {
	multicaster := NewMulticaster()
	ui1 := NewFakeUI()
	ui2 := NewFakeUI()
	multicaster.Register(ui1, false)
	multicaster.Register(ui2, false)

	assert.Equal(t, []UserInterface{ui1, ui2}, multicaster.uiList)
}

func Test_multicasting_ui_messages(t *testing.T) {
	var multicaster *Multicaster
	testFlags := []struct {
		desc      string
		operation func()
		expected  Call
	}{
		{
			"start",
			func() { multicaster.Start() },
			CallStart,
		},
		{
			"show running mode",
			func() { multicaster.ShowRunningMode(nil) },
			CallShowRunningMode,
		},
		{
			"show session info",
			func() { multicaster.ShowSessionInfo() },
			CallShowSessionInfo,
		},
		{
			"start reporting",
			func() { multicaster.StartReporting() },
			CallStartReporting,
		},
		{
			"stop reporting",
			func() { multicaster.StopReporting() },
			CallStopReporting,
		},
		{
			"mute desktop notifications",
			func() { multicaster.MuteDesktopNotifications(true) },
			CallMuteDesktopNotifications,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			multicaster = NewMulticaster()
			uis := []*FakeUI{NewFakeUI(), NewFakeUI()}
			for i, ui := range uis {
				multicaster.Register(ui, i == 0)
			}

			tt.operation()

			for _, ui := range uis {
				assert.Contains(t, ui.GetCallHistory(), tt.expected)
			}
		})
	}
}

func Test_sending_confirm_message_to_primary_ui_only(t *testing.T) {
	multicaster := NewMulticaster()
	primary := NewFakeUI()
	secondary1 := NewFakeUI()
	secondary2 := NewFakeUI()
	multicaster.Register(primary, true)
	multicaster.Register(secondary1, false)
	multicaster.Register(secondary2, false)

	multicaster.Confirm("some question", true)

	assert.Contains(t, primary.GetCallHistory(), CallConfirm)
	assert.NotContains(t, secondary1.GetCallHistory(), CallConfirm)
	assert.NotContains(t, secondary2.GetCallHistory(), CallConfirm)
}

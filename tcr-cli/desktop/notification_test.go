/*
Copyright (c) 2021 Murex

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

package desktop

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_show_notification(t *testing.T) {
	var testFlags = []struct {
		desc        string
		level       NotificationLevel
		title       string
		message     string
		muted       bool
		returnError error
	}{
		{
			desc:        "Normal Level",
			level:       NormalLevel,
			title:       "some normal level title",
			message:     "some normal level message",
			muted:       false,
			returnError: nil,
		},
		{
			desc:        "High Level",
			level:       HighLevel,
			title:       "some high level title",
			message:     "some high level message",
			muted:       false,
			returnError: nil,
		},
		{
			desc:        "With Error",
			level:       NormalLevel,
			title:       "some title",
			message:     "some message",
			muted:       false,
			returnError: errors.New("Some Error"),
		},
		{
			desc:        "When Muted",
			level:       HighLevel,
			title:       "some title",
			message:     "some message",
			muted:       true,
			returnError: nil,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			fake := &FakeNotifier{ReturnError: tt.returnError}
			desktop := NewDesktop(fake)

			var expectedTitle, expectedMessage string
			var expectedLevel NotificationLevel
			var expectedError error
			if tt.muted {
				desktop.MuteNotifications()
				expectedLevel = 0
				expectedTitle = ""
				expectedMessage = ""
				expectedError = nil
			} else {
				desktop.UnmuteNotifications()
				expectedLevel = tt.level
				expectedTitle = tt.title
				expectedMessage = tt.message
				expectedError = tt.returnError
			}
			err := desktop.ShowNotification(tt.level, tt.title, tt.message)
			assert.Equal(t, expectedLevel, fake.LastLevel)
			assert.Equal(t, expectedTitle, fake.LastTitle)
			assert.Equal(t, expectedMessage, fake.LastMessage)
			assert.Equal(t, expectedError, err)
		})
	}
}

func Test_mute_notifications(t *testing.T) {
	desktop := NewDesktop(nil)
	desktop.MuteNotifications()
	assert.True(t, desktop.IsMuted())
}

func Test_unmute_notifications(t *testing.T) {
	desktop := NewDesktop(nil)
	desktop.UnmuteNotifications()
	assert.False(t, desktop.IsMuted())
}

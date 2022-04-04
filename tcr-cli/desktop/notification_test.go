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
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeNotifier struct {
	lastLevel   NotificationLevel
	lastTitle   string
	lastMessage string
}

func (n fakeNotifier) normalLevelNotification(title string, message string) error {
	n.lastLevel = NormalLevel
	n.lastTitle = title
	n.lastMessage = message
	return nil
}

func (n fakeNotifier) highLevelNotification(title string, message string) error {
	n.lastLevel = HighLevel
	n.lastTitle = title
	n.lastMessage = message
	return nil
}

// TODO Figure out a way to test notifications without GUI elements displayed

func Test_show_notification(t *testing.T) {
	var testFlags = []struct {
		desc    string
		level   NotificationLevel
		title   string
		message string
	}{
		{
			desc:    "Normal Level",
			level:   NormalLevel,
			title:   "some normal level title",
			message: "some normal level message",
		},
		//{
		//	desc:    "High Level",
		//	level:   HighLevel,
		//	title:   "some high level title",
		//	message: "some high level message",
		//},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			fake := fakeNotifier{}
			notifier = fake
			ShowNotification(tt.level, tt.title, tt.message)
			assert.Equal(t, tt.level, fake.lastLevel)
			assert.Equal(t, tt.title, fake.lastTitle)
			assert.Equal(t, tt.message, fake.lastMessage)
		})
	}
}

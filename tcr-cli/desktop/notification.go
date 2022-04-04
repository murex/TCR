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
	"github.com/gen2brain/beeep"
	"github.com/murex/tcr/tcr-engine/report"
)

type (
	// NotificationLevel is the level of desktop notification messages. It can be either
	// normal or high
	NotificationLevel int

	// notifierInterface provides the interface to display desktop notifications
	notifierInterface interface {
		normalLevelNotification(title string, message string) error
		highLevelNotification(title string, message string) error
	}

	// beeepNotifier provides a wrapper around beeep 3rd-party module for desktop notifications
	beeepNotifier struct {
	}
)

// List of possible values for desktop notification level
const (
	NormalLevel NotificationLevel = iota
	HighLevel
)

var (
	// The notifier by default is beeep (3rd-party)
	notifier = newBeeepNotifier()
)

// newBeeepNotifier creates a beeep notifier instance
func newBeeepNotifier() notifierInterface {
	return beeepNotifier{}
}

// ShowNotification shows a notification message on the desktop. Implementation depends on the underlying OS.
func ShowNotification(level NotificationLevel, title string, message string) {
	var err error

	switch level {
	case NormalLevel:
		err = notifier.normalLevelNotification(title, message)
	case HighLevel:
		err = notifier.highLevelNotification(title, message)
	}

	if err != nil {
		report.PostError("ShowNotification:", err)
	}
}

func (b beeepNotifier) highLevelNotification(title string, message string) error {
	return beeep.Alert(title, message, "")
}

func (b beeepNotifier) normalLevelNotification(title string, message string) error {
	return beeep.Notify(title, message, "")
}

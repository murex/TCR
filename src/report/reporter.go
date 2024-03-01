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

package report

import (
	"fmt"
	"github.com/imkira/go-observer"
	"github.com/murex/tcr/report/role_event"
	"github.com/murex/tcr/report/timer_event"
	"github.com/murex/tcr/role"
	"sync"
	"time"
)

// Category provides different categories for reported message
type Category int

// List of possible values for Category field
const (
	Normal Category = iota
	Info
	Title
	Success
	Warning
	Error
	RoleEvent
	TimerEvent
)

// MessageReporter provides the interface that any message listener needs to implement
type MessageReporter interface {
	ReportSimple(emphasis bool, a ...any)
	ReportInfo(emphasis bool, a ...any)
	ReportTitle(emphasis bool, a ...any)
	ReportSuccess(emphasis bool, a ...any)
	ReportWarning(emphasis bool, a ...any)
	ReportError(emphasis bool, a ...any)
	ReportRoleEvent(emphasis bool, a ...any)
	ReportTimerEvent(emphasis bool, a ...any)
}

// MessageType type used for message characterization
type MessageType struct {
	Category Category
	Emphasis bool
}

// Message is the placeholder for any reported message
type Message struct {
	Type      MessageType
	Text      string
	Timestamp time.Time
}

var msgProperty observer.Property

func init() {
	Reset()
}

// Reset resets the reporter pipeline
func Reset() {
	msgProperty = observer.NewProperty(Message{Type: MessageType{Category: Normal}, Text: ""})
}

// Subscribe allows a listener to subscribe to any posted message through the reporter.
// The listener must implement the MessageReporter interface. The returned channel
// shall be kept by the listener as this channel will be used for unsubscription
func Subscribe(reporter MessageReporter) chan bool {
	stream := msgProperty.Observe()

	msg, _ := stream.Value().(Message) //nolint:revive
	// fmt.Printf("initial value: %v\n", msg)

	unsubscribe := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(s observer.Stream) {
		wg.Done()
		for {
			select {
			// wait for changes
			case <-s.Changes():
				// advance to next value
				s.Next()
				msg, _ = s.Value().(Message) //nolint:revive
				// fmt.Printf("got new value: %v\n", msg)
				reportMessage(reporter, msg)
			case <-unsubscribe:
				return
			}
		}
	}(stream)
	wg.Wait()
	return unsubscribe
}

// reportMessage tells the reporter to report msg depending on its category
func reportMessage(reporter MessageReporter, msg Message) {
	report := map[Category]func(r MessageReporter, emphasis bool, a ...any){
		Info:       MessageReporter.ReportInfo,
		Normal:     MessageReporter.ReportSimple,
		Title:      MessageReporter.ReportTitle,
		Success:    MessageReporter.ReportSuccess,
		Warning:    MessageReporter.ReportWarning,
		Error:      MessageReporter.ReportError,
		RoleEvent:  MessageReporter.ReportRoleEvent,
		TimerEvent: MessageReporter.ReportTimerEvent,
	}
	report[msg.Type.Category](reporter, msg.Type.Emphasis, msg.Text)
}

// Unsubscribe unsubscribes the listener associated to the provided channel from being notified
// of any further reported message
func Unsubscribe(c chan bool) {
	c <- true
}

// Post posts some text for reporting. This is actually the same as PostText()
func Post(a ...any) {
	PostText(a...)
}

// PostText posts some text for reporting
func PostText(a ...any) {
	postMessage(MessageType{Category: Normal}, a...)
}

// PostInfo posts an information message for reporting
func PostInfo(a ...any) {
	postMessage(MessageType{Category: Info}, a...)
}

// PostTitle posts a title message for reporting
func PostTitle(a ...any) {
	postMessage(MessageType{Category: Title}, a...)
}

// PostWarning posts a warning message for reporting
func PostWarning(a ...any) {
	postMessage(MessageType{Category: Warning}, a...)
}

// PostError posts an error message for reporting
func PostError(a ...any) {
	postMessage(MessageType{Category: Error}, a...)
}

// PostRoleEvent posts a role event
func PostRoleEvent(trigger string, r role.Role) {
	msg := role_event.Message{
		Trigger: role_event.Trigger(trigger),
		Role:    r,
	}
	postMessage(MessageType{Category: RoleEvent, Emphasis: msg.WithEmphasis()},
		role_event.WrapMessage(msg))
}

// PostTimerEvent posts a timer event
func PostTimerEvent(eventType string, timeout time.Duration, elapsed time.Duration, remaining time.Duration) {
	msg := timer_event.Message{
		Trigger:   timer_event.Trigger(eventType),
		Timeout:   timeout,
		Elapsed:   elapsed,
		Remaining: remaining,
	}
	postMessage(MessageType{Category: TimerEvent, Emphasis: msg.WithEmphasis()},
		timer_event.WrapMessage(msg))
}

// PostSuccessWithEmphasis posts a success message for reporting
func PostSuccessWithEmphasis(a ...any) {
	postMessage(MessageType{Category: Success, Emphasis: true}, a...)
}

// PostWarningWithEmphasis posts a warning with emphasis
func PostWarningWithEmphasis(a ...any) {
	postMessage(MessageType{Category: Warning, Emphasis: true}, a...)
}

// PostErrorWithEmphasis posts an error message for reporting
func PostErrorWithEmphasis(a ...any) {
	postMessage(MessageType{Category: Error, Emphasis: true}, a...)
}

func postMessage(msgType MessageType, a ...any) {
	msgProperty.Update(NewMessage(msgType, a...))
}

// NewMessage returns a message with the specified type
func NewMessage(messageType MessageType, a ...any) Message {
	return Message{messageType, fmt.Sprint(a...), time.Now()}
}

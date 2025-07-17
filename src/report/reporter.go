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
	"github.com/imkira/go-observer"
	"github.com/murex/tcr/report/role_event"
	"github.com/murex/tcr/report/text"
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

// MessageReporter provides the interface that any message listener should implement
type MessageReporter interface {
	ReportSimple(emphasis bool, payload text.Message)
	ReportInfo(emphasis bool, payload text.Message)
	ReportTitle(emphasis bool, payload text.Message)
	ReportSuccess(emphasis bool, payload text.Message)
	ReportWarning(emphasis bool, payload text.Message)
	ReportError(emphasis bool, payload text.Message)
	ReportRoleEvent(emphasis bool, payload role_event.Message)
	ReportTimerEvent(emphasis bool, payload timer_event.Message)
}

// MessageType contains message characterization information
type MessageType struct {
	Category Category
	Emphasis bool
}

// MessagePayload provides the abstraction for different types of message contents
type MessagePayload interface {
	ToString() string
}

// Message is the placeholder for any reported message
type Message struct {
	Type      MessageType
	Payload   MessagePayload
	Timestamp time.Time
}

// Reporter encapsulates the message pipeline for isolation
type Reporter struct {
	msgProperty observer.Property
}

var defaultReporter *Reporter

func init() {
	defaultReporter = NewReporter()
}

// NewReporter creates a new Reporter instance
func NewReporter() *Reporter {
	r := &Reporter{}
	r.Reset()
	return r
}

// Reset resets the reporter pipeline (for default reporter)
func Reset() {
	defaultReporter.Reset()
}

// Reset resets the reporter pipeline (for Reporter instance)
func (r *Reporter) Reset() {
	r.msgProperty = observer.NewProperty(Message{Type: MessageType{Category: Normal}, Payload: text.New("")})
}

// Subscribe allows a listener to subscribe to any posted message through the reporter (default)
func Subscribe(reporter MessageReporter) chan bool {
	return defaultReporter.Subscribe(reporter)
}

// Subscribe allows a listener to subscribe to any posted message through the reporter (instance)
func (r *Reporter) Subscribe(reporter MessageReporter) chan bool {
	stream := r.msgProperty.Observe()

	msg, _ := stream.Value().(Message) //nolint:revive

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
	switch msg.Type.Category {
	case Info:
		MessageReporter.ReportInfo(reporter, msg.Type.Emphasis, msg.Payload.(text.Message))
	case Normal:
		MessageReporter.ReportSimple(reporter, msg.Type.Emphasis, msg.Payload.(text.Message))
	case Title:
		MessageReporter.ReportTitle(reporter, msg.Type.Emphasis, msg.Payload.(text.Message))
	case Success:
		MessageReporter.ReportSuccess(reporter, msg.Type.Emphasis, msg.Payload.(text.Message))
	case Warning:
		MessageReporter.ReportWarning(reporter, msg.Type.Emphasis, msg.Payload.(text.Message))
	case Error:
		MessageReporter.ReportError(reporter, msg.Type.Emphasis, msg.Payload.(text.Message))
	case RoleEvent:
		MessageReporter.ReportRoleEvent(reporter, msg.Type.Emphasis, msg.Payload.(role_event.Message))
	case TimerEvent:
		MessageReporter.ReportTimerEvent(reporter, msg.Type.Emphasis, msg.Payload.(timer_event.Message))
	}
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
	postMessage(MessageType{Category: Normal}, text.New(a...))
}

// PostInfo posts an information message for reporting
func PostInfo(a ...any) {
	postMessage(MessageType{Category: Info}, text.New(a...))
}

// PostTitle posts a title message for reporting
func PostTitle(a ...any) {
	postMessage(MessageType{Category: Title}, text.New(a...))
}

// PostWarning posts a warning message for reporting
func PostWarning(a ...any) {
	postMessage(MessageType{Category: Warning}, text.New(a...))
}

// PostError posts an error message for reporting
func PostError(a ...any) {
	postMessage(MessageType{Category: Error}, text.New(a...))
}

// PostSuccessWithEmphasis posts a success message for reporting
func PostSuccessWithEmphasis(a ...any) {
	postMessage(MessageType{Category: Success, Emphasis: true}, text.New(a...))
}

// PostWarningWithEmphasis posts a warning with emphasis
func PostWarningWithEmphasis(a ...any) {
	postMessage(MessageType{Category: Warning, Emphasis: true}, text.New(a...))
}

// PostErrorWithEmphasis posts an error message for reporting
func PostErrorWithEmphasis(a ...any) {
	postMessage(MessageType{Category: Error, Emphasis: true}, text.New(a...))
}

// PostRoleEvent posts a role event
func PostRoleEvent(trigger role_event.Trigger, r role.Role) {
	msg := role_event.New(trigger, r)
	postMessage(MessageType{Category: RoleEvent, Emphasis: msg.WithEmphasis()}, msg)
}

// PostTimerEvent posts a timer event
func PostTimerEvent(
	trigger timer_event.Trigger,
	timeout time.Duration,
	elapsed time.Duration,
	remaining time.Duration,
) {
	msg := timer_event.New(trigger, timeout, elapsed, remaining)
	postMessage(MessageType{Category: TimerEvent, Emphasis: msg.WithEmphasis()}, msg)
}

func postMessage(msgType MessageType, payload MessagePayload) {
	defaultReporter.msgProperty.Update(NewMessage(msgType, payload))
}

// NewMessage builds a new reporter message
func NewMessage(messageType MessageType, messagePayload MessagePayload) Message {
	return Message{
		Type:      messageType,
		Payload:   messagePayload,
		Timestamp: time.Now(),
	}
}

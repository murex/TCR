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

package report

import (
	"fmt"
	"github.com/imkira/go-observer"
	"sync"
	"time"
)

// Severity provides level of severity for message
type Severity int

// List of possible values for Severity field
const (
	Normal Severity = iota
	Info
	Title
	Timer
	Success
	Warning
	Error
)

// MessageReporter provides the interface that any message listener needs to implement
type MessageReporter interface {
	ReportSimple(emphasis bool, a ...interface{})
	ReportInfo(emphasis bool, a ...interface{})
	ReportTitle(emphasis bool, a ...interface{})
	ReportTimer(emphasis bool, a ...interface{})
	ReportSuccess(emphasis bool, a ...interface{})
	ReportWarning(emphasis bool, a ...interface{})
	ReportError(emphasis bool, a ...interface{})
}

// MessageType type used for message characterization
type MessageType struct {
	Severity Severity
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
	msgProperty = observer.NewProperty(Message{Type: MessageType{Severity: Normal}, Text: ""})
}

// Subscribe allows a listener to subscribe to any posted message through the reporter.
// The listener must implement the MessageReporter interface. The returned channel
// shall be kept by the listener as this channel will be used for unsubscription
func Subscribe(reporter MessageReporter) chan bool {
	stream := msgProperty.Observe()

	msg := stream.Value().(Message)
	//fmt.Printf("initial value: %v\n", msg)

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
				msg = s.Value().(Message)
				//fmt.Printf("got new value: %v\n", msg)
				reportMessage(reporter, msg)
			case <-unsubscribe:
				return
			}
		}
	}(stream)
	wg.Wait()
	return unsubscribe
}

// reportMessage tells the reporter to report msg depending on its severity
func reportMessage(reporter MessageReporter, msg Message) {
	report := map[Severity]func(r MessageReporter, emphasis bool, a ...interface{}){
		Info:    MessageReporter.ReportInfo,
		Normal:  MessageReporter.ReportSimple,
		Title:   MessageReporter.ReportTitle,
		Timer:   MessageReporter.ReportTimer,
		Success: MessageReporter.ReportSuccess,
		Warning: MessageReporter.ReportWarning,
		Error:   MessageReporter.ReportError,
	}
	report[msg.Type.Severity](reporter, msg.Type.Emphasis, msg.Text)
}

// Unsubscribe unsubscribes the listener associated to the provided channel from being notified
// of any further reported message
func Unsubscribe(c chan bool) {
	c <- true
}

// Post posts some text for reporting. This is actually the same as PostText()
func Post(a ...interface{}) {
	PostText(a...)
}

// PostText posts some text for reporting
func PostText(a ...interface{}) {
	postMessage(MessageType{Severity: Normal}, a...)
}

// PostInfo posts an information message for reporting
func PostInfo(a ...interface{}) {
	postMessage(MessageType{Severity: Info}, a...)
}

// PostTitle posts a title message for reporting
func PostTitle(a ...interface{}) {
	postMessage(MessageType{Severity: Title}, a...)
}

// PostWarning posts a warning message for reporting
func PostWarning(a ...interface{}) {
	postMessage(MessageType{Severity: Warning}, a...)
}

// PostError posts an error message for reporting
func PostError(a ...interface{}) {
	postMessage(MessageType{Severity: Error}, a...)
}

// PostInfoWithEmphasis posts an info message with emphasis
func PostInfoWithEmphasis(a ...interface{}) {
	postMessage(MessageType{Severity: Info, Emphasis: true}, a...)
}

// PostTimerWithEmphasis posts a timer message with emphasis
func PostTimerWithEmphasis(a ...interface{}) {
	postMessage(MessageType{Severity: Timer, Emphasis: true}, a...)
}

// PostSuccessWithEmphasis posts a success message for reporting
func PostSuccessWithEmphasis(a ...interface{}) {
	postMessage(MessageType{Severity: Success, Emphasis: true}, a...)
}

// PostWarningWithEmphasis posts a warning with emphasis
func PostWarningWithEmphasis(a ...interface{}) {
	postMessage(MessageType{Severity: Warning, Emphasis: true}, a...)
}

// PostErrorWithEmphasis posts an error message for reporting
func PostErrorWithEmphasis(a ...interface{}) {
	postMessage(MessageType{Severity: Error, Emphasis: true}, a...)
}

func postMessage(msgType MessageType, a ...interface{}) {
	msgProperty.Update(NewMessage(msgType, a...))
}

// NewMessage returns a message with the specified type
func NewMessage(messageType MessageType, a ...interface{}) Message {
	return Message{messageType, fmt.Sprint(a...), time.Now()}
}

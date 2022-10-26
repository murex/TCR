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
	Warning
	Error
)

type MessageReport interface {
	Notification(msgText string)
}

// MessageType type used for message characterization
type MessageType struct {
	Severity Severity
	Emphasis bool

	Report func(reporter MessageReport)
}

func (mt MessageType) compare(x MessageType) bool {
	return mt.Severity == x.Severity && mt.Emphasis == x.Emphasis
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
// onReport() will be called every time a new message is posted. The returned channel
// shall be kept by the listener as this channel will be used for unsubscription
func Subscribe(onReport func(msg Message)) chan bool {
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
				onReport(msg)
			case <-unsubscribe:
				return
			}
		}
	}(stream)
	wg.Wait()
	return unsubscribe
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

// PostNotification posts an event message for reporting
func PostNotification(a ...interface{}) {
	postMessage(MessageType{Severity: Info, Emphasis: true}, a...)
}

func postMessage(msgType MessageType, a ...interface{}) {
	message := Message{msgType, fmt.Sprint(a...), time.Now()}
	//fmt.Println("Reporting message:", message)
	msgProperty.Update(message)
}

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

package report

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_can_retrieve_reported_message(t *testing.T) {
	text := "Dummy Message"
	result := reportAndReceive(func() {
		Post(text)
	})
	assert.Equal(t, text, result.Text)
}

func Test_one_message_and_multiple_receivers(t *testing.T) {
	const nbListeners = 2
	text := "Dummy Message"
	var c [nbListeners]chan bool
	var stubs [nbListeners]MessageReporterStub

	for i := 0; i < nbListeners; i++ {
		go func(i int) {
			stubs[i] = NewMessageReporterStub(i)
			c[i] = Subscribe(&stubs[i])
		}(i)
	}

	// To make sure observers are ready to receive
	time.Sleep(1 * time.Millisecond)
	Post(text)

	for i := 0; i < nbListeners; i++ {
		iReceived := <-stubs[i].received
		Unsubscribe(c[iReceived])
		assert.Equal(t, text, stubs[iReceived].message.Text)
	}
}

func Test_multiple_messages_and_one_receiver(t *testing.T) {
	const nbMessages = 3

	stub := NewMessageReporterStub(0)
	c := Subscribe(&stub)

	// To make sure the observer is ready to receive
	time.Sleep(1 * time.Millisecond)
	for i := 0; i < nbMessages; i++ {
		text := fmt.Sprintf("Dummy Message %v", i)
		Post(text)
		<-stub.received
		assert.Equal(t, text, stub.message.Text)
	}
	Unsubscribe(c)
}

func Test_report_simple_message(t *testing.T) {
	text := "Normal Message"
	result := reportAndReceive(func() {
		PostText(text)
	})
	assertMessageMatch(t, text, MessageType{Severity: Normal}, result)
}

func Test_report_info_message(t *testing.T) {
	text := "info Message"
	result := reportAndReceive(func() {
		PostInfo(text)
	})
	assertMessageMatch(t, text, MessageType{Severity: Info}, result)
}

func Test_report_title_message(t *testing.T) {
	text := "Title Message"
	result := reportAndReceive(func() {
		PostTitle(text)
	})
	assertMessageMatch(t, text, MessageType{Severity: Title}, result)
}

func Test_report_warning_message(t *testing.T) {
	text := "Warning Message"
	result := reportAndReceive(func() {
		PostWarning(text)
	})
	assertMessageMatch(t, text, MessageType{Severity: Warning}, result)
}

func Test_report_error_message(t *testing.T) {
	text := "Error Message"
	result := reportAndReceive(func() {
		PostError(text)
	})
	assertMessageMatch(t, text, MessageType{Severity: Error}, result)
}

func Test_report_info_with_emphasis_message(t *testing.T) {
	text := "PostInfoWithEmphasis Message"
	result := reportAndReceive(func() {
		PostInfoWithEmphasis(text)
	})
	assertMessageMatch(t, text, MessageType{Severity: Info, Emphasis: true}, result)
}

func Test_report_warning_with_emphasis(t *testing.T) {
	text := "PostWarningWithEmphasis Message"
	result := reportAndReceive(func() {
		PostWarningWithEmphasis(text)
	})
	assertMessageMatch(t, text, MessageType{Severity: Warning, Emphasis: true}, result)
}

func assertMessageMatch(t *testing.T, text string, msgType MessageType, msg Message) {
	t.Helper()
	assert.Equal(t, text, msg.Text)
	assert.Equal(t, msgType, msg.Type)
	assert.NotZero(t, msg.Timestamp)
}

func reportAndReceive(report func()) Message {
	stub := NewMessageReporterStub(0)
	c := Subscribe(&stub)

	// To make sure the observer is ready to receive
	time.Sleep(1 * time.Millisecond)
	report()
	<-stub.received
	Unsubscribe(c)
	return stub.message
}

type MessageReporterStub struct {
	index    int
	received chan int
	message  Message
}

func NewMessageReporterStub(index int) MessageReporterStub {
	stub := MessageReporterStub{}
	stub.index = index
	stub.received = make(chan int)
	return stub
}

// ReportSimple reports simple messages
func (stub *MessageReporterStub) ReportSimple(emphasis bool, a ...interface{}) {
	stub.message = NewMessage(MessageType{Normal, emphasis}, a...)
	stub.received <- stub.index
}

// ReportInfo reports info messages
func (stub *MessageReporterStub) ReportInfo(emphasis bool, a ...interface{}) {
	stub.message = NewMessage(MessageType{Info, emphasis}, a...)
	stub.received <- stub.index
}

// ReportTitle reports title messages
func (stub *MessageReporterStub) ReportTitle(emphasis bool, a ...interface{}) {
	stub.message = NewMessage(MessageType{Title, emphasis}, a...)
	stub.received <- stub.index
}

// ReportWarning reports warning messages
func (stub *MessageReporterStub) ReportWarning(emphasis bool, a ...interface{}) {
	stub.message = NewMessage(MessageType{Warning, emphasis}, a...)
	stub.received <- stub.index
}

// ReportError reports error messages
func (stub *MessageReporterStub) ReportError(emphasis bool, a ...interface{}) {
	stub.message = NewMessage(MessageType{Error, emphasis}, a...)
	stub.received <- stub.index
}

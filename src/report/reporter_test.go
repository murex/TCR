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
	"github.com/murex/tcr/report/role_event"
	"github.com/murex/tcr/report/timer_event"
	"github.com/murex/tcr/role"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_can_retrieve_reported_message(t *testing.T) {
	text := "dummy message"
	result := reportAndReceive(func() {
		Post(text)
	})
	assert.Equal(t, text, result.Text)
}

func Test_one_message_and_multiple_receivers(t *testing.T) {
	const nbListeners = 2
	text := "dummy message"
	var c [nbListeners]chan bool
	var stubs [nbListeners]*messageReporterStub

	for i := 0; i < nbListeners; i++ {
		go func(i int) {
			stubs[i] = newMessageReporterStub(i)
			c[i] = Subscribe(stubs[i])
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

	stub := newMessageReporterStub(0)
	c := Subscribe(stub)

	// To make sure the observer is ready to receive
	time.Sleep(1 * time.Millisecond)
	for i := 0; i < nbMessages; i++ {
		text := fmt.Sprintf("dummy message %v", i)
		Post(text)
		<-stub.received
		assert.Equal(t, text, stub.message.Text)
	}
	Unsubscribe(c)
}

func Test_post_text_message_functions(t *testing.T) {
	testCases := []struct {
		text         string
		postFunction func(a ...any)
		expectedType MessageType
	}{
		{
			"normal message",
			PostText,
			MessageType{Normal, false},
		},
		{
			"info message",
			PostInfo,
			MessageType{Info, false},
		},
		{
			"title message",
			PostTitle,
			MessageType{Title, false},
		},
		{
			"warning message",
			PostWarning,
			MessageType{Warning, false},
		},
		{
			"error message",
			PostError,
			MessageType{Error, false},
		},
		{
			"success message with emphasis",
			PostSuccessWithEmphasis,
			MessageType{Success, true},
		},
		{
			"warning message with emphasis",
			PostWarningWithEmphasis,
			MessageType{Warning, true},
		},
		{
			"error message with emphasis",
			PostErrorWithEmphasis,
			MessageType{Error, true},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.text, func(t *testing.T) {
			result := reportAndReceive(func() {
				tt.postFunction(tt.text)
			})
			assert.Equal(t, tt.text, result.Text)
			assert.Equal(t, tt.expectedType, result.Type)
			assert.NotZero(t, result.Timestamp)
		})
	}
}

func Test_post_event_message_functions(t *testing.T) {
	testCases := []struct {
		text         string
		postFunction func()
		expectedType MessageType
		expectedText string
	}{
		{
			"role event message",
			func() {
				PostRoleEvent(string(role_event.TriggerStart), role.Navigator{})
			},
			MessageType{RoleEvent, false},
			"navigator:start",
		},
		{
			"timer event message",
			func() {
				PostTimerEvent(string(timer_event.TriggerStart), 0, 0, 0)
			},
			MessageType{TimerEvent, true},
			"start:0:0:0",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.text, func(t *testing.T) {
			result := reportAndReceive(func() {
				tt.postFunction()
			})
			assert.Equal(t, tt.expectedText, result.Text)
			assert.Equal(t, tt.expectedType, result.Type)
			assert.NotZero(t, result.Timestamp)
		})
	}
}

func reportAndReceive(report func()) Message {
	stub := newMessageReporterStub(0)
	c := Subscribe(stub)

	// To make sure the observer is ready to receive
	time.Sleep(1 * time.Millisecond)
	report()
	<-stub.received
	Unsubscribe(c)
	return stub.message
}

type messageReporterStub struct {
	index    int
	received chan int
	message  Message
}

func newMessageReporterStub(index int) *messageReporterStub {
	return &messageReporterStub{
		index:    index,
		received: make(chan int),
	}
}

func (stub *messageReporterStub) report(category Category, emphasis bool, a ...any) {
	stub.message = NewMessage(MessageType{category, emphasis}, a...)
	stub.received <- stub.index
}

// ReportSimple reports simple messages
func (stub *messageReporterStub) ReportSimple(emphasis bool, a ...any) {
	stub.report(Normal, emphasis, a...)
}

// ReportInfo reports info messages
func (stub *messageReporterStub) ReportInfo(emphasis bool, a ...any) {
	stub.report(Info, emphasis, a...)
}

// ReportTitle reports title messages
func (stub *messageReporterStub) ReportTitle(emphasis bool, a ...any) {
	stub.report(Title, emphasis, a...)
}

// ReportSuccess reports success messages
func (stub *messageReporterStub) ReportSuccess(emphasis bool, a ...any) {
	stub.report(Success, emphasis, a...)
}

// ReportWarning reports warning messages
func (stub *messageReporterStub) ReportWarning(emphasis bool, a ...any) {
	stub.report(Warning, emphasis, a...)
}

// ReportError reports error messages
func (stub *messageReporterStub) ReportError(emphasis bool, a ...any) {
	stub.report(Error, emphasis, a...)
}

// ReportTimerEvent reports role event messages
func (stub *messageReporterStub) ReportRoleEvent(emphasis bool, a ...any) {
	stub.report(RoleEvent, emphasis, a...)
}

// ReportTimerEvent reports timer event messages
func (stub *messageReporterStub) ReportTimerEvent(emphasis bool, a ...any) {
	stub.report(TimerEvent, emphasis, a...)
}

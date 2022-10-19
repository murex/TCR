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
	var result [nbListeners]Message
	var c [nbListeners]chan bool

	received := make(chan int, nbListeners)

	for i := 0; i < nbListeners; i++ {
		go func(i int) {
			c[i] = Subscribe(func(msg Message) {
				result[i] = msg
				received <- i
			})
		}(i)
	}

	// To make sure observers are ready to receive
	time.Sleep(1 * time.Millisecond)
	Post(text)

	for i := 0; i < nbListeners; i++ {
		iReceived := <-received
		Unsubscribe(c[iReceived])
		assert.Equal(t, text, result[iReceived].Text)
	}
}

func Test_multiple_messages_and_one_receiver(t *testing.T) {
	const nbMessages = 2
	received := make(chan Message)

	c := Subscribe(func(msg Message) {
		received <- msg
	})

	// To make sure the observer is ready to receive
	time.Sleep(1 * time.Millisecond)
	for i := 0; i < nbMessages; i++ {
		text := fmt.Sprintf("Dummy Message %v", i)
		Post(text)
		result := <-received
		assert.Equal(t, text, result.Text)
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

func Test_report_notification_message(t *testing.T) {
	text := "Notification Message"
	result := reportAndReceive(func() {
		PostNotification(text)
	})
	assertMessageMatch(t, text, MessageType{Severity: Notification, Emphasis: true}, result)
}

func assertMessageMatch(t *testing.T, text string, msgType MessageType, msg Message) {
	assert.Equal(t, text, msg.Text)
	assert.Equal(t, msgType, msg.Type)
	assert.NotZero(t, msg.Timestamp)
}

func reportAndReceive(report func()) Message {
	var result Message
	received := make(chan bool)

	c := Subscribe(func(msg Message) {
		result = msg
		received <- true
	})

	// To make sure the observer is ready to receive
	time.Sleep(1 * time.Millisecond)
	report()
	<-received
	Unsubscribe(c)
	return result
}

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
	"sync"
	"testing"
	"time"

	"github.com/murex/tcr/report/role_event"
	"github.com/murex/tcr/report/text"
	"github.com/murex/tcr/report/timer_event"
	"github.com/murex/tcr/role"
	"github.com/stretchr/testify/assert"
)

func Test_can_retrieve_reported_message(t *testing.T) {
	txt := "dummy message"
	TestWithIsolatedReporter(func(reporter *Reporter, sniffer *Sniffer) {
		reporter.Post(txt)
		// Give a small moment for the message to be processed
		time.Sleep(10 * time.Millisecond)
		sniffer.Stop()
		assert.Equal(t, 1, sniffer.GetMatchCount())
		assert.Equal(t, txt, sniffer.GetAllMatches()[0].Payload.ToString())
	})
}

func Test_one_message_and_multiple_receivers(t *testing.T) {
	const nbListeners = 2
	txt := "dummy message"

	TestWithIsolatedReporter(func(reporter *Reporter, sniffer *Sniffer) {
		var c [nbListeners]chan bool
		var stubs [nbListeners]*messageReporterStub

		// Create stubs and subscribe synchronously
		for i := 0; i < nbListeners; i++ {
			stubs[i] = newMessageReporterStub(i)
			c[i] = reporter.Subscribe(stubs[i])
		}

		// Post the message
		reporter.Post(txt)

		// Wait for responses and verify with better synchronization
		var wg sync.WaitGroup
		wg.Add(nbListeners)

		for i := range nbListeners {
			go func(idx int) {
				defer wg.Done()
				select {
				case iReceived := <-stubs[idx].received:
					assert.Equal(t, idx, iReceived)
					stubs[idx].mutex.RLock()
					msgPayload := stubs[idx].message.Payload.ToString()
					stubs[idx].mutex.RUnlock()
					assert.Equal(t, txt, msgPayload)
				case <-time.After(100 * time.Millisecond):
					t.Errorf("timeout waiting for receiver %d", idx)
				}
			}(i)
		}

		wg.Wait()

		// Cleanup - use goroutines to avoid blocking
		for i := 0; i < nbListeners; i++ {
			go func(ch chan bool) {
				select {
				case ch <- true:
				case <-time.After(10 * time.Millisecond):
					// Channel might be blocked, that's ok
				}
			}(c[i])
		}

		// Small delay to allow cleanup
		time.Sleep(10 * time.Millisecond)
	})
}

func Test_multiple_messages_and_one_receiver(t *testing.T) {
	const nbMessages = 3

	TestWithIsolatedReporter(func(reporter *Reporter, sniffer *Sniffer) {
		stub := newMessageReporterStub(0)
		c := reporter.Subscribe(stub)

		for i := range nbMessages {
			txt := fmt.Sprintf("dummy message %v", i)
			reporter.Post(txt)
			select {
			case <-stub.received:
				stub.mutex.RLock()
				msgPayload := stub.message.Payload.ToString()
				stub.mutex.RUnlock()
				assert.Equal(t, txt, msgPayload)
			case <-time.After(100 * time.Millisecond):
				t.Fatalf("timeout waiting for message %d", i)
			}
		}

		// Cleanup using goroutine to avoid blocking
		go func() {
			select {
			case c <- true:
			case <-time.After(10 * time.Millisecond):
				// Channel might be blocked, that's ok
			}
		}()

		// Small delay to allow cleanup
		time.Sleep(10 * time.Millisecond)
	})
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
			TestWithIsolatedReporter(func(reporter *Reporter, sniffer *Sniffer) {
				// Use the isolated reporter instance instead of global functions
				switch tt.expectedType.Category {
				case Normal:
					reporter.PostText(tt.text)
				case Info:
					reporter.PostInfo(tt.text)
				case Title:
					reporter.PostTitle(tt.text)
				case Warning:
					if tt.expectedType.Emphasis {
						reporter.PostWarningWithEmphasis(tt.text)
					} else {
						reporter.PostWarning(tt.text)
					}
				case Error:
					if tt.expectedType.Emphasis {
						reporter.PostErrorWithEmphasis(tt.text)
					} else {
						reporter.PostError(tt.text)
					}
				case Success:
					if tt.expectedType.Emphasis {
						reporter.PostSuccessWithEmphasis(tt.text)
					}
				}

				// Give time for message processing
				time.Sleep(10 * time.Millisecond)
				sniffer.Stop()
				assert.Equal(t, 1, sniffer.GetMatchCount())
				result := sniffer.GetAllMatches()[0]
				assert.Equal(t, text.New(tt.text), result.Payload)
				assert.Equal(t, tt.expectedType, result.Type)
				assert.NotZero(t, result.Timestamp)
			})
		})
	}
}

func Test_post_event_message_functions(t *testing.T) {
	testCases := []struct {
		text            string
		postFunction    func()
		expectedType    MessageType
		expectedPayload MessagePayload
	}{
		{
			"role event message",
			func() {
				PostRoleEvent(role_event.TriggerStart, role.Navigator{})
			},
			MessageType{RoleEvent, false},
			role_event.Message{
				Trigger: role_event.TriggerStart,
				Role:    role.Navigator{},
			},
		},
		{
			"timer event message",
			func() {
				PostTimerEvent(timer_event.TriggerCountdown, 3*time.Second, 2*time.Second, 1*time.Second)
			},
			MessageType{TimerEvent, true},
			timer_event.Message{
				Trigger:   timer_event.TriggerCountdown,
				Timeout:   3 * time.Second,
				Elapsed:   2 * time.Second,
				Remaining: 1 * time.Second,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.text, func(t *testing.T) {
			TestWithIsolatedReporter(func(reporter *Reporter, sniffer *Sniffer) {
				// Use the isolated reporter instance for event functions
				switch tt.text {
				case "role event message":
					reporter.PostRoleEvent(role_event.TriggerStart, role.Navigator{})
				case "timer event message":
					reporter.PostTimerEvent(timer_event.TriggerCountdown, 3*time.Second, 2*time.Second, 1*time.Second)
				}

				// Give time for message processing
				time.Sleep(10 * time.Millisecond)
				sniffer.Stop()
				assert.Equal(t, 1, sniffer.GetMatchCount())
				result := sniffer.GetAllMatches()[0]
				assert.Equal(t, tt.expectedPayload, result.Payload)
				assert.Equal(t, tt.expectedType, result.Type)
				assert.NotZero(t, result.Timestamp)
			})
		})
	}
}

type messageReporterStub struct {
	index    int
	received chan int
	message  Message
	mutex    sync.RWMutex
}

func newMessageReporterStub(index int) *messageReporterStub {
	return &messageReporterStub{
		index:    index,
		received: make(chan int),
	}
}

func (stub *messageReporterStub) report(category Category, emphasis bool, payload MessagePayload) {
	stub.mutex.Lock()
	stub.message = NewMessage(MessageType{category, emphasis}, payload)
	stub.mutex.Unlock()
	stub.received <- stub.index
}

// ReportSimple reports simple messages
func (stub *messageReporterStub) ReportSimple(emphasis bool, payload text.Message) {
	stub.report(Normal, emphasis, payload)
}

// ReportInfo reports info messages
func (stub *messageReporterStub) ReportInfo(emphasis bool, payload text.Message) {
	stub.report(Info, emphasis, payload)
}

// ReportTitle reports title messages
func (stub *messageReporterStub) ReportTitle(emphasis bool, payload text.Message) {
	stub.report(Title, emphasis, payload)
}

// ReportSuccess reports success messages
func (stub *messageReporterStub) ReportSuccess(emphasis bool, payload text.Message) {
	stub.report(Success, emphasis, payload)
}

// ReportWarning reports warning messages
func (stub *messageReporterStub) ReportWarning(emphasis bool, payload text.Message) {
	stub.report(Warning, emphasis, payload)
}

// ReportError reports error messages
func (stub *messageReporterStub) ReportError(emphasis bool, payload text.Message) {
	stub.report(Error, emphasis, payload)
}

// ReportTimerEvent reports role event messages
func (stub *messageReporterStub) ReportRoleEvent(emphasis bool, payload role_event.Message) {
	stub.report(RoleEvent, emphasis, payload)
}

// ReportTimerEvent reports timer event messages
func (stub *messageReporterStub) ReportTimerEvent(emphasis bool, payload timer_event.Message) {
	stub.report(TimerEvent, emphasis, payload)
}

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

package ws

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/report/role_event"
	"github.com/murex/tcr/report/timer_event"
	"github.com/murex/tcr/role"
	"github.com/stretchr/testify/assert"
)

type fakeHTTPServer struct {
	url url.URL
}

func newFakeHTTPServer(urlStr string) *fakeHTTPServer {
	u, _ := url.Parse(urlStr)
	return &fakeHTTPServer{url: *u}
}

// InDevMode indicates if the server is running in dev (development) mode
func (s *fakeHTTPServer) InDevMode() bool {
	return false
}

// GetServerAddress returns the TCP server address that the server is listening to.
func (s *fakeHTTPServer) GetServerAddress() string {
	return s.url.Host
}

// GetWebsocketTimeout returns the timeout after which inactive websocket connections
// should be closed
func (s *fakeHTTPServer) GetWebsocketTimeout() time.Duration {
	// We differentiate CI and local machine to optimize test speed execution when run on local machine
	// while not failing when run on CI (which runs slower)
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		fmt.Println("Using long websocket timeouts")
		return 1 * time.Second
	} else {
		fmt.Println("Using short websocket timeouts")
		return 100 * time.Millisecond
	}
}

func Test_websocket_report_messages(t *testing.T) {
	const messageText = "hello from TCR!"
	tests := []struct {
		desc     string
		action   func(reporter *report.Reporter)
		expected message
	}{
		{
			desc:     "report.Post",
			action:   func(reporter *report.Reporter) { reporter.Post(messageText) },
			expected: newMessage(messageTypeSimple, messageSeverityNormal, false, messageText),
		},
		{
			desc:     "report.PostText",
			action:   func(reporter *report.Reporter) { reporter.PostText(messageText) },
			expected: newMessage(messageTypeSimple, messageSeverityNormal, false, messageText),
		},
		{
			desc:     "report.PostInfo",
			action:   func(reporter *report.Reporter) { reporter.PostInfo(messageText) },
			expected: newMessage(messageTypeInfo, messageSeverityNormal, false, messageText),
		},
		{
			desc:     "report.PostTitle",
			action:   func(reporter *report.Reporter) { reporter.PostTitle(messageText) },
			expected: newMessage(messageTypeTitle, messageSeverityNormal, false, messageText),
		},
		{
			desc:     "report.PostWarning",
			action:   func(reporter *report.Reporter) { reporter.PostWarning(messageText) },
			expected: newMessage(messageTypeWarning, messageSeverityLow, false, messageText),
		},
		{
			desc:     "report.PostError",
			action:   func(reporter *report.Reporter) { reporter.PostError(messageText) },
			expected: newMessage(messageTypeError, messageSeverityHigh, false, messageText),
		},
		{
			desc:     "report.PostRoleEvent navigator start",
			action:   func(reporter *report.Reporter) { reporter.PostRoleEvent(role_event.TriggerStart, role.Navigator{}) },
			expected: newMessage(messageTypeRole, messageSeverityNormal, false, "navigator:start"),
		},
		{
			desc:     "report.PostRoleEvent driver start",
			action:   func(reporter *report.Reporter) { reporter.PostRoleEvent(role_event.TriggerStart, role.Driver{}) },
			expected: newMessage(messageTypeRole, messageSeverityNormal, false, "driver:start"),
		},
		{
			desc:     "report.PostRoleEvent navigator end",
			action:   func(reporter *report.Reporter) { reporter.PostRoleEvent(role_event.TriggerEnd, role.Navigator{}) },
			expected: newMessage(messageTypeRole, messageSeverityNormal, false, "navigator:end"),
		},
		{
			desc:     "report.PostRoleEvent driver end",
			action:   func(reporter *report.Reporter) { reporter.PostRoleEvent(role_event.TriggerEnd, role.Driver{}) },
			expected: newMessage(messageTypeRole, messageSeverityNormal, false, "driver:end"),
		},
		{
			desc:     "report.PostTimerEvent start",
			action:   func(reporter *report.Reporter) { reporter.PostTimerEvent(timer_event.TriggerStart, 0, 0, 0) },
			expected: newMessage(messageTypeTimer, messageSeverityNormal, true, "start:0:0:0"),
		},
		{
			desc:     "report.PostTimerEvent countdown",
			action:   func(reporter *report.Reporter) { reporter.PostTimerEvent(timer_event.TriggerCountdown, 0, 0, 0) },
			expected: newMessage(messageTypeTimer, messageSeverityNormal, true, "countdown:0:0:0"),
		},
		{
			desc:     "report.PostTimerEvent stop",
			action:   func(reporter *report.Reporter) { reporter.PostTimerEvent(timer_event.TriggerStop, 0, 0, 0) },
			expected: newMessage(messageTypeTimer, messageSeverityNormal, true, "stop:0:0:0"),
		},
		{
			desc:     "report.PostTimerEvent first timeout",
			action:   func(reporter *report.Reporter) { reporter.PostTimerEvent(timer_event.TriggerTimeout, 0, 0, 0) },
			expected: newMessage(messageTypeTimer, messageSeverityNormal, true, "timeout:0:0:0"),
		},
		{
			desc: "report.PostTimerEvent second timeout",
			action: func(reporter *report.Reporter) {
				reporter.PostTimerEvent(timer_event.TriggerTimeout, 0, 0, -1*time.Second)
			},
			expected: newMessage(messageTypeTimer, messageSeverityNormal, false, "timeout:0:0:-1"),
		},
		{
			desc:     "report.PostSuccessWithEmphasis",
			action:   func(reporter *report.Reporter) { reporter.PostSuccessWithEmphasis(messageText) },
			expected: newMessage(messageTypeSuccess, messageSeverityNormal, true, messageText),
		},
		{
			desc:     "report.PostWarningWithEmphasis",
			action:   func(reporter *report.Reporter) { reporter.PostWarningWithEmphasis(messageText) },
			expected: newMessage(messageTypeWarning, messageSeverityLow, true, messageText),
		},
		{
			desc:     "report.PostErrorWithEmphasis",
			action:   func(reporter *report.Reporter) { reporter.PostErrorWithEmphasis(messageText) },
			expected: newMessage(messageTypeError, messageSeverityHigh, true, messageText),
		},
	}

	// Create HTTP test server with the websocket connection handler.
	s := httptest.NewUnstartedServer(http.HandlerFunc(websocketConnectionHandler))
	var fakeServer tcrHTTPServer
	s.Config.BaseContext = func(l net.Listener) context.Context {
		fakeServer = newFakeHTTPServer(s.URL)
		return context.WithValue(context.Background(), serverContextKey, fakeServer)
	}
	s.Start()
	defer s.Close()

	// Build URL and header for websocket connection request
	u, _ := url.Parse(s.URL)
	u.Scheme = "ws"
	hd := http.Header{}
	hd.Add("Origin", s.URL)
	// Create the websocket connection
	var ws, _, err = websocket.DefaultDialer.Dial(u.String(), hd) //nolint:bodyclose
	defer func(ws *websocket.Conn) {
		_ = ws.Close()
	}(ws)
	if err != nil {
		t.Fatalf("%v", err)
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			// Run the test action (message posting through report)
			test.action(report.GetDefaultReporter())

			// Retrieve the message sent through the websocket and verify its contents
			var msg message

			// Set read deadline to prevent hanging
			_ = ws.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

			readErr := ws.ReadJSON(&msg)
			if readErr != nil {
				// Check if it's a close error, which is expected during shutdown
				if websocket.IsCloseError(readErr, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
					t.Logf("WebSocket closed during test: %v", readErr)
					return
				}
			}
			assert.NoError(t, readErr)
			assertMessagesMatch(t, test.expected, msg)
		})
	}

	// Send close frame before waiting for timeout
	_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	// Wait for the websocket connection to time out with some extra buffer
	time.Sleep(fakeServer.GetWebsocketTimeout() + 100*time.Millisecond)
}

func Test_websocket_upgrader_with_invalid_request_header(t *testing.T) {
	tests := []struct {
		desc      string
		hdBuilder func() http.Header
	}{
		{
			desc: "origin with no protocol",
			hdBuilder: func() http.Header {
				hd := http.Header{}
				hd.Add("Origin", "://127.0.0.1")
				return hd
			},
		},
		{
			desc: "origin with invalid hostname",
			hdBuilder: func() http.Header {
				hd := http.Header{}
				hd.Add("Origin", "http://dummy.url")
				return hd
			},
		},
		{
			desc: "origin with invalid port",
			hdBuilder: func() http.Header {
				hd := http.Header{}
				hd.Add("Origin", "http://127.0.0.1:9999")
				return hd
			},
		},
		{
			desc: "origin not set",
			hdBuilder: func() http.Header {
				return http.Header{}
			},
		},
	}

	// Create HTTP test server with the websocket connection handler.
	s := httptest.NewUnstartedServer(http.HandlerFunc(websocketConnectionHandler))
	var fakeServer tcrHTTPServer
	s.Config.BaseContext = func(l net.Listener) context.Context {
		fakeServer = newFakeHTTPServer(s.URL)
		return context.WithValue(context.Background(), serverContextKey, fakeServer)
	}
	s.Start()
	defer s.Close()

	u, _ := url.Parse(s.URL)
	u.Scheme = "ws"
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			hd := test.hdBuilder()
			// Try creating the websocket connection using this header
			_, _, err := websocket.DefaultDialer.Dial(u.String(), hd) //nolint:bodyclose
			// And verify that it gets rejected
			assert.Error(t, err)
		})
	}
}

// assertMessagesMatch checks that 2 message instance messages match.
// Used in place of assert.Equal() to ignore potential timestamp variations.
func assertMessagesMatch(t *testing.T, expected message, msg message) {
	t.Helper()
	assert.Equal(t, expected.Type, msg.Type)
	assert.Equal(t, expected.Severity, msg.Severity)
	assert.Equal(t, expected.Emphasis, msg.Emphasis)
	assert.Equal(t, expected.Text, msg.Text)
	expectedTS, _ := time.Parse(time.RFC3339, expected.Timestamp)
	msgTS, _ := time.Parse(time.RFC3339, msg.Timestamp)
	assert.WithinDuration(t, expectedTS, msgTS, 10*time.Second)
}

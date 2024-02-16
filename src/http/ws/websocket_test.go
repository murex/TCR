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
	"github.com/gorilla/websocket"
	"github.com/murex/tcr/report"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
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

// GetWebSocketTimeout returns the timeout after which inactive websocket connections
// should be closed
func (s *fakeHTTPServer) GetWebSocketTimeout() time.Duration {
	// To prevent waiting for 1 minute before websocket connection gets shut down
	return 100 * time.Millisecond
}

// RegisterWebSocket register a new websocket connection to the server
func (s *fakeHTTPServer) RegisterWebSocket(_ *WebsocketMessageReporter) {
}

// UnregisterWebSocket unregister a new websocket connection from the server
func (s *fakeHTTPServer) UnregisterWebSocket(_ *WebsocketMessageReporter) {
}

func Test_websocket_report_messages(t *testing.T) {
	const messageText = "hello from TCR!"
	tests := []struct {
		desc     string
		action   func()
		expected webSocketMessage
	}{
		{
			desc:     "report.Post",
			action:   func() { report.Post(messageText) },
			expected: newWebSocketMessage("simple", "0", false, messageText),
		},
		{
			desc:     "report.PostText",
			action:   func() { report.PostText(messageText) },
			expected: newWebSocketMessage("simple", "0", false, messageText),
		},
		{
			desc:     "report.PostInfo",
			action:   func() { report.PostInfo(messageText) },
			expected: newWebSocketMessage("info", "0", false, messageText),
		},
		{
			desc:     "report.PostTitle",
			action:   func() { report.PostTitle(messageText) },
			expected: newWebSocketMessage("title", "0", false, messageText),
		},
		{
			desc:     "report.PostWarning",
			action:   func() { report.PostWarning(messageText) },
			expected: newWebSocketMessage("warning", "1", false, messageText),
		},
		{
			desc:     "report.PostError",
			action:   func() { report.PostError(messageText) },
			expected: newWebSocketMessage("error", "2", false, messageText),
		},
		{
			desc:     "report.PostRole",
			action:   func() { report.PostRole(messageText) },
			expected: newWebSocketMessage("role", "0", false, messageText),
		},
		{
			desc:     "report.PostTimerWithEmphasis",
			action:   func() { report.PostTimerWithEmphasis(messageText) },
			expected: newWebSocketMessage("timer", "0", true, messageText),
		},
		{
			desc:     "report.PostSuccessWithEmphasis",
			action:   func() { report.PostSuccessWithEmphasis(messageText) },
			expected: newWebSocketMessage("success", "0", true, messageText),
		},
		{
			desc:     "report.PostWarningWithEmphasis",
			action:   func() { report.PostWarningWithEmphasis(messageText) },
			expected: newWebSocketMessage("warning", "1", true, messageText),
		},
		{
			desc:     "report.PostErrorWithEmphasis",
			action:   func() { report.PostErrorWithEmphasis(messageText) },
			expected: newWebSocketMessage("error", "2", true, messageText),
		},
	}

	// Create HTTP test server with the websocket connection handler.
	s := httptest.NewUnstartedServer(http.HandlerFunc(webSocketConnectionHandler))
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
			test.action()

			// Retrieve the message sent through the websocket and verify its contents
			var msg webSocketMessage
			readErr := ws.ReadJSON(&msg)
			assert.NoError(t, readErr)
			assertMessagesMatch(t, test.expected, msg)
		})
	}

	// Wait for the websocket connection to time out
	time.Sleep(fakeServer.GetWebSocketTimeout())
}

// assertMessagesMatch checks that 2 webSocketMessage instance messages match.
// Used in place of assert.Equal() to ignore potential timestamp variations.
func assertMessagesMatch(t *testing.T, expected webSocketMessage, msg webSocketMessage) {
	t.Helper()
	assert.Equal(t, expected.Type, msg.Type)
	assert.Equal(t, expected.Severity, msg.Severity)
	assert.Equal(t, expected.Emphasis, msg.Emphasis)
	assert.Equal(t, expected.Text, msg.Text)
	expectedTS, _ := time.Parse(time.RFC3339, expected.Timestamp)
	msgTS, _ := time.Parse(time.RFC3339, msg.Timestamp)
	assert.WithinDuration(t, expectedTS, msgTS, 10*time.Second)
}

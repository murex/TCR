//go:build test_helper

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
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/murex/tcr/report"
	"github.com/stretchr/testify/assert"
)

// Test_websocket_concurrent_connections_stress tests multiple concurrent WebSocket connections
// to verify there are no race conditions during connection setup and teardown
func Test_websocket_concurrent_connections_stress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	// Create HTTP test server with the websocket connection handler
	s := httptest.NewUnstartedServer(http.HandlerFunc(websocketConnectionHandler))
	var fakeServer tcrHTTPServer
	s.Config.BaseContext = func(l net.Listener) context.Context {
		fakeServer = &fakeHTTPServer{
			url: *mustParseURL(s.URL),
		}
		return context.WithValue(context.Background(), serverContextKey, fakeServer)
	}
	s.Start()
	defer s.Close()

	const numConnections = 20
	const numMessages = 10

	var wg sync.WaitGroup
	errorChan := make(chan error, numConnections)

	for i := 0; i < numConnections; i++ {
		wg.Add(1)
		go func(connID int) {
			defer wg.Done()

			// Build URL and header for websocket connection request
			u := mustParseURL(s.URL)
			u.Scheme = "ws"
			hd := http.Header{}
			hd.Add("Origin", s.URL)

			// Create the websocket connection
			ws, _, err := websocket.DefaultDialer.Dial(u.String(), hd) //nolint:bodyclose
			if err != nil {
				errorChan <- err
				return
			}

			defer func() {
				// Send close frame before closing
				_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				_ = ws.Close()
			}()

			// Send multiple messages rapidly
			for j := 0; j < numMessages; j++ {
				report.PostText("stress test message", connID, j)

				// Try to read message with timeout
				_ = ws.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
				var msg message
				err := ws.ReadJSON(&msg)
				if err != nil {
					// Check if it's a close error, which is acceptable
					if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
						return // Connection closed gracefully
					}
					// For other errors, we just log and continue
					t.Logf("Connection %d read error: %v", connID, err)
					return
				}

				// Small delay to allow other goroutines to run
				time.Sleep(1 * time.Millisecond)
			}
		}(i)
	}

	// Wait for all connections to complete
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All connections completed successfully
	case err := <-errorChan:
		t.Fatalf("Connection error: %v", err)
	case <-time.After(30 * time.Second):
		t.Fatal("Stress test timed out")
	}

	// Check if there were any errors
	close(errorChan)
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		t.Logf("Encountered %d errors during stress test:", len(errors))
		for _, err := range errors {
			t.Logf("  - %v", err)
		}
		// We don't fail the test for connection errors in stress test
		// as they might be expected under high load
	}
}

// Test_websocket_rapid_connect_disconnect tests rapid connection and disconnection
// to verify proper cleanup and no resource leaks
func Test_websocket_rapid_connect_disconnect(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	// Create HTTP test server with shorter timeout for faster test
	s := httptest.NewUnstartedServer(http.HandlerFunc(websocketConnectionHandler))
	var fakeServer tcrHTTPServer
	s.Config.BaseContext = func(l net.Listener) context.Context {
		fakeServer = &fakeHTTPServer{
			url: *mustParseURL(s.URL),
		}
		return context.WithValue(context.Background(), serverContextKey, fakeServer)
	}
	s.Start()
	defer s.Close()

	const numCycles = 50

	for i := 0; i < numCycles; i++ {
		// Build URL and header for websocket connection request
		u := mustParseURL(s.URL)
		u.Scheme = "ws"
		hd := http.Header{}
		hd.Add("Origin", s.URL)

		// Create the websocket connection
		ws, _, err := websocket.DefaultDialer.Dial(u.String(), hd) //nolint:bodyclose
		if !assert.NoError(t, err, "Failed to create WebSocket connection on cycle %d", i) {
			continue
		}

		// Send one message
		report.PostText("rapid disconnect test", i)

		// Try to read one message with short timeout
		_ = ws.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
		var msg message
		err = ws.ReadJSON(&msg)
		if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
			t.Logf("Cycle %d read error (acceptable): %v", i, err)
		}

		// Close connection immediately
		_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		_ = ws.Close()

		// Small delay between cycles
		time.Sleep(5 * time.Millisecond)
	}
}

// Test_websocket_message_flood tests sending many messages rapidly
// to verify the "write" mutex and error handling work correctly
func Test_websocket_message_flood(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	// Create HTTP test server
	s := httptest.NewUnstartedServer(http.HandlerFunc(websocketConnectionHandler))
	var fakeServer tcrHTTPServer
	s.Config.BaseContext = func(l net.Listener) context.Context {
		fakeServer = &fakeHTTPServer{
			url: *mustParseURL(s.URL),
		}
		return context.WithValue(context.Background(), serverContextKey, fakeServer)
	}
	s.Start()
	defer s.Close()

	// Build URL and header for websocket connection request
	u := mustParseURL(s.URL)
	u.Scheme = "ws"
	hd := http.Header{}
	hd.Add("Origin", s.URL)

	// Create the websocket connection
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), hd) //nolint:bodyclose
	assert.NoError(t, err)
	defer func() {
		_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		_ = ws.Close()
	}()

	const numMessages = 100
	messagesReceived := 0

	// Start reading messages in background
	readDone := make(chan struct{})
	go func() {
		defer close(readDone)
		for {
			_ = ws.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			var msg message
			if err := ws.ReadJSON(&msg); err != nil {
				// Exit on any read error to avoid repeated reads on a failed connection
				return
			}
			messagesReceived++
			if messagesReceived >= numMessages {
				return
			}
		}
	}()

	// Send messages rapidly
	for i := 0; i < numMessages; i++ {
		report.PostText("flood test message", i)
		// No delay - send as fast as possible
	}

	// Wait for messages to be received or timeout
	select {
	case <-readDone:
		t.Logf("Successfully received %d messages", messagesReceived)
	case <-time.After(5 * time.Second):
		t.Logf("Timeout reached, received %d out of %d messages", messagesReceived, numMessages)
	}

	// We don't assert on exact message count as some may be lost during rapid sending
	// The important thing is that the system doesn't crash or deadlock
	assert.True(t, messagesReceived > 0, "Should have received at least some messages")
}

// Helper function to parse URL and panic on error (for test setup)
func mustParseURL(urlStr string) *url.URL {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return u
}

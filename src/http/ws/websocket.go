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
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/report/role_event"
	"github.com/murex/tcr/report/text"
	"github.com/murex/tcr/report/timer_event"
)

// message is used to JSON-encode TCR report messages
type message struct {
	Type      string `json:"type"`
	Severity  string `json:"severity"`
	Text      string `json:"text"`
	Emphasis  bool   `json:"emphasis"`
	Timestamp string `json:"timestamp"`
}

type messageType string

const (
	messageTypeSimple  messageType = "simple"
	messageTypeInfo    messageType = "info"
	messageTypeTitle   messageType = "title"
	messageTypeSuccess messageType = "success"
	messageTypeWarning messageType = "warning"
	messageTypeError   messageType = "error"
	messageTypeRole    messageType = "role"
	messageTypeTimer   messageType = "timer"
)

type messageSeverity int

const (
	messageSeverityNormal = iota
	messageSeverityLow
	messageSeverityHigh
)

func newMessage(msgType messageType, severity messageSeverity, emphasis bool, a ...any) message {
	return message{
		Type:      string(msgType),
		Severity:  strconv.Itoa(int(severity)),
		Text:      fmt.Sprint(a...),
		Emphasis:  emphasis,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// MessageReporter is in charge of sending TCR report messages over a websocket
type MessageReporter struct {
	server           tcrHTTPServer
	reportingChannel chan bool
	conn             *websocket.Conn
	connMutex        sync.Mutex
}

func newMessageReporter(server tcrHTTPServer, conn *websocket.Conn) *MessageReporter {
	return &MessageReporter{
		server: server,
		conn:   conn,
	}
}

func (r *MessageReporter) startReporting() {
	r.reportingChannel = report.Subscribe(r)
}

func (r *MessageReporter) stopReporting() {
	if r.reportingChannel != nil {
		report.Unsubscribe(r.reportingChannel)
	}
}

// ReportSimple reports simple messages
func (r *MessageReporter) ReportSimple(emphasis bool, payload text.Message) {
	r.write(newMessage(messageTypeSimple, messageSeverityNormal, emphasis, payload.ToString()))
}

// ReportInfo reports info messages
func (r *MessageReporter) ReportInfo(emphasis bool, payload text.Message) {
	r.write(newMessage(messageTypeInfo, messageSeverityNormal, emphasis, payload.ToString()))
}

// ReportTitle reports title messages
func (r *MessageReporter) ReportTitle(emphasis bool, payload text.Message) {
	r.write(newMessage(messageTypeTitle, messageSeverityNormal, emphasis, payload.ToString()))
}

// ReportSuccess reports success messages
func (r *MessageReporter) ReportSuccess(emphasis bool, payload text.Message) {
	r.write(newMessage(messageTypeSuccess, messageSeverityNormal, emphasis, payload.ToString()))
}

// ReportWarning reports warning messages
func (r *MessageReporter) ReportWarning(emphasis bool, payload text.Message) {
	r.write(newMessage(messageTypeWarning, messageSeverityLow, emphasis, payload.ToString()))
}

// ReportError reports error messages
func (r *MessageReporter) ReportError(emphasis bool, payload text.Message) {
	r.write(newMessage(messageTypeError, messageSeverityHigh, emphasis, payload.ToString()))
}

// ReportRoleEvent reports role event messages
func (r *MessageReporter) ReportRoleEvent(emphasis bool, payload role_event.Message) {
	r.write(newMessage(messageTypeRole, messageSeverityNormal, emphasis, payload.ToString()))
}

// ReportTimerEvent reports timer event messages
func (r *MessageReporter) ReportTimerEvent(emphasis bool, payload timer_event.Message) {
	r.write(newMessage(messageTypeTimer, messageSeverityNormal, emphasis, payload.ToString()))
}

func (r *MessageReporter) write(msg message) {
	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	// Check if connection is still open before writing
	if r.conn == nil {
		return
	}

	// Set a write deadline to prevent hanging
	_ = r.conn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))

	// We deliberately ignore write errors, which could happen
	// every time a client closes their console browser page
	_ = r.conn.WriteJSON(msg)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		server := r.Context().Value(serverContextKey).(tcrHTTPServer)
		if server.InDevMode() {
			// server and client ports are different when running in devMode,
			// so we bypass any CORS restriction in this mode
			return true
		}
		origin := r.Header.Get("Origin")
		url, err := url.Parse(origin)
		if err != nil {
			report.PostWarning("invalid origin: \"", origin, "\" - ", err.Error())
			return false
		}
		if url.Host != server.GetServerAddress() {
			// Note: This policy is quite restrictive:
			// - can't use "localhost" in browser URL as server is listening on "127.0.0.1".
			// - will not work if we allow connections from any HTTP client (eg. listening on "0.0.0.0")
			// We may need to soften it a bit depending on intended usage
			report.PostWarning("client host not authorized: ", url.Host)
			return false
		}
		return true
	},
}

// WebsocketHandler is the entry point for handling websocket requests sent to the HTTP server
func WebsocketHandler(c *gin.Context) {
	// Converts the gin request into a "regular" http HandlerFunc
	websocketConnectionHandler(c.Writer, requestWithGinContext(c))
}

// websocketConnectionHandler is responsible for opening a new websocket connection request
// and keeping it alive until we reach the connection timeout
func websocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		report.PostWarning("failed to upgrade to a websocket connection: ", err.Error())
		return
	}

	server := r.Context().Value(serverContextKey).(tcrHTTPServer)
	reporter := newMessageReporter(server, conn)
	reporter.startReporting()

	// Set up a context for graceful shutdown
	ctx, cancel := context.WithTimeout(r.Context(), server.GetWebsocketTimeout())
	defer cancel()

	defer func() {
		// Stop reporting first to prevent new messages
		reporter.stopReporting()

		// Mark connection as closed in reporter
		reporter.connMutex.Lock()
		reporter.conn = nil
		reporter.connMutex.Unlock()

		// Send close frame and close connection gracefully
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		_ = conn.Close()
	}()

	// We kill the connection after a fixed period of time to avoid keeping sending
	// messages to clients that are no longer there.
	// This should not be an issue for clients that are still connected
	// as the webapp client will automatically open a new connection after this one is gone.
	<-ctx.Done()
}

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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/utils"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// webSocketMessage is used to JSON-encode TCR report messages
type webSocketMessage struct {
	Type      string `json:"type"`
	Severity  string `json:"severity"`
	Text      string `json:"text"`
	Emphasis  bool   `json:"emphasis"`
	Timestamp string `json:"timestamp"`
}

func newWebSocketMessage(msgType string, severity string, emphasis bool, a ...any) webSocketMessage {
	return webSocketMessage{
		Type:      msgType,
		Severity:  severity,
		Text:      fmt.Sprint(a...),
		Emphasis:  emphasis,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// WebsocketMessageReporter is in charge of sending TCR report messages over a websocket
type WebsocketMessageReporter struct {
	reportingChannel chan bool
	conn             *websocket.Conn
	connMutex        sync.Mutex
}

func newWebSocketMessageReporter(conn *websocket.Conn) *WebsocketMessageReporter {
	var reporter = &WebsocketMessageReporter{conn: conn}
	reporter.startReporting()
	return reporter
}

func (r *WebsocketMessageReporter) startReporting() {
	server.RegisterWebSocket(r)
	r.reportingChannel = report.Subscribe(r)
}

func (r *WebsocketMessageReporter) stopReporting() {
	server.UnregisterWebSocket(r)
	if r.reportingChannel != nil {
		report.Unsubscribe(r.reportingChannel)
	}
}

// ReportSimple reports simple messages
func (r *WebsocketMessageReporter) ReportSimple(emphasis bool, a ...any) {
	r.write(newWebSocketMessage("simple", "0", emphasis, a...))
}

// ReportInfo reports info messages
func (r *WebsocketMessageReporter) ReportInfo(emphasis bool, a ...any) {
	r.write(newWebSocketMessage("info", "0", emphasis, a...))
}

// ReportTitle reports title messages
func (r *WebsocketMessageReporter) ReportTitle(emphasis bool, a ...any) {
	r.write(newWebSocketMessage("title", "0", emphasis, a...))
}

// ReportRole reports role event messages
// Note: this function is not part of the reporter interface (should be added)
func (r *WebsocketMessageReporter) ReportRole(emphasis bool, a ...any) {
	r.write(newWebSocketMessage("role", "0", emphasis, a...))
}

// ReportTimer reports timer messages
func (r *WebsocketMessageReporter) ReportTimer(emphasis bool, a ...any) {
	r.write(newWebSocketMessage("timer", "0", emphasis, a...))
}

// ReportSuccess reports success messages
func (r *WebsocketMessageReporter) ReportSuccess(emphasis bool, a ...any) {
	r.write(newWebSocketMessage("success", "0", emphasis, a...))
}

// ReportWarning reports warning messages
func (r *WebsocketMessageReporter) ReportWarning(emphasis bool, a ...any) {
	r.write(newWebSocketMessage("warning", "1", emphasis, a...))
}

// ReportError reports error messages
func (r *WebsocketMessageReporter) ReportError(emphasis bool, a ...any) {
	r.write(newWebSocketMessage("error", "2", emphasis, a...))
}

func (r *WebsocketMessageReporter) write(msg webSocketMessage) {
	r.connMutex.Lock()
	err := r.conn.WriteJSON(msg)
	r.connMutex.Unlock()
	if err != nil {
		report.PostWarning("websocket message sending failure - ", err.Error())
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
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

// WebSocketHandler is the entry point for handling websocket requests sent to the HTTP server
func WebSocketHandler(c *gin.Context) {
	handleWebSocket(c.Writer, c.Request)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.Trace(err)
		return
	}

	reporter := newWebSocketMessageReporter(conn)

	defer func() {
		reporter.stopReporting()
		_ = conn.Close()
	}()

	// We kill the connection after a fixed period of time to avoid keeping sending
	// messages to clients that are no longer there.
	// This should not be an issue for clients that are still connected
	// as the webapp client will automatically open a new connection after this one is gone.
	time.Sleep(server.GetWebSocketTimeout())
}

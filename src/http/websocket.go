/*
Copyright (c) 2023 Murex

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

package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/utils"
	"net/http"
	"sync"
	"time"
)

// webSocketConnectionTimeout is the delay after which we close a websocket connection
const webSocketConnectionTimeout = 1 * time.Minute

type webSocketMessage struct {
	Type      string `json:"type"`
	Severity  string `json:"severity"`
	Text      string `json:"text"`
	Emphasis  bool   `json:"emphasis"`
	Timestamp string `json:"timestamp"`
}

type websocketMessageReporter struct {
	reportingChannel chan bool
	conn             *websocket.Conn
	connMutex        sync.Mutex
}

func newWebSocketMessageReporter(conn *websocket.Conn) *websocketMessageReporter {
	var reporter = &websocketMessageReporter{conn: conn}
	reporter.startReporting()
	return reporter
}

func (r *websocketMessageReporter) startReporting() {
	ServerInstance.registerWebSocket(r)
	r.reportingChannel = report.Subscribe(r)
}

func (r *websocketMessageReporter) stopReporting() {
	ServerInstance.unregisterWebSocket(r)
	if r.reportingChannel != nil {
		report.Unsubscribe(r.reportingChannel)
	}
}

// ReportSimple reports simple messages
func (r *websocketMessageReporter) ReportSimple(emphasis bool, a ...any) {
	msg := webSocketMessage{
		Type:      "simple",
		Severity:  "0",
		Text:      fmt.Sprint(a...),
		Emphasis:  emphasis,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	r.write(msg)
}

// ReportInfo reports info messages
func (r *websocketMessageReporter) ReportInfo(emphasis bool, a ...any) {
	msg := webSocketMessage{
		Type:      "info",
		Severity:  "0",
		Text:      fmt.Sprint(a...),
		Emphasis:  emphasis,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	r.write(msg)
}

// ReportTitle reports title messages
func (r *websocketMessageReporter) ReportTitle(emphasis bool, a ...any) {
	msg := webSocketMessage{
		Type:      "title",
		Severity:  "0",
		Text:      fmt.Sprint(a...),
		Emphasis:  emphasis,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	r.write(msg)
}

// ReportTimer reports timer messages
func (r *websocketMessageReporter) ReportTimer(emphasis bool, a ...any) {
	msg := webSocketMessage{
		Type:      "timer",
		Severity:  "0",
		Text:      fmt.Sprint(a...),
		Emphasis:  emphasis,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	r.write(msg)
}

// ReportSuccess reports success messages
func (r *websocketMessageReporter) ReportSuccess(emphasis bool, a ...any) {
	msg := webSocketMessage{
		Type:      "success",
		Severity:  "0",
		Text:      fmt.Sprint(a...),
		Emphasis:  emphasis,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	r.write(msg)
}

// ReportWarning reports warning messages
func (r *websocketMessageReporter) ReportWarning(emphasis bool, a ...any) {
	msg := webSocketMessage{
		Type:      "warning",
		Severity:  "1",
		Text:      fmt.Sprint(a...),
		Emphasis:  emphasis,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	r.write(msg)
}

// ReportError reports error messages
func (r *websocketMessageReporter) ReportError(emphasis bool, a ...any) {
	msg := webSocketMessage{
		Type:      "error",
		Severity:  "2",
		Text:      fmt.Sprint(a...),
		Emphasis:  emphasis,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	r.write(msg)
}

func (r *websocketMessageReporter) write(msg webSocketMessage) {
	r.connMutex.Lock()
	err := r.conn.WriteJSON(msg)
	r.connMutex.Unlock()
	if err != nil {
		// utils.Trace(err)
		// TODO handle case when client is gone
		// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
		//	var req *http.Request
		//	report.PostWarning("error: %v, user-agent: %v", err, req.Header.Get("User-Agent"))
		// }
		return
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// fmt.Println(r)
		// TODO enforce origin in production mode?
		return true
	},
}

func webSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Trace(err)
		return
	}

	r := newWebSocketMessageReporter(conn)

	defer func() {
		r.stopReporting()
		_ = conn.Close()
	}()

	// We kill the connection after a fixed period of time to avoid keeping sending
	// messages to clients that are no longer there.
	// This should not be an issue for clients that are still connected
	// as the webapp client will automatically open a new connection after this one is gone.
	time.Sleep(webSocketConnectionTimeout)
}

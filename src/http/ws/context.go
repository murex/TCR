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
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type tcrHTTPServer = interface {
	InDevMode() bool
	GetServerAddress() string
	GetWebsocketTimeout() time.Duration
}

type serverContextKeyType string

const serverContextKey serverContextKeyType = "tcr-http-server"

// HTTPServerMiddleware adds the HTTP server instance to gin context
// so that websocket handlers can interact with it.
func HTTPServerMiddleware(s tcrHTTPServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(string(serverContextKey), s)
		c.Next()
	}
}

// requestWithGinContext inserts gin context value set for
// HTTP server instance to the request context sent to websocket handler
func requestWithGinContext(c *gin.Context) *http.Request {
	ctx := context.WithValue(c.Request.Context(),
		serverContextKey, c.MustGet(string(serverContextKey)).(tcrHTTPServer))
	return c.Request.WithContext(ctx)
}

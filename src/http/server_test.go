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

package http

import (
	"github.com/gin-gonic/gin"
	"github.com/murex/tcr/engine"
	"github.com/murex/tcr/params"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_create_http_server(t *testing.T) {
	p := *params.AParamSet()
	tcr := engine.NewFakeTCREngine()
	srv := New(p, tcr)

	tests := []struct {
		desc     string
		asserter func(t *testing.T)
	}{
		{
			desc: "access to TCR instance",
			asserter: func(t *testing.T) {
				assert.Equal(t, tcr, srv.tcr)
			},
		},
		{
			desc: "host",
			asserter: func(t *testing.T) {
				assert.Equal(t, "127.0.0.1", srv.host)
			},
		},
		{
			desc: "development mode",
			asserter: func(t *testing.T) {
				assert.Equal(t, true, srv.devMode)
			},
		},
		{
			desc: "websocket connections timeout",
			asserter: func(t *testing.T) {
				assert.Equal(t, 1*time.Minute, srv.websocketTimeout)
			},
		},
		{
			desc: "websocket connection pool",
			asserter: func(t *testing.T) {
				assert.Equal(t, 0, len(*srv.websockets))
			},
		},
		{
			desc: "access to application parameters",
			asserter: func(t *testing.T) {
				assert.Equal(t, p, srv.params)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, test.asserter)
	}
}

func Test_cors_middleware_handler(t *testing.T) {

	tests := []struct {
		desc     string
		cors     bool
		origin   string
		expected string
	}{
		{
			desc:     "with CORS and same origin",
			cors:     true,
			origin:   "",
			expected: "",
		},
		{
			desc:     "with CORS and different origin",
			cors:     true,
			origin:   "http://some.other.origin",
			expected: "*",
		},
		{
			desc:     "without CORS",
			cors:     false,
			origin:   "http://some.other.origin",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			// Setup the router
			rPath := "/"
			router := gin.Default()

			// Add handlers
			if test.cors {
				router.Use(corsMiddleware())
			}
			router.GET(rPath, func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// Prepare the request, send it and capture the response
			req, _ := http.NewRequest("GET", rPath, nil)
			if test.origin != "" {
				req.Header.Add("Origin", test.origin)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, test.expected, w.Header().Get("Access-Control-Allow-Origin"))
		})
	}
}

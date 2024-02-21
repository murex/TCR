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
	"slices"
	"strings"
	"testing"
	"time"
)

func Test_create_web_ui_server(t *testing.T) {
	p := *params.AParamSet()
	tcr := engine.NewFakeTCREngine()
	wuis := New(p, tcr)

	tests := []struct {
		desc     string
		asserter func(t *testing.T)
	}{
		{
			desc: "access to TCR instance",
			asserter: func(t *testing.T) {
				assert.Equal(t, tcr, wuis.tcr)
			},
		},
		{
			desc: "host",
			asserter: func(t *testing.T) {
				assert.Equal(t, "127.0.0.1", wuis.host)
			},
		},
		{
			desc: "development mode",
			asserter: func(t *testing.T) {
				assert.Equal(t, true, wuis.devMode)
			},
		},
		{
			desc: "http server instance",
			asserter: func(t *testing.T) {
				assert.Nil(t, wuis.httpServer)
			},
		},
		{
			desc: "websocket connections timeout",
			asserter: func(t *testing.T) {
				assert.Equal(t, 1*time.Minute, wuis.websocketTimeout)
			},
		},
		{
			desc: "websocket connection pool",
			asserter: func(t *testing.T) {
				assert.Equal(t, 0, len(*wuis.websockets))
			},
		},
		{
			desc: "access to application parameters",
			asserter: func(t *testing.T) {
				assert.Equal(t, p, wuis.params)
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
			req, _ := http.NewRequest(http.MethodGet, rPath, nil)
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

func Test_init_gin_engine(t *testing.T) {
	tests := []struct {
		desc             string
		devMode          bool
		expectedHandlers int
	}{
		{
			desc:             "development mode",
			devMode:          true,
			expectedHandlers: 3, // gin.Recovery, gin.Logger and corsMiddleware
		},
		{
			desc:             "production mode",
			devMode:          false,
			expectedHandlers: 1, // gin.Recovery only
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			srv := New(*params.AParamSet(), engine.NewFakeTCREngine())
			srv.devMode = test.devMode
			srv.initGinEngine()
			assert.NotNil(t, srv.router)
			assert.Equal(t, test.expectedHandlers, len(srv.router.Handlers))
		})
	}
}

func Test_add_static_routes(t *testing.T) {
	// Setup the router
	wuis := New(*params.AParamSet(), engine.NewFakeTCREngine())
	wuis.initGinEngine()
	wuis.addStaticRoutes()

	// Prepare the request, send it and capture the response
	rPath := "/some_path"
	req, _ := http.NewRequest(http.MethodGet, rPath, nil)
	w := httptest.NewRecorder()
	wuis.router.ServeHTTP(w, req)

	// Note: we don't test the regular case where rPath = "/" (returning a StatusOK)
	// This is to avoid getting dangling test results depending whether
	// frontend files have been generated under static/webapp/browser
	assert.Equal(t, http.StatusMovedPermanently, w.Code)
}

type testRESTRouteParams struct {
	path    string   // REST request path
	methods []string // REST methods accepted for this path
}

func allRESTMethods() []string {
	return []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodPost,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodHead,
		http.MethodOptions,
		http.MethodTrace,
	}
}

func testRESTRoutes(t *testing.T, router *gin.Engine, tests []testRESTRouteParams) {
	t.Helper()
	for _, test := range tests {
		// Hierarchical test runner seems to get confused when there are "/" in description.
		// Replacing them with "\" allows to workaround this issue
		descPath := strings.Replace(test.path, "/", "\\", -1)
		t.Run(descPath, func(t *testing.T) {
			for _, method := range allRESTMethods() {
				methodSupported := slices.Contains(test.methods, method)
				descMethod := method
				if methodSupported {
					descMethod += "_(yes)"
				} else {
					descMethod += "_(no)"
				}
				t.Run(descMethod, func(t *testing.T) {
					// Prepare the request, send it and capture the response
					req, _ := http.NewRequest(method, test.path, nil)
					w := httptest.NewRecorder()
					router.ServeHTTP(w, req)
					// We only check at router level if we get a 404 or not.
					// Actual return codes should be tested at handlers levels
					assert.Equal(t, methodSupported, w.Code != http.StatusNotFound)
				})
			}
		})
	}
}

func Test_add_api_routes(t *testing.T) {
	tests := []testRESTRouteParams{
		{
			path:    "/api",
			methods: []string{},
		},
		{
			path:    "/api/build-info",
			methods: []string{http.MethodGet},
		},
		{
			path:    "/api/session-info",
			methods: []string{http.MethodGet},
		},
		{
			path:    "/api/roles",
			methods: []string{http.MethodGet},
		},
		{
			path:    "/api/roles/name",
			methods: []string{http.MethodGet},
		},
		{
			path:    "/api/roles/name/action",
			methods: []string{http.MethodPost},
		},
	}

	// Setup the router
	wuis := New(*params.AParamSet(), engine.NewFakeTCREngine())
	wuis.initGinEngine()
	wuis.addAPIRoutes()

	// check every route + REST method combination
	testRESTRoutes(t, wuis.router, tests)
}

func Test_add_websocket_routes(t *testing.T) {
	tests := []testRESTRouteParams{
		{
			path:    "/ws",
			methods: []string{http.MethodGet},
		},
		{
			path:    "/ws/any",
			methods: []string{""},
		},
	}

	// Setup the router
	wuis := New(*params.AParamSet(), engine.NewFakeTCREngine())
	wuis.initGinEngine()
	wuis.addWebsocketRoutes()

	// check every route + REST method combination
	testRESTRoutes(t, wuis.router, tests)
}

func Test_start_server(t *testing.T) {
	wuis := New(*params.AParamSet(), engine.NewFakeTCREngine())
	wuis.Start()
	t.Cleanup(func() {
		wuis.stopGinEngine()
	})

	// Check that the HTTP server instance is here
	assert.NotNil(t, wuis.httpServer)
	// Smoke test: Send a simple HTTP request and verify that we get a response
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	wuis.httpServer.Handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

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

package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/murex/tcr/engine"
	"github.com/stretchr/testify/assert"
)

func Test_controls_post_handler(t *testing.T) {
	tests := []struct {
		control              string
		expectedHTTPResponse int
		expectedCalls        []engine.TCRCall
	}{
		{
			control:              controlAbortCommand,
			expectedHTTPResponse: http.StatusAccepted,
			expectedCalls:        []engine.TCRCall{engine.TCRCallAbortCommand},
		},
		{
			control:              "unrecognized-control",
			expectedHTTPResponse: http.StatusBadRequest,
			expectedCalls:        nil,
		},
	}

	for _, test := range tests {
		subPath := test.control
		t.Run(subPath, func(t *testing.T) {
			// Setup the router
			router := gin.Default()
			// Plug it to a fake TCR engine
			tcr := engine.NewFakeTCREngine()
			router.Use(TCREngineMiddleware(tcr))
			// Add the route
			rPath := "/api/controls/:name"
			router.POST(rPath, ControlsPostHandler)

			// Prepare the request, send it and capture the response
			req, _ := http.NewRequest(http.MethodPost,
				strings.Replace(rPath, ":name", subPath, 1), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify HTTP response code
			assert.Equal(t, test.expectedHTTPResponse, w.Code)
			// Verify TCR call history
			assert.Equal(t, test.expectedCalls, tcr.GetCallHistory())
		})
	}
}

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
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/murex/tcr/engine"
	"github.com/murex/tcr/role"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_roles_get_handler(t *testing.T) {
	// Setup the router
	rPath := "/api/roles"
	router := gin.Default()
	tcr := engine.NewFakeTCREngine()
	router.Use(TCREngineMiddleware(tcr))
	router.GET(rPath, RolesGetHandler)

	// Prepare the request, send it and capture the response
	req, _ := http.NewRequest(http.MethodGet, rPath, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify the response's code, header and body
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	var expected []roleData
	current := tcr.GetCurrentRole()
	for _, r := range role.All() {
		expected = append(expected, roleData{
			Name:        r.Name(),
			Description: r.LongName(),
			Active:      r == current,
		})
	}
	var actual []roleData
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &actual))
	assert.Equal(t, expected, actual)
}

func Test_role_get_handler(t *testing.T) {
	// Setup the router
	rPath := "/api/roles/:name"
	router := gin.Default()
	tcr := engine.NewFakeTCREngine()
	router.Use(TCREngineMiddleware(tcr))
	router.GET(rPath, RoleGetHandler)

	for _, r := range role.All() {
		t.Run(r.Name(), func(t *testing.T) {
			// Prepare the request, send it and capture the response
			req, _ := http.NewRequest(http.MethodGet,
				strings.Replace(rPath, ":name", r.Name(), 1), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify the response's code, header and body
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
			expected := roleData{
				Name:        r.Name(),
				Description: r.LongName(),
				Active:      r == tcr.GetCurrentRole(),
			}
			var actual roleData
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &actual))
			assert.Equal(t, expected, actual)

		})
	}
}

func Test_role_get_handler_with_invalid_params(t *testing.T) {
	// Setup the router
	rPath := "/api/roles/:name"
	router := gin.Default()
	tcr := engine.NewFakeTCREngine()
	router.Use(TCREngineMiddleware(tcr))
	router.GET(rPath, RoleGetHandler)

	// Prepare the request, send it and capture the response
	req, _ := http.NewRequest(http.MethodGet,
		strings.Replace(rPath, ":name", "unknown-role", 1), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify the response's code
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_roles_post_handler(t *testing.T) {
	// Setup the router
	rPath := "/api/roles/:name/:action"
	router := gin.Default()
	tcr := engine.NewFakeTCREngine()
	router.Use(TCREngineMiddleware(tcr))
	router.POST(rPath, RolesPostHandler)

	tests := []struct {
		role     string
		action   string
		expected roleData
	}{
		{
			role:   role.Navigator{}.Name(),
			action: startAction,
			expected: roleData{
				Name:        role.Navigator{}.Name(),
				Description: role.Navigator{}.LongName(),
				Active:      true,
			},
		},
		{
			role:   role.Navigator{}.Name(),
			action: stopAction,
			expected: roleData{
				Name:        role.Navigator{}.Name(),
				Description: role.Navigator{}.LongName(),
				Active:      false,
			},
		},
		{
			role:   role.Driver{}.Name(),
			action: startAction,
			expected: roleData{
				Name:        role.Driver{}.Name(),
				Description: role.Driver{}.LongName(),
				Active:      true,
			},
		},
		{
			role:   role.Driver{}.Name(),
			action: stopAction,
			expected: roleData{
				Name:        role.Driver{}.Name(),
				Description: role.Driver{}.LongName(),
				Active:      false,
			},
		},
	}

	for _, test := range tests {
		subPath := test.role + "/" + test.action
		t.Run(subPath, func(t *testing.T) {
			// Prepare the request, send it and capture the response
			req, _ := http.NewRequest(http.MethodPost,
				strings.Replace(rPath, ":name/:action", subPath, 1), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify the response's code, header and body
			assert.Equal(t, http.StatusAccepted, w.Code)
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
			var actual roleData
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &actual))
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_roles_post_handler_with_invalid_params(t *testing.T) {
	// Setup the router
	rPath := "/api/roles/:name/:action"
	router := gin.Default()
	tcr := engine.NewFakeTCREngine()
	router.Use(TCREngineMiddleware(tcr))
	router.POST(rPath, RolesPostHandler)

	const invalidRole = "unknown-role"
	const invalidAction = "unknown-action"
	tests := []struct {
		role   string
		action string
	}{
		{role: invalidRole, action: startAction},
		{role: invalidRole, action: stopAction},
		{role: role.Navigator{}.Name(), action: invalidAction},
		{role: role.Driver{}.Name(), action: invalidAction},
	}

	for _, test := range tests {
		subPath := test.role + "/" + test.action
		t.Run(subPath, func(t *testing.T) {
			// Prepare the request, send it and capture the response
			req, _ := http.NewRequest(http.MethodPost,
				strings.Replace(rPath, ":name/:action", subPath, 1), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify the response's code
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

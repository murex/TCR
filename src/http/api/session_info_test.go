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
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_session_info_get_handler(t *testing.T) {
	// Setup the router
	rPath := "/api/session-info"
	router := gin.Default()
	tcr := engine.NewFakeTCREngine()
	router.Use(TCREngineMiddleware(tcr))
	router.GET(rPath, SessionInfoGetHandler)

	// Prepare the request, send it and capture the response
	req, _ := http.NewRequest(http.MethodGet, rPath, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify the response's code, header and body
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	info := tcr.GetSessionInfo()
	expected := sessionInfo{
		BaseDir:           info.BaseDir,
		WorkDir:           info.WorkDir,
		LanguageName:      info.LanguageName,
		ToolchainName:     info.ToolchainName,
		VCSName:           info.VCSName,
		VCSSessionSummary: info.VCSSessionSummary,
		Variant:           info.Variant,
		GitAutoPush:       info.GitAutoPush,
		MessageSuffix:     info.MessageSuffix,
	}
	var actual sessionInfo
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &actual))
	assert.Equal(t, expected, actual)
}

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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_http_server_middleware_updates_gin_context(t *testing.T) {
	server := newFakeHTTPServer("http://127.0.0.1")
	ctx := gin.Context{}

	HTTPServerMiddleware(server)(&ctx)

	assert.Equal(t, server, ctx.MustGet(string(serverContextKey)))
}

func Test_insert_gin_context_value_into_http_request(t *testing.T) {
	server := newFakeHTTPServer("http://127.0.0.1")
	req := http.Request{}
	ctx := gin.Context{Request: &req}
	ctx.Set(string(serverContextKey), server)

	newReq := requestWithGinContext(&ctx)

	assert.Equal(t, server, newReq.Context().Value(serverContextKey))
}

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
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeWebsocketWriter struct {
	operationCount int
}

func newFakeWebsocketWriter() *fakeWebsocketWriter {
	return &fakeWebsocketWriter{operationCount: 0}
}

func (f *fakeWebsocketWriter) ReportTitle(_ bool, _ ...any) {
	f.operationCount++
}

func (f *fakeWebsocketWriter) ReportRole(_ bool, _ ...any) {
	f.operationCount++
}

func Test_connection_pool_registration(t *testing.T) {
	cp := NewConnectionPool()
	ws1 := newFakeWebsocketWriter()
	ws2 := newFakeWebsocketWriter()

	assert.Empty(t, *cp)

	cp.Register(ws1)
	assert.Len(t, *cp, 1)

	cp.Register(ws2)
	assert.Len(t, *cp, 2)

	cp.Unregister(ws1)
	assert.Len(t, *cp, 1)

	cp.Unregister(ws2)
	assert.Len(t, *cp, 0)
}

func Test_connection_pool_dispatch(t *testing.T) {
	cp := NewConnectionPool()
	ws1 := newFakeWebsocketWriter()
	cp.Register(ws1)
	ws2 := newFakeWebsocketWriter()
	cp.Register(ws2)

	assert.Equal(t, 0, ws1.operationCount)
	assert.Equal(t, 0, ws2.operationCount)
	cp.Dispatch(func(w WebsocketWriter) {
		w.ReportTitle(false, "")
	})
	assert.Equal(t, 1, ws1.operationCount)
	assert.Equal(t, 1, ws2.operationCount)
}

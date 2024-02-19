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

// ConnectionPool is in charge of managing websocket connections
type ConnectionPool []WebsocketWriter

// NewConnectionPool returns a new instance of websocket connection pool
func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{}
}

// Register registers a new websocket with the connection pool
func (p *ConnectionPool) Register(w WebsocketWriter) {
	*p = append(*p, w)
}

// Unregister unregisters a websocket from the connection pool
func (p *ConnectionPool) Unregister(w WebsocketWriter) {
	for i, registered := range *p {
		if w == registered {
			*p = append((*p)[:i], (*p)[i+1:]...)
			return
		}
	}
}

// Dispatch dispatches the provided operation to all registered websockets
func (p *ConnectionPool) Dispatch(operation func(w WebsocketWriter)) {
	for _, registered := range *p {
		operation(registered)
	}
}

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

package utils

import (
	"fmt"
	"github.com/murex/tcr/settings"
	"io"
)

// simpleTraceWriter is the writer used by Trace() to write trace messages to io.Writer
var simpleTraceWriter io.Writer

// SetSimpleTrace sets the writer used by Trace()
func SetSimpleTrace(w io.Writer) {
	if w != nil {
		simpleTraceWriter = w
	}
}

// Trace writes simple trace messages
func Trace(a ...any) {
	if simpleTraceWriter != nil {
		_, _ = fmt.Fprintln(simpleTraceWriter, "["+settings.ApplicationName+"]", fmt.Sprint(a...))
	}
}

// TraceConfigValue writes simple trace messages for a configuration key/value pair
func TraceConfigValue(key string, value any) {
	Trace("- ", key, ": ", value)
}

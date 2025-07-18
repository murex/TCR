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

package helpers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_simple_trace_when_writer_is_not_set(t *testing.T) {
	assert.NotPanics(t, func() {
		SetSimpleTrace(nil)
		Trace("Some dummy message")
	})
}

func Test_simple_trace_format(t *testing.T) {
	msg := "Some dummy message"
	AssertSimpleTrace(t, []string{msg},
		func() {
			Trace(msg)
		},
	)
}

func Test_simple_trace_key_value_format(t *testing.T) {
	k, v := "some-key", "some-value"
	expected := []string{fmt.Sprintf("- %s: %s", k, v)}
	AssertSimpleTrace(t, expected,
		func() {
			TraceKeyValue(k, v)
		},
	)
}

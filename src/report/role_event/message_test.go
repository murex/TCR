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

package role_event

import (
	"testing"

	"github.com/murex/tcr/role"
	"github.com/stretchr/testify/assert"
)

func Test_wrap_message(t *testing.T) {
	tests := []struct {
		message Message
		wrapped string
	}{
		{
			message: New(TriggerStart, role.Driver{}),
			wrapped: "driver:start",
		},
		{
			message: New(TriggerStart, role.Navigator{}),
			wrapped: "navigator:start",
		},
		{
			message: New(TriggerEnd, role.Driver{}),
			wrapped: "driver:end",
		},
		{
			message: New(TriggerEnd, role.Navigator{}),
			wrapped: "navigator:end",
		},
	}

	for _, test := range tests {
		t.Run(test.wrapped, func(t *testing.T) {
			str := test.message.ToString()
			assert.Equal(t, test.wrapped, str)
		})
	}
}

func Test_message_emphasis(t *testing.T) {
	tests := []struct {
		desc     string
		message  Message
		expected bool
	}{
		{
			desc:     "start trigger",
			message:  New(TriggerStart, role.Driver{}),
			expected: false,
		},
		{
			desc:     "end trigger",
			message:  New(TriggerEnd, role.Driver{}),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert.Equal(t, test.expected, test.message.WithEmphasis())
		})
	}
}

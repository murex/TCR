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

package timer_event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_wrap_event_message(t *testing.T) {
	tests := []struct {
		message Message
		wrapped string
	}{
		{
			message: New(TriggerStart, 10*time.Second, 0*time.Second, 10*time.Second),
			wrapped: "start:10:0:10",
		},
		{
			message: New(TriggerCountdown, 10*time.Second, 1*time.Second, 9*time.Second),
			wrapped: "countdown:10:1:9",
		},
		{
			message: New(TriggerStop, 10*time.Second, 4*time.Second, 0*time.Second),
			wrapped: "stop:10:4:0",
		},
		{
			message: New(TriggerTimeout, 10*time.Second, 11*time.Second, -1*time.Second),
			wrapped: "timeout:10:11:-1",
		},
	}

	for _, test := range tests {
		t.Run(test.wrapped, func(t *testing.T) {
			str := test.message.ToString()
			assert.Equal(t, test.wrapped, str)
		})
	}
}

func Test_event_message_emphasis(t *testing.T) {
	tests := []struct {
		desc     string
		message  Message
		expected bool
	}{
		{
			desc:     "start trigger",
			message:  New(TriggerStart, 10*time.Second, 0*time.Second, 10*time.Second),
			expected: true,
		},
		{
			desc:     "countdown trigger",
			message:  New(TriggerCountdown, 10*time.Second, 1*time.Second, 9*time.Second),
			expected: true,
		},
		{
			desc:     "stop trigger",
			message:  New(TriggerStop, 10*time.Second, 4*time.Second, 0*time.Second),
			expected: true,
		},
		{
			desc:     "first timeout trigger",
			message:  New(TriggerTimeout, 10*time.Second, 10*time.Second, 0*time.Second),
			expected: true,
		},
		{
			desc:     "later timeout trigger",
			message:  New(TriggerTimeout, 10*time.Second, 11*time.Second, -1*time.Second),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert.Equal(t, test.expected, test.message.WithEmphasis())
		})
	}
}

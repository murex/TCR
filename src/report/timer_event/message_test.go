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
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_wrap_unwrap_event_message(t *testing.T) {
	tests := []struct {
		wrapped string
		message Message
	}{
		{
			wrapped: "start:10:0:10",
			message: Message{
				Trigger:   TriggerStart,
				Timeout:   10 * time.Second,
				Elapsed:   0 * time.Second,
				Remaining: 10 * time.Second,
			},
		},
		{
			wrapped: "countdown:10:1:9",
			message: Message{
				Trigger:   TriggerCountdown,
				Timeout:   10 * time.Second,
				Elapsed:   1 * time.Second,
				Remaining: 9 * time.Second,
			},
		},
		{
			wrapped: "stop:10:4:0",
			message: Message{
				Trigger:   TriggerStop,
				Timeout:   10 * time.Second,
				Elapsed:   4 * time.Second,
				Remaining: 0 * time.Second,
			},
		},
		{
			wrapped: "timeout:10:0:-5",
			message: Message{
				Trigger:   TriggerTimeout,
				Timeout:   10 * time.Second,
				Elapsed:   0 * time.Second,
				Remaining: -5 * time.Second,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.wrapped, func(t *testing.T) {
			result := WrapMessage(test.message)
			assert.Equal(t, test.wrapped, result)
			assert.Equal(t, test.message, UnwrapMessage(result))
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
			message:  Message{Trigger: TriggerStart},
			expected: true,
		},
		{
			desc:     "countdown trigger",
			message:  Message{Trigger: TriggerCountdown},
			expected: true,
		},
		{
			desc:     "stop trigger",
			message:  Message{Trigger: TriggerStop},
			expected: true,
		},
		{
			desc:     "first timeout trigger",
			message:  Message{Trigger: TriggerTimeout, Remaining: 0},
			expected: true,
		},
		{
			desc:     "later timeout trigger",
			message:  Message{Trigger: TriggerTimeout, Remaining: -1},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert.Equal(t, test.expected, test.message.WithEmphasis())
		})
	}
}

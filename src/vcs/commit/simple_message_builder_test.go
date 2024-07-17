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

package commit

import (
	"github.com/murex/tcr/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_simple_message_builder(t *testing.T) {
	sampleEvent := events.ATcrEvent(events.WithCommandStatus(events.StatusPass))
	textFormattedEvent := sampleEvent.ToYAML()

	tests := []struct {
		description string
		header      string
		event       *events.TCREvent
		suffix      string
		expected    []string
	}{
		{
			description: "message with header only",
			header:      "commit message header",
			event:       nil,
			suffix:      "",
			expected:    []string{"commit message header"},
		},
		{
			description: "message with event",
			header:      "commit message header",
			event:       sampleEvent,
			suffix:      "",
			expected:    []string{"commit message header", textFormattedEvent},
		},
		{
			description: "message with suffix",
			header:      "commit message header",
			event:       nil,
			suffix:      "my suffix",
			expected:    []string{"commit message header", "\nmy suffix"},
		},
		{
			description: "message with event and suffix",
			header:      "commit message header",
			event:       sampleEvent,
			suffix:      "my suffix",
			expected:    []string{"commit message header", textFormattedEvent, "\nmy suffix"},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			builder := NewSimpleMessageBuilder(test.header, test.event, test.suffix)
			message, _ := builder.GenerateMessage()
			assert.Equal(t, test.expected, message)
		})
	}
}

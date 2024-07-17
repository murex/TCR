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
)

// SimpleMessageBuilder provides implementation for a simple commit message builder. This
// is the default message builder used by TCR.
type SimpleMessageBuilder struct {
	header string
	event  *events.TCREvent
	suffix string
}

// NewSimpleMessageBuilder creates a new SimpleMessageBuilder instance
func NewSimpleMessageBuilder(header string, event *events.TCREvent, suffix string) *SimpleMessageBuilder {
	return &SimpleMessageBuilder{
		header: header,
		event:  event,
		suffix: suffix,
	}
}

// GenerateMessage generates a commit message
func (smb *SimpleMessageBuilder) GenerateMessage() ([]string, error) {
	messages := []string{smb.header}
	if smb.event != nil {
		messages = append(messages, smb.event.ToYAML())
	}
	if smb.suffix != "" {
		messages = append(messages, "\n"+smb.suffix)
	}
	return messages, nil
}

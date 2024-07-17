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
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_external_message_builder(t *testing.T) {
	tests := []struct {
		description string
		commandPath string
		commandArgs []string
		expectedMsg []string
		expectErr   bool
	}{
		{
			description: "command with no output",
			commandPath: "true",
			commandArgs: []string{},
			expectedMsg: []string{},
			expectErr:   false,
		},
		{
			description: "single-line echo command",
			commandPath: "echo",
			commandArgs: []string{"Hello World!"},
			expectedMsg: []string{"Hello World!"},
			expectErr:   false,
		},
		{
			description: "multi-line echo command",
			commandPath: "echo",
			commandArgs: []string{"'Hello\nWorld!'"},
			expectedMsg: []string{"Hello", "World!"},
			expectErr:   false,
		},
		{
			description: "unknown command should return an error",
			commandPath: "invalid",
			commandArgs: []string{},
			expectedMsg: []string{},
			expectErr:   true,
		},
		{
			description: "failing command should return an error",
			commandPath: "false",
			commandArgs: []string{},
			expectedMsg: []string{},
			expectErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			builder := NewExternalMessageBuilder(test.commandPath, test.commandArgs...)
			msg, err := builder.GenerateMessage()
			assert.Equal(t, test.expectedMsg, msg)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

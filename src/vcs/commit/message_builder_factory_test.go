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
	"os"
	"testing"
)

func Test_builder_factory_default(t *testing.T) {
	_ = os.Unsetenv(envVariableName)
	builder := NewMessageBuilder("", nil, "")
	assert.IsType(t, &SimpleMessageBuilder{}, builder)
}

func Test_builder_factory_when_env_variable_is_set(t *testing.T) {
	t.Setenv(envVariableName, "some-path")
	builder := NewMessageBuilder("", nil, "")
	assert.IsType(t, &ExternalMessageBuilder{}, builder)
}

func Test_extract_path_and_args(t *testing.T) {
	tests := []struct {
		description  string
		input        string
		expectedPath string
		expectedArgs []string
	}{
		{
			description:  "empty value",
			input:        "",
			expectedPath: "",
			expectedArgs: []string{},
		},
		{
			description:  "command with no arguments",
			input:        "echo",
			expectedPath: "echo",
			expectedArgs: []string{},
		},
		{
			description:  "command with leading and trailing spaces",
			input:        " \t\r\necho \r \t \n",
			expectedPath: "echo",
			expectedArgs: []string{},
		},
		{
			description:  "command with one argument",
			input:        "echo hello",
			expectedPath: "echo",
			expectedArgs: []string{"hello"},
		},
		{
			description:  "command with two argument",
			input:        "echo hello world",
			expectedPath: "echo",
			expectedArgs: []string{"hello", "world"},
		},
		{
			description:  "command with arguments and extra spaces",
			input:        "\t\n echo\t \nhello\tworld\r\n",
			expectedPath: "echo",
			expectedArgs: []string{"hello", "world"},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			path, args := extractPathAndArgs(test.input)
			assert.Equal(t, test.expectedPath, path)
			assert.Equal(t, test.expectedArgs, args)
		})
	}
}

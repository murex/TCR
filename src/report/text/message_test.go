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

package text

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_wrap_unwrap_message(t *testing.T) {
	tests := []struct {
		message Message
		wrapped string
	}{
		{
			message: New("single string"),
			wrapped: "single string",
		},
		{
			message: New("two", " strings"),
			wrapped: "two strings",
		},
		{
			message: New(2, " different types"),
			wrapped: "2 different types",
		},
	}

	for _, test := range tests {
		t.Run(test.wrapped, func(t *testing.T) {
			str := test.message.ToString()
			assert.Equal(t, test.wrapped, str)
			assert.Equal(t, test.message, UnwrapMessage(str))
		})
	}
}

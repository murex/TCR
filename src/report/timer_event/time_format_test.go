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

func Test_format_duration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{duration: 0 * time.Second, expected: "0s"},
		{duration: 59 * time.Second, expected: "59s"},
		{duration: 60 * time.Second, expected: "1m"},
		{duration: 61 * time.Second, expected: "1m1s"},
		{duration: 120 * time.Second, expected: "2m"},
		{duration: 59 * time.Minute, expected: "59m"},
		{duration: 60 * time.Minute, expected: "1h"},
		{duration: 60*time.Minute + 1*time.Second, expected: "1h0m1s"},
		{duration: 61 * time.Minute, expected: "1h1m"},
		{duration: 120 * time.Minute, expected: "2h"},
	}

	for _, test := range tests {
		t.Run(test.duration.String(), func(t *testing.T) {
			assert.Equal(t, test.expected, FormatDuration(test.duration))
		})
	}
}

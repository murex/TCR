/*
Copyright (c) 2022 Murex

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

package xunit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_retrieve_xunit_test_counters(t *testing.T) {
	var parser *Parser
	testFlags := []struct {
		desc              string
		extractorFunction func() int
		expected          int
	}{
		{"total tests", func() int { return parser.getTotalTests() }, 3},
		{"total tests passed", func() int { return parser.getTotalTestsPassed() }, 1},
		{"total tests failed", func() int { return parser.getTotalTestsFailed() }, 1},
		{"total tests in error", func() int { return parser.getTotalTestsInError() }, 0},
		{"total tests skipped", func() int { return parser.getTotalTestsSkipped() }, 1},
		{"total tests run", func() int { return parser.getTotalTestsRun() }, 2},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			parser = NewParser()
			_ = parser.parse(xunitSample)
			assert.Equal(t, tt.expected, tt.extractorFunction())
		})
	}
}

func Test_retrieve_xunit_test_duration(t *testing.T) {
	parser := NewParser()
	_ = parser.parse(xunitSample)
	assert.Equal(t, sampleTotalsSuite0.Duration+sampleTotalsSuite1.Duration, parser.getTotalTestDuration())
}

func Test_parsing_invalid_data(t *testing.T) {
	testFlags := []struct {
		desc        string
		xunitData   []byte
		expectError bool
	}{
		{"no data",
			nil,
			false,
		},
		{"empty data",
			[]byte(""),
			false,
		},
		{"header only data",
			[]byte(`<?xml version="1.0" encoding="UTF-8"?>`),
			false,
		},
		{"truncated data",
			[]byte(`<?xml version="1.0" encoding="UTF-8"?><testsuites`),
			true,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			err := NewParser().parse(tt.xunitData)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

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
	"github.com/joshdk/go-junit"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var xunitSample = []byte(`
    <?xml version="1.0" encoding="UTF-8"?>
    <testsuites>
        <testsuite name="JUnitXmlReporter" errors="0" tests="0" failures="0" time="0" timestamp="2013-05-24T10:23:58" />
        <testsuite name="JUnitXmlReporter.constructor" errors="0" skipped="1" tests="3" failures="1" time="0.006" timestamp="2013-05-24T10:23:58">
            <properties>
                <property name="java.vendor" value="Sun Microsystems Inc." />
                <property name="compiler.debug" value="on" />
                <property name="project.jdk.classpath" value="jdk.classpath.1.6" />
            </properties>
            <testcase classname="JUnitXmlReporter.constructor" name="should default path to an empty string" time="0.006">
                <failure message="test failure">Assertion failed</failure>
            </testcase>
            <testcase classname="JUnitXmlReporter.constructor" name="should default consolidate to true" time="0">
                <skipped />
            </testcase>
            <testcase classname="JUnitXmlReporter.constructor" name="should default useDotNotation to true" time="0" />
        </testsuite>
    </testsuites>
`)

func Test_parse_xunit_sample(t *testing.T) {
	suites, err := junit.Ingest(xunitSample)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(suites))

	assert.Equal(t, junit.Totals{Tests: 0, Passed: 0, Skipped: 0, Failed: 0, Error: 0, Duration: 0}, suites[0].Totals)
	assert.Equal(t, 0, len(suites[0].Tests))

	assert.Equal(t, junit.Totals{Tests: 3, Passed: 1, Skipped: 1, Failed: 1, Error: 0, Duration: 6 * time.Millisecond}, suites[1].Totals)
	assert.Equal(t, 3, len(suites[1].Tests))
	assert.Equal(t, junit.StatusFailed, suites[1].Tests[0].Status)
	assert.Equal(t, junit.StatusSkipped, suites[1].Tests[1].Status)
	assert.Equal(t, junit.StatusPassed, suites[1].Tests[2].Status)
}

func Test_retrieve_xunit_test_stats(t *testing.T) {
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

func Test_parsing_invalid_data(t *testing.T) {
	var parser *Parser
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
			parser = NewParser()
			err := parser.parse(tt.xunitData)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

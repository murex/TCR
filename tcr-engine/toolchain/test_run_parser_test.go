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

package toolchain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_extracting_a_test_info_from_a_maven_test_build_output(t *testing.T) {
	testFlags := []struct {
		desc            string
		mavenTestOutput string
		expected        TestResults
	}{
		{
			"legitimate test output",
			"[INFO] Results:\r" +
				"[INFO]\r" +
				"[WARNING] Tests run: 26, Failures: 1, Errors: 3, Skipped: 4\r" +
				"[INFO]\r",
			NewTestResults(26, 18, 1, 4, 3, 0),
		},
		{
			"another legitimate test output",
			"[INFO] Results:\r" +
				"[INFO]\r" +
				"[WARNING] Tests run: 30, Failures: 5, Errors: 4, Skipped: 2\r" +
				"[INFO]\r",
			NewTestResults(30, 19, 5, 2, 4, 0),
		},
		{
			"it takes the report line in the results section of a legitimate output",
			"[INFO] Tests run: 5, Failures: 1, Errors: 0, Skipped: 0, Time elapsed: 0.1 s - in com.tcr\r" +
				"[INFO]\r" +
				"[INFO] Results:\r" +
				"[INFO]\r" +
				"[WARNING] Tests run: 30, Failures: 5, Errors: 4, Skipped: 2\r" +
				"[INFO]\r",
			NewTestResults(30, 19, 5, 2, 4, 0),
		},
		{
			"incomplete test output returns a default TestResults object",
			"[WARNING] Tests run: 26, Failures: 0",
			NewTestResults(0, 0, 0, 0, 0, 0),
		},
		{
			"an empty build output return a default TestResults object",
			"",
			NewTestResults(0, 0, 0, 0, 0, 0),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			info := extractTestResults(tt.mavenTestOutput)
			assert.Equal(t, tt.expected, info)
		})
	}
}

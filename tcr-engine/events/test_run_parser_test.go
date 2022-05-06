package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_extracting_a_test_info_from_a_maven_test_build_output(t *testing.T) {
	testFlags := []struct {
		desc            string
		mavenTestOutput string
		expected        TestRunInformation
	}{
		{
			"legitimate test output",
			"[INFO] Results:\r" +
				"[INFO]\r" +
				"[WARNING] Tests run: 26, Failures: 1, Errors: 3, Skipped: 4\r" +
				"[INFO]\r",
			NewTestRunInformation(26, 18, 1, 4, 3),
		},
		{
			"another legitimate test output",
			"[INFO] Results:\r" +
				"[INFO]\r" +
				"[WARNING] Tests run: 30, Failures: 5, Errors: 4, Skipped: 2\r" +
				"[INFO]\r",
			NewTestRunInformation(30, 19, 5, 2, 4),
		},
		{
			"it takes the report line in the results section of a legitimate output",
			"[INFO] Tests run: 5, Failures: 1, Errors: 0, Skipped: 0, Time elapsed: 0.1 s - in com.tcr\r" +
				"[INFO]\r" +
				"[INFO] Results:\r" +
				"[INFO]\r" +
				"[WARNING] Tests run: 30, Failures: 5, Errors: 4, Skipped: 2\r" +
				"[INFO]\r",
			NewTestRunInformation(30, 19, 5, 2, 4),
		},
		{
			"incomplete test output returns a default TestRunInformation object",
			"[WARNING] Tests run: 26, Failures: 0",
			NewTestRunInformation(0, 0, 0, 0, 0),
		},
		{
			"an empty build output return a default TestRunInformation object",
			"",
			NewTestRunInformation(0, 0, 0, 0, 0),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			info := ExtractTestRunInformation(tt.mavenTestOutput)
			assert.Equal(t, tt.expected, info)
		})
	}
}

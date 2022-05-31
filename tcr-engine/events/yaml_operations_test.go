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

package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	emptyTcrEvent    = *ATcrEvent(WithTimestamp(time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)))
	sampleDate       = time.Date(2022, 4, 11, 15, 52, 3, 0, time.UTC)
	sampleDateString = "2022-04-11T15:52:03Z" // YAML uses ISO-8601 standard to express dates
	zeroDateString   = "0001-01-01T00:00:00Z"
)

func Test_convert_a_yaml_string_to_a_tcr_event(t *testing.T) {
	testFlags := []struct {
		desc       string
		yamlString string
		expected   TcrEvent
	}{
		{
			"timestamp in UTC",
			buildYamlString(sampleDateString, "0", "0", "0", "0", "0", "0", "0", "0", "0"),
			*ATcrEvent(WithTimestamp(sampleDate)),
		},
		{
			"modified source lines",
			buildYamlString(sampleDateString, "2", "0", "0", "0", "0", "0", "0", "0", "0"),
			*ATcrEvent(WithTimestamp(sampleDate), WithModifiedSrcLines(2)),
		},
		{
			"modified test lines",
			buildYamlString(sampleDateString, "0", "3", "0", "0", "0", "0", "0", "0", "0"),
			*ATcrEvent(WithTimestamp(sampleDate), WithModifiedTestLines(3)),
		},
		{
			"with build passed",
			buildYamlString(sampleDateString, "0", "0", "1", "0", "0", "0", "0", "0", "0"),
			*ATcrEvent(WithTimestamp(sampleDate), WithPassingBuild()),
		},
		{
			"with build failed",
			buildYamlString(sampleDateString, "0", "0", "2", "0", "0", "0", "0", "0", "0"),
			*ATcrEvent(WithTimestamp(sampleDate), WithFailingBuild()),
		},
		{
			"with tests passed",
			buildYamlString(sampleDateString, "0", "0", "0", "1", "0", "0", "0", "0", "0"),
			*ATcrEvent(WithTimestamp(sampleDate), WithPassingTests()),
		},
		{
			"with tests failed",
			buildYamlString(sampleDateString, "0", "0", "0", "2", "0", "0", "0", "0", "0"),
			*ATcrEvent(WithTimestamp(sampleDate), WithFailingTests()),
		},
		{
			"total test cases run",
			buildYamlString(sampleDateString, "0", "0", "0", "0", "4", "0", "0", "0", "0"),
			*ATcrEvent(WithTimestamp(sampleDate), WithTotalTestsRun(4)),
		},
		{
			"passed test cases",
			buildYamlString(sampleDateString, "0", "0", "0", "0", "0", "3", "0", "0", "0"),
			*ATcrEvent(WithTimestamp(sampleDate), WithTestsPassed(3)),
		},
		{
			"failed test cases",
			buildYamlString(sampleDateString, "0", "0", "0", "0", "0", "0", "2", "0", "0"),
			*ATcrEvent(WithTimestamp(sampleDate), WithTestsFailed(2)),
		},
		{
			"skipped test cases",
			buildYamlString(sampleDateString, "0", "0", "0", "0", "0", "0", "0", "5", "0"),
			*ATcrEvent(WithTimestamp(sampleDate), WithTestsSkipped(5)),
		},
		{
			"test cases with errors",
			buildYamlString(sampleDateString, "0", "0", "0", "0", "0", "0", "0", "0", "4"),
			*ATcrEvent(WithTimestamp(sampleDate), WithTestsWithErrors(4)),
		},
		{
			"empty yaml string",
			"",
			emptyTcrEvent,
		},
		{
			"yaml with empty values",
			buildYamlString("", "", "", "", "", "", "", "", "", ""),
			emptyTcrEvent,
		},
		{
			"yaml with invalid timestamp value",
			buildYamlString("wrong timestamp", "", "", "", "", "", "", "", "", ""),
			emptyTcrEvent,
		},
		{
			"yaml with invalid int value",
			buildYamlString("", "wrong", "", "", "", "", "", "", "", ""),
			emptyTcrEvent,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			yamlRecord := yamlToTcrEvent(tt.yamlString)
			assert.Equal(t, tt.expected, yamlRecord)
		})
	}
}

func Test_convert_a_tcr_event_to_a_yaml_string(t *testing.T) {
	testFlags := []struct {
		desc     string
		event    TcrEvent
		expected string
	}{
		{
			"timestamp in UTC",
			*ATcrEvent(WithTimestamp(sampleDate)),
			buildYamlString(sampleDateString, "0", "0", "0", "0", "0", "0", "0", "0", "0"),
		},
		{
			"modified source lines",
			*ATcrEvent(WithTimestamp(sampleDate), WithModifiedSrcLines(2)),
			buildYamlString(sampleDateString, "2", "0", "0", "0", "0", "0", "0", "0", "0"),
		},
		{
			"modified test lines",
			*ATcrEvent(WithTimestamp(sampleDate), WithModifiedTestLines(3)),
			buildYamlString(sampleDateString, "0", "3", "0", "0", "0", "0", "0", "0", "0"),
		},
		{
			"with build passed",
			*ATcrEvent(WithTimestamp(sampleDate), WithPassingBuild()),
			buildYamlString(sampleDateString, "0", "0", "1", "0", "0", "0", "0", "0", "0"),
		},
		{
			"with build failed",
			*ATcrEvent(WithTimestamp(sampleDate), WithFailingBuild()),
			buildYamlString(sampleDateString, "0", "0", "2", "0", "0", "0", "0", "0", "0"),
		},
		{
			"with tests passed",
			*ATcrEvent(WithTimestamp(sampleDate), WithPassingTests()),
			buildYamlString(sampleDateString, "0", "0", "0", "1", "0", "0", "0", "0", "0"),
		},
		{
			"with tests failed",
			*ATcrEvent(WithTimestamp(sampleDate), WithFailingTests()),
			buildYamlString(sampleDateString, "0", "0", "0", "2", "0", "0", "0", "0", "0"),
		},
		{
			"total test cases run",
			*ATcrEvent(WithTimestamp(sampleDate), WithTotalTestsRun(4)),
			buildYamlString(sampleDateString, "0", "0", "0", "0", "4", "0", "0", "0", "0"),
		},
		{
			"passed test cases",
			*ATcrEvent(WithTimestamp(sampleDate), WithTestsPassed(3)),
			buildYamlString(sampleDateString, "0", "0", "0", "0", "0", "3", "0", "0", "0"),
		},
		{
			"failed test cases",
			*ATcrEvent(WithTimestamp(sampleDate), WithTestsFailed(2)),
			buildYamlString(sampleDateString, "0", "0", "0", "0", "0", "0", "2", "0", "0"),
		},
		{
			"skipped test cases",
			*ATcrEvent(WithTimestamp(sampleDate), WithTestsSkipped(5)),
			buildYamlString(sampleDateString, "0", "0", "0", "0", "0", "0", "0", "5", "0"),
		},
		{
			"test cases with errors",
			*ATcrEvent(WithTimestamp(sampleDate), WithTestsWithErrors(4)),
			buildYamlString(sampleDateString, "0", "0", "0", "0", "0", "0", "0", "0", "4"),
		},
		{
			"empty TCR event",
			emptyTcrEvent,
			buildYamlString(zeroDateString, "0", "0", "0", "0", "0", "0", "0", "0", "0"),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			yaml := tcrEventToYaml(tt.event)
			assert.Equal(t, tt.expected, yaml)
		})
	}
}

func buildYamlString(timestamp, srcLines, testLines, buildStatus, testStatus, totalTests, testsPassed, testsFailed, testsSkipped, testsWithErrors string) string {
	return buildYamlLine("timestamp", timestamp) +
		buildYamlLine("modified-src-lines", srcLines) +
		buildYamlLine("modified-test-lines", testLines) +
		buildYamlLine("build-status", buildStatus) +
		buildYamlLine("tests-status", testStatus) +
		buildYamlLine("total-tests-run", totalTests) +
		buildYamlLine("tests-passed", testsPassed) +
		buildYamlLine("tests-failed", testsFailed) +
		buildYamlLine("tests-skipped", testsSkipped) +
		buildYamlLine("tests-with-errors", testsWithErrors)
}

func buildYamlLine(key, value string) string {
	return key + ": " + value + "\n"
}

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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_convert_a_yaml_string_to_a_tcr_event(t *testing.T) {
	testFlags := []struct {
		desc       string
		yamlString string
		expected   TCREvent
	}{
		{
			"modified source lines",
			buildYAMLString("2", "0", "0", "0", "0", "0", "0", "0s"),
			*ATcrEvent(WithModifiedSrcLines(2)),
		},
		{
			"modified test lines",
			buildYAMLString("0", "3", "0", "0", "0", "0", "0", "0s"),
			*ATcrEvent(WithModifiedTestLines(3)),
		},
		{
			"total test cases run",
			buildYAMLString("0", "0", "4", "0", "0", "0", "0", "0s"),
			*ATcrEvent(WithTotalTestsRun(4)),
		},
		{
			"passed test cases",
			buildYAMLString("0", "0", "0", "3", "0", "0", "0", "0s"),
			*ATcrEvent(WithTestsPassed(3)),
		},
		{
			"failed test cases",
			buildYAMLString("0", "0", "0", "0", "2", "0", "0", "0s"),
			*ATcrEvent(WithTestsFailed(2)),
		},
		{
			"skipped test cases",
			buildYAMLString("0", "0", "0", "0", "0", "5", "0", "0s"),
			*ATcrEvent(WithTestsSkipped(5)),
		},
		{
			"test cases with errors",
			buildYAMLString("0", "0", "0", "0", "0", "0", "4", "0s"),
			*ATcrEvent(WithTestsWithErrors(4)),
		},
		{
			"test duration",
			buildYAMLString("0", "0", "0", "0", "0", "0", "0", "20s"),
			*ATcrEvent(WithTestsDuration(20 * time.Second)),
		},
		{
			"empty yaml string",
			"",
			*ATcrEvent(),
		},
		{
			"yaml with empty values",
			buildYAMLString("", "", "", "", "", "", "", ""),
			*ATcrEvent(),
		},
		{
			"yaml with invalid timestamp value",
			buildYAMLString("", "", "", "", "", "", "", ""),
			*ATcrEvent(),
		},
		{
			"yaml with invalid int value",
			buildYAMLString("wrong", "", "", "", "", "", "", ""),
			*ATcrEvent(),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			yamlRecord := yamlToTCREvent(tt.yamlString)
			assert.Equal(t, tt.expected, yamlRecord)
		})
	}
}

func Test_convert_a_tcr_event_to_a_yaml_string(t *testing.T) {
	testFlags := []struct {
		desc     string
		event    TCREvent
		expected string
	}{
		{
			"modified source lines",
			*ATcrEvent(WithModifiedSrcLines(2)),
			buildYAMLString("2", "0", "0", "0", "0", "0", "0", "0s"),
		},
		{
			"modified test lines",
			*ATcrEvent(WithModifiedTestLines(3)),
			buildYAMLString("0", "3", "0", "0", "0", "0", "0", "0s"),
		},
		{
			"total test cases run",
			*ATcrEvent(WithTotalTestsRun(4)),
			buildYAMLString("0", "0", "4", "0", "0", "0", "0", "0s"),
		},
		{
			"passed test cases",
			*ATcrEvent(WithTestsPassed(3)),
			buildYAMLString("0", "0", "0", "3", "0", "0", "0", "0s"),
		},
		{
			"failed test cases",
			*ATcrEvent(WithTestsFailed(2)),
			buildYAMLString("0", "0", "0", "0", "2", "0", "0", "0s"),
		},
		{
			"skipped test cases",
			*ATcrEvent(WithTestsSkipped(5)),
			buildYAMLString("0", "0", "0", "0", "0", "5", "0", "0s"),
		},
		{
			"test cases with errors",
			*ATcrEvent(WithTestsWithErrors(4)),
			buildYAMLString("0", "0", "0", "0", "0", "0", "4", "0s"),
		},
		{
			"test duration",
			*ATcrEvent(WithTestsDuration(20 * time.Second)),
			buildYAMLString("0", "0", "0", "0", "0", "0", "0", "20s"),
		},
		{
			"empty TCR event",
			*ATcrEvent(),
			buildYAMLString("0", "0", "0", "0", "0", "0", "0", "0s"),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			yaml := tcrEventToYAML(tt.event)
			assert.Equal(t, tt.expected, yaml)
		})
	}
}

func buildYAMLString(srcLines, testLines, totalTests, testsPassed, testsFailed, testsSkipped, testsWithErrors, testDuration string) string {
	return buildYAMLSectionLine("changed-lines") +
		buildYAMLKeyValueLine("src", srcLines) +
		buildYAMLKeyValueLine("test", testLines) +
		buildYAMLSectionLine("test-stats") +
		buildYAMLKeyValueLine("run", totalTests) +
		buildYAMLKeyValueLine("passed", testsPassed) +
		buildYAMLKeyValueLine("failed", testsFailed) +
		buildYAMLKeyValueLine("skipped", testsSkipped) +
		buildYAMLKeyValueLine("error", testsWithErrors) +
		buildYAMLKeyValueLine("duration", testDuration)
}

func buildYAMLSectionLine(section string) string {
	return section + ":" + "\n"
}

func buildYAMLKeyValueLine(key, value string) string {
	return "    " + key + ": " + value + "\n"
}

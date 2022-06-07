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
	"regexp"
	"strconv"
	"strings"
)

const regexBuildOutput = `Tests run: [0-9]+, Failures: [0-9]+, Errors: [0-9]+, Skipped: [0-9]+[\r\n]+`
const regexTestsRun = "Tests run: [0-9]+"
const regexFailures = "Failures: [0-9]+"
const regexErrors = "Errors: [0-9]+"
const regexSkipped = "Skipped: [0-9]+"

func extractTestResults(mvnTestRunOutput string) TestResults {
	testResult := extractTestInfo(mvnTestRunOutput)
	totalTests := extractSectionInfo(testResult, regexTestsRun)
	testsFailed := extractSectionInfo(testResult, regexFailures)
	testsSkipped := extractSectionInfo(testResult, regexSkipped)
	testsWithErrors := extractSectionInfo(testResult, regexErrors)
	testsPassed := totalTests - (testsFailed + testsSkipped + testsWithErrors)
	return NewTestResults(totalTests, testsPassed, testsFailed, testsSkipped, testsWithErrors, 0)
}

func extractTestInfo(testRunOutput string) string {
	re := regexp.MustCompile(regexBuildOutput)
	return re.FindString(testRunOutput)
}

func extractSectionInfo(testInfo, regex string) int {
	return extractNumberFrom(extractMatchingString(testInfo, regex))
}

func extractMatchingString(input, regex string) string {
	testsRunRegEx := regexp.MustCompile(regex)
	return testsRunRegEx.FindString(input)
}

func extractNumberFrom(testRunResult string) (num int) {
	txtSplit := strings.Split(testRunResult, ":")
	if len(txtSplit) == 2 {
		num, _ = strconv.Atoi(strings.TrimSpace(txtSplit[1]))
	}
	return
}

package events

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

// ExtractTestResults extracts an instance of TestResults from the buildOutput
func ExtractTestResults(mvnTestRunOutput string) TestResults {
	testResult := extractTestInfo(mvnTestRunOutput)
	totalTests := extractSectionInfo(testResult, regexTestsRun)
	testsFailed := extractSectionInfo(testResult, regexFailures)
	testsSkipped := extractSectionInfo(testResult, regexSkipped)
	testsWithErrors := extractSectionInfo(testResult, regexErrors)
	testsPassed := totalTests - (testsFailed + testsSkipped + testsWithErrors)
	return NewTestResults(totalTests, testsPassed, testsFailed, testsSkipped, testsWithErrors)
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

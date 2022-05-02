package events

// TestRunInformation is the structure containing information of the test run
type TestRunInformation struct {
	TotalTestsRan   int
	TestsPassed     int
	TestsFailed     int
	TestsSkipped    int
	TestsWithErrors int
}

// NewTestRunInformation create a new
func NewTestRunInformation(totalTestsRan, testsPassed,
	testsFailed, testsSkipped, testsWithErrors int) TestRunInformation {
	return TestRunInformation{
		TotalTestsRan:   totalTestsRan,
		TestsFailed:     testsFailed,
		TestsPassed:     testsPassed,
		TestsSkipped:    testsSkipped,
		TestsWithErrors: testsWithErrors,
	}
}

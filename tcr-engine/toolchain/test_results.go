package toolchain

// TestResults is the structure containing information of the test run
type TestResults struct {
	TotalRun   int
	Passed     int
	Failed     int
	Skipped    int
	WithErrors int
}

// NewTestResults create a new instance of the TestResults class
func NewTestResults(totalRun, passed, failed, skipped, withErrors int) TestResults {
	return TestResults{
		TotalRun:   totalRun,
		Failed:     failed,
		Passed:     passed,
		Skipped:    skipped,
		WithErrors: withErrors,
	}
}

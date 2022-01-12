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

package engine

// Status is used for describing the current TCR engine status
type Status struct {
	rc int
}

// List of TCR engine possible status values
var (
	StatusOk          = Status{rc: 0} // Build and Test Passed and changes were committed with no error
	StatusBuildFailed = Status{rc: 1} // Build failed
	StatusTestFailed  = Status{rc: 2} // Build passed, one or more test failed, and changes were reverted
	StatusConfigError = Status{rc: 3} // Error in configuration or parameters
	StatusGitError    = Status{rc: 4} // Git error
	StatusOtherError  = Status{rc: 5} // Any other error
)

var currentState Status

func recordState(state Status) {
	currentState = state
}

func getCurrentState() Status {
	return currentState
}

func getReturnCode() int {
	return currentState.rc
}

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

package status

// Status is used for describing the current TCR engine status
type Status struct {
	rc int
}

// List of TCR engine possible status values
var (
	Ok          = NewStatus(0) // Build and Test Passed and changes were committed with no error
	BuildFailed = NewStatus(1) // Build failed
	TestFailed  = NewStatus(2) // Build passed, one or more test failed, and changes were reverted
	ConfigError = NewStatus(3) // Error in configuration or parameters
	GitError    = NewStatus(4) // Git error
	OtherError  = NewStatus(5) // Any other error
)

var currentState Status

// NewStatus creates a new application status
func NewStatus(rc int) Status {
	return Status{rc: rc}
}

// RecordState records the state to be returned as return code by the application
func RecordState(state Status) {
	currentState = state
}

// GetCurrentState returns the current application state
func GetCurrentState() Status {
	return currentState
}

// GetReturnCode provides the application's return code
func GetReturnCode() int {
	return GetCurrentState().rc
}

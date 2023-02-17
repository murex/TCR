/*
Copyright (c) 2023 Murex

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

package model

import "github.com/murex/tcr/status"

// CheckStatus provides the return status for a specific check
type CheckStatus int

// Check status values
const (
	CheckStatusOk      CheckStatus = 0 // Check status is OK
	CheckStatusWarning CheckStatus = 1 // Check status is Warning
	CheckStatusError   CheckStatus = 2 // Build status is Error
)

// UpdateReturnState updates the application's return state according to CheckGroup's status
func UpdateReturnState(results *CheckGroup) {
	if int(results.GetStatus()) > status.GetReturnCode() {
		RecordCheckState(results.GetStatus())
	}
}

// RecordCheckState sets the application's return state with the provided check status
func RecordCheckState(s CheckStatus) {
	status.RecordState(status.NewStatus(int(s)))
}

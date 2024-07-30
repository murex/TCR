/*
Copyright (c) 2024 Murex

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

package command

import (
	"fmt"
	"github.com/murex/tcr/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_command_result_outcome(t *testing.T) {
	testCases := []struct {
		status         Status
		expectedPassed bool
		expectedFailed bool
	}{
		{"pass", true, false},
		{"fail", false, true},
		{"unknown", false, false},
	}
	for _, tt := range testCases {
		t.Run(fmt.Sprint(tt.status, "_status"), func(t *testing.T) {
			result := Result{Status: tt.status}
			assert.Equal(t, tt.expectedPassed, result.Passed())
			assert.Equal(t, tt.expectedFailed, result.Failed())
		})
	}
}

func Test_run_command(t *testing.T) {
	testCases := []struct {
		desc           string
		command        Command
		expectedStatus Status
	}{
		{
			"unknown command",
			Command{Path: "unknown-command"},
			StatusFail,
		},
		{
			"passing command",
			Command{Path: "true"},
			StatusPass,
		},
		{
			"failing command",
			Command{Path: "false"},
			StatusFail,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			result := GetRunner().Run("", &tt.command)
			assert.Equal(t, tt.expectedStatus, result.Status)
		})
	}
}

func Test_abort_command(t *testing.T) {
	// this test fails randomly on Windows for an unexplained reason.
	// This seems to be related to command.Process never being set when running a command,
	// but no explanation why this happens on Windows while this works as expected on Unix OS's.
	utils.SkipOnWindows(t)

	testCases := []struct {
		desc     string
		command  *Command
		expected bool
	}{
		{
			"with no running command",
			nil,
			false,
		},
		{
			"with running command",
			&Command{Path: "sleep", Arguments: []string{"5"}},
			true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.command != nil {
				go GetRunner().Run("", tt.command)
				// wait for the new command process to start
				for GetRunner().command == nil {
					time.Sleep(10 * time.Millisecond)
				}
			}
			result := GetRunner().AbortRunningCommand()
			assert.Equal(t, tt.expected, result)
		})
	}
}

//go:build test_helper

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

package utils

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
)

// SlowTestTag is a test utility function for marking tests that take a long time to be executed.
// When added at the beginning of a test, the corresponding test is skipped when running tests with -short flag
func SlowTestTag(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

// AssertSimpleTrace is a utility function to assert simple trace messages
func AssertSimpleTrace(t *testing.T, expected []string, operation func()) {
	t.Helper()
	var output bytes.Buffer
	SetSimpleTrace(&output)
	operation()
	var expectedWithWrapping string
	for _, line := range expected {
		expectedWithWrapping += "[TCR] " + line + "\n"
	}
	assert.Equal(t, expectedWithWrapping, output.String())
}

// SkipOnWindowsCI allows to prevent running a test when on Windows CI when called at the beginning of the test
func SkipOnWindowsCI(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" && runtime.GOOS == "windows" {
		t.Skip("test skipped on windows CI")
	}
}

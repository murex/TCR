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

// Package built_in_test_data provides test fixture data for TCR's built-in
// toolchains, used by toolchain tests to verify build and test command configuration.
package built_in_test_data //nolint:revive

// BuiltInTestData provides test data for a built-in toolchain.
type BuiltInTestData struct {
	Name             string
	BuildCommandPath string
	BuildCommandArgs []string
	TestCommandPath  string
	TestCommandArgs  []string
	TestResultDir    string
}

// BuiltInTests is the slice containing all built-in toolchain test data.
var BuiltInTests []BuiltInTestData

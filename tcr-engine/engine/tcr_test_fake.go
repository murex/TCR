//go:build test_helper

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

// FakeTcrEngine is a TCR engine fake. Used mainly for testing peripheral packages
// such as cli.
type FakeTcrEngine struct {
	TcrEngine
}

// NewFakeTcrEngine creates a FakeToolchain instance
func NewFakeTcrEngine() *FakeTcrEngine {
	return &FakeTcrEngine{}
}

// GetSessionInfo returns a SessionInfo struct filled with "fake" values
func (fake *FakeTcrEngine) GetSessionInfo() SessionInfo {
	return SessionInfo{
		BaseDir:       "fake",
		WorkDir:       "fake",
		LanguageName:  "fake",
		ToolchainName: "fake",
		AutoPush:      false,
		BranchName:    "fake",
	}
}

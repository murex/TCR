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

package filesystem

// FakeSourceTree is a fake implementation Source Tree interface
type FakeSourceTree struct {
	baseDir string
}

// NewFakeSourceTree provides a fake source tree instance
func NewFakeSourceTree(dir string) *FakeSourceTree {
	return &FakeSourceTree{dir}
}

// GetBaseDir returns the base directory for the source tree instance
func (fst FakeSourceTree) GetBaseDir() string {
	return fst.baseDir
}

// IsValid indicates that the source tree instance is valid (forced to true)
func (fst FakeSourceTree) IsValid() bool {
	return true
}

// Watch is a fake implementation of Watch command (not usable as is)
func (fst FakeSourceTree) Watch(_ []string, _ func(filename string) bool, _ <-chan bool) bool {
	fakeChannel := make(chan bool)
	return <-fakeChannel
}

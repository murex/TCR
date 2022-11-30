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

package vcs

type (
	// FileDiff is a structure containing diff information for a file
	FileDiff struct {
		Path         string
		addedLines   int
		removedLines int
	}

	// FileDiffs contains a set of FileDiff data in a slice
	FileDiffs []FileDiff
)

// NewFileDiff creates a new instance of FileDiff
func NewFileDiff(filename string, added int, removed int) FileDiff {
	return FileDiff{Path: filename, addedLines: added, removedLines: removed}
}

// ChangedLines returns the number of changed lines for this file
func (fd FileDiff) ChangedLines() int {
	return fd.addedLines + fd.removedLines
}

// ChangedLines returns the total number of changed lines for this set of files
func (diffs FileDiffs) ChangedLines(predicate func(filepath string) bool) (changes int) {
	for _, fd := range diffs {
		if predicate == nil || predicate(fd.Path) {
			changes += fd.ChangedLines()
		}
	}
	return changes
}

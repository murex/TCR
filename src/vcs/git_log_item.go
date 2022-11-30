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

import (
	"sort"
	"time"
)

type (
	// GitLogItem contains git log information for a commit
	GitLogItem struct {
		Hash      string
		Timestamp time.Time
		Message   string
	}

	// GitLogItems contains a set of git log items in a slice
	GitLogItems []GitLogItem
)

// NewGitLogItem creates a new git log item instance
func NewGitLogItem(hash string, timestamp time.Time, message string) GitLogItem {
	return GitLogItem{hash, timestamp, message}
}

func (items *GitLogItems) sortByDate() {
	sort.Slice(*items, func(i, j int) bool {
		return (*items)[i].Timestamp.Before((*items)[j].Timestamp)
	})
}

func (items *GitLogItems) add(d GitLogItem) {
	*items = append(*items, d)
}

func (items *GitLogItems) len() int {
	return len(*items)
}

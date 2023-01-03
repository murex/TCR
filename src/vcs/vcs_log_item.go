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
	// LogItem contains VCS log information for a commit
	LogItem struct {
		Hash      string
		Timestamp time.Time
		Message   string
	}

	// LogItems contains a set of VCS log items in a slice
	LogItems []LogItem
)

// NewLogItem creates a new VCS log item instance
func NewLogItem(hash string, timestamp time.Time, message string) LogItem {
	return LogItem{hash, timestamp, message}
}

func (items *LogItems) sortByDate() {
	sort.Slice(*items, func(i, j int) bool {
		return (*items)[i].Timestamp.Before((*items)[j].Timestamp)
	})
}

// Add adds a LogItem to the LogItems collection
func (items *LogItems) Add(d LogItem) {
	*items = append(*items, d)
}

// Len returns the length of the items array
func (items *LogItems) Len() int {
	return len(*items)
}

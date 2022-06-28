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
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_add_regular_git_log_item(t *testing.T) {
	var items GitLogItems
	item := NewGitLogItem("xxx", time.Now(), "some message")
	items.add(item)
	assert.Len(t, items, 1)
	assert.Contains(t, items, item)
}

func Test_add_empty_git_log_item(t *testing.T) {
	var items GitLogItems
	item := NewGitLogItem("", time.Time{}, "")
	items.add(item)
	assert.Len(t, items, 1)
}

func Test_sort_empty_git_log_items_does_not_panic(t *testing.T) {
	var items GitLogItems
	assert.NotPanics(t, func() {
		items.sort()
	})
	assert.Len(t, items, 0)
}

func Test_sort_already_sorted_git_log_items(t *testing.T) {
	item1 := NewGitLogItem("xxx1", time.Now(), "first commit")
	item2 := NewGitLogItem("xxx2", time.Now().Add(1*time.Second), "second commit")
	items := GitLogItems{item1, item2}
	items.sort()
	assert.Equal(t, GitLogItems{item1, item2}, items)
}

func Test_sort_unsorted_git_log_items(t *testing.T) {
	item1 := NewGitLogItem("xxx1", time.Now(), "first commit")
	item2 := NewGitLogItem("xxx2", time.Now().Add(1*time.Second), "second commit")
	items := GitLogItems{item2, item1}
	items.sort()
	assert.Equal(t, GitLogItems{item1, item2}, items)
}

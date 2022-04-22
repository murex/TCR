/*
Copyright (c) 2021 Murex

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
	"path/filepath"
	"testing"
)

// push-enabling flag

func Test_git_auto_push_is_disabled_default(t *testing.T) {
	git, _ := New(".")
	assert.Zero(t, git.IsPushEnabled())
}

func Test_git_enable_disable_push(t *testing.T) {
	git, _ := New(".")
	git.EnablePush(true)
	assert.NotZero(t, git.IsPushEnabled())
	git.EnablePush(false)
	assert.Zero(t, git.IsPushEnabled())
}

// Working Branch

func Test_init_fails_when_working_dir_is_not_in_a_git_repo(t *testing.T) {
	git, err := New("/")
	assert.Zero(t, git)
	assert.NotZero(t, err)
}

func Test_can_retrieve_working_branch(t *testing.T) {
	git, _ := New(".")
	assert.NotZero(t, git.GetWorkingBranch())
}

func Test_git_diff_command(t *testing.T) {
	testFlags := []struct {
		desc          string
		gitDiffOutput string
		expected      []FileDiff
	}{
		{"0 file changed",
			"",
			nil,
		},
		{"1 file changed",
			"1\t1\tsome-file.txt\n",
			[]FileDiff{
				{"some-file.txt", 1, 1},
			},
		},
		{"2 files changed",
			"1\t1\tfile1.txt\n1\t1\tfile2.txt\n",
			[]FileDiff{
				{"file1.txt", 1, 1},
				{"file2.txt", 1, 1},
			},
		},
		{"file changed in sub-directory",
			"1\t1\tsome-dir/some-file.txt\n",
			[]FileDiff{
				{filepath.Join("some-dir", "some-file.txt"), 1, 1},
			},
		},
		{"1 file changed with added lines only",
			"15\t0\tsome-file.txt\n",
			[]FileDiff{
				{"some-file.txt", 15, 0},
			},
		},
		{"1 file changed with removed lines only",
			"0\t7\tsome-file.txt\n",
			[]FileDiff{
				{"some-file.txt", 0, 7},
			},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			git := &GitImpl{
				// git command calls are faked
				runGitFunction: func(_ []string) (output []byte, err error) {
					return []byte(tt.gitDiffOutput), nil
				},
			}
			fileDiffs, _ := git.Diff()
			assert.Equal(t, tt.expected, fileDiffs)
		})
	}
}

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
	"errors"
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
	assert.Error(t, err)
}

func Test_can_retrieve_working_branch(t *testing.T) {
	git, _ := New(".")
	assert.NotZero(t, git.GetWorkingBranch())
}

func Test_git_diff_command(t *testing.T) {
	testFlags := []struct {
		desc          string
		gitDiffOutput string
		gitDiffError  error
		expectError   bool
		expectedArgs  []string
		expectedDiff  []FileDiff
	}{
		{"git command arguments",
			"",
			errors.New("git diff error"),
			true,
			[]string{
				"diff", "--numstat", "--ignore-cr-at-eol", "--ignore-all-space", "--ignore-blank-lines", "HEAD"},
			nil,
		},
		{"git error",
			"",
			errors.New("git diff error"),
			true,
			nil,
			nil,
		},
		{"0 file changed",
			"",
			nil,
			false,
			nil,
			nil,
		},
		{"1 file changed",
			"1\t1\tsome-file.txt\n",
			nil,
			false,
			nil,
			[]FileDiff{
				{"some-file.txt", 1, 1},
			},
		},
		{"2 files changed",
			"1\t1\tfile1.txt\n1\t1\tfile2.txt\n",
			nil,
			false,
			nil,
			[]FileDiff{
				{"file1.txt", 1, 1},
				{"file2.txt", 1, 1},
			},
		},
		{"file changed in sub-directory",
			"1\t1\tsome-dir/some-file.txt\n",
			nil,
			false,
			nil,
			[]FileDiff{
				{filepath.Join("some-dir", "some-file.txt"), 1, 1},
			},
		},
		{"1 file changed with added lines only",
			"15\t0\tsome-file.txt\n",
			nil,
			false,
			nil,
			[]FileDiff{
				{"some-file.txt", 15, 0},
			},
		},
		{"1 file changed with removed lines only",
			"0\t7\tsome-file.txt\n",
			nil,
			false,
			nil,
			[]FileDiff{
				{"some-file.txt", 0, 7},
			},
		},
		{"noise in output trace",
			"warning: LF will be replaced by CRLF in some-file.txt.\n" +
				"The file will have its original line endings in your working directory\n" +
				"1\t1\tsome-file.txt\n",
			nil,
			false,
			nil,
			[]FileDiff{
				{"some-file.txt", 1, 1},
			},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			git := &GitImpl{
				// git command calls are faked
				runGitFunction: func(args []string) (output []byte, err error) {
					actualArgs = args[2:]
					return []byte(tt.gitDiffOutput), tt.gitDiffError
				},
			}
			fileDiffs, err := git.Diff()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.expectedArgs != nil {
				assert.Equal(t, tt.expectedArgs, actualArgs)
			}
			assert.Equal(t, tt.expectedDiff, fileDiffs)
		})
	}
}

func Test_git_push_command(t *testing.T) {
	testFlags := []struct {
		desc                 string
		pushEnabled          bool
		gitError             error
		expectError          bool
		expectBranchOnRemote bool
	}{
		{
			"push enabled and no git error",
			true,
			nil,
			false,
			true,
		},
		{
			"push enabled and git error",
			true,
			errors.New("git push error"),
			true,
			false,
		},
		{
			"push disabled and no git error",
			false,
			nil,
			false,
			false,
		},
		{
			"push disabled and git error",
			false,
			errors.New("git push error"),
			false,
			false,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			git := &GitImpl{
				// git command calls are faked
				traceGitFunction: func(_ []string) (err error) {
					return tt.gitError
				},
				pushEnabled: tt.pushEnabled,
			}
			err := git.Push()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectBranchOnRemote, git.workingBranchExistsOnRemote)
		})
	}
}

func Test_git_pull_command(t *testing.T) {
	testFlags := []struct {
		desc           string
		branchOnRemote bool
		gitError       error
		expectError    bool
	}{
		{
			"branch on remote and no git error",
			true,
			nil,
			false,
		},
		{
			"branch on remote and git error",
			true,
			errors.New("git push error"),
			true,
		},
		{
			"no branch on remote and no git error",
			false,
			nil,
			false,
		},
		{
			"no branch on remote and git error",
			false,
			errors.New("git push error"),
			false,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			git := &GitImpl{
				// git command calls are faked
				traceGitFunction: func(_ []string) (err error) {
					return tt.gitError
				},
				workingBranchExistsOnRemote: tt.branchOnRemote,
			}
			err := git.Pull()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_git_add_command(t *testing.T) {
	testFlags := []struct {
		desc         string
		paths        []string
		gitError     error
		expectError  bool
		expectedArgs []string
	}{
		{
			"no git error",
			[]string{"some-path"},
			nil,
			false,
			[]string{"add", "some-path"},
		},
		{
			"git error",
			[]string{"some-path"},
			errors.New("git add error"),
			true,
			[]string{"add", "some-path"},
		},
		{
			"default path",
			nil,
			nil,
			false,
			[]string{"add", "."},
		},
		{
			"multiple paths",
			[]string{"path1", "path2"},
			nil,
			false,
			[]string{"add", "path1", "path2"},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			git := &GitImpl{
				// git command calls are faked
				traceGitFunction: func(args []string) (err error) {
					actualArgs = args[2:]
					return tt.gitError
				},
			}
			err := git.Add(tt.paths...)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedArgs, actualArgs)
		})
	}
}

func Test_git_commit_command(t *testing.T) {
	testFlags := []struct {
		desc         string
		message      string
		amend        bool
		gitError     error
		expectError  bool
		expectedArgs []string
	}{
		{
			"no git error",
			"some message", false,
			nil,
			false,
			[]string{"commit", "--no-gpg-sign", "-m", "some message"},
		},
		{
			"git error",
			"some message", false,
			errors.New("git commit error"),
			false, // We currently ignore git commit errors to handle the case when there's nothing to commit
			[]string{"commit", "--no-gpg-sign", "-m", "some message"},
		},
		{
			"with amend option",
			"some message", true,
			errors.New("git commit error"),
			false, // We currently ignore git commit errors to handle the case when there's nothing to commit
			[]string{"commit", "--no-gpg-sign", "--amend", "-m", "some message"},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			git := &GitImpl{
				// git command calls are faked
				traceGitFunction: func(args []string) (err error) {
					actualArgs = args[2:]
					return tt.gitError
				},
			}
			err := git.Commit(tt.amend, tt.message)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedArgs, actualArgs)
		})
	}
}

func Test_git_restore_command(t *testing.T) {
	testFlags := []struct {
		desc        string
		gitError    error
		expectError bool
	}{
		{
			"no git error",
			nil,
			false,
		},
		{
			"git error",
			errors.New("git restore error"),
			true,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			git := &GitImpl{
				// git command calls are faked
				traceGitFunction: func(_ []string) (err error) {
					return tt.gitError
				},
			}
			err := git.Restore("some-path")
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_git_revert_command(t *testing.T) {
	testFlags := []struct {
		desc         string
		gitError     error
		expectError  bool
		expectedArgs []string
	}{
		{
			"no git error",
			nil,
			false,
			[]string{"revert", "--no-gpg-sign", "--no-edit", "HEAD"},
		},
		{
			"git error",
			errors.New("git revert error"),
			true,
			[]string{"revert", "--no-gpg-sign", "--no-edit", "HEAD"},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			git := &GitImpl{
				// git command calls are faked
				traceGitFunction: func(args []string) (err error) {
					actualArgs = args[2:]
					return tt.gitError
				},
			}
			err := git.Revert()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedArgs, actualArgs)
		})
	}
}

func Test_git_stash_command(t *testing.T) {
	testFlags := []struct {
		desc         string
		message      string
		gitError     error
		expectError  bool
		expectedArgs []string
	}{
		{
			"no git error",
			"some message",
			nil,
			false,
			[]string{"stash", "push", "--quiet", "--include-untracked", "--message", "some message"},
		},
		{
			"git error",
			"some message",
			errors.New("git stash push error"),
			true,
			[]string{"stash", "push", "--quiet", "--include-untracked", "--message", "some message"},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			git := &GitImpl{
				// git command calls are faked
				traceGitFunction: func(args []string) (err error) {
					actualArgs = args[2:]
					return tt.gitError
				},
			}
			err := git.Stash(tt.message)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedArgs, actualArgs)
		})
	}
}

func Test_git_unstash_command(t *testing.T) {
	testFlags := []struct {
		desc         string
		keep         bool
		gitError     error
		expectError  bool
		expectedArgs []string
	}{
		{
			"keep stash and no git error",
			true,
			nil,
			false,
			[]string{"stash", "apply", "--quiet"},
		},
		{
			"keep stash and git error",
			true,
			errors.New("git stash apply error"),
			true,
			[]string{"stash", "apply", "--quiet"},
		},
		{
			"remove stash and no git error",
			false,
			nil,
			false,
			[]string{"stash", "pop", "--quiet"},
		},
		{
			"remove stash and git error",
			false,
			errors.New("git stash pop error"),
			true,
			[]string{"stash", "pop", "--quiet"},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			git := &GitImpl{
				// git command calls are faked
				traceGitFunction: func(args []string) (err error) {
					actualArgs = args[2:]
					return tt.gitError
				},
			}
			err := git.UnStash(tt.keep)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedArgs, actualArgs)
		})
	}
}

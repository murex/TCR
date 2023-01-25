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

package git

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/murex/tcr/vcs"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func slowTestTag(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

func Test_git_auto_push_is_disabled_default(t *testing.T) {
	g, _ := newGitImpl(inMemoryRepoInit, "")
	assert.Zero(t, g.IsPushEnabled())
}

func Test_git_enable_disable_push(t *testing.T) {
	g, _ := newGitImpl(inMemoryRepoInit, "")
	g.EnablePush(true)
	assert.NotZero(t, g.IsPushEnabled())
	g.EnablePush(false)
	assert.Zero(t, g.IsPushEnabled())
}

func Test_init_fails_when_working_dir_is_not_in_a_git_repo(t *testing.T) {
	g, err := New("/")
	assert.Zero(t, g)
	assert.Error(t, err)
}

func Test_can_retrieve_working_branch_on_in_memory_repo(t *testing.T) {
	g, _ := newGitImpl(inMemoryRepoInit, "")
	assert.Equal(t, "master", g.GetWorkingBranch())
}

func Test_can_retrieve_working_branch_on_current_repo(t *testing.T) {
	g, _ := New(".")
	assert.NotEmpty(t, g.GetWorkingBranch())
}

func Test_check_remote_access_on_current_repo(t *testing.T) {
	slowTestTag(t)
	g, _ := New(".")
	// Depending on where the test is run (local machine or CI) and which is the current branch,
	// CheckRemoteAccess() result can vary.
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("skipped when on GitHub Actions")
	}
	if g.IsRemoteEnabled() {
		assert.True(t, g.CheckRemoteAccess())
	} else {
		assert.False(t, g.CheckRemoteAccess())
	}
}

func Test_check_remote_access_on_in_memory_repo(t *testing.T) {
	g, _ := newGitImpl(inMemoryRepoInit, "")
	assert.False(t, g.CheckRemoteAccess())
}

func Test_git_diff_command(t *testing.T) {
	testFlags := []struct {
		desc          string
		gitDiffOutput string
		gitDiffError  error
		expectError   bool
		expectedArgs  []string
		expectedDiff  vcs.FileDiffs
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
			vcs.FileDiffs{
				{Path: filepath.Join("/", "some-file.txt"), AddedLines: 1, RemovedLines: 1},
			},
		},
		{"2 files changed",
			"1\t1\tfile1.txt\n1\t1\tfile2.txt\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: filepath.Join("/", "file1.txt"), AddedLines: 1, RemovedLines: 1},
				{Path: filepath.Join("/", "file2.txt"), AddedLines: 1, RemovedLines: 1},
			},
		},
		{"file changed in sub-directory",
			"1\t1\tsome-dir/some-file.txt\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: filepath.Join("/", "some-dir", "some-file.txt"), AddedLines: 1, RemovedLines: 1},
			},
		},
		{"1 file changed with added lines only",
			"15\t0\tsome-file.txt\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: filepath.Join("/", "some-file.txt"), AddedLines: 15, RemovedLines: 0},
			},
		},
		{"1 file changed with removed lines only",
			"0\t7\tsome-file.txt\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: filepath.Join("/", "some-file.txt"), AddedLines: 0, RemovedLines: 7},
			},
		},
		{"noise in output trace",
			"warning: LF will be replaced by CRLF in some-file.txt.\n" +
				"The file will have its original line endings in your working directory\n" +
				"1\t1\tsome-file.txt\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: filepath.Join("/", "some-file.txt"), AddedLines: 1, RemovedLines: 1},
			},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			g, _ := newGitImpl(inMemoryRepoInit, "")
			g.runGitFunction = func(args ...string) (output []byte, err error) {
				actualArgs = args[2:]
				return []byte(tt.gitDiffOutput), tt.gitDiffError
			}
			fileDiffs, err := g.Diff()
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
			g, _ := newGitImpl(inMemoryRepoInit, "")
			g.traceGitFunction = func(_ ...string) (err error) {
				return tt.gitError
			}
			g.pushEnabled = tt.pushEnabled
			g.remoteEnabled = true
			g.remoteName = DefaultRemoteName

			err := g.Push()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectBranchOnRemote, g.workingBranchExistsOnRemote)
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
			g, _ := newGitImpl(inMemoryRepoInit, "")
			g.traceGitFunction = func(_ ...string) (err error) {
				return tt.gitError
			}
			g.remoteEnabled = true
			g.remoteName = DefaultRemoteName
			g.workingBranchExistsOnRemote = tt.branchOnRemote

			err := g.Pull()
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
			g, _ := newGitImpl(inMemoryRepoInit, "")
			g.traceGitFunction = func(args ...string) (err error) {
				actualArgs = args[2:]
				return tt.gitError
			}

			err := g.Add(tt.paths...)
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
		messages     []string
		amend        bool
		gitError     error
		expectError  bool
		expectedArgs []string
	}{
		{
			"no git error",
			[]string{"some message"}, false,
			nil,
			false,
			[]string{"commit", "--no-gpg-sign", "-m", "some message"},
		},
		{
			"git error",
			[]string{"some message"}, false,
			errors.New("git commit error"),
			false, // We currently ignore git commit errors to handle the case when there's nothing to commit
			[]string{"commit", "--no-gpg-sign", "-m", "some message"},
		},
		{
			"with multiple messages",
			[]string{"main message", "additional message"}, false,
			nil,
			false,
			[]string{"commit", "--no-gpg-sign", "-m", "main message", "-m", "additional message"},
		},
		{
			"with multi-line messages",
			[]string{"main message", "- line 1\n- line 2"}, false,
			nil,
			false,
			[]string{"commit", "--no-gpg-sign", "-m", "main message", "-m", "- line 1\n- line 2"},
		},
		{
			"with amend option",
			[]string{"some message"}, true,
			errors.New("git commit error"),
			false,
			[]string{"commit", "--no-gpg-sign", "--amend", "-m", "some message"},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			g, _ := newGitImpl(inMemoryRepoInit, "")
			g.traceGitFunction = func(args ...string) error {
				actualArgs = args[2:]
				return tt.gitError
			}
			err := g.Commit(tt.amend, tt.messages...)
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
			g, _ := newGitImpl(inMemoryRepoInit, "")
			g.traceGitFunction = func(_ ...string) (err error) {
				return tt.gitError
			}

			err := g.Restore("some-path")
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
			g, _ := newGitImpl(inMemoryRepoInit, "")
			g.traceGitFunction = func(args ...string) (err error) {
				actualArgs = args[2:]
				return tt.gitError
			}

			err := g.Revert()
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
			g, _ := newGitImpl(inMemoryRepoInit, "")
			g.traceGitFunction = func(args ...string) (err error) {
				actualArgs = args[2:]
				return tt.gitError
			}

			err := g.Stash(tt.message)
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
			g, _ := newGitImpl(inMemoryRepoInit, "")
			g.traceGitFunction = func(args ...string) (err error) {
				actualArgs = args[2:]
				return tt.gitError
			}

			err := g.UnStash(tt.keep)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedArgs, actualArgs)
		})
	}
}

func Test_git_log_command(t *testing.T) {
	// Note: this test may break if for any reason the TCR repository initial commit is altered
	tcrInitialCommit := vcs.LogItem{
		Hash:      "a823c098187455bade90ee44874d2b41c7ef96d9",
		Timestamp: time.Date(2021, time.June, 16, 15, 29, 41, 0, time.UTC),
		Message:   "Initial commit",
	}

	testFlags := []struct {
		desc     string
		filter   func(msg string) bool
		asserter func(t *testing.T, items vcs.LogItems)
	}{
		{
			"filter matching no item",
			func(_ string) bool { return false },
			func(t *testing.T, items vcs.LogItems) {
				assert.Equal(t, 0, items.Len())
			},
		},
		{
			"filter matching all items",
			nil,
			func(t *testing.T, items vcs.LogItems) {
				assert.Greater(t, items.Len(), 100)
			},
		},
		{
			"filter matching one single item",
			func(msg string) bool { return tcrInitialCommit.Message == msg },
			func(t *testing.T, items vcs.LogItems) {
				assert.Equal(t, 1, items.Len())
				assert.Equal(t, tcrInitialCommit, items[0])
			},
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			g, _ := New(".")
			items, err := g.Log(tt.filter)
			assert.NoError(t, err)
			tt.asserter(t, items)
		})
	}
}

func Test_nothing_to_commit(t *testing.T) {
	testFlags := []struct {
		desc           string
		gitInitializer func() *gitImpl
		expected       bool
	}{
		{
			"with empty working tree",
			func() *gitImpl {
				g, _ := newGitImpl(inMemoryRepoInit, "")
				return g
			},
			true,
		},
		{
			"with a new file",
			func() *gitImpl {
				g, _ := newGitImpl(inMemoryRepoInit, "")
				filePath := "my-file.txt"
				newFile, _ := g.filesystem.Create(filePath)
				_, _ = newFile.Write([]byte("My new file\n"))
				_ = newFile.Close()
				return g
			},
			false,
		},
		{
			"with a new added uncommitted file",
			func() *gitImpl {
				g, _ := newGitImpl(inMemoryRepoInit, "")
				filePath := "my-file.txt"
				newFile, _ := g.filesystem.Create(filePath)
				_, _ = newFile.Write([]byte("My new file\n"))
				_ = newFile.Close()
				worktree, _ := g.repository.Worktree()
				_, _ = worktree.Add(filePath)
				return g
			},
			false,
		},
		{
			"with a committed file",
			func() *gitImpl {
				g, _ := newGitImpl(inMemoryRepoInit, "")
				filePath := "my-file.txt"
				newFile, _ := g.filesystem.Create(filePath)
				_, _ = newFile.Write([]byte("My new file\n"))
				_ = newFile.Close()
				worktree, _ := g.repository.Worktree()
				_, _ = worktree.Add(filePath)
				_, _ = worktree.Commit("", &git.CommitOptions{
					Author: &object.Signature{Name: "John Doe", Email: "john@doe.org", When: time.Now()}})
				return g
			},
			true,
		},
		{
			"with a modified file",
			func() *gitImpl {
				g, _ := newGitImpl(inMemoryRepoInit, "")
				filePath := "my-file.txt"
				newFile, _ := g.filesystem.Create(filePath)
				_, _ = newFile.Write([]byte("My new file\n"))
				_ = newFile.Close()
				worktree, _ := g.repository.Worktree()
				_, _ = worktree.Add(filePath)
				_, _ = worktree.Commit("", &git.CommitOptions{
					Author: &object.Signature{Name: "John Doe", Email: "john@doe.org", When: time.Now()}})
				openFile, _ := g.filesystem.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
				_, _ = openFile.Write([]byte("A new line\n"))
				_ = openFile.Close()
				return g
			},
			false,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			g := tt.gitInitializer()
			assert.Equal(t, tt.expected, g.nothingToCommit())

		})
	}
}

func Test_is_on_root_branch(t *testing.T) {
	testFlags := []struct {
		desc     string
		name     string
		expected bool
	}{
		{"master branch", "master", true},
		{"main branch", "main", true},
		{"other branch", "xxx", false},
		{"empty branch", "", false},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			g, _ := newGitImpl(inMemoryRepoInit, "")
			g.workingBranch = tt.name
			assert.Equal(t, tt.expected, g.IsOnRootBranch())
		})
	}
}

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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/murex/tcr/vcs"
	"github.com/stretchr/testify/assert"
)

//func slowTestTag(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping test in short mode.")
//	}
//}

// inMemoryRepoInit initializes a brand new repository in memory (for use in tests)
func inMemoryRepoInit(_ string) (repo *git.Repository, fs billy.Filesystem, err error) {
	fs = memfs.New()
	repo, err = git.Init(memory.NewStorage(), fs)
	return
}

func Test_get_vcs_name(t *testing.T) {
	g, _ := newGitImpl(inMemoryRepoInit, "", "")
	assert.Equal(t, "git", g.Name())
}

func Test_get_vcs_session_summary(t *testing.T) {
	tests := []struct {
		desc          string
		remoteName    string
		remoteEnabled bool
		expected      string
	}{
		{
			desc:          "remote not set",
			remoteName:    "",
			remoteEnabled: false,
			expected:      "git branch \"master\"",
		},
		{
			desc:          "remote origin enabled",
			remoteName:    "origin",
			remoteEnabled: true,
			expected:      "git branch \"origin/master\"",
		},
		{
			desc:          "remote origin disabled",
			remoteName:    "origin",
			remoteEnabled: false,
			expected:      "git branch \"master\"",
		},
		{
			desc:          "remote with custom name",
			remoteName:    "custom-origin",
			remoteEnabled: true,
			expected:      "git branch \"custom-origin/master\"",
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			g, _ := newGitImpl(inMemoryRepoInit, "", test.remoteName)
			// we need to overwrite remoteEnabled flag when running with an in-memory repo
			// (this flag is automatically set depending on the existence of this remote alias)
			g.remoteEnabled = test.remoteEnabled
			assert.Equal(t, test.expected, g.SessionSummary())
		})
	}

}

func Test_git_auto_push_is_disabled_default(t *testing.T) {
	g, _ := newGitImpl(inMemoryRepoInit, "", "")
	assert.Zero(t, g.IsAutoPushEnabled())
}

func Test_git_supports_emojis(t *testing.T) {
	g, _ := newGitImpl(inMemoryRepoInit, "", "")
	assert.True(t, g.SupportsEmojis())
}

func Test_git_enable_disable_push(t *testing.T) {
	g, _ := newGitImpl(inMemoryRepoInit, "", "")
	g.EnableAutoPush(true)
	assert.NotZero(t, g.IsAutoPushEnabled())
	g.EnableAutoPush(false)
	assert.Zero(t, g.IsAutoPushEnabled())
}

func Test_init_fails_when_working_dir_is_not_in_a_git_repo(t *testing.T) {
	g, err := New("/", "")
	assert.Zero(t, g)
	assert.Error(t, err)
}

func Test_can_retrieve_working_branch_on_in_memory_repo(t *testing.T) {
	g, _ := newGitImpl(inMemoryRepoInit, "", "")
	// go-git's in memory repository default branch is "master"
	assert.Equal(t, "master", g.GetWorkingBranch())
}

func Test_can_retrieve_working_branch_on_current_repo(t *testing.T) {
	g, _ := New(".", "")
	assert.NotEmpty(t, g.GetWorkingBranch())
}

func Test_check_remote_access_on_in_memory_repo(t *testing.T) {
	g, _ := newGitImpl(inMemoryRepoInit, "", "")
	assert.False(t, g.CheckRemoteAccess())
}

func Test_git_diff(t *testing.T) {
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
			nil,
			false,
			[]string{
				"diff", "--numstat", "--ignore-cr-at-eol", "--ignore-blank-lines", "HEAD"},
			nil,
		},
		{"git diff command call fails",
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
			g, _ := newGitImpl(inMemoryRepoInit, "", "")
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

func Test_git_push(t *testing.T) {
	testFlags := []struct {
		desc                 string
		autoPushEnabled      bool
		gitError             error
		expectError          bool
		expectBranchOnRemote bool
	}{
		{
			"auto push enabled and git push command call succeeds",
			true,
			nil,
			false,
			true,
		},
		{
			"auto push enabled and git push command call fails",
			true,
			errors.New("git push error"),
			true,
			false,
		},
		{
			"auto push disabled and git push command call succeeds",
			false,
			nil,
			false,
			true,
		},
		{
			"auto push disabled and git push command call fails",
			false,
			errors.New("git push error"),
			true,
			false,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			g, _ := newGitImpl(inMemoryRepoInit, "", "")
			g.traceGitFunction = func(_ ...string) (err error) {
				return tt.gitError
			}
			g.autoPushEnabled = tt.autoPushEnabled
			g.remoteEnabled = true

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

func Test_git_pull(t *testing.T) {
	testFlags := []struct {
		desc           string
		branchOnRemote bool
		gitError       error
		expectError    bool
	}{
		{
			"branch on remote and git pull command call succeeds",
			true,
			nil,
			false,
		},
		{
			"branch on remote and git pull command call fails",
			true,
			errors.New("git pull error"),
			true,
		},
		{
			"no branch on remote and git pull command call succeeds",
			false,
			nil,
			false,
		},
		{
			"no branch on remote and git pull command call fails",
			false,
			errors.New("git pull error"),
			false,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			g, _ := newGitImpl(inMemoryRepoInit, "", "")
			g.traceGitFunction = func(_ ...string) (err error) {
				return tt.gitError
			}
			g.remoteEnabled = true
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

func Test_git_add(t *testing.T) {
	testFlags := []struct {
		desc         string
		paths        []string
		gitError     error
		expectError  bool
		expectedArgs []string
	}{
		{
			"git add command call succeeds",
			[]string{"some-path"},
			nil,
			false,
			[]string{"add", "some-path"},
		},
		{
			"git add command call fails",
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
			g, _ := newGitImpl(inMemoryRepoInit, "", "")
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

func Test_git_commit(t *testing.T) {
	testFlags := []struct {
		desc         string
		messages     []string
		gitError     error
		expectError  bool
		expectedArgs []string
	}{
		{
			"git commit command call succeeds",
			[]string{"some message"},
			nil,
			false,
			[]string{"commit", "--no-gpg-sign", "-m", "some message"},
		},
		{
			"git commit command call fails",
			[]string{"some message"},
			errors.New("git commit error"),
			false, // We currently ignore git commit errors to handle the case when there's nothing to commit
			[]string{"commit", "--no-gpg-sign", "-m", "some message"},
		},
		{
			"with multiple messages",
			[]string{"main message", "additional message"},
			nil,
			false,
			[]string{"commit", "--no-gpg-sign", "-m", "main message", "-m", "additional message"},
		},
		{
			"with multi-line messages",
			[]string{"main message", "- line 1\n- line 2"},
			nil,
			false,
			[]string{"commit", "--no-gpg-sign", "-m", "main message", "-m", "- line 1\n- line 2"},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			g, _ := newGitImpl(inMemoryRepoInit, "", "")
			g.traceGitFunction = func(args ...string) error {
				actualArgs = args[2:]
				return tt.gitError
			}
			err := g.Commit(tt.messages...)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedArgs, actualArgs)
		})
	}
}

func Test_git_revert_local(t *testing.T) {
	testFlags := []struct {
		desc        string
		gitError    error
		expectError bool
	}{
		{
			"git checkout command call succeeds",
			nil,
			false,
		},
		{
			"git checkout command call fails",
			errors.New("git checkout error"),
			true,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			g, _ := newGitImpl(inMemoryRepoInit, "", "")
			g.traceGitFunction = func(_ ...string) (err error) {
				return tt.gitError
			}

			err := g.RevertLocal("some-path")
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_git_rollback_last_commit(t *testing.T) {
	testFlags := []struct {
		desc         string
		gitError     error
		expectError  bool
		expectedArgs []string
	}{
		{
			"git revert command call succeeds",
			nil,
			false,
			[]string{"revert", "--no-gpg-sign", "--no-edit", "--no-commit", "HEAD"},
		},
		{
			"git revert command call fails",
			errors.New("git revert error"),
			true,
			[]string{"revert", "--no-gpg-sign", "--no-edit", "--no-commit", "HEAD"},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			g, _ := newGitImpl(inMemoryRepoInit, "", "")
			g.traceGitFunction = func(args ...string) (err error) {
				actualArgs = args[2:]
				return tt.gitError
			}

			err := g.RollbackLastCommit()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedArgs, actualArgs)
		})
	}
}

func Test_git_log(t *testing.T) {
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
			g, _ := newGitImpl(inMemoryRepoInit, ".", "")
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
				g, _ := newGitImpl(inMemoryRepoInit, "", "")
				return g
			},
			true,
		},
		{
			"with a new file",
			func() *gitImpl {
				g, _ := newGitImpl(inMemoryRepoInit, "", "")
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
				g, _ := newGitImpl(inMemoryRepoInit, "", "")
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
				g, _ := newGitImpl(inMemoryRepoInit, "", "")
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
				g, _ := newGitImpl(inMemoryRepoInit, "", "")
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

func Test_check_remote_access(t *testing.T) {
	testFlags := []struct {
		desc           string
		remoteEnabled  bool
		remoteName     string
		branchName     string
		gitError       error
		expectedResult bool
		expectedArgs   []string
	}{
		{
			"git push dry-run command call succeeds",
			true, "origin", "main",
			nil,
			true,
			[]string{"push", "--dry-run", "origin", "main"},
		},
		{
			"git push dry-run command call fails",
			true, "origin", "main",
			errors.New("git push --dry-run error"),
			false,
			[]string{"push", "--dry-run", "origin", "main"},
		},
		{
			"git remote is disabled",
			false, "origin", "main",
			nil,
			false,
			nil,
		},
		{
			"git remote is not set",
			true, "", "main",
			nil,
			false,
			nil,
		},
		{
			"git branch is not set",
			true, "origin", "",
			nil,
			false,
			nil,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			g, _ := newGitImpl(inMemoryRepoInit, "", "")
			g.remoteEnabled = tt.remoteEnabled
			g.remoteName = tt.remoteName
			g.workingBranch = tt.branchName
			g.runGitFunction = func(args ...string) (out []byte, err error) {
				actualArgs = args[2:]
				return nil, tt.gitError
			}
			assert.Equal(t, tt.expectedResult, g.CheckRemoteAccess())
			assert.Equal(t, tt.expectedArgs, actualArgs)
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
			g, _ := newGitImpl(inMemoryRepoInit, "", "")
			g.workingBranch = tt.name
			assert.Equal(t, tt.expected, g.IsOnRootBranch())
		})
	}
}

func Test_git_run_command_global_parameters(t *testing.T) {
	g, _ := newGitImpl(inMemoryRepoInit, filepath.FromSlash("/basedir"), "")
	var cmdParams []string
	g.runGitFunction = func(params ...string) (out []byte, err error) {
		cmdParams = params
		return nil, nil
	}
	_, _ = g.runGit()
	assert.Equal(t, []string{"-C", filepath.FromSlash("/")}, cmdParams)
}

func Test_git_trace_command_global_parameters(t *testing.T) {
	g, _ := newGitImpl(inMemoryRepoInit, filepath.FromSlash("/basedir"), "")
	var cmdParams []string
	g.traceGitFunction = func(params ...string) (err error) {
		cmdParams = params
		return nil
	}
	_ = g.traceGit()
	assert.Equal(t, []string{"-C", filepath.FromSlash("/")}, cmdParams)
}

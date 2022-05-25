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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/go-git/go-billy/v5/helper/chroot"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/murex/tcr/tcr-engine/report"
	"path/filepath"
	"strconv"
	"strings"
)

// GitImpl provides the implementation of the git interface
type GitImpl struct {
	baseDir                     string
	rootDir                     string
	remoteName                  string
	workingBranch               string
	workingBranchExistsOnRemote bool
	pushEnabled                 bool
	runGitFunction              func(params []string) (output []byte, err error)
	traceGitFunction            func(params []string) (err error)
}

// New initializes the git implementation based on the provided directory from local clone
func New(dir string) (GitInterface, error) {
	var gitImpl = GitImpl{
		baseDir:          dir,
		remoteName:       DefaultRemoteName,
		pushEnabled:      DefaultPushEnabled,
		runGitFunction:   runGitCommand,
		traceGitFunction: traceGitCommand,
	}

	plainOpenOptions := git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: false,
	}
	repo, err := git.PlainOpenWithOptions(gitImpl.baseDir, &plainOpenOptions)
	if err != nil {
		return nil, err
	}
	r, _ := findRootDir(repo)
	gitImpl.rootDir = filepath.Dir(r)

	remoteOk, err := isRemoteDefined(repo, gitImpl.GetRemoteName())
	if !remoteOk {
		return nil, err
	}

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}
	gitImpl.workingBranch = head.Name().Short()
	gitImpl.workingBranchExistsOnRemote, err = isBranchOnRemote(repo, gitImpl.GetWorkingBranch(), gitImpl.GetRemoteName())

	return &gitImpl, err
}

// isRemoteDefined returns true is the provided remote is defined in the repository
func isRemoteDefined(repo *git.Repository, remoteName string) (bool, error) {
	_, err := repo.Remote(remoteName)
	if err != nil {
		return false, errors.New("git remote " + remoteName + " is not set")
	}
	return true, nil
}

// isBranchOnRemote returns true is the provided branch exists in the repository for the provided remote
func isBranchOnRemote(repo *git.Repository, branchName, remoteName string) (bool, error) {
	branches, err := remoteBranches(repo.Storer)
	if err != nil {
		return false, err
	}

	remoteBranchName := fmt.Sprintf("%v/%v", remoteName, branchName)
	var found = false
	_ = branches.ForEach(func(branch *plumbing.Reference) error {
		found = found || strings.HasSuffix(branch.Name().Short(), remoteBranchName)
		return nil
	})
	return found, nil
}

// remoteBranches returns the list of known remote branches
func remoteBranches(s storer.ReferenceStorer) (storer.ReferenceIter, error) {
	refs, err := s.IterReferences()
	if err != nil {
		return nil, err
	}

	// We keep only remote branches, and ignore symbolic references
	return storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return ref.Name().IsRemote() && ref.Type() != plumbing.SymbolicReference
	}, refs), nil
}

// findRootDir returns the local clone's root directory of provided repository
func findRootDir(r *git.Repository) (dir string, err error) {
	// Try to grab the repository Storer
	storage, ok := r.Storer.(*filesystem.Storage)
	if !ok {
		err = errors.New("repository storage is not filesystem.Storage")
		return
	}

	// Try to get the underlying billy.Filesystem
	fs, ok := storage.Filesystem().(*chroot.ChrootHelper)
	if !ok {
		err = errors.New("filesystem is not chroot.ChrootHelper")
		return
	}
	dir = fs.Root()
	return
}

// GetRootDir returns the root directory path
func (g *GitImpl) GetRootDir() string {
	return g.rootDir
}

// GetRemoteName returns the current git remote name
func (g *GitImpl) GetRemoteName() string {
	return g.remoteName
}

// GetWorkingBranch returns the current git working branch
func (g *GitImpl) GetWorkingBranch() string {
	return g.workingBranch
}

// Add adds the listed paths to git index.
// Current implementation uses a direct call to git
func (g *GitImpl) Add(paths ...string) error {
	gitArgs := []string{"add"}
	if len(paths) == 0 {
		gitArgs = append(gitArgs, ".")
	} else {
		gitArgs = append(gitArgs, paths...)
	}
	return g.traceGit(gitArgs...)
}

// Commit commits changes to git index.
// Current implementation uses a direct call to git
func (g *GitImpl) Commit(amend bool, message string) error {
	gitArgs := []string{"commit", "--no-gpg-sign"}
	if amend {
		gitArgs = append(gitArgs, "--amend")
	}
	gitArgs = append(gitArgs, "-m", message)
	_ = g.traceGit(gitArgs...)
	// We ignore return code on purpose to prevent raising an error
	// when there is nothing to commit
	// TODO find a way to check beforehand if there is something to commit
	// ("git diff --exit-code --quiet HEAD" seems to do the trick)
	return nil
}

// Restore restores to last commit for the provided path.
// Current implementation uses a direct call to git
func (g *GitImpl) Restore(path string) error {
	report.PostWarning("Reverting ", path)
	return g.traceGit("checkout", "HEAD", "--", path)
}

// Revert runs a git revert operation.
// Current implementation uses a direct call to git
func (g *GitImpl) Revert() error {
	report.PostInfo("Reverting changes")
	return g.traceGit("revert", "--no-gpg-sign", "--no-edit", "HEAD")
}

// Push runs a git push operation.
// Current implementation uses a direct call to git
func (g *GitImpl) Push() error {
	if !g.IsPushEnabled() {
		// There's nothing to do in this case
		return nil
	}

	report.PostInfo("Pushing changes to ", g.GetRemoteName(), "/", g.GetWorkingBranch())
	err := g.traceGit("push", "--no-recurse-submodules", g.GetRemoteName(), g.GetWorkingBranch())
	if err == nil {
		g.workingBranchExistsOnRemote = true
	}
	return err
}

// Pull runs a git pull operation.
// Current implementation uses a direct call to git
func (g *GitImpl) Pull() error {
	if !g.workingBranchExistsOnRemote {
		report.PostInfo("Working locally on branch ", g.GetWorkingBranch())
		return nil
	}
	report.PostInfo("Pulling latest changes from ", g.GetRemoteName(), "/", g.GetWorkingBranch())
	return g.traceGit("pull", "--no-recurse-submodules", g.GetRemoteName(), g.GetWorkingBranch())
}

// Stash creates a git stash.
// Current implementation uses a direct call to git
func (g *GitImpl) Stash(message string) error {
	report.PostInfo("Stashing changes")
	return g.traceGit("stash", "push", "--quiet", "--include-untracked", "--message", message)
	//	"--", "java/src/main", "java/src/test")
}

// UnStash applies a git stash. Depending on the keep argument value, either a "stash apply" or a "stash pop"
// command is executed under the hood.
// Current implementation uses a direct call to git
func (g *GitImpl) UnStash(keep bool) error {
	report.PostInfo("Applying stashed changes")
	stashAction := "pop"
	if keep {
		stashAction = "apply"
	}
	return g.traceGit("stash", stashAction, "--quiet")
}

// Diff returns the list of files modified since last commit with diff info for each file
// Current implementation uses a direct call to git
func (g *GitImpl) Diff() (diffs []FileDiff, err error) {
	var gitOutput []byte
	gitOutput, err = g.runGit("diff", "--numstat", "--ignore-cr-at-eol",
		"--ignore-all-space", "--ignore-blank-lines", "HEAD")
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(gitOutput))
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		if len(fields) == 3 {
			added, _ := strconv.Atoi(fields[0])
			removed, _ := strconv.Atoi(fields[1])
			filename := filepath.Join(g.rootDir, fields[2])
			diffs = append(diffs, NewFileDiff(filename, added, removed))
		}
	}
	return
}

// EnablePush sets a flag allowing to turn on/off git push operations
func (g *GitImpl) EnablePush(flag bool) {
	if g.pushEnabled == flag {
		return
	}
	g.pushEnabled = flag
	autoPushStr := "off"
	if g.pushEnabled {
		autoPushStr = "on"
	}
	report.PostInfo(fmt.Sprintf("Git auto-push is turned %v", autoPushStr))
}

// IsPushEnabled indicates if git push operations are turned on
func (g *GitImpl) IsPushEnabled() bool {
	return g.pushEnabled
}

// CheckRemoteAccess returns true if git remote can be accessed. This is currently done through
// checking the return value of "git push --dry-run". This very likely does not guarantee that
// git remote commands will work, but already gives an indication.
func (g *GitImpl) CheckRemoteAccess() bool {
	_, err := g.runGit("push", "--dry-run", g.GetRemoteName(), g.GetWorkingBranch())
	return err == nil
}

// traceGit runs a git command and traces its output.
// The command is launched from the git root directory
func (g *GitImpl) traceGit(args ...string) error {
	return g.traceGitFunction(append([]string{"-C", g.GetRootDir()}, args...))
}

// runGit calls git command in a separate process and returns its output traces
// The command is launched from the git root directory
func (g *GitImpl) runGit(args ...string) (output []byte, err error) {
	return g.runGitFunction(append([]string{"-C", g.GetRootDir()}, args...))
}

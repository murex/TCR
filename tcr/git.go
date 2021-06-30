package tcr

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/mengdaming/tcr/trace"
	"os"
	"strings"
	"time"
)

const (
	GitRemoteName = "origin"
)

var (
	gitWorkingBranch               string
	gitWorkingBranchExistsOnRemote bool
)

func detectGitWorkingBranch() {
	cwd, err := os.Getwd()
	if err != nil {
		trace.Error("os.Getwd(): ", err)
	}
	//trace.Info("Current Working Directory: ", cwd)

	repo, err := git.PlainOpen(cwd)
	if err != nil {
		trace.Error("git.PlainOpen(): ", err)
	}

	head, err := repo.Head()
	if err != nil {
		trace.Error("repo.Head(): ", err)
	}

	gitWorkingBranch = head.Name().Short()
	trace.Info("Git Working Branch: ", gitWorkingBranch)

	gitWorkingBranchExistsOnRemote = isBranchOnRemote(repo, gitWorkingBranch, GitRemoteName)
	trace.Info("Git Branch exists on ", GitRemoteName,": ", gitWorkingBranchExistsOnRemote)
}

func isBranchOnRemote(repo *git.Repository, branch, remote string) bool {
	remoteName := fmt.Sprintf("%v/%v", remote, branch)
	branches, err := remoteBranches(repo.Storer)
	if err != nil {
		trace.Error("remoteBranches(): ", err)
	}

	var found = false
	_ = branches.ForEach(func(branch *plumbing.Reference) error {
		found = found || strings.HasSuffix(branch.Name().Short(), remoteName)
		return nil
	})

	return found
}

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

func push() {
	trace.Info("Pushing changes to origin/", gitWorkingBranch)
	time.Sleep(1 * time.Second)
	// TODO Call to git push --no-recurse-submodules origin "${GIT_WORKING_BRANCH}"
	// TODO [ ${git_rc} -eq 0 ] && GIT_WORKING_BRANCH_EXISTS_ON_ORIGIN=1
	// TODO	return ${git_rc}
}

func pull() {
	if gitWorkingBranchExistsOnRemote {
		trace.Info("Pulling latest changes from origin/", gitWorkingBranch)
		time.Sleep(1 * time.Second)
		// TODO Call to git pull --no-recurse-submodules origin "${GIT_WORKING_BRANCH}"
	} else {
		trace.Info("Working locally on branch ", gitWorkingBranch)
		time.Sleep(1 * time.Second)
	}
}

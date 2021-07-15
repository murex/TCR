package tcr

import (
	"errors"
	"fmt"
	"github.com/go-git/go-billy/v5/helper/chroot"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/mengdaming/tcr/trace"
	"os"
	"path/filepath"
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

func detectGitWorkingBranch(dir string) {
	gitOptions := git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: false,
	}
	repo, err := git.PlainOpenWithOptions(dir, &gitOptions)
	if err != nil {
		trace.Error("git.PlainOpen(): ", err)
	}
	r, _ := root(repo)
	trace.Info("Repository Root: ", filepath.Dir(r))

	head, err := repo.Head()
	if err != nil {
		trace.Error("repo.Head(): ", err)
	}

	gitWorkingBranch = head.Name().Short()
	trace.Info("Git Working Branch: ", gitWorkingBranch)

	gitWorkingBranchExistsOnRemote = isBranchOnRemote(repo, gitWorkingBranch, GitRemoteName)
	trace.Info("Git Branch exists on ", GitRemoteName, ": ", gitWorkingBranchExistsOnRemote)
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

func root(r *git.Repository) (string, error) {
	// Try to grab the repository Storer
	s, ok := r.Storer.(*filesystem.Storage)
	if !ok {
		return "", errors.New("repository storage is not filesystem.Storage")
	}

	// Try to get the underlying billy.Filesystem
	fs, ok := s.Filesystem().(*chroot.ChrootHelper)
	if !ok {
		return "", errors.New("filesystem is not chroot.ChrootHelper")
	}

	return fs.Root(), nil
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

		dir, _ := os.Getwd()
		gitOptions := git.PlainOpenOptions{
			DetectDotGit:          true,
			EnableDotGitCommonDir: false,
		}
		repo, err := git.PlainOpenWithOptions(dir, &gitOptions)
		if err != nil {
			trace.Error("git.PlainOpen(): ", err)
		}

		worktree, err := repo.Worktree()
		if err != nil {
			trace.Error("repo.Worktree(): ", err)
		}

		err = worktree.Pull(&git.PullOptions{
			RemoteName:        GitRemoteName,
			ReferenceName:     plumbing.ReferenceName(gitWorkingBranch),
			SingleBranch:      true,
			RecurseSubmodules: git.NoRecurseSubmodules},
		)

		printLastCommit(repo)
	} else {
		trace.Info("Working locally on branch ", gitWorkingBranch)
	}
}

func printLastCommit(repo *git.Repository) {
	// Print the latest commit that was just pulled
	head, err := repo.Head()
	if err != nil {
		trace.Error("repo.Head(): ", err)
	}
	commit, err := repo.CommitObject(head.Hash())
	trace.Echo(commit)
}

func restore(dir string) {
	// TODO Call to git checkout HEAD -- ${SRC_DIRS}
	trace.Info("Restoring ", dir)
	//	dir, _ := os.Getwd()
	gitOptions := git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: false,
	}
	repo, err := git.PlainOpenWithOptions(dir, &gitOptions)
	if err != nil {
		trace.Error("git.PlainOpen(): ", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		trace.Error("repo.Worktree(): ", err)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.HEAD,
		Create: false,
		Keep:   false,
	})
	if err != nil {
		trace.Error("worktree.Checkout(): ", err)
	}
}

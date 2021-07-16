package tcr

import (
	"errors"
	"fmt"
	"github.com/codeskyblue/go-sh"
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
	DefaultRemoteName = "origin"
)

type GitInterface interface {
	WorkingBranch() string
	Commit()
	Restore(dir string)
	Push()
	Pull()
}

type GoGit struct {
	remoteName                  string
	workingBranch               string
	workingBranchExistsOnRemote bool
}

func NewGoGit(dir string) GitInterface {
	var goGit GoGit = GoGit{
		remoteName: DefaultRemoteName,
	}

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

	goGit.workingBranch = head.Name().Short()
	trace.Info("Git Working Branch: ", goGit.workingBranch)

	goGit.workingBranchExistsOnRemote = isBranchOnRemote(repo, goGit.workingBranch, goGit.remoteName)
	trace.Info("Git Branch exists on ",
		goGit.remoteName, ": ", goGit.workingBranchExistsOnRemote)

	return goGit
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

func (g GoGit) WorkingBranch() string {
	return g.workingBranch
}

func (g GoGit) Commit() {
	time.Sleep(1 * time.Second)
	// TODO Call to git commit -am TCR
}

func (g GoGit) Restore(dir string) {
	// Currently, go-git does not allow to checkout a subset of the
	// files related to a branch or commit.
	// There's a pending PR that should allow this, that we could use
	// once it's merged and packaged into go-git.
	// Cf. https://github.com/go-git/go-git/pull/343
	// In the meantime, we use direct call to git checkout HEAD
	// TODO When available, replace git call with go-git restore function

	trace.Info("Restoring ", dir)

	err := gitCommand([]string{"checkout", "HEAD", "--", dir})
	if err != nil {
		trace.Error(err)
	}
}

func (g GoGit) Push() {
	trace.Info("Pushing changes to origin/", g.workingBranch)
	time.Sleep(1 * time.Second)
	// TODO Call to git push --no-recurse-submodules origin "${GIT_WORKING_BRANCH}"
	// TODO [ ${git_rc} -eq 0 ] && GIT_WORKING_BRANCH_EXISTS_ON_ORIGIN=1
	// TODO	return ${git_rc}
}

func (g GoGit) Pull() {
	if !g.workingBranchExistsOnRemote {
		trace.Info("Working locally on branch ", g.workingBranch)
		return
	}

	trace.Info("Pulling latest changes from ",
		g.remoteName, "/", g.workingBranch)

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
		RemoteName:        g.remoteName,
		ReferenceName:     plumbing.ReferenceName(g.workingBranch),
		SingleBranch:      true,
		RecurseSubmodules: git.NoRecurseSubmodules},
	)

	printLastCommit(repo)
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

func gitCommand(params []string) error {
	gitCommand := "git"
	output, err := sh.Command(gitCommand, params).Output()
	if output != nil {
		trace.Echo(string(output))
	}
	if err != nil {
		trace.Warning(err)
	}
	return err
}

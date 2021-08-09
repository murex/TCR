package vcs

import (
	"errors"
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/go-git/go-billy/v5/helper/chroot"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/mengdaming/tcr/tcr/report"
	"path/filepath"
	"strings"
)

type GitImpl struct {
	baseDir                     string
	rootDir                     string
	remoteName                  string
	workingBranch               string
	workingBranchExistsOnRemote bool
	commitMessage               string
	pushEnabled                 bool
}

// New initializes the git implementation based on the provided directory from local clone
func New(dir string) GitInterface {
	var gitImpl = GitImpl{
		baseDir:       dir,
		remoteName:    DefaultRemoteName,
		commitMessage: DefaultCommitMessage,
		pushEnabled:   DefaultPushEnabled,
	}

	plainOpenOptions := git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: false,
	}
	repo, err := git.PlainOpenWithOptions(gitImpl.baseDir, &plainOpenOptions)
	if err != nil {
		report.PostError("git.PlainOpenWithOptions(): ", err)
		return nil
	}
	r, _ := rootDir(repo)
	gitImpl.rootDir = filepath.Dir(r)

	head, err := repo.Head()
	if err != nil {
		report.PostError("repo.Head(): ", err)
		return nil
	}

	gitImpl.workingBranch = head.Name().Short()

	gitImpl.workingBranchExistsOnRemote = isBranchOnRemote(repo, gitImpl.workingBranch, gitImpl.remoteName)

	return &gitImpl
}

// isBranchOnRemote returns true is the provided branch exists on provided remote
func isBranchOnRemote(repo *git.Repository, branch, remote string) bool {
	remoteName := fmt.Sprintf("%v/%v", remote, branch)
	branches, err := remoteBranches(repo.Storer)
	if err != nil {
		report.PostError("remoteBranches(): ", err)
	}

	var found = false
	_ = branches.ForEach(func(branch *plumbing.Reference) error {
		found = found || strings.HasSuffix(branch.Name().Short(), remoteName)
		return nil
	})

	return found
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

// rootDir Returns the local clone's root directory of provided repository
func rootDir(r *git.Repository) (string, error) {
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

// WorkingBranch returns the current git working branch
func (g *GitImpl) WorkingBranch() string {
	return g.workingBranch
}

// Commit restores to last commit.
// Current implementation uses a direct call to git
func (g *GitImpl) Commit() {
	// go-git Add function does not use .gitignore contents
	// when applied on a directory.
	// Until this is implemented we rely instead on a direct
	// git command call
	// TODO When gitignore rules are implemented, use go-git.Commit()

	// We ignore return code on purpose to prevent exiting
	// when there is nothing to commit
	_ = runGitCommand([]string{"commit", "-am", g.commitMessage})
}

// Restore restores to last commit for everything under dir.
// Current implementation uses a direct call to git
func (g *GitImpl) Restore(dir string) {
	// Currently, go-git does not allow to checkout a subset of the
	// files related to a branch or commit.
	// There's a pending PR that should allow this, that we could use
	// once it's merged and packaged into go-git.
	// Cf. https://github.com/go-git/go-git/pull/343
	// In the meantime, we use direct call to git checkout HEAD
	// TODO When available, replace git call with go-git restore function

	report.PostInfo("Restoring ", dir)

	err := runGitCommand([]string{"checkout", "HEAD", "--", dir})
	if err != nil {
		report.PostError(err)
	}
}

// Push runs a git push.
// Current implementation uses a direct call to git
func (g *GitImpl) Push() {
	if g.IsPushEnabled() {
		report.PostInfo("Pushing changes to origin/", g.workingBranch)

		// Solution below works but requires to provide username
		// and password, which is not acceptable here. Until we
		// find a way to reuse git credentials, we'll use a direct
		// git command call instead
		// TODO Look if there is a way to reuse git credentials

		//gitOptions := git.PlainOpenOptions{
		//	DetectDotGit:          true,
		//	EnableDotGitCommonDir: false,
		//}
		//repo, err := git.PlainOpenWithOptions(g.baseDir, &gitOptions)
		//if err != nil {
		//	report.PostError("git.PlainOpenWithOptions(): ", err)
		//}
		//
		//err = repo.Push(&git.PushOptions{
		//	RemoteName: g.remoteName,
		//	Auth: &http.BasicAuth{
		//		Username: "xxx",
		//		Password: "xxx",
		//	},
		//})
		//if err != nil {
		//	report.PostError("repo.Push(): ", err)
		//} else {
		//	g.workingBranchExistsOnRemote = isBranchOnRemote(
		//		repo, g.workingBranch, g.remoteName)
		//}

		err := runGitCommand([]string{
			"push",
			"--no-recurse-submodules",
			g.remoteName,
			g.workingBranch,
		})
		if err != nil {
			report.PostError(err)
		} else {
			g.workingBranchExistsOnRemote = true
		}
	}
}

// Pull runs a git pull operation.
// Current implementation uses a direct call to git
func (g *GitImpl) Pull() {
	if !g.workingBranchExistsOnRemote {
		report.PostInfo("Working locally on branch ", g.workingBranch)
		return
	}

	report.PostInfo("Pulling latest changes from ",
		g.remoteName, "/", g.workingBranch)

	// Solution below works but requires to provide username
	// and password, which is not acceptable here. Until we
	// find a way to reuse git credentials, we'll use a direct
	// git command call instead
	// TODO Look if there is a way to reuse git credentials

	//gitOptions := git.PlainOpenOptions{
	//	DetectDotGit:          true,
	//	EnableDotGitCommonDir: false,
	//}
	//repo, err := git.PlainOpenWithOptions(g.baseDir, &gitOptions)
	//if err != nil {
	//	report.PostError("git.PlainOpenWithOptions(): ", err)
	//}
	//
	//worktree, err := repo.Worktree()
	//if err != nil {
	//	report.PostError("repo.Worktree(): ", err)
	//}
	//
	//err = worktree.Pull(&git.PullOptions{
	//	RemoteName:        g.remoteName,
	//	ReferenceName:     plumbing.ReferenceName(g.workingBranch),
	//	SingleBranch:      true,
	//	RecurseSubmodules: git.NoRecurseSubmodules},
	//)

	//report.PostEcho("From ", g.remoteName)
	//report.PostEcho(" * branch\t", g.workingBranch, " -> FETCH_HEAD")
	//switch err {
	//case git.NoErrAlreadyUpToDate:
	//	report.PostEcho("Already up to date.")
	//case nil:
	//	printLastCommit(repo)
	//default:
	//	report.PostWarning("Pull(): ", err)
	//}

	err := runGitCommand([]string{
		"pull",
		"--no-recurse-submodules",
		g.remoteName,
		g.workingBranch,
	})
	if err != nil {
		report.PostError(err)
	}
}

// EnablePush Set a flag allowing to turn on/off git push operations
func (g *GitImpl) EnablePush(flag bool) {
	g.pushEnabled = flag
}

// IsPushEnabled Indicates if git push operations are turned on
func (g *GitImpl) IsPushEnabled() bool {
	return g.pushEnabled
}

// runGitCommand Calls git command in a separate process
func runGitCommand(params []string) error {
	gitCommand := "git"
	output, err := sh.Command(gitCommand, params).CombinedOutput()

	if output != nil {
		report.PostText(string(output))
	}
	return err
}

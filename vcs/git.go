package vcs

const (
	// DefaultRemoteName is the alias used by default for the git remote repository
	DefaultRemoteName = "origin"
	// DefaultCommitMessage is the message used by default by TCR every time it does a git commit
	DefaultCommitMessage = "TCR"
	// DefaultPushEnabled indicates the default state for git auto-push option
	DefaultPushEnabled = false
)

// GitInterface provides the interface that a git implementation must satisfy for TCR engine to be
// able to interact with git
type GitInterface interface {
	WorkingBranch() string
	Commit() error
	Restore(dir string) error
	Push() error
	Pull() error
	EnablePush(flag bool)
	IsPushEnabled() bool
}

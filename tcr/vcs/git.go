package vcs

const (
	DefaultRemoteName    = "origin"
	DefaultCommitMessage = "TCR"
	DefaultPushEnabled   = false
)

type GitInterface interface {
	WorkingBranch() string
	Commit() error
	Restore(dir string) error
	Push() error
	Pull() error
	EnablePush(flag bool)
	IsPushEnabled() bool
}

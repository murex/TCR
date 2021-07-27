package vcs

const (
	DefaultRemoteName    = "origin"
	DefaultCommitMessage = "TCR"
	DefaultPushEnabled   = false
)

type GitInterface interface {
	WorkingBranch() string
	Commit()
	Restore(dir string)
	Push()
	Pull()
	EnablePush(flag bool)
	IsPushEnabled() bool
}

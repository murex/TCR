package vcs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// push-enabling flag

func Test_git_auto_push_is_disabled_default(t *testing.T) {
	git, _ := New(".")
	assert.Zero(t, git.IsPushEnabled())
}

func Test_git_enable_disable_push(t *testing.T) {
	git, _ := New(".")
	git.EnablePush(true)
	assert.NotZero(t, git.IsPushEnabled())
	git.EnablePush(false)
	assert.Zero(t, git.IsPushEnabled())
}

// Working Branch

func Test_init_fails_when_working_dir_is_not_in_a_git_repo(t *testing.T) {
	git, err := New("/")
	assert.Zero(t, git)
	assert.NotZero(t, err)
}

func Test_can_retrieve_working_branch(t *testing.T) {
	git, _ := New(".")
	assert.NotZero(t, git.WorkingBranch())
}
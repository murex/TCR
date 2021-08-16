package filesystem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Source Tree initialization

func Test_init_source_tree_with_missing_directory_fails(t *testing.T) {
	tree, err := New("/dummy")
	assert.Zero(t, tree)
	assert.NotZero(t, err)
}

func Test_init_source_tree_with_existing_directory_passes(t *testing.T) {
	tree, err := New(".")
	assert.NotZero(t, tree)
	assert.Zero(t, err)
}

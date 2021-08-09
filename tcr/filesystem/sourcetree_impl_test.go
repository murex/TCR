package filesystem

import (
	"github.com/mengdaming/tcr/tcr/trace"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Prevent trace.Error() from triggering os.Exit()
	trace.SetTestMode()
	os.Exit(m.Run())
}

// Source Tree initialization

func Test_init_source_tree_with_missing_directory_fails(t *testing.T) {
	tree, err := New("/dummy")
	assert.Zero(t, tree)
	assert.NotZero(t, err)
	//assert.NotZero(t, trace.GetExitReturnCode())
}

func Test_init_source_tree_with_existing_directory_passes(t *testing.T) {
	tree, err := New(".")
	assert.NotZero(t, tree)
	assert.Zero(t, err)
}

package cli

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

func Test_confirm_with_default_answer_to_yes(t *testing.T) {
	assertConfirmBehaviour(t, []byte{enterKey}, true, true)
}

func Test_confirm_with_default_answer_to_no(t *testing.T) {
	assertConfirmBehaviour(t, []byte{enterKey}, false, false)
}

func Test_confirm_with_a_yes_answer(t *testing.T) {
	assertConfirmBehaviour(t, []byte{'y'}, false, true)
	assertConfirmBehaviour(t, []byte{'Y'}, false, true)
}

func Test_confirm_with_a_no_answer(t *testing.T) {
	assertConfirmBehaviour(t, []byte{'n'}, true, false)
	assertConfirmBehaviour(t, []byte{'N'}, true, false)
}

func Test_confirm_question_with_default_answer_to_no(t *testing.T) {
	assert.Equal(t, "[y/N]", yesOrNoAdvice(false))
}

func Test_confirm_question_with_default_answer_to_yes(t *testing.T) {
	assert.Equal(t, "[Y/n]", yesOrNoAdvice(true))
}

func assertConfirmBehaviour(t *testing.T, input []byte, defaultValue bool, expected bool) {
	stdin := os.Stdin
	stdout := os.Stdout
	// Restore stdin and stdout right after the test.
	defer func() { os.Stdin = stdin; os.Stdout = stdout }()
	// We mock stdin so that we can simulate a key press
	os.Stdin = mockStdin(t, input)
	// Displayed info on stdout is useless for the test
	os.Stdout = os.NewFile(0, os.DevNull)

	term := New()
	assert.Equal(t, expected, term.Confirm("", defaultValue))
}

func mockStdin(t *testing.T, input []byte) *os.File {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}
	_ = w.Close()
	return r
}

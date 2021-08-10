package cli

import (
	"github.com/kami-zh/go-capturer"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	ansiEscape   = "\x1b[0m"
	ansiRedFg    = "\x1b[31m"
	ansiYellowFg = "\x1b[33m"
	ansiCyanFg   = "\x1b[36m"
	newline      = "\n"
)

func Test_change_line_prefix(t *testing.T) {
	prefix := "PREFIX"
	msg := "Message"
	setLinePrefix(prefix)
	out := capturer.CaptureStdout(func() {
		printInCyan(msg)
	})
	assert.Contains(t, out, prefix)
}

func Test_print_untouched_function_does_not_alter_data(t *testing.T) {
	msg := "Dummy Message"
	out := capturer.CaptureStdout(func() {
		printUntouched(msg)
	})
	assert.Equal(t, msg+"\n", out)
}

func Test_print_in_cyan_function_formatting(t *testing.T) {
	prefix := "PREFIX"
	msg := "Message"
	setLinePrefix(prefix)
	out := capturer.CaptureStdout(func() {
		printInCyan(msg)
	})
	expected := ansiCyanFg + prefix + ansiEscape + " " + ansiCyanFg + msg + ansiEscape + newline
	assert.Equal(t, expected, out)
}

func Test_print_in_yellow_function_formatting(t *testing.T) {
	prefix := "PREFIX"
	msg := "Message"
	setLinePrefix(prefix)
	out := capturer.CaptureStdout(func() {
		printInYellow(msg)
	})
	expected := ansiYellowFg + prefix + ansiEscape + " " + ansiYellowFg + msg + ansiEscape + newline
	assert.Equal(t, expected, out)
}

func Test_print_in_red_function_formatting(t *testing.T) {
	prefix := "PREFIX"
	msg := "Message"
	setLinePrefix(prefix)
	out := capturer.CaptureStdout(func() {
		printInRed(msg)
	})
	expected := ansiRedFg + prefix + ansiEscape + " " + ansiRedFg + msg + ansiEscape + newline
	assert.Equal(t, expected, out)
}

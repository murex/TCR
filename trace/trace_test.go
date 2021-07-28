package trace

import (
	"github.com/kami-zh/go-capturer"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	ansiReset    = "\x1b[0m"
	ansiRedFg    = "\x1b[31m"
	ansiYellowFg = "\x1b[33m"
	ansiCyanFg   = "\x1b[36m"
	newline      = "\n"
)

func TestMain(m *testing.M) {
	// Prevent trace.Error() from triggering os.Exit()
	SetTestMode()
	os.Exit(m.Run())
}

func Test_change_line_prefix(t *testing.T) {
	prefix := "PREFIX"
	msg := "Message"
	SetLinePrefix(prefix)
	out := capturer.CaptureStdout(func() {
		Info(msg)
	})
	assert.Contains(t, out, prefix)
}

func Test_echo_function_does_not_alter_data(t *testing.T) {
	msg := "Dummy Message"
	out := capturer.CaptureStdout(func() {
		Echo(msg)
	})
	assert.Equal(t, msg+"\n", out)
}

func Test_info_function_formatting(t *testing.T) {
	prefix := "PREFIX"
	msg := "Message"
	SetLinePrefix(prefix)
	out := capturer.CaptureStdout(func() {
		Info(msg)
	})
	expected := ansiCyanFg + prefix + ansiReset + " " + ansiCyanFg + msg + ansiReset + newline
	assert.Equal(t, expected, out)
}

func Test_warning_function_formatting(t *testing.T) {
	prefix := "PREFIX"
	msg := "Message"
	SetLinePrefix(prefix)
	out := capturer.CaptureStdout(func() {
		Warning(msg)
	})
	expected := ansiYellowFg + prefix + ansiReset + " " + ansiYellowFg + msg + ansiReset + newline
	assert.Equal(t, expected, out)
}

func Test_error_function_formatting(t *testing.T) {
	prefix := "PREFIX"
	msg := "Message"
	SetLinePrefix(prefix)
	out := capturer.CaptureStdout(func() {
		Error(msg)
	})
	expected := ansiRedFg + prefix + ansiReset + " " + ansiRedFg + msg + ansiReset + newline
	assert.Equal(t, expected, out)
}

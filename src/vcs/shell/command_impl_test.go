/*
Copyright (c) 2023 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package shell

import (
	"bytes"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/vcs"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

func Test_is_in_path_for_a_valid_command(t *testing.T) {
	assert.True(t, NewCommandFunc("ls").IsInPath())
}

func Test_is_in_path_for_an_invalid_command(t *testing.T) {
	assert.False(t, NewCommandFunc("unknown-command").IsInPath())
}

func Test_get_full_path_for_a_valid_command(t *testing.T) {
	base := filepath.Base(NewCommandFunc("ls").GetFullPath())
	assert.Equal(t, strings.TrimSuffix(base, ".exe"), "ls")
}

func Test_get_full_path_for_an_invalid_command(t *testing.T) {
	assert.Zero(t, NewCommandFunc("unknown-command").GetFullPath())
}

func Test_run_valid_command_with_initial_parameters(t *testing.T) {
	output, err := NewCommandFunc("echo", "hello world!").Run()
	assert.NoError(t, err)
	trimmed := string(bytes.TrimRight(output, "\r\n"))
	assert.Equal(t, "hello world!", trimmed)
}

func Test_run_valid_command_with_additional_parameters(t *testing.T) {
	output, err := NewCommandFunc("echo").Run("hello world!")
	assert.NoError(t, err)
	trimmed := string(bytes.TrimRight(output, "\r\n"))
	assert.Equal(t, "hello world!", trimmed)
}

func Test_run_invalid_command(t *testing.T) {
	output, err := NewCommandFunc("unknown-command").Run()
	assert.Error(t, err)
	assert.Zero(t, output)
}

func Test_trace_valid_command_with_initial_parameters(t *testing.T) {
	sniffer := report.NewSniffer()
	err := NewCommandFunc("echo", "hello world!").Trace()
	sniffer.Stop()
	assert.NoError(t, err)
	assert.Equal(t, 1, sniffer.GetMatchCount())
	trimmed := strings.TrimRight(sniffer.GetAllMatches()[0].Text, "\r\n")
	assert.Equal(t, "hello world!", trimmed)
}

func Test_trace_valid_command_with_additional_parameters(t *testing.T) {
	sniffer := report.NewSniffer()
	err := NewCommandFunc("echo").Trace("hello world!")
	sniffer.Stop()
	assert.NoError(t, err)
	assert.Equal(t, 1, sniffer.GetMatchCount())
	trimmed := strings.TrimRight(sniffer.GetAllMatches()[0].Text, "\r\n")
	assert.Equal(t, "hello world!", trimmed)
}

func Test_trace_invalid_command(t *testing.T) {
	sniffer := report.NewSniffer()
	err := NewCommandFunc("unknown-command").Trace()
	sniffer.Stop()
	assert.Error(t, err)
	assert.Equal(t, 0, sniffer.GetMatchCount())
}

func Test_run_pipe_valid_commands_with_initial_parameters(t *testing.T) {
	output, err := NewCommandFunc("echo", "hello\tworld!").RunAndPipe(
		NewCommandFunc("cut", "-f", "1"))
	assert.NoError(t, err)
	trimmed := string(bytes.TrimRight(output, "\r\n"))
	assert.Equal(t, "hello", trimmed)
}

func Test_run_pipe_valid_commands_with_additional_parameters(t *testing.T) {
	output, err := NewCommandFunc("echo").RunAndPipe(
		NewCommandFunc("cut", "-f", "1"),
		"hello\tworld!")
	assert.NoError(t, err)
	trimmed := string(bytes.TrimRight(output, "\r\n"))
	assert.Equal(t, "hello", trimmed)
}

func Test_run_pipe_with_first_command_invalid(t *testing.T) {
	output, err := NewCommandFunc("unknown-command").RunAndPipe(
		NewCommandFunc("cut", "-f", "1"))
	assert.Error(t, err)
	assert.Zero(t, output)
}

func Test_run_pipe_with_second_command_invalid(t *testing.T) {
	output, err := NewCommandFunc("echo").RunAndPipe(
		NewCommandFunc("unknown-command"))
	assert.Error(t, err)
	assert.Zero(t, output)
}

func Test_trace_pipe_valid_commands_with_initial_parameters(t *testing.T) {
	sniffer := report.NewSniffer()
	err := NewCommandFunc("echo", "hello\tworld!").TraceAndPipe(
		NewCommandFunc("cut", "-f", "1"))
	sniffer.Stop()
	assert.NoError(t, err)
	assert.Equal(t, 1, sniffer.GetMatchCount())
	trimmed := strings.TrimRight(sniffer.GetAllMatches()[0].Text, "\r\n")
	assert.Equal(t, "hello", trimmed)
}

func Test_trace_pipe_valid_commands_with_additional_parameters(t *testing.T) {
	sniffer := report.NewSniffer()
	err := NewCommandFunc("echo").TraceAndPipe(
		NewCommandFunc("cut", "-f", "1"),
		"hello\tworld!")
	sniffer.Stop()
	assert.NoError(t, err)
	assert.Equal(t, 1, sniffer.GetMatchCount())
	trimmed := strings.TrimRight(sniffer.GetAllMatches()[0].Text, "\r\n")
	assert.Equal(t, "hello", trimmed)
}

func Test_trace_pipe_with_first_command_invalid(t *testing.T) {
	sniffer := report.NewSniffer()
	err := NewCommandFunc("unknown-command").TraceAndPipe(
		NewCommandFunc("cut", "-f", "1"))
	sniffer.Stop()
	assert.Error(t, err)
	assert.Equal(t, 0, sniffer.GetMatchCount())
}

func Test_trace_pipe_with_second_command_invalid(t *testing.T) {
	sniffer := report.NewSniffer()
	err := NewCommandFunc("echo").TraceAndPipe(
		NewCommandFunc("unknown-command"))
	sniffer.Stop()
	assert.Error(t, err)
	assert.Equal(t, 0, sniffer.GetMatchCount())
}

func Test_command_as_string(t *testing.T) {
	tests := []struct {
		desc      string
		command   Command
		extraArgs []string
		expected  string
	}{
		{
			"0 arg and 0 extra arg",
			NewCommandImpl("command"),
			nil,
			"command",
		},
		{
			"1 arg and 0 extra arg",
			NewCommandImpl("command", "arg1"),
			nil,
			"command arg1",
		},
		{
			"1 arg and 1 extra arg",
			NewCommandImpl("command", "arg1"),
			[]string{"extra1"},
			"command arg1 extra1",
		},
		{
			"0 arg and 1 extra arg",
			NewCommandImpl("command"),
			[]string{"extra1"},
			"command extra1",
		},
		{
			"2 args and 3 extra args",
			NewCommandImpl("command", "arg1", "arg2"),
			[]string{"extra1", "extra2", "extra3"},
			"command arg1 arg2 extra1 extra2 extra3",
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert.Equal(t, test.expected, test.command.String(test.extraArgs...))
		})
	}
}

func Test_run_command_with_vcs_trace_enabled(t *testing.T) {
	sniffer := report.NewSniffer(func(msg report.Message) bool {
		return msg.Type.Category == report.Warning && strings.Index(msg.Text, "echo") == 0
	})
	vcs.SetTrace(true)
	_, err := NewCommandFunc("echo", "hello world!").Run()
	vcs.SetTrace(false)
	sniffer.Stop()
	assert.NoError(t, err)
	assert.Equal(t, 1, sniffer.GetMatchCount())
	trimmed := strings.TrimRight(sniffer.GetAllMatches()[0].Text, "\r\n")
	assert.Equal(t, "echo hello world!", trimmed)
}

func Test_run_piped_command_with_vcs_trace_enabled(t *testing.T) {
	sniffer := report.NewSniffer(func(msg report.Message) bool {
		return msg.Type.Category == report.Warning && strings.Index(msg.Text, "echo") == 0
	})
	vcs.SetTrace(true)
	_, err := NewCommandFunc("echo", "hello\tworld!").RunAndPipe(
		NewCommandFunc("cut", "-f", "1"))
	vcs.SetTrace(false)
	sniffer.Stop()
	assert.NoError(t, err)
	assert.Equal(t, 1, sniffer.GetMatchCount())
	trimmed := strings.TrimRight(sniffer.GetAllMatches()[0].Text, "\r\n")
	assert.Equal(t, "echo hello\tworld! | cut -f 1", trimmed)
}

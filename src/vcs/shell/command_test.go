/*
Copyright (c) 2022 Murex

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
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

func Test_is_in_path_for_a_valid_command(t *testing.T) {
	assert.True(t, NewCommand("ls").IsInPath())
}

func Test_is_in_path_for_an_invalid_command(t *testing.T) {
	assert.False(t, NewCommand("unknown-command").IsInPath())
}

func Test_get_full_path_for_a_valid_command(t *testing.T) {
	base := filepath.Base(NewCommand("ls").GetFullPath())
	assert.Equal(t, strings.TrimSuffix(base, ".exe"), "ls")
}

func Test_get_full_path_for_an_invalid_command(t *testing.T) {
	assert.Zero(t, NewCommand("unknown-command").GetFullPath())
}

func Test_run_valid_command_with_initial_parameters(t *testing.T) {
	output, err := NewCommand("echo", "hello world!").Run()
	assert.NoError(t, err)
	trimmed := string(bytes.TrimRight(output, "\r\n"))
	assert.Equal(t, "hello world!", trimmed)
}

func Test_run_valid_command_with_additional_parameters(t *testing.T) {
	output, err := NewCommand("echo").Run("hello world!")
	assert.NoError(t, err)
	trimmed := string(bytes.TrimRight(output, "\r\n"))
	assert.Equal(t, "hello world!", trimmed)
}

func Test_run_invalid_command(t *testing.T) {
	output, err := NewCommand("unknown-command").Run()
	assert.Error(t, err)
	assert.Zero(t, output)
}

func Test_trace_valid_command_with_initial_parameters(t *testing.T) {
	sniffer := report.NewSniffer()
	err := NewCommand("echo", "hello world!").Trace()
	sniffer.Stop()
	assert.NoError(t, err)
	assert.Equal(t, 1, sniffer.GetMatchCount())
	trimmed := strings.TrimRight(sniffer.GetAllMatches()[0].Text, "\r\n")
	assert.Equal(t, "hello world!", trimmed)
}

func Test_trace_valid_command_with_additional_parameters(t *testing.T) {
	sniffer := report.NewSniffer()
	err := NewCommand("echo").Trace("hello world!")
	sniffer.Stop()
	assert.NoError(t, err)
	assert.Equal(t, 1, sniffer.GetMatchCount())
	trimmed := strings.TrimRight(sniffer.GetAllMatches()[0].Text, "\r\n")
	assert.Equal(t, "hello world!", trimmed)
}

func Test_trace_invalid_command(t *testing.T) {
	sniffer := report.NewSniffer()
	err := NewCommand("unknown-command").Trace()
	sniffer.Stop()
	assert.Error(t, err)
	assert.Equal(t, 0, sniffer.GetMatchCount())
}

func Test_run_pipe_valid_commands_with_initial_parameters(t *testing.T) {
	output, err := NewCommand("echo", "hello\tworld!").RunAndPipe(
		NewCommand("cut", "-f", "1"))
	assert.NoError(t, err)
	trimmed := string(bytes.TrimRight(output, "\r\n"))
	assert.Equal(t, "hello", trimmed)
}

func Test_run_pipe_valid_commands_with_additional_parameters(t *testing.T) {
	output, err := NewCommand("echo").RunAndPipe(
		NewCommand("cut", "-f", "1"),
		"hello\tworld!")
	assert.NoError(t, err)
	trimmed := string(bytes.TrimRight(output, "\r\n"))
	assert.Equal(t, "hello", trimmed)
}

func Test_run_pipe_with_first_command_invalid(t *testing.T) {
	output, err := NewCommand("unknown-command").RunAndPipe(
		NewCommand("cut", "-f", "1"))
	assert.Error(t, err)
	assert.Zero(t, output)
}

func Test_run_pipe_with_second_command_invalid(t *testing.T) {
	output, err := NewCommand("echo").RunAndPipe(
		NewCommand("unknown-command"))
	assert.Error(t, err)
	assert.Zero(t, output)
}

func Test_trace_pipe_valid_commands_with_initial_parameters(t *testing.T) {
	sniffer := report.NewSniffer()
	err := NewCommand("echo", "hello\tworld!").TraceAndPipe(
		NewCommand("cut", "-f", "1"))
	sniffer.Stop()
	assert.NoError(t, err)
	assert.Equal(t, 1, sniffer.GetMatchCount())
	trimmed := strings.TrimRight(sniffer.GetAllMatches()[0].Text, "\r\n")
	assert.Equal(t, "hello", trimmed)
}

func Test_trace_pipe_valid_commands_with_additional_parameters(t *testing.T) {
	sniffer := report.NewSniffer()
	err := NewCommand("echo").TraceAndPipe(
		NewCommand("cut", "-f", "1"),
		"hello\tworld!")
	sniffer.Stop()
	assert.NoError(t, err)
	assert.Equal(t, 1, sniffer.GetMatchCount())
	trimmed := strings.TrimRight(sniffer.GetAllMatches()[0].Text, "\r\n")
	assert.Equal(t, "hello", trimmed)
}

func Test_trace_pipe_with_first_command_invalid(t *testing.T) {
	sniffer := report.NewSniffer()
	err := NewCommand("unknown-command").TraceAndPipe(
		NewCommand("cut", "-f", "1"))
	sniffer.Stop()
	assert.Error(t, err)
	assert.Equal(t, 0, sniffer.GetMatchCount())
}

func Test_trace_pipe_with_second_command_invalid(t *testing.T) {
	sniffer := report.NewSniffer()
	err := NewCommand("echo").TraceAndPipe(
		NewCommand("unknown-command"))
	sniffer.Stop()
	assert.Error(t, err)
	assert.Equal(t, 0, sniffer.GetMatchCount())
}

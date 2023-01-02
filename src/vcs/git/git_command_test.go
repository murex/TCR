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

package git

import (
	"github.com/murex/tcr/report"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_is_git_command_available(t *testing.T) {
	assert.True(t, IsGitCommandAvailable())
}

func Test_get_git_command_path(t *testing.T) {
	assert.NotZero(t, GetGitCommandPath())
}

func Test_run_git_command(t *testing.T) {
	output, err := runGitCommand("status")
	assert.NoError(t, err)
	assert.NotZero(t, output)
}

func Test_trace_git_command(t *testing.T) {
	sniffer := report.NewSniffer()
	err := traceGitCommand("status")
	sniffer.Stop()
	assert.NoError(t, err)
	assert.NotZero(t, sniffer.GetMatchCount())
}

func Test_get_git_username(t *testing.T) {
	assert.NotZero(t, GetGitUserName())
}

func Test_get_git_config_value_with_undefined_key(t *testing.T) {
	assert.Equal(t, "not set", getGitConfigValue("undefined-config-value"))
}

func Test_get_git_command_version(t *testing.T) {
	assert.NotZero(t, GetGitCommandVersion())
}

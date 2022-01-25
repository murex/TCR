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

package vcs

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

func Test_check_command_is_available_for_a_valid_command(t *testing.T) {
	assert.True(t, isCommandAvailable("ls"))
}

func Test_check_command_is_available_for_an_invalid_command(t *testing.T) {
	assert.False(t, isCommandAvailable("unknown-command"))
}

func Test_check_command_path_for_a_valid_command(t *testing.T) {
	base := filepath.Base(getCommandPath("ls"))
	assert.Equal(t, strings.TrimSuffix(base, ".exe"), "ls")
}

func Test_check_command_path_for_an_invalid_command(t *testing.T) {
	assert.Zero(t, getCommandPath("unknown-command"))
}

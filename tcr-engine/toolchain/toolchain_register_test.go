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

package toolchain

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_does_not_support_empty_toolchain_name(t *testing.T) {
	assert.False(t, isSupported(""))
}

func Test_does_not_support_unregistered_toolchain_name(t *testing.T) {
	assert.False(t, isSupported("unregistered-toolchain"))
}

func Test_unrecognized_toolchain_name(t *testing.T) {
	toolchain, err := GetToolchain("dummy-toolchain")
	assert.Error(t, err)
	assert.Zero(t, toolchain)
}

func Test_can_add_a_built_in_toolchain(t *testing.T) {
	const name = "new-built-in-toolchain"
	assert.False(t, isBuiltIn(name))
	assert.NoError(t, addBuiltIn(*AToolchain(WithName(name))))
	assert.True(t, isBuiltIn(name))
}

func Test_cannot_add_a_built_in_toolchain_with_no_name(t *testing.T) {
	assert.Error(t, addBuiltIn(*AToolchain(WithName(""))))
}

func Test_toolchain_name_is_case_insensitive(t *testing.T) {
	const name = "miXeD-CasE"
	_ = Register(*AToolchain(WithName(name)))
	assertNameIsNotCaseSensitive(t, name)
}

func Test_can_register_a_toolchain(t *testing.T) {
	const name = "new-toolchain"
	assert.False(t, isSupported(name))
	assert.NoError(t, Register(*AToolchain(WithName(name))))
	assert.True(t, isSupported(name))
}

func Test_cannot_register_a_toolchain_with_no_name(t *testing.T) {
	assert.Error(t, Register(*AToolchain(WithName(""))))
}

func Test_cannot_register_a_toolchain_with_no_build_command(t *testing.T) {
	const name = "no-build-command"
	assert.Error(t, Register(*AToolchain(WithName(name), WithNoBuildCommand())))
	assert.False(t, isSupported(name))
}

func Test_cannot_register_a_toolchain_with_no_test_command(t *testing.T) {
	const name = "no-test-command"
	assert.Error(t, Register(*AToolchain(WithName(name), WithNoTestCommand())))
	assert.False(t, isSupported(name))
}

func assertIsABuiltInToolchain(t *testing.T, name string) {
	assert.True(t, isBuiltIn(name))
}

func assertIsSupported(t *testing.T, name string) {
	assert.True(t, isSupported(name))
}

func assertNameIsNotCaseSensitive(t *testing.T, name string) {
	assert.True(t, isSupported(name))
	assert.True(t, isSupported(strings.ToUpper(name)))
	assert.True(t, isSupported(strings.ToLower(name)))
	assert.True(t, isSupported(strings.Title(name)))
}

func assertToolchainInitialization(t *testing.T, name string) {
	toolchain, err := GetToolchain(name)
	assert.NoError(t, err)
	assert.Equal(t, name, toolchain.GetName())
}

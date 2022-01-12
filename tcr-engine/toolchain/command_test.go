/*
Copyright (c) 2021 Murex

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
	"testing"
)

func Test_os_darwin_is_recognized(t *testing.T) {
	assert.Contains(t, GetAllOsNames(), OsName("darwin"))
}

func Test_os_linux_is_recognized(t *testing.T) {
	assert.Contains(t, GetAllOsNames(), OsName("linux"))
}

func Test_os_windows_is_recognized(t *testing.T) {
	assert.Contains(t, GetAllOsNames(), OsName("windows"))
}

func Test_arch_386_is_recognized(t *testing.T) {
	assert.Contains(t, GetAllArchNames(), ArchName("386"))
}

func Test_arch_amd64_is_recognized(t *testing.T) {
	assert.Contains(t, GetAllArchNames(), ArchName("amd64"))
}

func Test_arch_arm64_is_recognized(t *testing.T) {
	assert.Contains(t, GetAllArchNames(), ArchName("arm64"))
}

func Test_unrecognized_os(t *testing.T) {
	assert.False(t, ACommand().runsWithOs(OsName("dummy")))
}

func Test_unrecognized_architecture(t *testing.T) {
	assert.False(t, ACommand().runsWithArch(ArchName("dummy")))
}

func Test_unrecognized_platform(t *testing.T) {
	dummyOs, dummyArch := OsName("dummy"), ArchName("dummy")

	assert.False(t, ACommand().runsOnPlatform(dummyOs, dummyArch))

	assert.False(t, ACommand().runsOnPlatform(OsDarwin, dummyArch))
	assert.False(t, ACommand().runsOnPlatform(OsWindows, dummyArch))
	assert.False(t, ACommand().runsOnPlatform(OsLinux, dummyArch))

	assert.False(t, ACommand().runsOnPlatform(dummyOs, Arch386))
	assert.False(t, ACommand().runsOnPlatform(dummyOs, ArchAmd64))
	assert.False(t, ACommand().runsOnPlatform(dummyOs, ArchArm64))
}

func Test_find_command_must_match_both_os_and_arch(t *testing.T) {
	myOs, myArch := OsName("my-os"), ArchName("my-arch")
	anotherOs, anotherArch := OsName("another-os"), ArchName("another-arch")
	myCommand := ACommand(WithOs(myOs), WithArch(myArch))
	commands := []Command{*myCommand}

	assert.Equal(t, findCommand(commands, myOs, myArch), myCommand)
	assert.Zero(t, findCommand(commands, anotherOs, anotherArch))
	assert.Zero(t, findCommand(commands, myOs, anotherArch))
	assert.Zero(t, findCommand(commands, anotherOs, myArch))
}

func Test_command_path_cannot_be_empty(t *testing.T) {
	assert.Error(t, ACommand(WithPath("")).check())
}

func Test_command_os_list_cannot_be_empty(t *testing.T) {
	assert.Error(t, ACommand(WithNoOs()).check())
}

func Test_a_command_os_cannot_be_empty(t *testing.T) {
	assert.Error(t, ACommand(WithOs("")).check())
}

func Test_command_arch_list_cannot_be_empty(t *testing.T) {
	assert.Error(t, ACommand(WithNoArch()).check())
}

func Test_a_command_arch_cannot_be_empty(t *testing.T) {
	assert.Error(t, ACommand(WithArch("")).check())
}

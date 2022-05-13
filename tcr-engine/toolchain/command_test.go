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
	"fmt"
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
	dummyOs, dummyArch := OsName("dummy_os"), ArchName("dummy_arch")

	cases := []struct {
		os   OsName
		arch ArchName
	}{
		{dummyOs, dummyArch},
		{OsDarwin, dummyArch},
		{OsWindows, dummyArch},
		{OsLinux, dummyArch},
		{dummyOs, Arch386},
		{dummyOs, ArchAmd64},
		{dummyOs, ArchArm64},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprint(tt.os, "-", tt.arch), func(t *testing.T) {
			assert.False(t, ACommand().runsOnPlatform(tt.os, tt.arch))
		})
	}
}

func Test_find_command_must_match_both_os_and_arch(t *testing.T) {
	myOs, myArch := OsName("my-os"), ArchName("my-arch")
	anotherOs, anotherArch := OsName("another-os"), ArchName("another-arch")
	myCommand := ACommand(WithOs(myOs), WithArch(myArch))
	commands := []Command{*myCommand}

	cases := []struct {
		desc     string
		os       OsName
		arch     ArchName
		expected *Command
	}{
		{"both os and arch match", myOs, myArch, myCommand},
		{"both os and arch do not match", anotherOs, anotherArch, nil},
		{"os matches but arch does not match", myOs, anotherArch, nil},
		{"arch matches but os does not match", anotherOs, myArch, nil},
	}

	for _, tt := range cases {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expected, findCommand(commands, tt.os, tt.arch))
		})
	}
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

func Test_a_valid_command_should_have_path_os_and_arch_non_empty(t *testing.T) {
	assert.NoError(t, ACommand(WithPath("some-path"), WithOs("some-os"), WithArch("some-arch")).check())
}

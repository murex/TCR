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

package p4

import (
	"errors"
	"fmt"
	"github.com/murex/tcr/vcs"
	"github.com/murex/tcr/vcs/shell"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/charmap"
	"path/filepath"
	"testing"
)

// inMemoryDepotInit initializes a brand new depot in memory (for use in tests)
func inMemoryDepotInit() afero.Fs {
	return afero.NewMemMapFs()
}

func Test_get_vcs_name(t *testing.T) {
	p, _ := newP4Impl(inMemoryDepotInit, "", true)
	assert.Equal(t, "p4", p.Name())
}

func Test_get_vcs_session_summary(t *testing.T) {
	p, _ := newP4Impl(inMemoryDepotInit, "", true)
	assert.Equal(t, "p4 client \"test\"", p.SessionSummary())
}

func Test_p4_auto_push_is_always_enabled(t *testing.T) {
	p, _ := newP4Impl(inMemoryDepotInit, "", true)
	assert.True(t, p.IsPushEnabled())
}

func Test_p4_enable_disable_push_has_no_effect(t *testing.T) {
	p, _ := newP4Impl(inMemoryDepotInit, "", true)
	p.EnablePush(true)
	assert.True(t, p.IsPushEnabled())
	p.EnablePush(false)
	assert.True(t, p.IsPushEnabled())
}

func Test_p4_init_fails_when_working_dir_is_not_in_a_depot(t *testing.T) {
	p, err := New("/")
	assert.Zero(t, p)
	assert.Error(t, err)
}

func Test_p4_working_branch_is_always_empty(t *testing.T) {
	p, _ := newP4Impl(inMemoryDepotInit, "", true)
	// We currently don't use p4 branches. This may change in the future
	assert.Equal(t, "", p.GetWorkingBranch())
}

func Test_p4_is_on_root_branch_is_always_false(t *testing.T) {
	p, _ := newP4Impl(inMemoryDepotInit, "", true)
	assert.False(t, p.IsOnRootBranch())
}

func Test_p4_is_always_remote_enabled(t *testing.T) {
	p, _ := newP4Impl(inMemoryDepotInit, "", true)
	assert.True(t, p.IsRemoteEnabled())
}

func Test_p4_check_remote_access_is_always_true(t *testing.T) {
	p, _ := newP4Impl(inMemoryDepotInit, "", true)
	assert.True(t, p.CheckRemoteAccess())
}

func Test_p4_get_remote_name_is_always_empty(t *testing.T) {
	p, _ := newP4Impl(inMemoryDepotInit, "", true)
	assert.Equal(t, "", p.GetRemoteName())
}

func Test_convert_line(t *testing.T) {
	tests := []struct {
		desc     string
		input    rune
		encoding *charmap.Charmap
		expected string
	}{
		{"UTC-8 smiley", '🙂', nil, "f09f9982"},
		{"ISO8859-1 smiley", '🙂', charmap.ISO8859_1, "c3b0c29fc299c282"},
		{"WINDOWS-1252 smiley", '🙂', charmap.Windows1252, "c3b0c5b8e284a2e2809a"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			//t.Logf("%x - %x", test.input, convertLine(test.encoding, string(test.input)))
			assert.Equal(t, test.expected, fmt.Sprintf("%x", convertLine(test.encoding, string(test.input))))
		})
	}
}

func Test_p4_diff(t *testing.T) {
	testFlags := []struct {
		desc         string
		p4DiffOutput string
		p4DiffError  error
		expectError  bool
		expectedArgs []string
		expectedDiff vcs.FileDiffs
	}{
		{"p4 command arguments",
			"",
			nil,
			false,
			[]string{"diff", "-f", "-Od", "-dw", "-ds", filepath.Clean("/...")},
			nil,
		},
		{"p4 diff command call fails",
			"",
			errors.New("p4 diff error"),
			true,
			nil,
			nil,
		},
		{"0 file changed",
			"",
			nil,
			false,
			nil,
			nil,
		},
		{"1 file changed",
			"==== //some-depot/some-file.txt#1 - C:\\some-path\\some-file.txt ====\r\n" +
				"add 1 chunks 1 lines\r\n" +
				"deleted 1 chunks 1 lines\r\n" +
				"changed 0 chunks 0 / 0 lines\r\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: "C:\\some-path\\some-file.txt", AddedLines: 1, RemovedLines: 1},
			},
		},
		{"2 files changed",
			"==== //some-depot/file1.txt#1 - C:\\some-path\\file1.txt ====\r\n" +
				"add 1 chunks 1 lines\r\n" +
				"deleted 1 chunks 1 lines\r\n" +
				"changed 1 chunks 1 / 1 lines\r\n" +
				"==== //some-depot/file2.txt#1 - C:\\some-path\\file2.txt ====\r\n" +
				"add 1 chunks 1 lines\r\n" +
				"deleted 1 chunks 1 lines\r\n" +
				"changed 1 chunks 1 / 1 lines\r\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: "C:\\some-path\\file1.txt", AddedLines: 1 + 1, RemovedLines: 1 + 1},
				{Path: "C:\\some-path\\file2.txt", AddedLines: 1 + 1, RemovedLines: 1 + 1},
			},
		},
		{"file changed in sub-directory",
			"==== //some-depot/some-dir/some-file.txt#1 - C:\\some-path\\some-dir\\some-file.txt ====\r\n" +
				"add 1 chunks 1 lines\r\n" +
				"deleted 1 chunks 1 lines\r\n" +
				"changed 0 chunks 0 / 0 lines\r\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: "C:\\some-path\\some-dir\\some-file.txt", AddedLines: 1, RemovedLines: 1},
			},
		},
		{"1 file changed with added lines only",
			"==== //some-depot/some-file.txt#1 - C:\\some-path\\some-file.txt ====\r\n" +
				"add 1 chunks 15 lines\r\n" +
				"deleted 0 chunks 0 lines\r\n" +
				"changed 0 chunks 0 / 0 lines\r\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: "C:\\some-path\\some-file.txt", AddedLines: 15, RemovedLines: 0},
			},
		},
		{"1 file changed with removed lines only",
			"==== //some-depot/some-file.txt#1 - C:\\some-path\\some-file.txt ====\r\n" +
				"add 0 chunks 0 lines\r\n" +
				"deleted 0 chunks 7 lines\r\n" +
				"changed 0 chunks 0 / 0 lines\r\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: "C:\\some-path\\some-file.txt", AddedLines: 0, RemovedLines: 7},
			},
		},
		{"1 file changed with changed lines only",
			"==== //some-depot/some-file.txt#1 - C:\\some-path\\some-file.txt ====\r\n" +
				"add 0 chunks 0 lines\r\n" +
				"deleted 0 chunks 0 lines\r\n" +
				"changed 1 chunks 2 / 5 lines\r\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: "C:\\some-path\\some-file.txt", AddedLines: 5, RemovedLines: 2},
			},
		},
		{"1 file changed with added, deleted and changed lines",
			"==== //some-depot/some-file.txt#1 - C:\\some-dir\\some-file.txt ====\r\n" +
				"add 2 chunks 10 lines\r\n" +
				"deleted 1 chunks 5 lines\r\n" +
				"changed 3 chunks 15 / 2 lines\r\n",
			nil,
			false,
			nil,
			vcs.FileDiffs{
				{Path: "C:\\some-dir\\some-file.txt", AddedLines: 10 + 2, RemovedLines: 5 + 15},
			},
		},
		{"noise in output trace",
			"Usage: diff [ -d<flags> -f -m max -Od -s<flag> -t ] [files...]\r\n" +
				"Invalid option: -xxx.\r\n",
			nil,
			true,
			nil,
			nil,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			p, _ := newP4Impl(inMemoryDepotInit, "", true)
			p.rootDir = ""
			p.runP4Function = func(args ...string) (output []byte, err error) {
				actualArgs = args[4:]
				return []byte(tt.p4DiffOutput), tt.p4DiffError
			}
			fileDiffs, err := p.Diff()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.expectedArgs != nil {
				assert.Equal(t, tt.expectedArgs, actualArgs)
			}
			assert.Equal(t, tt.expectedDiff, fileDiffs)
		})
	}
}

func Test_p4_push(t *testing.T) {
	p, _ := newP4Impl(inMemoryDepotInit, "", true)
	// p4 push does nothing, thus is never expected to return an error
	assert.NoError(t, p.Push())
}

func Test_p4_pull(t *testing.T) {
	testFlags := []struct {
		desc         string
		rootDir      string
		dir          string
		clientName   string
		p4Error      error
		expectError  bool
		expectedArgs []string
	}{
		{
			"p4 sync command call succeeds",
			filepath.FromSlash("/p4root"),
			filepath.FromSlash("/p4root/base_dir"),
			"test_client",
			nil,
			false,
			[]string{"sync", "//test_client/base_dir/..."},
		},
		{
			"p4 sync command call fails",
			filepath.FromSlash("/p4root"),
			filepath.FromSlash("/p4root/base_dir"),
			"test_client",
			errors.New("p4 sync error"),
			true,
			[]string{"sync", "//test_client/base_dir/..."},
		},
		{
			"base directory is empty",
			filepath.FromSlash("/p4root"),
			"",
			"test_client",
			nil,
			true,
			nil,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			p, _ := newP4Impl(inMemoryDepotInit, tt.dir, true)
			p.rootDir = tt.rootDir
			p.clientName = tt.clientName
			p.traceP4Function = func(args ...string) (err error) {
				actualArgs = args[4:]
				return tt.p4Error
			}
			err := p.Pull()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedArgs, actualArgs)
		})
	}
}

func Test_p4_add(t *testing.T) {
	testFlags := []struct {
		desc         string
		paths        []string
		p4Error      error
		expectError  bool
		expectedArgs []string
	}{
		{
			"p4 reconcile command call succeeds",
			[]string{"some-path"},
			nil,
			false,
			[]string{"reconcile", "-a", "-e", "-d", "some-path"},
		},
		{
			"p4 reconcile command call fails",
			[]string{"some-path"},
			errors.New("p4 reconcile error"),
			true,
			[]string{"reconcile", "-a", "-e", "-d", "some-path"},
		},
		{
			"default path",
			nil,
			nil,
			false,
			[]string{"reconcile", "-a", "-e", "-d", filepath.Clean("/...")},
		},
		{
			"multiple paths",
			[]string{"path1", "path2"},
			nil,
			false,
			[]string{"reconcile", "-a", "-e", "-d", "path1", "path2"},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var actualArgs []string
			p, _ := newP4Impl(inMemoryDepotInit, "", true)
			p.traceP4Function = func(args ...string) (err error) {
				actualArgs = args[4:]
				return tt.p4Error
			}

			err := p.Add(tt.paths...)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedArgs, actualArgs)
		})
	}
}

func Test_p4_commit(t *testing.T) {
	testFlags := []struct {
		desc                 string
		messages             []string
		amend                bool
		p4ChangeError        error
		p4ChangeOutput       string
		p4SubmitError        error
		p4SubmitExpectedArgs []string
		expectError          bool
	}{
		{
			"p4 change and p4 submit command calls succeed",
			[]string{"some message"}, false,
			nil, "change 1234567 created ...",
			nil, []string{"submit", "-c", "1234567"},
			false,
		},
		{
			"p4 change command call fails",
			[]string{"some message"}, false,
			errors.New("p4 change error"), "",
			nil, nil,
			true,
		},
		{
			"p4 submit command call fails",
			[]string{"some message"}, false,
			nil, "change 1234567 created ...",
			errors.New("p4 submit error"), []string{"submit", "-c", "1234567"},
			true,
		},
		{
			"with multiple messages",
			[]string{"main message", "additional message"}, false,
			nil, "change 1234567 created ...",
			nil, []string{"submit", "-c", "1234567"},
			false,
		},
		{
			"with multi-line messages",
			[]string{"main message", "- line 1\n- line 2"}, false,
			nil, "change 1234567 created ...",
			nil, []string{"submit", "-c", "1234567"},
			false,
		},
		{
			"with amend option",
			[]string{"some message"}, true,
			nil, "change 1234567 created ...",
			nil, []string{"submit", "-c", "1234567"},
			false,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var p4SubmitActualArgs []string
			p, _ := newP4Impl(inMemoryDepotInit, "", true)
			p.runPipedP4Function = func(toCmd shell.Command, args ...string) (output []byte, err error) {
				// Stub for the call to "p4 change ... -o | p4 change -i"
				return []byte(tt.p4ChangeOutput), tt.p4ChangeError
			}
			p.traceP4Function = func(args ...string) (err error) {
				// Stub for the call to "p4 submit -c <cl_number>"
				p4SubmitActualArgs = args[4:]
				return tt.p4SubmitError
			}

			err := p.Commit(tt.amend, tt.messages...)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.p4SubmitExpectedArgs, p4SubmitActualArgs)
		})
	}
}

func Test_p4_submit(t *testing.T) {
	testFlags := []struct {
		desc                 string
		p4Changelist         *changeList
		p4SubmitError        error
		p4SubmitExpectedArgs []string
		expectError          bool
	}{
		{
			"p4 submit command call succeeds",
			&changeList{number: "1234567"},
			nil, []string{"submit", "-c", "1234567"},
			false,
		},
		{
			"p4 submit command call fails",
			&changeList{number: "1234567"},
			errors.New("p4 submit error"), []string{"submit", "-c", "1234567"},
			true,
		},
		{
			"empty changelist",
			nil,
			errors.New("p4 submit error"), nil,
			true,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var p4SubmitActualArgs []string
			p, _ := newP4Impl(inMemoryDepotInit, "", true)
			p.traceP4Function = func(args ...string) (err error) {
				// Stub for the call to "p4 submit -c <cl_number>"
				p4SubmitActualArgs = args[4:]
				return tt.p4SubmitError
			}

			err := p.submitChangeList(tt.p4Changelist)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.p4SubmitExpectedArgs, p4SubmitActualArgs)
		})
	}
}

func Test_p4_restore(t *testing.T) {
	testFlags := []struct {
		desc         string
		p4Error      error
		expectedArgs []string
		expectError  bool
	}{
		{
			"p4 revert arguments",
			nil,
			[]string{"revert", "some-path"},
			false,
		},
		{
			"p4 revert command call succeeds",
			nil,
			nil,
			false,
		},
		{
			"p4 revert command call fails",
			errors.New("p4 revert error"),
			nil,
			true,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			p, _ := newP4Impl(inMemoryDepotInit, "", true)
			var actualArgs []string
			p.traceP4Function = func(args ...string) (err error) {
				actualArgs = args[4:]
				return tt.p4Error
			}

			err := p.Restore("some-path")
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.expectedArgs != nil {
				assert.Equal(t, tt.expectedArgs, actualArgs)
			}
		})
	}
}

func Test_convert_to_p4_client_path(t *testing.T) {
	testFlags := []struct {
		desc          string
		rootDir       string
		clientName    string
		dir           string
		expectedError error
		expected      string
	}{
		{
			"Dir is the root mount directory",
			filepath.FromSlash("/p4root"),
			"test_client",
			filepath.FromSlash("/p4root"),
			nil,
			"//test_client/...",
		},
		{
			"Dir is under the root directory",
			filepath.FromSlash("/p4root"),
			"test_client",
			filepath.FromSlash("/p4root/sub_dir"),
			nil,
			"//test_client/sub_dir/...",
		},
		{
			"Root Dir has trailing separators",
			filepath.FromSlash("/p4root//"),
			"test_client",
			filepath.FromSlash("/p4root/sub_dir"),
			nil,
			"//test_client/sub_dir/...",
		},
		{
			"Dir has extra separators",
			filepath.FromSlash("/p4root"),
			"test_client",
			filepath.FromSlash("/p4root//sub_dir///"),
			nil,
			"//test_client/sub_dir/...",
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			p, _ := newP4Impl(inMemoryDepotInit, tt.rootDir, true)
			p.clientName = tt.clientName
			clientPath, err := p.toP4ClientPath(tt.dir)
			assert.Equal(t, tt.expected, clientPath)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func Test_p4_run_command_global_parameters(t *testing.T) {
	p, _ := newP4Impl(inMemoryDepotInit, filepath.FromSlash("/basedir"), true)
	p.clientName = "test_client"
	var cmdParams []string
	p.runP4Function = func(params ...string) (out []byte, err error) {
		cmdParams = params
		return nil, nil
	}
	_, _ = p.runP4()
	assert.Equal(t, []string{"-d", filepath.FromSlash("/basedir"), "-c", "test_client"}, cmdParams)
}

func Test_p4_trace_command_global_parameters(t *testing.T) {
	p, _ := newP4Impl(inMemoryDepotInit, filepath.FromSlash("/basedir"), true)
	p.clientName = "test_client"
	var cmdParams []string
	p.traceP4Function = func(params ...string) (err error) {
		cmdParams = params
		return nil
	}
	_ = p.traceP4()
	assert.Equal(t, []string{"-d", filepath.FromSlash("/basedir"), "-c", "test_client"}, cmdParams)
}

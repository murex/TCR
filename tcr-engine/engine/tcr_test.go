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

package engine

import (
	"fmt"
	"github.com/murex/tcr/tcr-engine/language"
	"github.com/murex/tcr/tcr-engine/runmode"
	"github.com/murex/tcr/tcr-engine/toolchain"
	"github.com/murex/tcr/tcr-engine/ui"
	"github.com/murex/tcr/tcr-engine/vcs"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	failure  int
	failures []failure
)

const (
	failBuild failure = iota
	failTest
	failCommit
	failRestore
	failPush
	failPull
	failDiff
)

func (fs failures) contains(f failure) bool {
	for _, set := range fs {
		if f == set {
			return true
		}
	}
	return false
}

func Test_run_build_command_with_no_failure(t *testing.T) {
	initTcrEngineWithFakes(failures{})
	assertCommandEndState(t, build, StatusOk, false)
}

func Test_run_build_command_with_failure(t *testing.T) {
	initTcrEngineWithFakes(failures{failBuild})
	assertCommandEndState(t, build, StatusBuildFailed, true)
}

func Test_run_test_command_with_no_failure(t *testing.T) {
	initTcrEngineWithFakes(failures{})
	assertCommandEndState(t, test, StatusOk, false)
}

func Test_run_test_command_with_failure(t *testing.T) {
	initTcrEngineWithFakes(failures{failTest})
	assertCommandEndState(t, test, StatusTestFailed, true)
}

func Test_run_commit_operation_with_no_failure(t *testing.T) {
	initTcrEngineWithFakes(failures{})
	assertOperationEndState(t, commit, StatusOk)
}

func Test_run_commit_operation_with_commit_failure(t *testing.T) {
	initTcrEngineWithFakes(failures{failCommit})
	assertOperationEndState(t, commit, StatusGitError)
}

func Test_run_commit_operation_with_push_failure(t *testing.T) {
	initTcrEngineWithFakes(failures{failPush})
	assertOperationEndState(t, commit, StatusGitError)
}

func Test_run_revert_operation_with_no_changes_in_src_files(t *testing.T) {
	initTcrEngineWithFakes(failures{})
	assertOperationEndState(t, revert, StatusOk)
}

func Test_run_revert_operation_with_changes_in_src_files(t *testing.T) {
	initTcrEngineWithFakes(failures{})
	assertOperationEndState(t, revert, StatusOk)
}

func Test_run_revert_operation_with_diff_failure(t *testing.T) {
	initTcrEngineWithFakes(failures{failDiff})
	assertOperationEndState(t, revert, StatusGitError)
}

func Test_run_revert_operation_with_restore_failure(t *testing.T) {
	initTcrEngineWithFakes(failures{failRestore})
	assertOperationEndState(t, revert, StatusGitError)
}

func Test_run_tcr_cycle_with_no_failure(t *testing.T) {
	initTcrEngineWithFakes(failures{})
	assertOperationEndState(t, RunTCRCycle, StatusOk)
}

func Test_run_tcr_cycle_with_build_failing(t *testing.T) {
	initTcrEngineWithFakes(failures{failBuild})
	assertOperationEndState(t, RunTCRCycle, StatusBuildFailed)
}

func Test_run_tcr_cycle_with_test_failing(t *testing.T) {
	initTcrEngineWithFakes(failures{failTest})
	assertOperationEndState(t, RunTCRCycle, StatusOk)
}

func Test_run_tcr_cycle_with_commit_failing(t *testing.T) {
	initTcrEngineWithFakes(failures{failCommit})
	assertOperationEndState(t, RunTCRCycle, StatusGitError)
}

func Test_run_tcr_cycle_with_push_failing(t *testing.T) {
	initTcrEngineWithFakes(failures{failPush})
	assertOperationEndState(t, RunTCRCycle, StatusGitError)
}

func Test_run_tcr_cycle_with_test_and_diff_failing(t *testing.T) {
	initTcrEngineWithFakes(failures{failTest, failDiff})
	assertOperationEndState(t, RunTCRCycle, StatusGitError)
}

func Test_run_tcr_cycle_with_test_and_restore_failing(t *testing.T) {
	initTcrEngineWithFakes(failures{failTest, failRestore})
	assertOperationEndState(t, RunTCRCycle, StatusGitError)
}

func initTcrEngineWithFakes(f failures) {
	tchn := registerFakeToolchain(f.contains(failBuild), f.contains(failTest))
	lang := registerFakeLanguage(tchn)
	Init(ui.NewFakeUI(), Params{Language: lang, Toolchain: tchn, Mode: runmode.OneShot{}})
	replaceGitImplWithFake(
		f.contains(failCommit),
		f.contains(failRestore),
		f.contains(failPush),
		f.contains(failPull),
		f.contains(failDiff),
	)
}

func registerFakeToolchain(failingBuild, failingTest bool) string {
	fake := toolchain.NewFakeToolchain(failingBuild, failingTest)
	if err := toolchain.Register(fake); err != nil {
		fmt.Println(err)
	}
	return fake.GetName()
}

func registerFakeLanguage(toolchainName string) string {
	fake := language.NewFakeLanguage(toolchainName)
	if err := language.Register(fake); err != nil {
		fmt.Println(err)
	}
	return fake.GetName()
}

func replaceGitImplWithFake(failingCommit, failingRestore, failingPush, failingPull, failDiff bool) {
	git, _ = vcs.NewGitFake(failingCommit, failingRestore, failingPush, failingPull, failDiff, []string{"fake-src"})
}

func assertCommandEndState(t *testing.T, operation func() error, endState Status, expectError bool) {
	assert.Equal(t, StatusOk, getCurrentState())
	err := operation()
	if expectError {
		assert.Error(t, err)
	} else {
		assert.Zero(t, err)
	}
	assert.Equal(t, endState, getCurrentState())
}

func assertOperationEndState(t *testing.T, operation func(), endState Status) {
	assert.Equal(t, StatusOk, getCurrentState())
	operation()
	assert.Equal(t, endState, getCurrentState())
}

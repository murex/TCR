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
	"github.com/murex/tcr/tcr-engine/params"
	"github.com/murex/tcr/tcr-engine/runmode"
	"github.com/murex/tcr/tcr-engine/status"
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
	tcr := initTcrEngineWithFakes(failures{})
	assertCommandEndState(t, tcr.build, status.Ok, false)
}

func Test_run_build_command_with_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{failBuild})
	assertCommandEndState(t, tcr.build, status.BuildFailed, true)
}

func Test_run_test_command_with_no_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{})
	assertCommandEndState(t, tcr.test, status.Ok, false)
}

func Test_run_test_command_with_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{failTest})
	assertCommandEndState(t, tcr.test, status.TestFailed, true)
}

func Test_run_commit_operation_with_no_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{})
	assertOperationEndState(t, tcr.commit, status.Ok)
}

func Test_run_commit_operation_with_commit_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{failCommit})
	assertOperationEndState(t, tcr.commit, status.GitError)
}

func Test_run_commit_operation_with_push_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{failPush})
	assertOperationEndState(t, tcr.commit, status.GitError)
}

func Test_run_revert_operation_with_no_changes_in_src_files(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{})
	assertOperationEndState(t, tcr.revert, status.Ok)
}

func Test_run_revert_operation_with_changes_in_src_files(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{})
	assertOperationEndState(t, tcr.revert, status.Ok)
}

func Test_run_revert_operation_with_diff_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{failDiff})
	assertOperationEndState(t, tcr.revert, status.GitError)
}

func Test_run_revert_operation_with_restore_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{failRestore})
	assertOperationEndState(t, tcr.revert, status.GitError)
}

func Test_run_tcr_cycle_with_no_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{})
	assertOperationEndState(t, tcr.RunTCRCycle, status.Ok)
}

func Test_run_tcr_cycle_with_build_failing(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{failBuild})
	assertOperationEndState(t, tcr.RunTCRCycle, status.BuildFailed)
}

func Test_run_tcr_cycle_with_test_failing(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{failTest})
	assertOperationEndState(t, tcr.RunTCRCycle, status.Ok)
}

func Test_run_tcr_cycle_with_commit_failing(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{failCommit})
	assertOperationEndState(t, tcr.RunTCRCycle, status.GitError)
}

func Test_run_tcr_cycle_with_push_failing(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{failPush})
	assertOperationEndState(t, tcr.RunTCRCycle, status.GitError)
}

func Test_run_tcr_cycle_with_test_and_diff_failing(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{failTest, failDiff})
	assertOperationEndState(t, tcr.RunTCRCycle, status.GitError)
}

func Test_run_tcr_cycle_with_test_and_restore_failing(t *testing.T) {
	tcr := initTcrEngineWithFakes(failures{failTest, failRestore})
	assertOperationEndState(t, tcr.RunTCRCycle, status.GitError)
}

func initTcrEngineWithFakes(f failures) TcrInterface {
	tchn := registerFakeToolchain(f.contains(failBuild), f.contains(failTest))
	lang := registerFakeLanguage(tchn)
	tcr := NewTcrEngine()
	tcr.Init(ui.NewFakeUI(), params.Params{Language: lang, Toolchain: tchn, Mode: runmode.OneShot{}})
	replaceGitImplWithFake(tcr,
		f.contains(failCommit),
		f.contains(failRestore),
		f.contains(failPush),
		f.contains(failPull),
		f.contains(failDiff),
	)
	return tcr
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

func replaceGitImplWithFake(tcr TcrInterface, failingCommit, failingRestore, failingPush, failingPull, failDiff bool) {
	fake, _ := vcs.NewGitFake(failingCommit, failingRestore, failingPush, failingPull, failDiff, []string{"fake-src"})
	tcr.setVcs(fake)
}

func assertCommandEndState(t *testing.T, operation func() error, endState status.Status, expectError bool) {
	assert.Equal(t, status.Ok, status.GetCurrentState())
	err := operation()
	if expectError {
		assert.Error(t, err)
	} else {
		assert.Zero(t, err)
	}
	assert.Equal(t, endState, status.GetCurrentState())
}

func assertOperationEndState(t *testing.T, operation func(), endState status.Status) {
	assert.Equal(t, status.Ok, status.GetCurrentState())
	operation()
	assert.Equal(t, endState, status.GetCurrentState())
}

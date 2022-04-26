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
	"github.com/murex/tcr/tcr-engine/events"
	"github.com/murex/tcr/tcr-engine/language"
	"github.com/murex/tcr/tcr-engine/params"
	"github.com/murex/tcr/tcr-engine/report"
	"github.com/murex/tcr/tcr-engine/role"
	"github.com/murex/tcr/tcr-engine/runmode"
	"github.com/murex/tcr/tcr-engine/settings"
	"github.com/murex/tcr/tcr-engine/status"
	"github.com/murex/tcr/tcr-engine/toolchain"
	"github.com/murex/tcr/tcr-engine/ui"
	"github.com/murex/tcr/tcr-engine/vcs"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
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
	tcr := initTcrEngineWithFakes(nil, failures{})
	assertCommandEndState(t, tcr.build, status.Ok, false)
}

func Test_run_build_command_with_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{failBuild})
	assertCommandEndState(t, tcr.build, status.BuildFailed, true)
}

func Test_run_test_command_with_no_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{})
	assertCommandEndState(t, tcr.test, status.Ok, false)
}

func Test_run_test_command_with_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{failTest})
	assertCommandEndState(t, tcr.test, status.TestFailed, true)
}

func Test_run_commit_operation_with_no_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{})
	assertOperationEndState(t, tcr.commit, status.Ok)
}

func Test_run_commit_operation_with_commit_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{failCommit})
	assertOperationEndState(t, tcr.commit, status.GitError)
}

func Test_run_commit_operation_with_push_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{failPush})
	assertOperationEndState(t, tcr.commit, status.GitError)
}

func Test_run_revert_operation_with_no_changes_in_src_files(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{})
	assertOperationEndState(t, tcr.revert, status.Ok)
}

func Test_run_revert_operation_with_changes_in_src_files(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{})
	assertOperationEndState(t, tcr.revert, status.Ok)
}

func Test_run_revert_operation_with_diff_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{failDiff})
	assertOperationEndState(t, tcr.revert, status.GitError)
}

func Test_run_revert_operation_with_restore_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{failRestore})
	assertOperationEndState(t, tcr.revert, status.GitError)
}

func Test_run_tcr_cycle_with_no_failure(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{})
	assertOperationEndState(t, tcr.RunTCRCycle, status.Ok)
}

func Test_run_tcr_cycle_with_build_failing(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{failBuild})
	assertOperationEndState(t, tcr.RunTCRCycle, status.BuildFailed)
}

func Test_run_tcr_cycle_with_test_failing(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{failTest})
	assertOperationEndState(t, tcr.RunTCRCycle, status.Ok)
}

func Test_run_tcr_cycle_with_commit_failing(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{failCommit})
	assertOperationEndState(t, tcr.RunTCRCycle, status.GitError)
}

func Test_run_tcr_cycle_with_push_failing(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{failPush})
	assertOperationEndState(t, tcr.RunTCRCycle, status.GitError)
}

func Test_run_tcr_cycle_with_test_and_diff_failing(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{failTest, failDiff})
	assertOperationEndState(t, tcr.RunTCRCycle, status.GitError)
}

func Test_run_tcr_cycle_with_test_and_restore_failing(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{failTest, failRestore})
	assertOperationEndState(t, tcr.RunTCRCycle, status.GitError)
}

func initTcrEngineWithFakes(p *params.Params, f failures) TcrInterface {
	tchn := registerFakeToolchain(f.contains(failBuild), f.contains(failTest))
	lang := registerFakeLanguage(tchn)
	events.EventRepository = &events.TcrEventInMemoryRepository{}

	var parameters params.Params
	if p == nil {
		parameters = *params.AParamSet(
			params.WithLanguage(lang),
			params.WithToolchain(tchn),
			params.WithRunMode(runmode.OneShot{}),
		)
	} else {
		parameters = *params.AParamSet(
			params.WithConfigDir(p.ConfigDir),
			params.WithBaseDir(p.BaseDir),
			params.WithWorkDir(p.WorkDir),
			params.WithLanguage(lang),
			params.WithToolchain(tchn),
			params.WithMobTimerDuration(p.MobTurnDuration),
			params.WithAutoPush(p.AutoPush),
			params.WithPollingPeriod(p.PollingPeriod),
			params.WithRunMode(p.Mode),
		)
	}

	tcr := NewTcrEngine()
	tcr.Init(ui.NewFakeUI(), parameters)
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
	fake, _ := vcs.NewGitFake(failingCommit, failingRestore, failingPush, failingPull, failDiff,
		[]vcs.FileDiff{vcs.NewFileDiff("fake-src", 1, 1)})
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

func Test_run_as_role_methods(t *testing.T) {
	var tcr TcrInterface
	testFlags := []struct {
		role        role.Role
		runAsMethod func()
	}{
		{role.Driver{}, func() { tcr.RunAsDriver() }},
		{role.Navigator{}, func() { tcr.RunAsNavigator() }},
	}
	for _, tt := range testFlags {
		t.Run(tt.role.LongName(), func(t *testing.T) {
			tcr = initTcrEngineWithFakes(params.AParamSet(params.WithRunMode(runmode.Mob{})), failures{})
			tt.runAsMethod()
			time.Sleep(10 * time.Millisecond)
			assert.Equal(t, tt.role, tcr.GetCurrentRole())
		})
	}
}

func Test_generate_tcr_event_on_build_fail(t *testing.T) {
	tcr := initTcrEngineWithFakes(params.AParamSet(params.WithRunMode(runmode.OneShot{})), failures{failBuild})
	tcr.RunTCRCycle()
	assert.Equal(t, events.StatusFailed, events.EventRepository.Get().BuildStatus)
}

func Test_generate_tcr_event_on_build_pass_and_tests_pass(t *testing.T) {
	tcr := initTcrEngineWithFakes(params.AParamSet(params.WithRunMode(runmode.OneShot{})), failures{})
	tcr.RunTCRCycle()
	assert.Equal(t, events.StatusPassed, events.EventRepository.Get().BuildStatus)
	assert.Equal(t, events.StatusPassed, events.EventRepository.Get().TestsStatus)
}

func Test_generate_tcr_event_on_build_pass_and_tests_fail(t *testing.T) {
	tcr := initTcrEngineWithFakes(params.AParamSet(params.WithRunMode(runmode.OneShot{})), failures{failTest})
	tcr.RunTCRCycle()
	assert.Equal(t, events.StatusPassed, events.EventRepository.Get().BuildStatus)
	assert.Equal(t, events.StatusFailed, events.EventRepository.Get().TestsStatus)
}

func Test_set_auto_push(t *testing.T) {
	var tcr TcrInterface
	testFlags := []struct {
		desc  string
		state bool
	}{
		{"Turn on", true},
		{"Turn off", false},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			tcr = initTcrEngineWithFakes(nil, failures{})
			tcr.SetAutoPush(tt.state)
			assert.Equal(t, tt.state, tcr.GetSessionInfo().AutoPush)
		})
	}
}

func Test_toggle_auto_push(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{})
	tcr.SetAutoPush(false)
	assert.Equal(t, false, tcr.GetSessionInfo().AutoPush)
	tcr.ToggleAutoPush()
	assert.Equal(t, true, tcr.GetSessionInfo().AutoPush)
	tcr.ToggleAutoPush()
	assert.Equal(t, false, tcr.GetSessionInfo().AutoPush)
}

func Test_get_session_info(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, failures{})
	currentDir, _ := os.Getwd()
	expected := SessionInfo{
		BaseDir:       currentDir,
		WorkDir:       currentDir,
		LanguageName:  "fake-language",
		ToolchainName: "fake-toolchain",
		AutoPush:      false,
		BranchName:    "",
	}
	assert.Equal(t, expected, tcr.GetSessionInfo())
}

func Test_mob_timer_duration_trace_at_startup(t *testing.T) {
	var tcr TcrInterface
	testFlags := []struct {
		role        role.Role
		runAsMethod func()
		timer       time.Duration
	}{
		{role.Driver{}, func() { tcr.RunAsDriver() }, 1 * time.Minute},
		//	{role.Navigator{}, func() { tcr.RunAsNavigator() }, 3 * time.Minute},
		{role.Driver{}, func() { tcr.RunAsDriver() }, 2 * time.Hour},
	}
	for _, tt := range testFlags {
		t.Run(tt.role.LongName()+" "+tt.timer.String(), func(t *testing.T) {
			settings.EnableMobTimer = true
			sniffer := report.NewFilteringSniffer(
				func(msg report.Message) bool {
					return msg.Type == report.Info && msg.Text == "Timer duration is "+tt.timer.String()
				},
			)
			tcr = initTcrEngineWithFakes(
				params.AParamSet(
					params.WithRunMode(runmode.Mob{}),
					params.WithMobTimerDuration(tt.timer),
				),
				failures{})
			tcr.SetRunMode(runmode.Mob{})
			tt.runAsMethod()
			time.Sleep(1 * time.Millisecond)
			sniffer.Stop()
			//fmt.Println(sniffer.GetAllMatches())
			assert.Equal(t, 1, sniffer.GetMatchCount())
		})
	}
}

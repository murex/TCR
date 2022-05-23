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

func Test_tcr_command_end_state(t *testing.T) {
	testFlags := []struct {
		desc           string
		command        func() error
		expectedStatus status.Status
		expectingError bool
	}{
		{
			"build with no failure",
			initTcrEngineWithFakes(nil, failures{}).build,
			status.Ok, false,
		},
		{
			"build with failure",
			initTcrEngineWithFakes(nil, failures{failBuild}).build,
			status.BuildFailed, true,
		},
		{
			"test with no failure",
			func() error {
				_, err := initTcrEngineWithFakes(nil, failures{}).test()
				return err
			},
			status.Ok, false,
		},
		{
			"test with failure",
			func() error {
				_, err := initTcrEngineWithFakes(nil, failures{failTest}).test()
				return err
			},

			status.TestFailed, true,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			status.RecordState(status.Ok)
			err := tt.command()
			if tt.expectingError {
				assert.Error(t, err)
			} else {
				assert.Zero(t, err)
			}
			assert.Equal(t, tt.expectedStatus, status.GetCurrentState())
		})
	}
}

func Test_tcr_operation_end_state(t *testing.T) {
	testFlags := []struct {
		desc           string
		operation      func()
		expectedStatus status.Status
	}{
		{
			"commit with no failure",
			initTcrEngineWithFakes(nil, failures{}).commit,
			status.Ok,
		},
		{
			"commit with git commit failure",
			initTcrEngineWithFakes(nil, failures{failCommit}).commit,
			status.GitError,
		},
		{
			"commit with git push failure",
			initTcrEngineWithFakes(nil, failures{failPush}).commit,
			status.GitError,
		},
		{
			"revert with no failure",
			initTcrEngineWithFakes(nil, failures{}).revert,
			status.Ok,
		},
		{
			"revert with git diff failure",
			initTcrEngineWithFakes(nil, failures{failDiff}).revert,
			status.GitError,
		},
		{
			"revert with git restore failure",
			initTcrEngineWithFakes(nil, failures{failRestore}).revert,
			status.GitError,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			status.RecordState(status.Ok)
			tt.operation()
			assert.Equal(t, tt.expectedStatus, status.GetCurrentState())
		})
	}
}

func Test_tcr_cycle_end_state(t *testing.T) {
	testFlags := []struct {
		desc           string
		failures       failures
		expectedStatus status.Status
	}{
		{"with no failure", failures{}, status.Ok},
		{"with build failure", failures{failBuild}, status.BuildFailed},
		{"with test failure", failures{failTest}, status.Ok},
		{"with git commit failure", failures{failCommit}, status.GitError},
		{"with git push failure", failures{failPush}, status.GitError},
		{"with test and git diff failure", failures{failTest, failDiff}, status.GitError},
		{"with test and git restore failure", failures{failTest, failRestore}, status.GitError},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			status.RecordState(status.Ok)
			initTcrEngineWithFakes(nil, tt.failures).RunTCRCycle()
			assert.Equal(t, tt.expectedStatus, status.GetCurrentState())
		})
	}
}

func initTcrEngineWithFakes(p *params.Params, f failures) TcrInterface {
	tchn := registerFakeToolchain(f.contains(failBuild), f.contains(failTest), toolchain.TestResults{})
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

func registerFakeToolchain(failingBuild, failingTest bool, testResults toolchain.TestResults) string {
	fake := toolchain.NewFakeToolchain(failingBuild, failingTest, testResults)
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

func Test_generate_events_from_tcr_cycle(t *testing.T) {
	testsFlags := []struct {
		desc                string
		failures            failures
		expectedBuildStatus events.TcrEventStatus
		expectedTestStatus  events.TcrEventStatus
	}{
		{
			"with build failing",
			failures{failBuild},
			events.StatusFailed,
			events.StatusUnknown,
		},
		{
			"with tests failing",
			failures{failTest},
			events.StatusPassed,
			events.StatusFailed,
		},
		{
			"with tests passing",
			failures{},
			events.StatusPassed,
			events.StatusPassed,
		},
	}

	for _, tt := range testsFlags {
		t.Run(tt.desc, func(t *testing.T) {
			tcr := initTcrEngineWithFakes(params.AParamSet(params.WithRunMode(runmode.OneShot{})), tt.failures)
			tcr.RunTCRCycle()
			assert.Equal(t, tt.expectedBuildStatus, events.EventRepository.GetLast().BuildStatus)
			assert.Equal(t, tt.expectedTestStatus, events.EventRepository.GetLast().TestsStatus)
		})
	}
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
			tcr = initTcrEngineWithFakes(params.AParamSet(
				params.WithRunMode(runmode.Mob{}),
				params.WithMobTimerDuration(tt.timer),
			), failures{})
			tcr.SetRunMode(runmode.Mob{})
			tt.runAsMethod()
			time.Sleep(1 * time.Millisecond)
			tcr.Stop()
			sniffer.Stop()
			//fmt.Println(sniffer.GetAllMatches())
			assert.Equal(t, 1, sniffer.GetMatchCount())
		})
	}
}

func Test_mob_timer_should_not_start_in_solo_mode(t *testing.T) {
	settings.EnableMobTimer = true
	sniffer := report.NewFilteringSniffer(
		func(msg report.Message) bool {
			return msg.Type == report.Info && msg.Text == "Mob Timer is off"
		},
	)
	tcr := initTcrEngineWithFakes(params.AParamSet(
		params.WithRunMode(runmode.Solo{}),
	), failures{})

	tcr.RunAsDriver()
	time.Sleep(1 * time.Millisecond)
	tcr.ReportMobTimerStatus()
	tcr.Stop()
	sniffer.Stop()
	//fmt.Println(sniffer.GetAllMatches())
	assert.Equal(t, 1, sniffer.GetMatchCount())
}

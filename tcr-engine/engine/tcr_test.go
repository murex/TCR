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

func Test_tcr_command_end_state(t *testing.T) {
	testFlags := []struct {
		desc           string
		command        func() error
		expectedStatus status.Status
		expectingError bool
	}{
		{
			"build with no failure",
			func() error {
				return initTcrEngineWithFakes(nil, nil, nil).build()
			},
			status.Ok, false,
		},
		{
			"build with failure",
			func() error {
				return initTcrEngineWithFakes(nil, toolchain.Operations{toolchain.BuildOperation}, nil).build()
			},
			status.BuildFailed, true,
		},
		{
			"test with no failure",
			func() error {
				_, err := initTcrEngineWithFakes(nil, nil, nil).test()
				return err
			},
			status.Ok, false,
		},
		{
			"test with failure",
			func() error {
				_, err := initTcrEngineWithFakes(nil, toolchain.Operations{toolchain.TestOperation}, nil).test()
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
			func() {
				tcr := initTcrEngineWithFakes(nil, nil, nil)
				tcr.commit(events.TcrEvent{})
			},
			status.Ok,
		},
		{
			"commit with git commit failure",
			func() {
				tcr := initTcrEngineWithFakes(nil, nil, vcs.GitCommands{vcs.CommitCommand})
				tcr.commit(events.TcrEvent{})
			},
			status.GitError,
		},
		{
			"commit with git push failure",
			func() {
				tcr := initTcrEngineWithFakes(nil, nil, vcs.GitCommands{vcs.PushCommand})
				tcr.commit(events.TcrEvent{})
			},
			status.GitError,
		},
		{
			"revert with no failure",
			func() {
				tcr := initTcrEngineWithFakes(nil, nil, nil)
				tcr.revert(events.TcrEvent{})
			},
			status.Ok,
		},
		{
			"revert with git diff failure",
			func() {
				tcr := initTcrEngineWithFakes(nil, nil, vcs.GitCommands{vcs.DiffCommand})
				tcr.revert(events.TcrEvent{})
			},
			status.GitError,
		},
		{
			"revert with git restore failure",
			func() {
				tcr := initTcrEngineWithFakes(nil, nil, vcs.GitCommands{vcs.RestoreCommand})
				tcr.revert(events.TcrEvent{})
			},
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

func Test_tcr_revert_end_state_with_commit_on_fail_enabled(t *testing.T) {
	testFlags := []struct {
		desc           string
		gitFailures    vcs.GitCommands
		expectedStatus status.Status
	}{
		{
			"no failure",
			nil,
			status.Ok,
		},
		{
			"git stash failure",
			vcs.GitCommands{vcs.StashCommand},
			status.GitError,
		},
		{
			"git un-stash failure",
			vcs.GitCommands{vcs.UnStashCommand},
			status.GitError,
		},
		{
			"git add failure",
			vcs.GitCommands{vcs.AddCommand},
			status.GitError,
		},
		{
			"git commit failure",
			vcs.GitCommands{vcs.CommitCommand},
			status.GitError,
		},
		{
			"git revert failure",
			vcs.GitCommands{vcs.RevertCommand},
			status.GitError,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			status.RecordState(status.Ok)
			tcr := initTcrEngineWithFakes(nil, nil, tt.gitFailures)
			tcr.SetCommitOnFail(true)
			tcr.revert(events.TcrEvent{})
			assert.Equal(t, tt.expectedStatus, status.GetCurrentState())
		})
	}
}

func Test_tcr_cycle_end_state(t *testing.T) {
	testFlags := []struct {
		desc              string
		toolchainFailures toolchain.Operations
		gitFailures       vcs.GitCommands
		expectedStatus    status.Status
	}{
		{
			"with no failure",
			nil, nil,
			status.Ok,
		},
		{
			"with build failure",
			toolchain.Operations{toolchain.BuildOperation}, nil,
			status.BuildFailed,
		},
		{
			"with test failure",
			toolchain.Operations{toolchain.TestOperation}, nil,
			status.Ok,
		},
		{
			"with git add failure",
			nil, vcs.GitCommands{vcs.AddCommand},
			status.GitError,
		},
		{
			"with git commit failure",
			nil, vcs.GitCommands{vcs.CommitCommand},
			status.GitError,
		},
		{
			"with git push failure",
			nil, vcs.GitCommands{vcs.PushCommand},
			status.GitError,
		},
		{
			"with test and git diff failure",
			toolchain.Operations{toolchain.TestOperation}, vcs.GitCommands{vcs.DiffCommand},
			status.GitError,
		},
		{
			"with test and git restore failure",
			toolchain.Operations{toolchain.TestOperation}, vcs.GitCommands{vcs.RestoreCommand},
			status.GitError,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			status.RecordState(status.Ok)
			initTcrEngineWithFakes(nil, tt.toolchainFailures, tt.gitFailures).RunTCRCycle()
			assert.Equal(t, tt.expectedStatus, status.GetCurrentState())
		})
	}
}

func initTcrEngineWithFakes(p *params.Params, toolchainFailures toolchain.Operations, gitFailures vcs.GitCommands) TcrInterface {
	tchn := registerFakeToolchain(toolchainFailures)
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
			params.WithCommitFailures(p.CommitFailures),
			params.WithPollingPeriod(p.PollingPeriod),
			params.WithRunMode(p.Mode),
		)
	}

	tcr := NewTcrEngine()
	tcr.Init(ui.NewFakeUI(), parameters)
	replaceGitImplWithFake(tcr, gitFailures)
	return tcr
}

func registerFakeToolchain(failures toolchain.Operations) string {
	fake := toolchain.NewFakeToolchain(failures, toolchain.TestResults{})
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

func replaceGitImplWithFake(tcr TcrInterface, failures vcs.GitCommands) {
	fake, _ := vcs.NewGitFake(failures, []vcs.FileDiff{vcs.NewFileDiff("fake-src", 1, 1)})
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
			tcr = initTcrEngineWithFakes(params.AParamSet(params.WithRunMode(runmode.Mob{})), nil, nil)
			tt.runAsMethod()
			time.Sleep(10 * time.Millisecond)
			assert.Equal(t, tt.role, tcr.GetCurrentRole())
		})
	}
}

func Test_generate_events_from_tcr_cycle(t *testing.T) {
	testsFlags := []struct {
		desc                string
		toolchainFailures   toolchain.Operations
		expectedBuildStatus events.TcrEventStatus
		expectedTestStatus  events.TcrEventStatus
	}{
		{
			"with build failing",
			toolchain.Operations{toolchain.BuildOperation},
			events.StatusFailed,
			events.StatusUnknown,
		},
		{
			"with tests failing",
			toolchain.Operations{toolchain.TestOperation},
			events.StatusPassed,
			events.StatusFailed,
		},
		{
			"with tests passing",
			nil,
			events.StatusPassed,
			events.StatusPassed,
		},
	}

	for _, tt := range testsFlags {
		t.Run(tt.desc, func(t *testing.T) {
			tcr := initTcrEngineWithFakes(
				params.AParamSet(params.WithRunMode(runmode.OneShot{})),
				tt.toolchainFailures, nil,
			)
			tcr.RunTCRCycle()
			assert.Equal(t, tt.expectedBuildStatus, events.EventRepository.Get().BuildStatus)
			assert.Equal(t, tt.expectedTestStatus, events.EventRepository.Get().TestsStatus)
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
			tcr = initTcrEngineWithFakes(nil, nil, nil)
			tcr.SetAutoPush(tt.state)
			assert.Equal(t, tt.state, tcr.GetSessionInfo().AutoPush)
		})
	}
}

func Test_set_commit_on_fail(t *testing.T) {
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
			tcr = initTcrEngineWithFakes(nil, nil, nil)
			tcr.SetCommitOnFail(tt.state)
			assert.Equal(t, tt.state, tcr.GetSessionInfo().CommitOnFail)
		})
	}
}

func Test_toggle_auto_push(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, nil, nil)
	tcr.SetAutoPush(false)
	assert.Equal(t, false, tcr.GetSessionInfo().AutoPush)
	tcr.ToggleAutoPush()
	assert.Equal(t, true, tcr.GetSessionInfo().AutoPush)
	tcr.ToggleAutoPush()
	assert.Equal(t, false, tcr.GetSessionInfo().AutoPush)
}

func Test_get_session_info(t *testing.T) {
	tcr := initTcrEngineWithFakes(nil, nil, nil)
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
	t.Skip("Dangling test on GitHub actions")
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
			), nil, nil)
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
	tcr := initTcrEngineWithFakes(
		params.AParamSet(params.WithRunMode(runmode.Solo{})),
		nil, nil,
	)
	tcr.RunAsDriver()
	time.Sleep(1 * time.Millisecond)
	tcr.ReportMobTimerStatus()
	tcr.Stop()
	sniffer.Stop()
	//fmt.Println(sniffer.GetAllMatches())
	assert.Equal(t, 1, sniffer.GetMatchCount())
}

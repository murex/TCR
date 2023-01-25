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
	"github.com/murex/tcr/events"
	"github.com/murex/tcr/language"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/role"
	"github.com/murex/tcr/runmode"
	"github.com/murex/tcr/settings"
	"github.com/murex/tcr/status"
	"github.com/murex/tcr/toolchain"
	"github.com/murex/tcr/ui"
	"github.com/murex/tcr/vcs"
	"github.com/murex/tcr/vcs/git"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"time"
)

func slowTestTag(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

func Test_tcr_command_end_state(t *testing.T) {
	testFlags := []struct {
		desc              string
		command           func() toolchain.CommandResult
		expectedCmdStatus toolchain.CommandStatus
		expectedAppStatus status.Status
	}{
		{
			"build with no failure",
			func() toolchain.CommandResult {
				tcr, _ := initTCREngineWithFakes(nil, nil, nil, nil)
				return tcr.build()
			},
			toolchain.CommandStatusPass, status.Ok,
		},
		{
			"build with failure",
			func() toolchain.CommandResult {
				tcr, _ := initTCREngineWithFakes(nil, toolchain.Operations{toolchain.BuildOperation}, nil, nil)
				return tcr.build()
			},
			toolchain.CommandStatusFail, status.BuildFailed,
		},
		{
			"test with no failure",
			func() toolchain.CommandResult {
				tcr, _ := initTCREngineWithFakes(nil, nil, nil, nil)
				result := tcr.test()
				return result.CommandResult
			},
			toolchain.CommandStatusPass, status.Ok,
		},
		{
			"test with failure",
			func() toolchain.CommandResult {
				tcr, _ := initTCREngineWithFakes(nil, toolchain.Operations{toolchain.TestOperation}, nil, nil)
				result := tcr.test()
				return result.CommandResult
			},
			toolchain.CommandStatusFail, status.TestFailed,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			status.RecordState(status.Ok)
			assert.Equal(t, tt.expectedCmdStatus, tt.command().Status)
			assert.Equal(t, tt.expectedAppStatus, status.GetCurrentState())
		})
	}
}

func Test_tcr_reports_and_emphasises(t *testing.T) {
	testFlags := []struct {
		desc              string
		failAt            toolchain.Operation
		isExpectedMessage func(report.Message) bool
	}{
		{
			desc:   "reports build failures as warnings",
			failAt: toolchain.BuildOperation,
			isExpectedMessage: func(msg report.Message) bool {
				return msg.Text == buildFailureMessage && msg.Type.Severity == report.Warning
			},
		},
		{
			desc:   "emphasises build failures",
			failAt: toolchain.BuildOperation,
			isExpectedMessage: func(msg report.Message) bool {
				return msg.Text == buildFailureMessage && msg.Type.Emphasis
			},
		},
		{
			desc:   "reports test successes as success",
			failAt: toolchain.Never,
			isExpectedMessage: func(msg report.Message) bool {
				return msg.Text == testSuccessMessage && msg.Type.Severity == report.Success
			},
		},
		{
			desc:   "emphasises test successes",
			failAt: toolchain.Never,
			isExpectedMessage: func(msg report.Message) bool {
				return msg.Text == testSuccessMessage && msg.Type.Emphasis
			},
		},
		{
			desc:   "reports test failures as errors",
			failAt: toolchain.TestOperation,
			isExpectedMessage: func(msg report.Message) bool {
				return msg.Text == testFailureMessage && msg.Type.Severity == report.Error
			},
		},
		{
			desc:   "emphasises test failures",
			failAt: toolchain.TestOperation,
			isExpectedMessage: func(msg report.Message) bool {
				return msg.Text == testFailureMessage && msg.Type.Emphasis
			},
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			status.RecordState(status.Ok)

			sniffer := report.NewSniffer(tt.isExpectedMessage)

			tcr, _ := initTCREngineWithFakes(nil, toolchain.Operations{tt.failAt}, nil, nil)
			tcr.build()
			tcr.test()

			time.Sleep(1 * time.Millisecond)
			sniffer.Stop()

			assert.Equal(t, 1, sniffer.GetMatchCount())
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
				tcr, _ := initTCREngineWithFakes(nil, nil, nil, nil)
				tcr.commit(events.TCREvent{})
			},
			status.Ok,
		},
		{
			"commit with VCS commit failure",
			func() {
				tcr, _ := initTCREngineWithFakes(nil, nil, git.Commands{git.CommitCommand}, nil)
				tcr.commit(events.TCREvent{})
			},
			status.VCSError,
		},
		{
			"commit with VCS push failure",
			func() {
				tcr, _ := initTCREngineWithFakes(nil, nil, git.Commands{git.PushCommand}, nil)
				tcr.commit(events.TCREvent{})
			},
			status.VCSError,
		},
		{
			"revert with no failure",
			func() {
				tcr, _ := initTCREngineWithFakes(nil, nil, nil, nil)
				tcr.revert(events.TCREvent{})
			},
			status.Ok,
		},
		{
			"revert with VCS diff failure",
			func() {
				tcr, _ := initTCREngineWithFakes(nil, nil, git.Commands{git.DiffCommand}, nil)
				tcr.revert(events.TCREvent{})
			},
			status.VCSError,
		},
		{
			"revert with VCS restore failure",
			func() {
				tcr, _ := initTCREngineWithFakes(nil, nil, git.Commands{git.RestoreCommand}, nil)
				tcr.revert(events.TCREvent{})
			},
			status.VCSError,
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
		vcsFailures    git.Commands
		expectedStatus status.Status
	}{
		{
			"no failure",
			nil,
			status.Ok,
		},
		{
			"VCS stash failure",
			git.Commands{git.StashCommand},
			status.VCSError,
		},
		{
			"VCS un-stash failure",
			git.Commands{git.UnStashCommand},
			status.VCSError,
		},
		{
			"VCS add failure",
			git.Commands{git.AddCommand},
			status.VCSError,
		},
		{
			"VCS commit failure",
			git.Commands{git.CommitCommand},
			status.VCSError,
		},
		{
			"VCS revert failure",
			git.Commands{git.RevertCommand},
			status.VCSError,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			status.RecordState(status.Ok)
			tcr, _ := initTCREngineWithFakes(nil, nil, tt.vcsFailures, nil)
			tcr.SetCommitOnFail(true)
			tcr.revert(events.TCREvent{})
			assert.Equal(t, tt.expectedStatus, status.GetCurrentState())
		})
	}
}

func Test_tcr_cycle_end_state(t *testing.T) {
	testFlags := []struct {
		desc              string
		toolchainFailures toolchain.Operations
		vcsFailures       git.Commands
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
			"with VCS add failure",
			nil, git.Commands{git.AddCommand},
			status.VCSError,
		},
		{
			"with VCS commit failure",
			nil, git.Commands{git.CommitCommand},
			status.VCSError,
		},
		{
			"with VCS push failure",
			nil, git.Commands{git.PushCommand},
			status.VCSError,
		},
		{
			"with test and VCS diff failure",
			toolchain.Operations{toolchain.TestOperation}, git.Commands{git.DiffCommand},
			status.VCSError,
		},
		{
			"with test and VCS restore failure",
			toolchain.Operations{toolchain.TestOperation}, git.Commands{git.RestoreCommand},
			status.VCSError,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			status.RecordState(status.Ok)
			tcr, _ := initTCREngineWithFakes(nil, tt.toolchainFailures, tt.vcsFailures, nil)
			tcr.RunTCRCycle()
			assert.Equal(t, tt.expectedStatus, status.GetCurrentState())
		})
	}
}

func initTCREngineWithFakes(
	p *params.Params,
	toolchainFailures toolchain.Operations,
	vcsFailures git.Commands,
	logItems vcs.LogItems,
) (TCRInterface, *git.Fake) {
	tchn := registerFakeToolchain(toolchainFailures)
	lang := registerFakeLanguage(tchn)

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

	tcr := NewTCREngine()
	tcr.Init(ui.NewFakeUI(), parameters)
	vcsFake := replaceVCSImplWithFake(tcr, vcsFailures, logItems)
	return tcr, vcsFake
}

func registerFakeToolchain(failures toolchain.Operations) string {
	fake := toolchain.NewFakeToolchain(failures, toolchain.TestStats{})
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

func replaceVCSImplWithFake(tcr TCRInterface, failures git.Commands, logItems vcs.LogItems) *git.Fake {
	fakeSettings := git.FakeSettings{
		FailingCommands: failures,
		ChangedFiles:    vcs.FileDiffs{vcs.NewFileDiff("fake-src", 1, 1)},
		Logs:            logItems,
	}
	fake, _ := git.NewFake(fakeSettings)
	tcr.setVCS(fake)
	return fake
}

func Test_run_as_role_methods(t *testing.T) {
	var tcr TCRInterface
	testFlags := []struct {
		role        role.Role
		runAsMethod func()
	}{
		{role.Driver{}, func() { tcr.RunAsDriver() }},
		{role.Navigator{}, func() { tcr.RunAsNavigator() }},
	}
	for _, tt := range testFlags {
		t.Run(tt.role.LongName(), func(t *testing.T) {
			tcr, _ = initTCREngineWithFakes(params.AParamSet(params.WithRunMode(runmode.Mob{})), nil, nil, nil)
			tt.runAsMethod()
			time.Sleep(10 * time.Millisecond)
			assert.Equal(t, tt.role, tcr.GetCurrentRole())
		})
	}
}

func Test_set_auto_push(t *testing.T) {
	var tcr TCRInterface
	testFlags := []struct {
		desc  string
		state bool
	}{
		{"Turn on", true},
		{"Turn off", false},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			tcr, _ = initTCREngineWithFakes(nil, nil, nil, nil)
			tcr.SetAutoPush(tt.state)
			assert.Equal(t, tt.state, tcr.GetSessionInfo().AutoPush)
		})
	}
}

func Test_set_commit_on_fail(t *testing.T) {
	var tcr TCRInterface
	testFlags := []struct {
		desc  string
		state bool
	}{
		{"Turn on", true},
		{"Turn off", false},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			tcr, _ = initTCREngineWithFakes(nil, nil, nil, nil)
			tcr.SetCommitOnFail(tt.state)
			assert.Equal(t, tt.state, tcr.GetSessionInfo().CommitOnFail)
		})
	}
}

func Test_vcs_pull_calls_vcs_command(t *testing.T) {
	tcr, vcsFake := initTCREngineWithFakes(nil, nil, nil, nil)
	tcr.VCSPull()
	assert.Equal(t, git.PullCommand, vcsFake.GetLastCommand())
}

func Test_vcs_pull_highlights_errors(t *testing.T) {
	sniffer := report.NewSniffer(
		func(msg report.Message) bool {
			return msg.Type.Severity == report.Error && msg.Text == "VCS pull command failed!"
		},
	)
	tcr, _ := initTCREngineWithFakes(nil, nil, git.Commands{git.PullCommand}, nil)
	tcr.VCSPull()
	time.Sleep(1 * time.Millisecond)
	sniffer.Stop()
	assert.Equal(t, 1, sniffer.GetMatchCount())
}

func Test_vcs_push_highlights_errors(t *testing.T) {
	sniffer := report.NewSniffer(
		func(msg report.Message) bool {
			return msg.Type.Severity == report.Error && msg.Text == "VCS push command failed!"
		},
	)
	tcr, _ := initTCREngineWithFakes(nil, nil, git.Commands{git.PushCommand}, nil)
	tcr.VCSPush()
	time.Sleep(1 * time.Millisecond)
	sniffer.Stop()
	assert.Equal(t, 1, sniffer.GetMatchCount())
}

func Test_toggle_auto_push(t *testing.T) {
	tcr, _ := initTCREngineWithFakes(nil, nil, nil, nil)
	tcr.SetAutoPush(false)
	assert.Equal(t, false, tcr.GetSessionInfo().AutoPush)
	tcr.ToggleAutoPush()
	assert.Equal(t, true, tcr.GetSessionInfo().AutoPush)
	tcr.ToggleAutoPush()
	assert.Equal(t, false, tcr.GetSessionInfo().AutoPush)
}

func Test_get_session_info(t *testing.T) {
	tcr, _ := initTCREngineWithFakes(nil, nil, nil, nil)
	currentDir, _ := os.Getwd()
	expected := SessionInfo{
		BaseDir:       currentDir,
		WorkDir:       currentDir,
		LanguageName:  "fake-language",
		ToolchainName: "fake-toolchain",
		AutoPush:      false,
		BranchName:    "master",
	}
	assert.Equal(t, expected, tcr.GetSessionInfo())
}

func Test_mob_timer_duration_trace_at_startup(t *testing.T) {
	slowTestTag(t)
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("disabled on CI due to dangling results")
	}
	var tcr TCRInterface
	testFlags := []struct {
		timer time.Duration
	}{
		{1 * time.Minute},
		{2 * time.Hour},
	}
	for _, tt := range testFlags {
		t.Run("duration "+tt.timer.String(), func(t *testing.T) {
			settings.EnableMobTimer = true
			sniffer := report.NewSniffer(
				func(msg report.Message) bool {
					return msg.Type.Severity == report.Info && msg.Text == "Timer duration is "+tt.timer.String()
				},
			)
			tcr, _ = initTCREngineWithFakes(params.AParamSet(
				params.WithRunMode(runmode.Mob{}),
				params.WithMobTimerDuration(tt.timer),
			), nil, nil, nil)
			tcr.RunAsDriver()
			time.Sleep(1 * time.Millisecond)
			tcr.ReportMobTimerStatus()
			tcr.Stop()
			sniffer.Stop()
			//fmt.Println(sniffer.GetAllMatches())
			assert.Equal(t, 1, sniffer.GetMatchCount())
		})
	}
}

func Test_mob_timer_should_not_start_in_solo_mode(t *testing.T) {
	slowTestTag(t)
	settings.EnableMobTimer = true
	sniffer := report.NewSniffer(
		func(msg report.Message) bool {
			return msg.Type.Severity == report.Info && msg.Text == "Mob Timer is off"
		},
	)
	tcr, _ := initTCREngineWithFakes(params.AParamSet(params.WithRunMode(runmode.Solo{})), nil, nil, nil)
	tcr.RunAsDriver()
	time.Sleep(1 * time.Millisecond)
	tcr.ReportMobTimerStatus()
	tcr.Stop()
	sniffer.Stop()
	//fmt.Println(sniffer.GetAllMatches())
	assert.Equal(t, 1, sniffer.GetMatchCount())
}

func Test_tcr_print_log(t *testing.T) {
	slowTestTag(t)
	now := time.Now()
	sampleItems := vcs.LogItems{
		vcs.NewLogItem("1111", now, "✅ TCR - tests passing"),
		vcs.NewLogItem("2222", now, "❌ TCR - tests failing"),
		vcs.NewLogItem("3333", now, "⏪ TCR - revert changes"),
		vcs.NewLogItem("4444", now, "other commit message"),
	}
	testFlags := []struct {
		desc            string
		filter          func(msg report.Message) bool
		logItems        vcs.LogItems
		expectedMatches int
	}{
		{
			desc: "TCR passing commits are kept",
			filter: func(msg report.Message) bool {
				return msg.Type.Severity == report.Info && strings.Index(msg.Text, "message:   ✅ TCR - tests passing") == 0
			},
			logItems:        sampleItems,
			expectedMatches: 1,
		},
		{
			desc: "TCR failing commits are kept",
			filter: func(msg report.Message) bool {
				return msg.Type.Severity == report.Info && strings.Index(msg.Text, "message:   ❌ TCR - tests failing") == 0
			},
			logItems:        sampleItems,
			expectedMatches: 1,
		},
		{
			desc: "TCR revert commits are dropped",
			filter: func(msg report.Message) bool {
				return msg.Type.Severity == report.Info && strings.Index(msg.Text, "message:   ⏪ TCR - revert changes") == 0
			},
			logItems:        sampleItems,
			expectedMatches: 0,
		},
		{
			desc: "non-TCR commits are dropped",
			filter: func(msg report.Message) bool {
				return msg.Type.Severity == report.Info && strings.Index(msg.Text, "message:   other commit message") == 0
			},
			logItems:        sampleItems,
			expectedMatches: 0,
		},
		{
			desc: "commit hashtag is printed",
			filter: func(msg report.Message) bool {
				return msg.Type.Severity == report.Title && strings.Index(msg.Text, "commit:    1111") == 0
			},
			logItems:        sampleItems,
			expectedMatches: 1,
		},
		{
			desc: "commit timestamp is printed",
			filter: func(msg report.Message) bool {
				return msg.Type.Severity == report.Info && strings.Index(msg.Text, "timestamp: "+now.String()) == 0
			},
			logItems:        sampleItems,
			expectedMatches: 2,
		},
		{
			desc: "warning when no record found",
			filter: func(msg report.Message) bool {
				return msg.Type.Severity == report.Warning && strings.Index(msg.Text, "no TCR commit found in branch") == 0
			},
			logItems:        nil,
			expectedMatches: 1,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			sniffer := report.NewSniffer(tt.filter)
			p := params.AParamSet(params.WithRunMode(runmode.Log{}))
			tcr, _ := initTCREngineWithFakes(p, nil, nil, tt.logItems)
			tcr.PrintLog(*p)
			time.Sleep(1 * time.Millisecond)
			sniffer.Stop()
			//fmt.Println(sniffer.GetAllMatches())
			assert.Equal(t, tt.expectedMatches, sniffer.GetMatchCount())
		})
	}
}

func Test_parse_commit_message(t *testing.T) {
	testFlags := []struct {
		desc          string
		commitMessage string
		expected      events.TCREvent
	}{
		{
			desc:          "empty commit message",
			commitMessage: "",
			expected:      events.TCREvent{Status: events.StatusUnknown},
		},
		{
			desc: "test-passing commit",
			commitMessage: "✅ TCR - tests passing\n" +
				"\n" +
				"changed-lines:\n" +
				"    src: 2\n" +
				"    test: 7\n" +
				"test-stats:\n" +
				"    run: 1\n" +
				"    passed: 1\n" +
				"    failed: 0\n" +
				"    skipped: 1\n" +
				"    error: 0\n" +
				"    duration: 2ms\n" +
				"\n",
			expected: events.NewTCREvent(
				events.StatusPass,
				events.NewChangedLines(2, 7),
				events.NewTestStats(1, 1, 0, 1, 0, 2*time.Millisecond),
			),
		},
		{
			desc: "test-failing commit",
			commitMessage: "❌ TCR - tests failing\n" +
				"\n" +
				"changed-lines:\n" +
				"    src: 1\n" +
				"    test: 3\n" +
				"test-stats:\n" +
				"    run: 10\n" +
				"    passed: 8\n" +
				"    failed: 2\n" +
				"    skipped: 1\n" +
				"    error: 0\n" +
				"    duration: 40ms\n" +
				"\n",
			expected: events.NewTCREvent(
				events.StatusFail,
				events.NewChangedLines(1, 3),
				events.NewTestStats(10, 8, 2, 1, 0, 40*time.Millisecond),
			),
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			event := parseCommitMessage(tt.commitMessage)
			assert.Equal(t, tt.expected, event)
		})
	}
}

func Test_compute_number_of_changed_lines(t *testing.T) {
	testFlags := []struct {
		desc         string
		diffs        vcs.FileDiffs
		expectedSrc  int
		expectedTest int
	}{
		{"0 file", nil, 0, 0},
		{
			"1 src file - 0 test file",
			vcs.FileDiffs{
				vcs.NewFileDiff("src/main/Hello.java", 1, 2),
			},
			3, 0,
		},
		{
			"2 src files - 0 test file",
			vcs.FileDiffs{
				vcs.NewFileDiff("src/main/Hello1.java", 1, 2),
				vcs.NewFileDiff("src/main/Hello2.java", 3, 4),
			},
			10, 0,
		},
		{
			"1 src file - 1 test file",
			vcs.FileDiffs{
				vcs.NewFileDiff("src/main/Hello.java", 1, 2),
				vcs.NewFileDiff("src/test/HelloTest.java", 3, 4),
			},
			3, 7,
		},
		{
			"0 src file - 1 test file",
			vcs.FileDiffs{
				vcs.NewFileDiff("src/test/HelloTest.java", 1, 2),
			},
			0, 3,
		},
		{
			"0 src file - 2 test files",
			vcs.FileDiffs{
				vcs.NewFileDiff("src/test/HelloTest1.java", 1, 2),
				vcs.NewFileDiff("src/test/HelloTest2.java", 3, 4),
			},
			0, 10,
		},
	}

	lang, _ := language.Get("java")

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expectedSrc, tt.diffs.ChangedLines(lang.IsSrcFile))
			assert.Equal(t, tt.expectedTest, tt.diffs.ChangedLines(lang.IsTestFile))
		})
	}
}

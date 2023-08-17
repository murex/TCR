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

package timer

import (
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/runmode"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_best_tick_period_for_timeout(t *testing.T) {
	var tickTests = []struct {
		timeout  time.Duration
		expected time.Duration
	}{
		{1 * time.Second, 1 * time.Second},
		{10 * time.Second, 1 * time.Second},
		{11 * time.Second, 10 * time.Second},
		{1 * time.Minute, 10 * time.Second},
		{1*time.Minute + 1*time.Second, 1 * time.Minute},
		{10 * time.Minute, 1 * time.Minute},
	}

	for _, tt := range tickTests {
		t.Run(tt.timeout.String()+"->"+tt.expected.String(), func(t *testing.T) {
			assert.Equal(t, tt.expected, findBestTickPeriodFor(tt.timeout))
		})
	}
}

func Test_format_duration(t *testing.T) {
	var fmtTests = []struct {
		d        time.Duration
		expected string
	}{
		{0 * time.Second, "0s"},
		{59 * time.Second, "59s"},
		{60 * time.Second, "1m"},
		{61 * time.Second, "1m1s"},
		{1 * time.Hour, "1h"},
		{1*time.Hour + 1*time.Second, "1h0m1s"},
		{1*time.Hour + 1*time.Minute, "1h1m"},
	}
	for _, tt := range fmtTests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, fmtDuration(tt.d))
		})
	}
}

func Test_mob_turn_countdown_creation_in_mob_runmode(t *testing.T) {
	assert.NotZero(t, NewMobTurnCountdown(runmode.Mob{}, defaultTimeout))
}

func Test_mob_turn_countdown_creation_in_solo_runmode(t *testing.T) {
	assert.Zero(t, NewMobTurnCountdown(runmode.Solo{}, defaultTimeout))
}

func Test_mob_turn_countdown_creation_in_check_runmode(t *testing.T) {
	assert.Zero(t, NewMobTurnCountdown(runmode.Check{}, defaultTimeout))
}

func Test_mob_turn_countdown_creation_in_one_shot_runmode(t *testing.T) {
	assert.Zero(t, NewMobTurnCountdown(runmode.OneShot{}, defaultTimeout))
}

func Test_report_count_down_status_when_mob_is_not_started(t *testing.T) {
	sniffer := report.NewSniffer()
	ReportCountDownStatus(nil)
	sniffer.Stop()
	assert.Equal(t, 1, sniffer.GetMatchCount())
	assert.Equal(t, "Mob Timer is off", sniffer.GetAllMatches()[0].Text)
}

func Test_report_count_down_status_when_mob_is_started(t *testing.T) {
	tests := []struct {
		desc          string
		state         ReminderState
		expectedTrace string
	}{
		{
			"timer not started",
			NotStarted,
			"Mob Timer is not started",
		},
		{
			"timer running",
			Running,
			"Mob Timer: 0s done, 5m to go",
		},
		{
			"after timeout",
			AfterTimeOut,
			"Mob Timer has timed out: 5m over!",
		},
		{
			"timer stopped after interruption",
			StoppedAfterInterruption,
			"Mob Timer was interrupted",
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			sniffer := report.NewSniffer()
			reminder := NewPeriodicReminder(0, 0, func(ctx ReminderContext) {})
			reminder.state = test.state
			reminder.startTime = time.Now()
			ReportCountDownStatus(reminder)
			sniffer.Stop()
			assert.Equal(t, 1, sniffer.GetMatchCount())
			assert.Equal(t, test.expectedTrace, sniffer.GetAllMatches()[0].Text)
		})
	}
}

func Test_mob_turn_count_down(t *testing.T) {
	sniffer := report.NewSniffer()
	reminder := NewMobTurnCountdown(runmode.Mob{}, 2*time.Second)

	reminder.Start()
	time.Sleep(3200 * time.Millisecond)
	reminder.Stop()

	sniffer.Stop()
	assert.Equal(t, "(Mob Timer) Starting 2s countdown", sniffer.GetAllMatches()[0].Text)
	assert.Equal(t, "(Mob Timer) Your turn ends in 1s", sniffer.GetAllMatches()[1].Text)
	assert.Equal(t, "(Mob Timer) Time's up. Time to rotate! You are 0s over!", sniffer.GetAllMatches()[2].Text)
	assert.Equal(t, "(Mob Timer) Time's up. Time to rotate! You are 1s over!", sniffer.GetAllMatches()[3].Text)
	assert.Equal(t, "(Mob Timer) Stopping countdown after 3s", sniffer.GetAllMatches()[4].Text)
}

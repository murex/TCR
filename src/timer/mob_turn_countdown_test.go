/*
Copyright (c) 2024 Murex

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
	"testing"
	"time"

	"github.com/murex/tcr/report"
	"github.com/murex/tcr/runmode"
	"github.com/stretchr/testify/assert"
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

func Test_mob_turn_count_down(t *testing.T) {
	report.TestWithIsolatedReporter(func(reporter *report.Reporter, sniffer *report.Sniffer) {
		reminder := NewMobTurnCountdown(runmode.Mob{}, 2*time.Second)
		reminder.Start()
		time.Sleep(3200 * time.Millisecond)
		reminder.Stop()

		sniffer.Stop()

		expected := []struct {
			text     string
			category report.Category
			emphasis bool
		}{
			{"start:2:0:2", report.TimerEvent, true},
			{"countdown:2:1:1", report.TimerEvent, true},
			{"timeout:2:2:0", report.TimerEvent, true},
			{"timeout:2:3:-1", report.TimerEvent, false},
			{"stop:2:3:0", report.TimerEvent, true},
		}
		assert.Equal(t, len(expected), sniffer.GetMatchCount())
		for i, e := range expected {
			msg := sniffer.GetAllMatches()[i]
			assert.Equal(t, report.MessageType{Category: e.category, Emphasis: e.emphasis}, msg.Type)
			assert.Equal(t, e.text, msg.Payload.ToString())
		}
	})
}

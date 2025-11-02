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

package checker

import (
	"testing"
	"time"

	"github.com/murex/tcr/checker/model"
	"github.com/murex/tcr/params"
	"github.com/stretchr/testify/assert"
)

func Test_check_mob_configuration(t *testing.T) {
	assertCheckGroupRunner(t,
		checkMobConfiguration,
		&checkMobConfigurationRunners,
		*params.AParamSet(),
		"mob configuration")
}

func Test_check_mob_timer(t *testing.T) {
	tests := []struct {
		desc     string
		mobTimer time.Duration
		expected []model.CheckPoint
	}{
		{
			"0s turned off", 0,
			[]model.CheckPoint{
				model.OkCheckPoint("mob timer duration is set to 0s"),
				model.OkCheckPoint("mob timer is turned off"),
			},
		},
		{
			"2m too short", 2 * time.Minute,
			[]model.CheckPoint{
				model.OkCheckPoint("mob timer duration is set to 2m0s"),
				model.WarningCheckPoint("mob timer duration is quite short (under 3m0s)"),
			},
		},
		{
			"3m low threshold", 3 * time.Minute,
			[]model.CheckPoint{
				model.OkCheckPoint("mob timer duration is set to 3m0s"),
				model.OkCheckPoint("mob timer duration is in the recommended range"),
			},
		},
		{
			"10m in range", 10 * time.Minute,
			[]model.CheckPoint{
				model.OkCheckPoint("mob timer duration is set to 10m0s"),
				model.OkCheckPoint("mob timer duration is in the recommended range"),
			},
		},
		{
			"15m high threshold", 15 * time.Minute,
			[]model.CheckPoint{
				model.OkCheckPoint("mob timer duration is set to 15m0s"),
				model.OkCheckPoint("mob timer duration is in the recommended range"),
			},
		},
		{
			"20m too long", 20 * time.Minute,
			[]model.CheckPoint{
				model.OkCheckPoint("mob timer duration is set to 20m0s"),
				model.WarningCheckPoint("mob timer duration is quite long (above 15m0s)"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithMobTimerDuration(test.mobTimer))
			assert.Equal(t, test.expected, checkMobTimer(p))
		})
	}
}

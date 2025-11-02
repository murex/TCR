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
	"github.com/murex/tcr/vcs/git"
	"github.com/murex/tcr/vcs/p4"
	"github.com/stretchr/testify/assert"
)

func Test_check_vcs_configuration(t *testing.T) {
	assertCheckGroupRunner(t,
		checkVCSConfiguration,
		&checkVCSRunners,
		*params.AParamSet(),
		"VCS configuration")
}

func Test_check_vcs_selection(t *testing.T) {
	tests := []struct {
		desc     string
		vcsName  string
		expected []model.CheckPoint
	}{
		{
			"empty", "",
			[]model.CheckPoint{
				model.ErrorCheckPoint("no VCS is selected"),
			},
		},
		{
			"unknown", "unknown-vcs",
			[]model.CheckPoint{
				model.ErrorCheckPoint("selected VCS is not supported: \"unknown-vcs\""),
			},
		},
		{
			"git", git.Name,
			[]model.CheckPoint{
				model.OkCheckPoint("selected VCS is git"),
			},
		},
		{
			"GIT", "GIT",
			[]model.CheckPoint{
				model.OkCheckPoint("selected VCS is git"),
			},
		},
		{
			"p4", p4.Name,
			[]model.CheckPoint{
				model.OkCheckPoint("selected VCS is p4"),
			},
		},
		{
			"P4", "P4",
			[]model.CheckPoint{
				model.OkCheckPoint("selected VCS is p4"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithVCS(test.vcsName))
			assert.Equal(t, test.expected, checkVCSSelection(p))
		})
	}
}

func Test_check_vcs_polling_period(t *testing.T) {
	tests := []struct {
		desc          string
		pollingPeriod time.Duration
		expected      []model.CheckPoint
	}{
		{
			"0s turned off", 0,
			[]model.CheckPoint{
				model.OkCheckPoint("polling period is set to 0s"),
				model.OkCheckPoint("code refresh (for navigator role) is turned off"),
			},
		},
		{
			"1s too fast", 1 * time.Second,
			[]model.CheckPoint{
				model.OkCheckPoint("polling period is set to 1s"),
				model.WarningCheckPoint("polling is very fast (below 2s-period)"),
			},
		},
		{
			"2s low threshold", 2 * time.Second,
			[]model.CheckPoint{
				model.OkCheckPoint("polling period is set to 2s"),
				model.OkCheckPoint("polling period is in the recommended range"),
			},
		},
		{
			"30s in range", 30 * time.Second,
			[]model.CheckPoint{
				model.OkCheckPoint("polling period is set to 30s"),
				model.OkCheckPoint("polling period is in the recommended range"),
			},
		},
		{
			"1m high threshold", 1 * time.Minute,
			[]model.CheckPoint{
				model.OkCheckPoint("polling period is set to 1m0s"),
				model.OkCheckPoint("polling period is in the recommended range"),
			},
		},
		{
			"2m too slow", 2 * time.Minute,
			[]model.CheckPoint{
				model.OkCheckPoint("polling period is set to 2m0s"),
				model.WarningCheckPoint("polling is very slow (above 1m0s-period)"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithPollingPeriod(test.pollingPeriod))
			assert.Equal(t, test.expected, checkVCSPollingPeriod(p))
		})
	}
}

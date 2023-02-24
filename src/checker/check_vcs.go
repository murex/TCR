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
	"github.com/murex/tcr/checker/model"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/vcs/git"
	"github.com/murex/tcr/vcs/p4"
	"strings"
	"time"
)

var checkVCSRunners []checkPointRunner

func init() {
	checkVCSRunners = []checkPointRunner{
		checkVCSSelection,
		checkVCSPollingPeriod,
	}
}

func checkVCSConfiguration(p params.Params) (cg *model.CheckGroup) {
	cg = model.NewCheckGroup("VCS configuration")
	for _, runner := range checkVCSRunners {
		cg.Add(runner(p)...)
	}
	return cg
}

func checkVCSSelection(p params.Params) (cp []model.CheckPoint) {
	switch strings.ToLower(p.VCS) {
	case git.Name, p4.Name:
		cp = append(cp, model.OkCheckPoint("selected VCS is ", p.VCS))
	case "":
		cp = append(cp, model.ErrorCheckPoint("no VCS is selected"))
	default:
		cp = append(cp, model.ErrorCheckPoint("selected VCS is not supported: \"", p.VCS, "\""))
	}
	return cp
}

const (
	pollingPeriodLowThreshold  = 2 * time.Second
	pollingPeriodHighThreshold = 1 * time.Minute
)

func checkVCSPollingPeriod(p params.Params) (cp []model.CheckPoint) {
	cp = append(cp, model.OkCheckPoint("polling period is set to ", p.PollingPeriod.String()))
	if p.PollingPeriod == 0 {
		cp = append(cp, model.OkCheckPoint("code refresh (for navigator role) is turned off"))
	} else if p.PollingPeriod > pollingPeriodHighThreshold {
		cp = append(cp,
			model.WarningCheckPoint("polling is very slow (above ", pollingPeriodHighThreshold, "-period)"))
	} else if p.PollingPeriod < pollingPeriodLowThreshold {
		cp = append(cp,
			model.WarningCheckPoint("polling is very fast (below ", pollingPeriodLowThreshold, "-period)"))
	} else {
		cp = append(cp, model.OkCheckPoint("polling period is in the recommended range"))
	}
	return cp
}

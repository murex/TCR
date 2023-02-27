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
	"time"
)

var checkMobConfigurationRunners []checkPointRunner

func init() {
	checkMobConfigurationRunners = []checkPointRunner{
		checkMobTimer,
	}
}

func checkMobConfiguration(p params.Params) (cg *model.CheckGroup) {
	cg = model.NewCheckGroup("mob configuration")
	for _, runner := range checkMobConfigurationRunners {
		cg.Add(runner(p)...)
	}
	return cg
}

const (
	mobTimerLowThreshold  = 3 * time.Minute
	mobTimerHighThreshold = 15 * time.Minute
)

func checkMobTimer(p params.Params) (cp []model.CheckPoint) {
	timer := p.MobTurnDuration
	cp = append(cp, model.OkCheckPoint("mob timer duration is set to ", timer.String()))
	switch {
	case timer == 0:
		cp = append(cp, model.OkCheckPoint("mob timer is turned off"))
	case timer < mobTimerLowThreshold:
		cp = append(cp, model.WarningCheckPoint(
			"mob timer duration is quite short (under ", mobTimerLowThreshold, ")"))
	case timer > mobTimerHighThreshold:
		cp = append(cp, model.WarningCheckPoint(
			"mob timer duration is quite long (above ", mobTimerHighThreshold, ")"))
	default:
		cp = append(cp, model.OkCheckPoint("mob timer duration is in the recommended range"))
	}
	return cp
}

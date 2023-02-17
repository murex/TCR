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

package checker

import (
	"github.com/murex/tcr/checker/model"
	"github.com/murex/tcr/params"
	"time"
)

const (
	mobTimerLowThreshold  = 3 * time.Minute
	mobTimerHighThreshold = 15 * time.Minute
)

func checkMobTimer(p params.Params) (cg *model.CheckGroup) {
	cg = model.NewCheckGroup("mob timer")
	cg.Ok("mob timer duration is ", p.MobTurnDuration.String())
	if p.MobTurnDuration == 0 {
		cg.Warning("mob timer is turned off")
	} else if p.MobTurnDuration < mobTimerLowThreshold {
		cg.Warning("mob timer duration is quite short (under ", mobTimerLowThreshold, ")")
	} else if p.MobTurnDuration > mobTimerHighThreshold {
		cg.Warning("mob timer duration is quite long (above ", mobTimerHighThreshold, ")")
	} else {
		cg.Ok("mob timer duration is in the recommended range")
	}
	return cg
}

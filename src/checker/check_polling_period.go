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
	"github.com/murex/tcr/params"
	"time"
)

const (
	pollingPeriodLowThreshold  = 2 * time.Second
	pollingPeriodHighThreshold = 1 * time.Minute
)

func checkPollingPeriod(p params.Params) (cr *CheckResults) {
	cr = NewCheckResults("polling period")
	cr.ok("polling period is ", p.PollingPeriod.String())
	if p.PollingPeriod == 0 {
		cr.warning("code refresh (for navigator role) is turned off")
	} else if p.PollingPeriod > pollingPeriodHighThreshold {
		cr.warning("polling period is very slow (above ", pollingPeriodHighThreshold, ")")
	} else if p.PollingPeriod < pollingPeriodLowThreshold {
		cr.warning("polling period is very fast (below ", pollingPeriodLowThreshold, ")")
	} else {
		cr.ok("polling period is in the recommended range")
	}
	return cr
}

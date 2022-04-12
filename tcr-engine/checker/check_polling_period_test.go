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
	"github.com/murex/tcr/tcr-engine/params"
	"testing"
	"time"
)

func Test_check_polling_period_returns_warning_when_set_to_0(t *testing.T) {
	assertWarning(t, checkPollingPeriod, *params.AParamSet(params.WithPollingPeriod(0)))
}

func Test_check_polling_period_returns_warning_when_set_to_1s(t *testing.T) {
	assertWarning(t, checkPollingPeriod, *params.AParamSet(params.WithPollingPeriod(1 * time.Second)))
}

func Test_check_polling_period_returns_ok_when_set_to_2s(t *testing.T) {
	assertOk(t, checkPollingPeriod, *params.AParamSet(params.WithPollingPeriod(2 * time.Second)))
}

func Test_check_polling_period_returns_warning_when_set_to_2m(t *testing.T) {
	assertWarning(t, checkPollingPeriod, *params.AParamSet(params.WithPollingPeriod(2 * time.Minute)))
}

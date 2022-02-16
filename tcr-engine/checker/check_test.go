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
	"github.com/murex/tcr/tcr-engine/engine"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
	"time"
)

const (
	testDataRootDir = "../testdata"
)

var (
	testDataDirJava = filepath.Join(testDataRootDir, "java")
	//testDataDirGo = filepath.Join(testDataRootDir, "go")
	//testDataDirCpp  = filepath.Join(testDataRootDir, "cpp")
)

// Assert utility functions

func assertStatus(t *testing.T, expected CheckStatus, checker func(params engine.Params) (cr *CheckResults), params engine.Params) {
	initCheckEnv(params)
	assert.Equal(t, expected, checker(params).getStatus())
}

func assertOk(t *testing.T, checker func(params engine.Params) (cr *CheckResults), params engine.Params) {
	assertStatus(t, CheckStatusOk, checker, params)
}

func assertWarning(t *testing.T, checker func(params engine.Params) (cr *CheckResults), params engine.Params) {
	assertStatus(t, CheckStatusWarning, checker, params)
}

func assertError(t *testing.T, checker func(params engine.Params) (cr *CheckResults), params engine.Params) {
	assertStatus(t, CheckStatusError, checker, params)
}

// Return code for check subcommand

func Test_checker_should_return_0_if_no_error_or_warning(t *testing.T) {
	t.Skip("need to provide fake configuration settings for tests")
	Run(*engine.AParamSet(
		engine.WithConfigDir(testDataDirJava),
		engine.WithBaseDir(testDataDirJava),
		engine.WithWorkDir(testDataDirJava),
		engine.WithMobTimerDuration(mobTimerLowThreshold),
		engine.WithPollingPeriod(pollingPeriodLowThreshold),
	))
	assert.Equal(t, 0, engine.GetReturnCode())
}

func Test_checker_should_return_2_if_one_or_more_errors(t *testing.T) {
	Run(*engine.AParamSet(
		engine.WithConfigDir("invalid-dir"),
		engine.WithBaseDir("invalid-dir"),
		engine.WithWorkDir("invalid-dir"),
	))
	assert.Equal(t, 2, engine.GetReturnCode())
}

func Test_checker_should_return_1_if_one_or_more_warnings(t *testing.T) {
	// The warning is triggered by the mob timer duration being under the min threshold
	Run(*engine.AParamSet(
		engine.WithConfigDir(testDataDirJava),
		engine.WithBaseDir(testDataDirJava),
		engine.WithWorkDir(testDataDirJava),
		engine.WithMobTimerDuration(1*time.Second),
	))
	assert.Equal(t, 1, engine.GetReturnCode())
}

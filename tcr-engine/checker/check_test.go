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
	"github.com/murex/tcr/tcr-engine/status"
	"github.com/murex/tcr/tcr-engine/vcs"
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

func initTestCheckEnv(params params.Params) {
	initCheckEnv(params)
	// We replace git implementation with a fake so that we bypass real git access
	checkEnv.git, checkEnv.gitErr = vcs.NewGitFake(vcs.GitFakeSettings{})
}

func assertStatus(t *testing.T, expected CheckStatus, checker func(params params.Params) (cr *CheckResults), params params.Params) {
	initTestCheckEnv(params)
	assert.Equal(t, expected, checker(params).getStatus())
}

func assertOk(t *testing.T, checker func(params params.Params) (cr *CheckResults), params params.Params) {
	assertStatus(t, CheckStatusOk, checker, params)
}

func assertWarning(t *testing.T, checker func(params params.Params) (cr *CheckResults), params params.Params) {
	assertStatus(t, CheckStatusWarning, checker, params)
}

func assertError(t *testing.T, checker func(params params.Params) (cr *CheckResults), params params.Params) {
	assertStatus(t, CheckStatusError, checker, params)
}

// Return code for check subcommand

func Test_checker_should_return_0_if_no_error_or_warning(t *testing.T) {
	t.Skip("need to provide fake configuration settings for tests")
	Run(*params.AParamSet(
		params.WithConfigDir(testDataDirJava),
		params.WithBaseDir(testDataDirJava),
		params.WithWorkDir(testDataDirJava),
		params.WithMobTimerDuration(mobTimerLowThreshold),
		params.WithPollingPeriod(pollingPeriodLowThreshold),
	))
	assert.Equal(t, 0, status.GetReturnCode())
}

func Test_checker_should_return_1_if_one_or_more_warnings(t *testing.T) {
	// The warning is triggered by the mob timer duration being under the min threshold
	t.Skip("disabled due to git remote access check failing when running from CI")
	Run(*params.AParamSet(
		params.WithConfigDir(testDataDirJava),
		params.WithBaseDir(testDataDirJava),
		params.WithWorkDir(testDataDirJava),
		params.WithMobTimerDuration(1*time.Second),
	))
	assert.Equal(t, 1, status.GetReturnCode())
}

func Test_checker_should_return_2_if_one_or_more_errors(t *testing.T) {
	const invalidDir = "invalid-dir"
	Run(*params.AParamSet(
		params.WithConfigDir(invalidDir),
		params.WithBaseDir(invalidDir),
		params.WithWorkDir(invalidDir),
	))
	assert.Equal(t, 2, status.GetReturnCode())
}

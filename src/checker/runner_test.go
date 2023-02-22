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
	"github.com/murex/tcr/status"
	"github.com/murex/tcr/vcs"
	"github.com/murex/tcr/vcs/factory"
	"github.com/murex/tcr/vcs/fake"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
	"time"
)

const (
	testDataRootDir = "../testdata"
)

// stub checkpoint runners used for testing

var (
	checkPointRunnerOkStub = func(p params.Params) []model.CheckPoint {
		return []model.CheckPoint{model.OkCheckPoint("")}
	}
	checkPointRunnerWarningStub = func(p params.Params) []model.CheckPoint {
		return []model.CheckPoint{model.WarningCheckPoint("")}
	}
	checkPointRunnerErrorStub = func(p params.Params) []model.CheckPoint {
		return []model.CheckPoint{model.ErrorCheckPoint("")}
	}
)

// Assert utility functions

var (
	testDataDirJava = filepath.Join(testDataRootDir, "java")
)

func initTestCheckEnv(params params.Params) {
	// Replace VCS factory initializer in order to use a VCS fake instead of the real thing
	factory.InitVCS = func(_ string, _ string) (vcs.Interface, error) {
		return fake.NewVCSFake(fake.Settings{}), nil
	}
	initCheckEnv(params)
}

func assertStatus(t *testing.T, expected model.CheckStatus, checker checkGroupRunner, params params.Params) {
	initTestCheckEnv(params)
	assert.Equal(t, expected, checker(params).GetStatus())
}

func assertOk(t *testing.T, checker checkGroupRunner, params params.Params) {
	t.Helper()
	assertStatus(t, model.CheckStatusOk, checker, params)
}

// Return code for check subcommand

func assertWarning(t *testing.T, checker checkGroupRunner, params params.Params) {
	t.Helper()
	assertStatus(t, model.CheckStatusWarning, checker, params)
}

func assertError(t *testing.T, checker checkGroupRunner, params params.Params) {
	t.Helper()
	assertStatus(t, model.CheckStatusError, checker, params)
}

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
	Run(*params.AParamSet(
		params.WithConfigDir(testDataDirJava),
		params.WithBaseDir(testDataDirJava),
		params.WithWorkDir(testDataDirJava),
		params.WithMobTimerDuration(1*time.Second),
	))
	assert.Equal(t, status.GetReturnCode(), 1)
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

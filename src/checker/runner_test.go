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

	"github.com/murex/tcr/checker/model"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/status"
	"github.com/murex/tcr/vcs"
	"github.com/murex/tcr/vcs/factory"
	"github.com/murex/tcr/vcs/fake"
	"github.com/stretchr/testify/assert"
)

func initTestCheckEnv(params params.Params) {
	// Replace VCS factory initializer in order to use a VCS fake instead of the real thing
	factory.InitVCS = func(_ string, _ string, _ string) (vcs.Interface, error) {
		return fake.NewVCSFake(fake.Settings{}), nil
	}
	initCheckEnv(params)
}

func assertCheckGroupRunner(t *testing.T,
	cgRunner checkGroupRunner,
	cpRunners *[]checkPointRunner,
	p params.Params,
	expectedTopic string) {
	t.Helper()
	tests := []struct {
		desc           string
		runnerStub     checkPointRunner
		expectedStatus model.CheckStatus
	}{
		{
			"ok",
			func(p params.Params) []model.CheckPoint {
				return []model.CheckPoint{model.OkCheckPoint("")}
			},
			model.CheckStatusOk,
		},
		{
			"warning",
			func(p params.Params) []model.CheckPoint {
				return []model.CheckPoint{model.WarningCheckPoint("")}
			},
			model.CheckStatusWarning,
		},
		{
			"error",
			func(p params.Params) []model.CheckPoint {
				return []model.CheckPoint{model.ErrorCheckPoint("")}
			},
			model.CheckStatusError,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Helper()
			if cpRunners != nil {
				cpRunnersBackup := *cpRunners
				t.Cleanup(func() { *cpRunners = cpRunnersBackup })
				*cpRunners = []checkPointRunner{test.runnerStub}
			}
			cg := cgRunner(p)
			assert.Equal(t, expectedTopic, cg.GetTopic())
			assert.Equal(t, test.expectedStatus, cg.GetStatus())
		})
	}
}

func Test_checker_run(t *testing.T) {
	okRunner := func(_ params.Params) *model.CheckGroup {
		cg := model.NewCheckGroup("ok runner")
		cg.Add(model.OkCheckPoint("always returns ok"))
		return cg
	}
	warningRunner := func(_ params.Params) *model.CheckGroup {
		cg := model.NewCheckGroup("warning runner")
		cg.Add(model.WarningCheckPoint("always returns warning"))
		return cg
	}
	errorRunner := func(_ params.Params) *model.CheckGroup {
		cg := model.NewCheckGroup("error runner")
		cg.Add(model.ErrorCheckPoint("always returns error"))
		return cg
	}

	tests := []struct {
		desc       string
		runners    []checkGroupRunner
		expectedRC int
	}{
		{"1 ok", []checkGroupRunner{okRunner}, 0},
		{"1 warning", []checkGroupRunner{warningRunner}, 1},
		{"1 error", []checkGroupRunner{errorRunner}, 2},
		{"1 ok 1 warning", []checkGroupRunner{okRunner, warningRunner}, 1},
		{"1 ok 1 error", []checkGroupRunner{okRunner, errorRunner}, 2},
		{"1 warning 1 error", []checkGroupRunner{warningRunner, errorRunner}, 2},
		{"1 ok 1 warning 1 error", []checkGroupRunner{okRunner, warningRunner, errorRunner}, 2},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			checkGroupRunners = test.runners
			Run(params.Params{})
			assert.Equal(t, test.expectedRC, status.GetReturnCode())
		})
	}
}

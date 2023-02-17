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
	"testing"
)

func Test_check_git_returns_warning_with_brand_new_repo(t *testing.T) {
	// Warning is triggered by default if branch (master) being a root branch
	assertWarning(t, checkGitEnvironment, *params.AParamSet())
}

func Test_check_git_auto_push(t *testing.T) {
	tests := []struct {
		desc     string
		value    bool
		expected model.CheckStatus
	}{
		{"enabled", true, model.CheckStatusOk},
		{"disabled", false, model.CheckStatusOk},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assertStatus(t, test.expected,
				func(p params.Params) (cg *model.CheckGroup) {
					cg = model.NewCheckGroup("git auto-push parameter")
					cg.Add(checkGitAutoPush(p)...)
					return cg
				},
				*params.AParamSet(params.WithAutoPush(test.value)))
		})
	}
}

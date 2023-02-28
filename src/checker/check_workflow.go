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
)

var checkWorkflowRunners []checkPointRunner

func init() {
	checkWorkflowRunners = []checkPointRunner{
		checkCommitFailures,
	}
}

func checkWorkflowConfiguration(p params.Params) (cg *model.CheckGroup) {
	cg = model.NewCheckGroup("TCR workflow configuration")
	for _, runner := range checkWorkflowRunners {
		cg.Add(runner(p)...)
	}
	return cg
}

func checkCommitFailures(p params.Params) (cp []model.CheckPoint) {
	switch p.CommitFailures {
	case true:
		cp = append(cp, model.OkCheckPoint(
			"commit-failures is turned on: test-breaking changes will be committed"))
	case false:
		cp = append(cp, model.OkCheckPoint(
			"commit-failures is turned off: test-breaking changes will not be committed"))
	}
	return cp
}

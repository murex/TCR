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
	"github.com/murex/tcr/helpers"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/vcs/p4"
	"strings"
)

var checkP4Runners []checkPointRunner

func init() {
	checkP4Runners = []checkPointRunner{
		checkP4Command,
		checkP4Config,
		checkP4Workspace,
	}
}

func checkP4Environment(p params.Params) (cg *model.CheckGroup) {
	cg = model.NewCheckGroup("perforce environment")
	// p4 environment is checked only when p4 is the selected VCS
	if strings.ToLower(p.VCS) == p4.Name {
		for _, runner := range checkP4Runners {
			cg.Add(runner(p)...)
		}
	}
	return cg
}

func checkP4Command(_ params.Params) (cp []model.CheckPoint) {
	if !p4.IsP4CommandAvailable() {
		cp = append(cp, model.ErrorCheckPoint("p4 command was not found on path"))
		return cp
	}
	cp = append(cp, model.OkCheckPoint("p4 command path is ", p4.GetP4CommandPath()))
	cp = append(cp, model.OkCheckPoint("p4 version is ", p4.GetP4CommandVersion()))
	return cp
}

func checkP4Config(_ params.Params) (cp []model.CheckPoint) {
	if p4.GetP4UserName() == "not set" {
		cp = append(cp, model.WarningCheckPoint("p4 username is not set"))
		return cp
	}
	cp = append(cp, model.OkCheckPoint("p4 username is ", p4.GetP4UserName()))
	return cp
}

func checkP4Workspace(p params.Params) (cp []model.CheckPoint) {
	if p4.GetP4ClientName() == "not set" {
		cp = append(cp, model.ErrorCheckPoint("p4 client name is not set"))
		return cp
	}
	cp = append(cp, model.OkCheckPoint("p4 client name is ", p4.GetP4ClientName()))

	p4RootDir, err := p4.GetP4RootDir()
	if err != nil {
		cp = append(cp, model.ErrorCheckPoint("p4 client root is not set"))
		return cp
	}
	cp = append(cp, model.OkCheckPoint("p4 client root is ", p4RootDir))

	if !helpers.IsSubPathOf(p.BaseDir, p4RootDir) {
		cp = append(cp, model.ErrorCheckPoint("TCR base dir is not under p4 client root dir"))
		return cp
	}

	if !helpers.IsSubPathOf(p.WorkDir, p4RootDir) {
		cp = append(cp, model.ErrorCheckPoint("TCR work dir is not under p4 client root dir"))
		return cp
	}

	return cp
}

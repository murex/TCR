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
)

var checkDirRunners []checkPointRunner

func init() {
	checkDirRunners = []checkPointRunner{
		checkBaseDirectory,
		checkWorkDirectory,
	}
}

func checkDirectories(p params.Params) (cg *model.CheckGroup) {
	cg = model.NewCheckGroup("directories")
	for _, runner := range checkDirRunners {
		cg.Add(runner(p)...)
	}
	return cg
}

func checkBaseDirectory(p params.Params) (cp []model.CheckPoint) {
	if p.BaseDir == "" {
		cp = append(cp, model.OkCheckPoint("base directory parameter is not set explicitly"))
	} else {
		cp = append(cp, model.OkCheckPoint("base directory parameter is ", p.BaseDir))
	}

	if checkEnv.sourceTreeErr != nil {
		cp = append(cp, model.CheckpointsForDirAccessError(p.BaseDir, checkEnv.sourceTreeErr)...)
	} else {
		cp = append(cp, model.OkCheckPoint(
			"base directory absolute path is ", checkEnv.sourceTree.GetBaseDir()))
	}
	return cp
}

func checkWorkDirectory(p params.Params) (cp []model.CheckPoint) {
	if p.WorkDir == "" {
		cp = append(cp, model.OkCheckPoint("work directory parameter is not set explicitly"))
	} else {
		cp = append(cp, model.OkCheckPoint("work directory parameter is ", p.WorkDir))
	}

	if checkEnv.workDirErr != nil {
		cp = append(cp, model.CheckpointsForDirAccessError(p.WorkDir, checkEnv.workDirErr)...)
	} else {
		cp = append(cp, model.OkCheckPoint("work directory absolute path is ", checkEnv.workDir))
	}
	return cp
}

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
	"github.com/murex/tcr/config"
	"github.com/murex/tcr/filesystem"
	"github.com/murex/tcr/language"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/toolchain"
	"github.com/murex/tcr/vcs"
	"github.com/murex/tcr/vcs/factory"
)

type checkGroupRunner func(params.Params) *model.CheckGroup

type checkPointRunner func(p params.Params) []model.CheckPoint

var checkEnv struct {
	configDir     string
	configDirErr  error
	workDir       string
	workDirErr    error
	sourceTree    filesystem.SourceTree
	sourceTreeErr error
	lang          language.LangInterface
	langErr       error
	tchn          toolchain.TchnInterface
	tchnErr       error
	vcs           vcs.Interface
	vcsErr        error
}

// Run goes through all configuration, parameters and local environment to check
// if TCR is ready to be used
func Run(p params.Params) {
	initCheckEnv(p)

	for _, checker := range initCheckers() {
		results := checker(p)
		results.Print()
		model.UpdateReturnState(results)
	}
	report.PostInfo("")
}

func initCheckEnv(p params.Params) {
	model.RecordCheckState(model.CheckStatusOk)

	checkEnv.configDir = config.GetConfigDirPath()
	checkEnv.sourceTree, checkEnv.sourceTreeErr = filesystem.New(p.BaseDir)

	if checkEnv.sourceTreeErr == nil {
		checkEnv.lang, checkEnv.langErr = language.GetLanguage(p.Language, checkEnv.sourceTree.GetBaseDir())
	} else {
		checkEnv.lang, checkEnv.langErr = language.Get(p.Language)
	}

	if checkEnv.langErr == nil {
		checkEnv.tchn, checkEnv.tchnErr = checkEnv.lang.GetToolchain(p.Toolchain)
	} else {
		checkEnv.tchn, checkEnv.tchnErr = toolchain.GetToolchain(p.Toolchain)
	}

	checkEnv.workDirErr = toolchain.SetWorkDir(p.WorkDir)
	checkEnv.workDir = toolchain.GetWorkDir()

	if checkEnv.sourceTreeErr == nil {
		checkEnv.vcs, checkEnv.vcsErr = factory.InitVCS(p.VCS, checkEnv.sourceTree.GetBaseDir())
	}
}

func initCheckers() []checkGroupRunner {
	return []checkGroupRunner{
		checkConfigDirectory,
		checkBaseDirectory,
		checkWorkDirectory,
		checkLanguage,
		checkToolchain,
		checkVCSConfiguration,
		checkGitEnvironment,
		checkP4Environment,
		checkWorkflowConfiguration,
		checkMobConfiguration,
	}
}

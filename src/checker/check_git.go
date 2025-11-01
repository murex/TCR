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
	"strings"

	"github.com/murex/tcr/checker/model"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/vcs/git"
)

var checkGitRunners []checkPointRunner

func init() {
	checkGitRunners = []checkPointRunner{
		checkGitCommand,
		checkGitConfig,
		checkGitRepository,
		checkGitRemote,
		checkGitAutoPush,
	}
}

func checkGitEnvironment(p params.Params) (cg *model.CheckGroup) {
	cg = model.NewCheckGroup("git environment")
	// git environment is checked only when git is the selected VCS
	if strings.ToLower(p.VCS) == git.Name {
		for _, runner := range checkGitRunners {
			cg.Add(runner(p)...)
		}
	}
	return cg
}

func checkGitCommand(_ params.Params) (cp []model.CheckPoint) {
	if !git.IsGitCommandAvailable() {
		cp = append(cp, model.ErrorCheckPoint("git command was not found on path"))
		return cp
	}
	cp = append(cp, model.OkCheckPoint("git command path is ", git.GetGitCommandPath()))
	cp = append(cp, model.OkCheckPoint("git version is ", git.GetGitCommandVersion()))
	// We could add here a check on git minimum version. No specific need for now.
	return cp
}

func checkGitConfig(_ params.Params) (cp []model.CheckPoint) {
	if git.GetGitUserName() == "not set" {
		cp = append(cp, model.WarningCheckPoint("git username is not set"))
		return cp
	}
	cp = append(cp, model.OkCheckPoint("git username is ", git.GetGitUserName()))
	return cp
}

func checkGitRepository(_ params.Params) (cp []model.CheckPoint) {
	if checkEnv.sourceTreeErr != nil {
		cp = append(cp, model.ErrorCheckPoint("cannot retrieve git repository information from base directory name"))
		return cp
	}
	if checkEnv.vcsErr != nil {
		cp = append(cp, model.ErrorCheckPoint(checkEnv.vcsErr))
		return cp
	}
	if checkEnv.vcs == nil {
		cp = append(cp, model.ErrorCheckPoint("git repository not properly initialized"))
		return cp
	}

	cp = append(cp, model.OkCheckPoint("git repository root is ", checkEnv.vcs.GetRootDir()))

	cp = append(cp, model.OkCheckPoint("git working branch is ", checkEnv.vcs.GetWorkingBranch()))
	if checkEnv.vcs.IsOnRootBranch() {
		cp = append(cp, model.WarningCheckPoint("running TCR from a root branch is not recommended"))
	}
	return cp
}

func checkGitRemote(_ params.Params) (cp []model.CheckPoint) {
	if checkEnv.vcs == nil || checkEnv.vcsErr != nil {
		// If git is not properly initialized, no point in trying to go further
		return []model.CheckPoint{}
	}

	if !checkEnv.vcs.IsRemoteEnabled() {
		if checkEnv.vcs.GetRemoteName() != "" {
			cp = append(cp, model.WarningCheckPoint("git remote not found: ", checkEnv.vcs.GetRemoteName()))
		}
		cp = append(cp, model.OkCheckPoint("git remote is disabled: all operations will be done locally"))
		return cp
	}

	cp = append(cp, model.OkCheckPoint("git remote name is ", checkEnv.vcs.GetRemoteName()))

	if checkEnv.vcs.CheckRemoteAccess() {
		cp = append(cp, model.OkCheckPoint("git remote access seems to be working"))
	} else {
		cp = append(cp, model.ErrorCheckPoint("git remote access does not seem to be working"))
	}
	return cp
}

func checkGitAutoPush(p params.Params) (cp []model.CheckPoint) {
	if p.AutoPush {
		cp = append(cp, model.OkCheckPoint("git auto-push is turned on: every commit will be pushed to origin"))
	} else {
		cp = append(cp, model.OkCheckPoint("git auto-push is turned off: commits will only be applied locally"))
	}
	return cp
}

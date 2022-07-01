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
	"github.com/murex/tcr/tcr-engine/vcs"
)

func checkGitEnvironment(_ params.Params) (cr *CheckResults) {
	cr = NewCheckResults("git environment")
	cr.add(checkGitCommand())
	cr.add(checkGitConfig())
	cr.add(checkGitRepository())
	cr.add(checkGitRemote())
	return
}

func checkGitCommand() (cp []CheckPoint) {
	if !vcs.IsGitCommandAvailable() {
		cp = append(cp, errorCheckPoint("git command was not found on path"))
		return
	}
	cp = append(cp, okCheckPoint("git command path is ", vcs.GetGitCommandPath()))
	cp = append(cp, okCheckPoint("git version is ", vcs.GetGitCommandVersion()))
	// TODO check git minimum version?
	return
}

func checkGitConfig() (cp []CheckPoint) {
	cp = append(cp, okCheckPoint("git username is ", vcs.GetGitUserName()))
	return
}

func checkGitRepository() (cp []CheckPoint) {
	if checkEnv.sourceTreeErr != nil {
		cp = append(cp, errorCheckPoint("cannot retrieve git repository information from base directory name"))
		return
	}
	if checkEnv.gitErr != nil {
		cp = append(cp, errorCheckPoint(checkEnv.gitErr))
		return
	}

	cp = append(cp, okCheckPoint("git repository root is ", checkEnv.git.GetRootDir()))

	branch := checkEnv.git.GetWorkingBranch()
	cp = append(cp, okCheckPoint("git working branch is ", branch))
	if vcs.IsRootBranch(branch) {
		cp = append(cp, warningCheckPoint("running TCR from a root branch is not recommended"))
	}
	return
}

func checkGitRemote() (cp []CheckPoint) {
	if checkEnv.git == nil {
		// If git is not properly initialized, no point in trying to go further
		return
	}

	if !checkEnv.git.IsRemoteEnabled() {
		cp = append(cp, okCheckPoint("git remote is disabled: all operations will be done locally"))
		return
	}

	cp = append(cp, okCheckPoint("git remote name is ", checkEnv.git.GetRemoteName()))

	if checkEnv.git.CheckRemoteAccess() {
		cp = append(cp, okCheckPoint("git remote access seems to be working"))
	} else {
		cp = append(cp, errorCheckPoint("git remote access does not seem to be working"))
	}
	return
}

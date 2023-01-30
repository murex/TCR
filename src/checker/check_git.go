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
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/vcs/git"
)

func checkGitEnvironment(p params.Params) (cr *CheckResults) {
	cr = NewCheckResults("git environment")
	cr.add(checkGitCommand())
	cr.add(checkGitConfig())
	cr.add(checkGitRepository())
	cr.add(checkGitRemote())
	cr.add(checkGitAutoPush(p))
	return cr
}

func checkGitCommand() (cp []CheckPoint) {
	if !git.IsGitCommandAvailable() {
		cp = append(cp, errorCheckPoint("git command was not found on path"))
		return cp
	}
	cp = append(cp, okCheckPoint("git command path is ", git.GetGitCommandPath()))
	cp = append(cp, okCheckPoint("git version is ", git.GetGitCommandVersion()))
	// We could add here a check on git minimum version. No specific need for now.
	return cp
}

func checkGitConfig() (cp []CheckPoint) {
	cp = append(cp, okCheckPoint("git username is ", git.GetGitUserName()))
	return cp
}

func checkGitRepository() (cp []CheckPoint) {
	if checkEnv.sourceTreeErr != nil {
		cp = append(cp, errorCheckPoint("cannot retrieve git repository information from base directory name"))
		return cp
	}
	if checkEnv.vcsErr != nil {
		cp = append(cp, errorCheckPoint(checkEnv.vcsErr))
		return cp
	}

	cp = append(cp, okCheckPoint("git repository root is ", checkEnv.vcs.GetRootDir()))

	cp = append(cp, okCheckPoint("git working branch is ", checkEnv.vcs.GetWorkingBranch()))
	if checkEnv.vcs.IsOnRootBranch() {
		cp = append(cp, warningCheckPoint("running TCR from a root branch is not recommended"))
	}
	return cp
}

func checkGitRemote() (cp []CheckPoint) {
	if checkEnv.vcs == nil {
		// If git is not properly initialized, no point in trying to go further
		return cp
	}

	if !checkEnv.vcs.IsRemoteEnabled() {
		cp = append(cp, okCheckPoint("git remote is disabled: all operations will be done locally"))
		return cp
	}

	cp = append(cp, okCheckPoint("git remote name is ", checkEnv.vcs.GetRemoteName()))

	if checkEnv.vcs.CheckRemoteAccess() {
		cp = append(cp, okCheckPoint("git remote access seems to be working"))
	} else {
		cp = append(cp, errorCheckPoint("git remote access does not seem to be working"))
	}
	return cp
}

func checkGitAutoPush(p params.Params) (cp []CheckPoint) {
	if p.AutoPush {
		cp = append(cp, okCheckPoint("git auto-push is turned on: every commit will be pushed to origin"))
	} else {
		cp = append(cp, okCheckPoint("git auto-push is turned off: commits will only be applied locally"))
	}
	return cp
}

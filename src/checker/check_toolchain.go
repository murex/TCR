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
	"github.com/murex/tcr/toolchain"
	"runtime"
)

func checkToolchain(p params.Params) (cg *model.CheckGroup) {
	cg = model.NewCheckGroup("toolchain")

	if p.Toolchain == "" {
		cg.Add(checkpointsWhenToolchainIsNotSet()...)
	} else {
		cg.Add(checkpointsWhenToolchainIsSet(p.Toolchain)...)
	}

	if checkEnv.tchn != nil {
		cg.Ok("local platform is "+toolchain.OsName(runtime.GOOS), "/", toolchain.ArchName(runtime.GOARCH))
		cg.Add(checkCommandLine("build", checkEnv.tchn.BuildCommandPath(), checkEnv.tchn.BuildCommandLine())...)
		cg.Add(checkCommandLine("test", checkEnv.tchn.TestCommandPath(), checkEnv.tchn.TestCommandLine())...)
		cg.Add(checkTestResultDir(checkEnv.tchn.GetTestResultDir(), checkEnv.tchn.GetTestResultPath())...)
	}
	return cg
}

func checkCommandLine(name string, cmdPath string, cmdLine string) (cp []model.CheckPoint) {
	cp = append(cp, model.OkCheckPoint(name, " command line: ", cmdLine))

	path, err := checkEnv.tchn.CheckCommandAccess(cmdPath)
	if err != nil {
		cp = append(cp, model.ErrorCheckPoint("cannot access ", name, " command: ", cmdPath))
	} else {
		cp = append(cp, model.OkCheckPoint(name, " command path: ", path))
	}
	return cp
}

func checkTestResultDir(dir string, path string) (cp []model.CheckPoint) {
	if dir == "" {
		cp = append(cp, model.WarningCheckPoint(
			"test result directory parameter is not set explicitly (default: work directory)"))
	} else {
		cp = append(cp, model.OkCheckPoint("test result directory parameter is ", dir))
	}

	cp = append(cp, model.OkCheckPoint("test result directory absolute path is ", path))
	return cp
}

func checkpointsWhenToolchainIsSet(name string) (cp []model.CheckPoint) {
	cp = append(cp, model.OkCheckPoint("toolchain parameter is set to ", name))
	if checkEnv.tchnErr != nil {
		cp = append(cp, model.ErrorCheckPoint(checkEnv.tchnErr))
	} else {
		cp = append(cp, model.OkCheckPoint(checkEnv.tchn.GetName(), " toolchain is valid"))
		if checkEnv.langErr == nil {
			cp = append(cp, model.OkCheckPoint(checkEnv.tchn.GetName(), " toolchain is compatible with ",
				checkEnv.lang.GetName(), " language"))
		} else {
			cp = append(cp, model.WarningCheckPoint("skipping toolchain/language compatibility check"))
		}
	}
	return cp
}

func checkpointsWhenToolchainIsNotSet() (cp []model.CheckPoint) {
	cp = append(cp, model.OkCheckPoint("toolchain parameter is not set explicitly"))

	if checkEnv.langErr != nil {
		cp = append(cp, model.WarningCheckPoint("language is unknown"))
		cp = append(cp, model.ErrorCheckPoint("cannot retrieve toolchain from an unknown language"))
	} else {
		cp = append(cp, model.OkCheckPoint("using language's default toolchain"))
		if checkEnv.tchnErr != nil {
			cp = append(cp, model.ErrorCheckPoint(checkEnv.tchnErr))
		} else {
			cp = append(cp, model.OkCheckPoint("default toolchain for ", checkEnv.lang.GetName(),
				" language is ", checkEnv.tchn.GetName()))
		}
	}
	return cp
}

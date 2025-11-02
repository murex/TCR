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
	"runtime"

	"github.com/murex/tcr/checker/model"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/toolchain/command"
)

var checkToolchainRunners []checkPointRunner

func init() {
	checkToolchainRunners = []checkPointRunner{
		checkToolchainParameter,
		checkToolchainLanguageCompatibility,
		checkToolchainDetection,
		checkToolchainPlatform,
		checkToolchainBuildCommand,
		checkToolchainTestCommand,
		checkToolchainTestResultDir,
	}
}

func checkToolchain(p params.Params) (cg *model.CheckGroup) {
	cg = model.NewCheckGroup("toolchain")
	for _, runner := range checkToolchainRunners {
		cg.Add(runner(p)...)
	}
	return cg
}

func checkToolchainParameter(p params.Params) (cp []model.CheckPoint) {
	if p.Toolchain == "" {
		return cp
	}
	cp = append(cp, model.OkCheckPoint("toolchain parameter is set to ", p.Toolchain))

	if checkEnv.tchnErr != nil {
		cp = append(cp, model.ErrorCheckPoint(checkEnv.tchnErr))
		return cp
	}

	cp = append(cp, model.OkCheckPoint(p.Toolchain, " toolchain is valid"))
	return cp
}

func checkToolchainLanguageCompatibility(p params.Params) (cp []model.CheckPoint) {
	if p.Toolchain == "" {
		return cp
	}

	if checkEnv.langErr != nil {
		cp = append(cp, model.WarningCheckPoint("skipping toolchain/language compatibility check"))
		return cp
	}

	cp = append(cp, model.OkCheckPoint(p.Toolchain, " toolchain is compatible with ",
		checkEnv.lang.GetName(), " language"))
	return cp
}

func checkToolchainDetection(p params.Params) (cp []model.CheckPoint) {
	if p.Toolchain != "" {
		return cp
	}

	cp = append(cp, model.OkCheckPoint("toolchain parameter is not set explicitly"))

	if checkEnv.langErr != nil {
		cp = append(cp, model.WarningCheckPoint("language is unknown"))
		cp = append(cp, model.ErrorCheckPoint("cannot retrieve toolchain from an unknown language"))
		return cp
	}

	cp = append(cp, model.OkCheckPoint("using language's default toolchain"))

	if checkEnv.tchnErr != nil {
		cp = append(cp, model.ErrorCheckPoint(checkEnv.tchnErr))
		return cp
	}

	cp = append(cp, model.OkCheckPoint("default toolchain for ",
		checkEnv.lang.GetName(), " language is ",
		checkEnv.lang.GetToolchains().Default))
	return cp
}

func checkToolchainPlatform(_ params.Params) (cp []model.CheckPoint) {
	if checkEnv.tchn == nil {
		return cp
	}
	cp = append(cp, model.OkCheckPoint("local platform is ",
		command.OsName(runtime.GOOS), "/", command.ArchName(runtime.GOARCH)))
	return cp
}

func checkToolchainBuildCommand(_ params.Params) (cp []model.CheckPoint) {
	if checkEnv.tchn == nil {
		return cp
	}
	return checkCommandLine("build", checkEnv.tchn.BuildCommandPath(), checkEnv.tchn.BuildCommandLine())
}

func checkToolchainTestCommand(_ params.Params) (cp []model.CheckPoint) {
	if checkEnv.tchn == nil {
		return cp
	}
	return checkCommandLine("test", checkEnv.tchn.TestCommandPath(), checkEnv.tchn.TestCommandLine())
}

func checkCommandLine(name string, cmdPath string, cmdLine string) (cp []model.CheckPoint) {
	cp = append(cp, model.OkCheckPoint(name, " command line: ", cmdLine))

	path, err := checkEnv.tchn.CheckCommandAccess(cmdPath)
	if err != nil {
		cp = append(cp, model.ErrorCheckPoint("cannot access ", name, " command: ", cmdPath))
		return cp
	}

	cp = append(cp, model.OkCheckPoint(name, " command path: ", path))
	return cp
}

func checkToolchainTestResultDir(_ params.Params) (cp []model.CheckPoint) {
	if checkEnv.tchn == nil {
		return cp
	}

	dir := checkEnv.tchn.GetTestResultDir()
	if dir == "" {
		cp = append(cp, model.WarningCheckPoint(
			"test result directory parameter is not set explicitly (default: work directory)"))
	} else {
		cp = append(cp, model.OkCheckPoint("test result directory parameter is ", dir))
	}

	cp = append(cp, model.OkCheckPoint(
		"test result directory absolute path is ", checkEnv.tchn.GetTestResultPath()))
	return cp
}

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
	"github.com/murex/tcr/toolchain"
	"runtime"
)

func checkToolchain(p params.Params) (cr *CheckResults) {
	cr = NewCheckResults("toolchain")

	if p.Toolchain == "" {
		cr.add(checkpointsWhenToolchainIsNotSet())
	} else {
		cr.add(checkpointsWhenToolchainIsSet(p.Toolchain))
	}

	if checkEnv.tchn != nil {
		cr.ok("local platform is "+toolchain.OsName(runtime.GOOS), "/", toolchain.ArchName(runtime.GOARCH))
		cr.add(checkCommandLine("build", checkEnv.tchn.BuildCommandPath(), checkEnv.tchn.BuildCommandLine()))
		cr.add(checkCommandLine("test", checkEnv.tchn.TestCommandPath(), checkEnv.tchn.TestCommandLine()))
		cr.add(checkTestResultDir(checkEnv.tchn.GetTestResultDir(), checkEnv.tchn.GetTestResultPath()))
	}
	return cr
}

func checkCommandLine(name string, cmdPath string, cmdLine string) (cp []CheckPoint) {
	cp = append(cp, okCheckPoint(name, " command line: ", cmdLine))

	path, err := checkEnv.tchn.CheckCommandAccess(cmdPath)
	if err != nil {
		cp = append(cp, errorCheckPoint("cannot access ", name, " command: ", cmdPath))
	} else {
		cp = append(cp, okCheckPoint(name, " command path: ", path))
	}
	return cp
}

func checkTestResultDir(dir string, path string) (cp []CheckPoint) {
	if dir == "" {
		cp = append(cp, warningCheckPoint(
			"test result directory parameter is not set explicitly (default: work directory)"))
	} else {
		cp = append(cp, okCheckPoint("test result directory parameter is ", dir))
	}

	cp = append(cp, okCheckPoint("test result directory absolute path is ", path))
	return cp
}

func checkpointsWhenToolchainIsSet(name string) (cp []CheckPoint) {
	cp = append(cp, okCheckPoint("toolchain parameter is set to ", name))
	if checkEnv.tchnErr != nil {
		cp = append(cp, errorCheckPoint(checkEnv.tchnErr))
	} else {
		cp = append(cp, okCheckPoint(checkEnv.tchn.GetName(), " toolchain is valid"))
		if checkEnv.langErr == nil {
			cp = append(cp, okCheckPoint(checkEnv.tchn.GetName(), " toolchain is compatible with ",
				checkEnv.lang.GetName(), " language"))
		} else {
			cp = append(cp, warningCheckPoint("skipping toolchain/language compatibility check"))
		}
	}
	return cp
}

func checkpointsWhenToolchainIsNotSet() (cp []CheckPoint) {
	cp = append(cp, okCheckPoint("toolchain parameter is not set explicitly"))

	if checkEnv.langErr != nil {
		cp = append(cp, warningCheckPoint("language is unknown"))
		cp = append(cp, errorCheckPoint("cannot retrieve toolchain from an unknown language"))
	} else {
		cp = append(cp, okCheckPoint("using language's default toolchain"))
		if checkEnv.tchnErr != nil {
			cp = append(cp, errorCheckPoint(checkEnv.tchnErr))
		} else {
			cp = append(cp, okCheckPoint("default toolchain for ", checkEnv.lang.GetName(),
				" language is ", checkEnv.tchn.GetName()))
		}
	}
	return cp
}

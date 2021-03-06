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
	"fmt"
	"github.com/murex/tcr/tcr-engine/config"
	"github.com/murex/tcr/tcr-engine/filesystem"
	"github.com/murex/tcr/tcr-engine/language"
	"github.com/murex/tcr/tcr-engine/params"
	"github.com/murex/tcr/tcr-engine/report"
	"github.com/murex/tcr/tcr-engine/status"
	"github.com/murex/tcr/tcr-engine/toolchain"
	"github.com/murex/tcr/tcr-engine/vcs"
	"os"
)

// CheckStatus provides the return status for a specific check
type CheckStatus int

// Check status values
const (
	CheckStatusOk      CheckStatus = 0 // Check status is OK
	CheckStatusWarning CheckStatus = 1 // Check status is Warning
	CheckStatusError   CheckStatus = 2 // Build status is Error
)

// CheckPoint is used for describing the result of a single check point
type CheckPoint struct {
	rc          CheckStatus
	description string
}

// CheckResults contains an aggregation of CheckPoint for a specific topic
type CheckResults struct {
	topic       string
	checkPoints []CheckPoint
}

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
	git           vcs.GitInterface
	gitErr        error
}

// Run goes through all configuration, parameters and local environment to check
// if TCR is ready to be used
func Run(p params.Params) {
	initCheckEnv(p)

	checkers := []func(params.Params) *CheckResults{
		checkConfigDirectory,
		checkBaseDirectory,
		checkWorkDirectory,
		checkLanguage,
		checkToolchain,
		checkGitEnvironment,
		checkAutoPush,
		checkCommitFailures,
		checkMobTimer,
		checkPollingPeriod,
	}

	for _, checker := range checkers {
		results := checker(p)
		results.print()
		updateReturnState(results)
	}
	report.PostInfo("")
}

func updateReturnState(results *CheckResults) {
	if int(results.getStatus()) > status.GetReturnCode() {
		recordCheckState(results.getStatus())
	}
}

func recordCheckState(s CheckStatus) {
	status.RecordState(status.NewStatus(int(s)))
}

func initCheckEnv(params params.Params) {
	recordCheckState(CheckStatusOk)

	checkEnv.configDir = config.GetConfigDirPath()
	checkEnv.sourceTree, checkEnv.sourceTreeErr = filesystem.New(params.BaseDir)

	if checkEnv.sourceTreeErr == nil {
		checkEnv.lang, checkEnv.langErr = language.GetLanguage(params.Language, checkEnv.sourceTree.GetBaseDir())
	} else {
		checkEnv.lang, checkEnv.langErr = language.Get(params.Language)
	}

	if checkEnv.langErr == nil {
		checkEnv.tchn, checkEnv.tchnErr = checkEnv.lang.GetToolchain(params.Toolchain)
	} else {
		checkEnv.tchn, checkEnv.tchnErr = toolchain.GetToolchain(params.Toolchain)
	}

	checkEnv.workDirErr = toolchain.SetWorkDir(params.WorkDir)
	checkEnv.workDir = toolchain.GetWorkDir()

	if checkEnv.sourceTreeErr == nil {
		checkEnv.git, checkEnv.gitErr = vcs.New(checkEnv.sourceTree.GetBaseDir())
	}
}

func checkpointsForDirAccessError(dir string, err error) []CheckPoint {
	var checkpoint CheckPoint
	if os.IsNotExist(err) {
		checkpoint = errorCheckPoint("directory not found: ", dir)
	} else if os.IsPermission(err) {
		checkpoint = errorCheckPoint("cannot access directory ", dir)
	} else {
		checkpoint = errorCheckPoint(err)
	}
	return []CheckPoint{checkpoint}
}

func checkpointsForList(headerMsg string, emptyMsg string, values []string) (cp []CheckPoint) {
	if len(values) == 0 {
		cp = append(cp, warningCheckPoint(emptyMsg))
		return
	}
	cp = append(cp, okCheckPoint(headerMsg))
	for _, value := range values {
		cp = append(cp, okCheckPoint("- ", value))
	}
	return
}

func okCheckPoint(a ...interface{}) CheckPoint {
	return CheckPoint{rc: CheckStatusOk, description: fmt.Sprint(a...)}
}

func warningCheckPoint(a ...interface{}) CheckPoint {
	return CheckPoint{rc: CheckStatusWarning, description: fmt.Sprint(a...)}
}

func errorCheckPoint(a ...interface{}) CheckPoint {
	return CheckPoint{rc: CheckStatusError, description: fmt.Sprint(a...)}
}

// NewCheckResults creates a new CheckResults instance
func NewCheckResults(topic string) *CheckResults {
	return &CheckResults{
		topic:       topic,
		checkPoints: []CheckPoint{},
	}
}

func (cr *CheckResults) addCheckPoint(checkpoint CheckPoint) {
	cr.checkPoints = append(cr.checkPoints, checkpoint)
}

func (cr *CheckResults) add(checkPoints []CheckPoint) {
	cr.checkPoints = append(cr.checkPoints, checkPoints...)
}

func (cr *CheckResults) ok(a ...interface{}) {
	cr.addCheckPoint(okCheckPoint(a...))
}

func (cr *CheckResults) warning(a ...interface{}) {
	cr.addCheckPoint(warningCheckPoint(a...))
}

//func (cr *CheckResults) error(a ...interface{}) {
//	cr.addCheckPoint(errorCheckPoint(a...))
//}

func (cr *CheckResults) getStatus() (status CheckStatus) {
	status = CheckStatusOk
	for _, check := range cr.checkPoints {
		if check.rc > status {
			status = check.rc
		}
	}
	return
}

func (cr *CheckResults) print() {
	report.PostInfo()
	const messagePrefix = "??? checking "
	switch cr.getStatus() {
	case CheckStatusOk:
		report.PostInfo(messagePrefix, cr.topic)
	case CheckStatusWarning:
		report.PostWarning(messagePrefix, cr.topic)
	case CheckStatusError:
		report.PostError(messagePrefix, cr.topic)
	}
	for _, check := range cr.checkPoints {
		check.print()
	}
}

func (cp CheckPoint) print() {
	switch cp.rc {
	case CheckStatusOk:
		report.PostInfo("\t??? ", cp.description)
	case CheckStatusWarning:
		report.PostWarning("\t??? ", cp.description)
	case CheckStatusError:
		report.PostError("\t??? ", cp.description)
		//fmt.Println("EEE - ", cp.description)
	}
}

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
	"github.com/murex/tcr/params"
)

var checkLanguageRunners []checkPointRunner

func init() {
	checkLanguageRunners = []checkPointRunner{
		checkLanguageParameter,
		checkLanguageDetection,
		checkLanguageSrcDirectories,
		checkLanguageSrcPatterns,
		checkLanguageSrcFiles,
		checkLanguageTestDirectories,
		checkLanguageTestPatterns,
		checkLanguageTestFiles,
	}
}

func checkLanguage(p params.Params) (cg *model.CheckGroup) {
	cg = model.NewCheckGroup("language")
	for _, runner := range checkLanguageRunners {
		cg.Add(runner(p)...)
	}
	return cg
}

func checkLanguageParameter(p params.Params) (cp []model.CheckPoint) {
	if p.Language == "" {
		return cp
	}
	cp = append(cp, model.OkCheckPoint("language parameter is set to ", p.Language))

	if checkEnv.langErr != nil {
		cp = append(cp, model.ErrorCheckPoint(checkEnv.langErr))
		return cp
	}

	cp = append(cp, model.OkCheckPoint(p.Language+" language is valid"))
	return cp
}

func checkLanguageDetection(p params.Params) (cp []model.CheckPoint) {
	if p.Language != "" {
		return cp
	}
	cp = append(cp, model.OkCheckPoint("language parameter is not set explicitly"))

	if checkEnv.sourceTree == nil || !checkEnv.sourceTree.IsValid() {
		cp = append(cp, model.ErrorCheckPoint("cannot retrieve language from base directory name"))
		return cp
	}

	cp = append(cp, model.OkCheckPoint("using base directory name for language detection"))
	cp = append(cp, model.OkCheckPoint("base directory is ", checkEnv.sourceTree.GetBaseDir()))

	if checkEnv.langErr != nil {
		cp = append(cp, model.ErrorCheckPoint(checkEnv.langErr))
		return cp
	}

	cp = append(cp, model.OkCheckPoint("language retrieved from base directory name: ", checkEnv.lang.GetName()))
	return cp
}

func checkLanguageSrcDirectories(_ params.Params) (cp []model.CheckPoint) {
	if checkEnv.lang == nil {
		return cp
	}
	return model.CheckpointsForList(
		"source directories:",
		"no source directory defined for "+languageAsText(),
		checkEnv.lang.GetSrcFileFilter().Directories...)
}

func checkLanguageSrcPatterns(_ params.Params) (cp []model.CheckPoint) {
	if checkEnv.lang == nil {
		return cp
	}
	return model.CheckpointsForList(
		"source filename matching patterns:",
		"no source filename pattern defined for "+languageAsText(),
		checkEnv.lang.GetSrcFileFilter().FilePatterns...)
}

func checkLanguageSrcFiles(_ params.Params) (cp []model.CheckPoint) {
	if checkEnv.lang == nil {
		return cp
	}
	srcFiles, err := checkEnv.lang.AllSrcFiles()
	if err != nil {
		cp = append(cp, model.ErrorCheckPoint(err))
		return cp
	}
	return model.CheckpointsForList(
		"matching source files found:",
		"no matching source file found",
		srcFiles...)
}

func checkLanguageTestDirectories(_ params.Params) (cp []model.CheckPoint) {
	if checkEnv.lang == nil {
		return cp
	}
	return model.CheckpointsForList(
		"test directories:",
		"no test directory defined for "+languageAsText(),
		checkEnv.lang.GetTestFileFilter().Directories...)
}

func checkLanguageTestPatterns(_ params.Params) (cp []model.CheckPoint) {
	if checkEnv.lang == nil {
		return cp
	}
	return model.CheckpointsForList(
		"test filename matching patterns:",
		"no test filename pattern defined for "+languageAsText(),
		checkEnv.lang.GetTestFileFilter().FilePatterns...)
}

func checkLanguageTestFiles(_ params.Params) (cp []model.CheckPoint) {
	if checkEnv.lang == nil {
		return cp
	}
	testFiles, err := checkEnv.lang.AllTestFiles()
	if err != nil {
		cp = append(cp, model.ErrorCheckPoint(err))
		return cp
	}
	return model.CheckpointsForList(
		"matching test files found:",
		"no matching test file found",
		testFiles...)
}

func languageAsText() string {
	return checkEnv.lang.GetName() + " language"
}

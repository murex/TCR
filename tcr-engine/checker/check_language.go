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
)

func checkLanguage(params params.Params) (cr *CheckResults) {
	cr = NewCheckResults("language")

	if params.Language == "" {
		cr.add(checkpointsWhenLanguageIsNotSet())
	} else {
		cr.add(checkpointsWhenLanguageIsSet(params.Language))
	}

	if checkEnv.lang != nil {
		cr.add(checkSrcDirectories())
		cr.add(checkSrcPatterns())
		cr.add(checkSrcFiles())

		cr.add(checkTestDirectories())
		cr.add(checkTestPatterns())
		cr.add(checkTestFiles())
	}
	return
}

func languageInText() string {
	return checkEnv.lang.GetName() + " language"
}

func checkpointsWhenLanguageIsSet(name string) (cp []CheckPoint) {
	cp = append(cp, okCheckPoint("language parameter is set to ", name))
	if checkEnv.langErr != nil {
		cp = append(cp, errorCheckPoint(checkEnv.langErr))
	} else {
		cp = append(cp, okCheckPoint(languageInText()+" is valid"))
	}
	return
}

func checkpointsWhenLanguageIsNotSet() (cp []CheckPoint) {
	cp = append(cp, okCheckPoint("language parameter is not set explicitly"))
	if checkEnv.sourceTree != nil && checkEnv.sourceTree.IsValid() {
		cp = append(cp, okCheckPoint("using base directory name for language detection"))
		cp = append(cp, okCheckPoint("base directory is ", checkEnv.sourceTree.GetBaseDir()))

		if checkEnv.langErr != nil {
			cp = append(cp, errorCheckPoint(checkEnv.langErr))
		} else {
			cp = append(cp, okCheckPoint("language retrieved from base directory name: ", checkEnv.lang.GetName()))
		}
	} else {
		cp = append(cp, errorCheckPoint("cannot retrieve language from base directory name"))
	}
	return
}

func checkSrcDirectories() []CheckPoint {
	return checkpointsForList(
		"source directories:",
		"no source directory defined for "+languageInText(),
		checkEnv.lang.GetSrcFileFilter().Directories,
	)
}

func checkSrcPatterns() []CheckPoint {
	return checkpointsForList(
		"source filename matching patterns:",
		"no source filename pattern defined for "+languageInText(),
		checkEnv.lang.GetSrcFileFilter().FilePatterns,
	)
}

func checkSrcFiles() (cp []CheckPoint) {
	srcFiles, err := checkEnv.lang.AllSrcFiles()
	if err != nil {
		cp = append(cp, errorCheckPoint(err))
	}
	cp = append(cp, checkpointsForList(
		"matching source files found:",
		"no matching source file found",
		srcFiles,
	)...)
	return
}

func checkTestDirectories() []CheckPoint {
	return checkpointsForList(
		"test directories:",
		"no test directory defined for "+languageInText(),
		checkEnv.lang.GetTestFileFilter().Directories,
	)
}

func checkTestPatterns() []CheckPoint {
	return checkpointsForList(
		"test filename matching patterns:",
		"no test filename pattern defined for "+languageInText(),
		checkEnv.lang.GetTestFileFilter().FilePatterns,
	)
}

func checkTestFiles() (cp []CheckPoint) {
	testFiles, err := checkEnv.lang.AllTestFiles()
	if err != nil {
		cp = append(cp, errorCheckPoint(err))
	}
	cp = append(cp, checkpointsForList(
		"matching test files found:",
		"no matching test file found",
		testFiles,
	)...)
	return
}

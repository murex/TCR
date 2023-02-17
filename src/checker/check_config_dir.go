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
	"github.com/murex/tcr/config"
	"github.com/murex/tcr/params"
	"path/filepath"
)

func checkConfigDirectory(p params.Params) (cr *CheckResults) {
	cr = NewCheckResults("configuration directory")

	// TODO see how we can handle incorrect configuration settings, knowing that conf is already loaded when we get here

	if p.ConfigDir == "" {
		cr.ok("configuration directory parameter is not set explicitly")
		cr.ok("using current directory as configuration directory")
	} else {
		cr.ok("configuration directory parameter is: ", p.ConfigDir)
	}
	tcrDirPath, _ := filepath.Abs(config.GetConfigDirPath())
	cr.ok("TCR configuration root directory is ", tcrDirPath)

	cr.add(checkLanguageConfig()...)
	cr.add(checkToolchainConfig()...)

	return cr
}

func checkLanguageConfig() []CheckPoint {
	return checkpointsForConfigSubDir("language",
		config.GetLanguageConfigDirPath(),
		config.GetLanguageConfigFileList())
}

func checkToolchainConfig() []CheckPoint {
	return checkpointsForConfigSubDir("toolchain",
		config.GetToolchainConfigDirPath(),
		config.GetToolchainConfigFileList())
}

func checkpointsForConfigSubDir(name string, path string, files []string) (cp []CheckPoint) {
	dirPath, _ := filepath.Abs(path)
	cp = append(cp, okCheckPoint(name+" configuration directory is ", dirPath))
	cp = append(cp, checkpointsForList(
		name+" configuration files:", "no "+name+" configuration file found", files...)...)
	return cp
}

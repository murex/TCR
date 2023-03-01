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
	"github.com/murex/tcr/config"
	"github.com/murex/tcr/params"
	"path/filepath"
)

var checkConfigRunners []checkPointRunner

type configSubDir struct {
	name        string
	getDirPath  func() string
	getFileList func() []string
}

var (
	languageConfigDir  configSubDir
	toolchainConfigDir configSubDir
)

func init() {
	checkConfigRunners = []checkPointRunner{
		checkConfigDirectory,
		checkLanguageConfig,
		checkToolchainConfig,
	}

	languageConfigDir = configSubDir{
		name:        "language",
		getDirPath:  config.GetLanguageConfigDirPath,
		getFileList: config.GetLanguageConfigFileList,
	}
	toolchainConfigDir = configSubDir{
		name:        "toolchain",
		getDirPath:  config.GetToolchainConfigDirPath,
		getFileList: config.GetToolchainConfigFileList,
	}
}

func checkConfigFiles(p params.Params) (cg *model.CheckGroup) {
	cg = model.NewCheckGroup("configuration files")
	for _, runner := range checkConfigRunners {
		cg.Add(runner(p)...)
	}
	return cg
}

func checkConfigDirectory(p params.Params) (cp []model.CheckPoint) {
	if p.ConfigDir == "" {
		cp = append(cp, model.OkCheckPoint("configuration directory parameter is not set explicitly"))
		cp = append(cp, model.OkCheckPoint("using current directory as configuration directory"))
	} else {
		cp = append(cp, model.OkCheckPoint("configuration directory parameter is ", p.ConfigDir))
	}
	tcrDirPath, _ := filepath.Abs(config.GetConfigDirPath())
	cp = append(cp, model.OkCheckPoint("TCR configuration root directory is ", tcrDirPath))
	return cp
}

func checkLanguageConfig(_ params.Params) (cp []model.CheckPoint) {
	return checkSubDirConfig(languageConfigDir)
}

func checkToolchainConfig(_ params.Params) (cp []model.CheckPoint) {
	return checkSubDirConfig(toolchainConfigDir)
}

func checkSubDirConfig(configDir configSubDir) (cp []model.CheckPoint) {
	dirPath, _ := filepath.Abs(configDir.getDirPath())
	cp = append(cp, model.OkCheckPoint(configDir.name+" configuration directory is ", dirPath))
	cp = append(cp, model.CheckpointsForList(configDir.name+" configuration files:",
		"no "+configDir.name+" configuration file found", configDir.getFileList()...)...)
	return cp
}

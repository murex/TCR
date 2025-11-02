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
	"errors"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/murex/tcr/checker/model"
	"github.com/murex/tcr/params"
	"github.com/stretchr/testify/assert"
)

func Test_check_directories(t *testing.T) {
	assertCheckGroupRunner(t,
		checkDirectories,
		&checkDirRunners,
		*params.AParamSet(),
		"directories")
}

func Test_check_base_directory(t *testing.T) {
	currentDir, _ := filepath.Abs(".")
	tests := []struct {
		desc          string
		value         string
		sourceTreeErr error
		expected      []model.CheckPoint
	}{
		{
			"not set", "", nil,
			[]model.CheckPoint{
				model.OkCheckPoint("base directory parameter is not set explicitly"),
				model.OkCheckPoint("base directory absolute path is ", currentDir),
			},
		},
		{
			"set and exists", ".", nil,
			[]model.CheckPoint{
				model.OkCheckPoint("base directory parameter is ."),
				model.OkCheckPoint("base directory absolute path is ", currentDir),
			},
		},
		{
			"set but does not exist", "missing-dir",
			fs.ErrNotExist,
			[]model.CheckPoint{
				model.OkCheckPoint("base directory parameter is missing-dir"),
				model.ErrorCheckPoint("directory not found: missing-dir"),
			},
		},
		{
			"set but insufficient permissions", "no-perm-dir",
			fs.ErrPermission,
			[]model.CheckPoint{
				model.OkCheckPoint("base directory parameter is no-perm-dir"),
				model.ErrorCheckPoint("cannot access directory no-perm-dir"),
			},
		},
		{
			"set but not a directory", "some-file",
			errors.New(filepath.Join(currentDir, "some-file") + " exists but is not a directory"),
			[]model.CheckPoint{
				model.OkCheckPoint("base directory parameter is some-file"),
				model.ErrorCheckPoint(filepath.Join(currentDir, "some-file"), " exists but is not a directory"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithBaseDir(test.value))
			initTestCheckEnv(p)
			checkEnv.sourceTreeErr = test.sourceTreeErr
			assert.Equal(t, test.expected, checkBaseDirectory(p))
		})
	}
}

func Test_check_work_directory(t *testing.T) {
	currentDir, _ := filepath.Abs(".")
	tests := []struct {
		desc       string
		value      string
		workDirErr error
		expected   []model.CheckPoint
	}{
		{
			"not set", "", nil,
			[]model.CheckPoint{
				model.OkCheckPoint("work directory parameter is not set explicitly"),
				model.OkCheckPoint("work directory absolute path is ", currentDir),
			},
		},
		{
			"set and exists", ".", nil,
			[]model.CheckPoint{
				model.OkCheckPoint("work directory parameter is ."),
				model.OkCheckPoint("work directory absolute path is ", currentDir),
			},
		},
		{
			"set but does not exist", "missing-dir",
			fs.ErrNotExist,
			[]model.CheckPoint{
				model.OkCheckPoint("work directory parameter is missing-dir"),
				model.ErrorCheckPoint("directory not found: missing-dir"),
			},
		},
		{
			"set but insufficient permissions", "no-perm-dir",
			fs.ErrPermission,
			[]model.CheckPoint{
				model.OkCheckPoint("work directory parameter is no-perm-dir"),
				model.ErrorCheckPoint("cannot access directory no-perm-dir"),
			},
		},
		{
			"set but not a directory", "some-file",
			errors.New(filepath.Join(currentDir, "some-file") + " exists but is not a directory"),
			[]model.CheckPoint{
				model.OkCheckPoint("work directory parameter is some-file"),
				model.ErrorCheckPoint(filepath.Join(currentDir, "some-file"), " exists but is not a directory"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithWorkDir(test.value))
			initTestCheckEnv(p)
			checkEnv.workDirErr = test.workDirErr
			assert.Equal(t, test.expected, checkWorkDirectory(p))
		})
	}
}

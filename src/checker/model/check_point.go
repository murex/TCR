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

package model

import (
	"fmt"
	"github.com/murex/tcr/report"
	"os"
)

// CheckPoint is used for describing the result of a single check point
type CheckPoint struct {
	rc          CheckStatus
	description string
}

// CheckpointsForDirAccessError is a utility function providing usual check points
// related to accessing a directory
func CheckpointsForDirAccessError(dir string, err error) []CheckPoint {
	var checkpoint CheckPoint
	if os.IsNotExist(err) {
		checkpoint = ErrorCheckPoint("directory not found: ", dir)
	} else if os.IsPermission(err) {
		checkpoint = ErrorCheckPoint("cannot access directory ", dir)
	} else {
		checkpoint = ErrorCheckPoint(err)
	}
	return []CheckPoint{checkpoint}
}

// CheckpointsForList is a utility function providing checkpoints for a list of values
func CheckpointsForList(headerMsg string, emptyMsg string, values ...string) (cp []CheckPoint) {
	if len(values) == 0 {
		cp = append(cp, WarningCheckPoint(emptyMsg))
		return cp
	}
	cp = append(cp, OkCheckPoint(headerMsg))
	for _, value := range values {
		cp = append(cp, OkCheckPoint("- ", value))
	}
	return cp
}

// OkCheckPoint creates a checkpoint with OK status
func OkCheckPoint(a ...any) CheckPoint {
	return CheckPoint{rc: CheckStatusOk, description: fmt.Sprint(a...)}
}

// WarningCheckPoint creates a checkpoint with Warning status
func WarningCheckPoint(a ...any) CheckPoint {
	return CheckPoint{rc: CheckStatusWarning, description: fmt.Sprint(a...)}
}

// ErrorCheckPoint creates a checkpoint with Error status
func ErrorCheckPoint(a ...any) CheckPoint {
	return CheckPoint{rc: CheckStatusError, description: fmt.Sprint(a...)}
}

// Print prints a checkpoint. Color and icon depend on the checkpoint's return code
func (cp CheckPoint) Print() {
	switch cp.rc {
	case CheckStatusOk:
		report.PostInfo("\t✔ ", cp.description)
	case CheckStatusWarning:
		report.PostWarning("\t● ", cp.description)
	case CheckStatusError:
		report.PostError("\t▼ ", cp.description)
	}
}

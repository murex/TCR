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

import "github.com/murex/tcr/report"

// CheckGroup contains an aggregation of CheckPoint for a specific topic
type CheckGroup struct {
	topic       string
	checkpoints []CheckPoint
}

// NewCheckGroup creates a new CheckGroup instance
func NewCheckGroup(topic string) *CheckGroup {
	return &CheckGroup{
		topic:       topic,
		checkpoints: []CheckPoint{},
	}
}

// Add adds the provided checkpoints to the CheckGroup instance
func (cg *CheckGroup) Add(checkpoints ...CheckPoint) {
	cg.checkpoints = append(cg.checkpoints, checkpoints...)
}

// Ok adds a checkpoint with status Ok to the CheckGroup instance
func (cg *CheckGroup) Ok(a ...any) {
	cg.Add(OkCheckPoint(a...))
}

// Warning adds a checkpoint with status Warning to the CheckGroup instance
func (cg *CheckGroup) Warning(a ...any) {
	cg.Add(WarningCheckPoint(a...))
}

// Error adds a checkpoint with status Error to the CheckGroup instance
func (cg *CheckGroup) Error(a ...any) {
	cg.Add(ErrorCheckPoint(a...))
}

// GetStatus returns the CheckGroup status. This status is the max value of all
// its contained checkpoints statuses (with Error > Warning > Ok)
func (cg *CheckGroup) GetStatus() (s CheckStatus) {
	s = CheckStatusOk
	if cg == nil {
		return CheckStatusOk
	}
	for _, checkpoint := range cg.checkpoints {
		if checkpoint.rc > s {
			s = checkpoint.rc
		}
	}
	return s
}

// Print prints a report for the CheckGroup and all its contained checkpoints.
// Print color depends on the CheckGroup's status
func (cg *CheckGroup) Print() {
	if cg == nil {
		return
	}
	report.PostInfo()
	const messagePrefix = "➤ checking "
	switch cg.GetStatus() {
	case CheckStatusOk:
		report.PostInfo(messagePrefix, cg.topic)
	case CheckStatusWarning:
		report.PostWarning(messagePrefix, cg.topic)
	case CheckStatusError:
		report.PostError(messagePrefix, cg.topic)
	}
	for _, checkpoint := range cg.checkpoints {
		checkpoint.Print()
	}
}

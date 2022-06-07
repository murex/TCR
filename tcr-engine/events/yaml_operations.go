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

package events

import (
	"bytes"
	"github.com/murex/tcr/tcr-engine/report"
	"gopkg.in/yaml.v3"
	"time"
)

func tcrEventToYaml(event TcrEvent) string {
	return newTcrEventYaml(event).marshal()
}

func yamlToTcrEvent(yaml string) TcrEvent {
	return unmarshal(yaml).toTcrEvent()
}

// TcrEventYaml provides the mapping between of a TCR event in yaml.
type TcrEventYaml struct {
	ModifiedSrcLines  int           `yaml:"modified-src-lines"`
	ModifiedTestLines int           `yaml:"modified-test-lines"`
	TotalTestsRun     int           `yaml:"total-tests-run"`
	TestsPassed       int           `yaml:"tests-passed"`
	TestsFailed       int           `yaml:"tests-failed"`
	TestsSkipped      int           `yaml:"tests-skipped"`
	TestsWithErrors   int           `yaml:"tests-with-errors"`
	TestsDuration     time.Duration `yaml:"tests-duration"`
}

func newTcrEventYaml(event TcrEvent) TcrEventYaml {
	return TcrEventYaml{
		ModifiedSrcLines:  event.ModifiedSrcLines,
		ModifiedTestLines: event.ModifiedTestLines,
		TotalTestsRun:     event.TotalTestsRun,
		TestsPassed:       event.TestsPassed,
		TestsFailed:       event.TestsFailed,
		TestsSkipped:      event.TestsSkipped,
		TestsWithErrors:   event.TestsWithErrors,
		TestsDuration:     event.TestsDuration,
	}
}

func (e TcrEventYaml) toTcrEvent() TcrEvent {
	return TcrEvent{
		ModifiedSrcLines:  e.ModifiedSrcLines,
		ModifiedTestLines: e.ModifiedTestLines,
		TotalTestsRun:     e.TotalTestsRun,
		TestsPassed:       e.TestsPassed,
		TestsFailed:       e.TestsFailed,
		TestsSkipped:      e.TestsSkipped,
		TestsWithErrors:   e.TestsWithErrors,
		TestsDuration:     e.TestsDuration,
	}
}

func (e TcrEventYaml) marshal() string {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(0)
	err := yamlEncoder.Encode(&e)
	if err != nil {
		report.PostWarning(err)
	}
	return b.String()
}

// unmarshal un-marshals a yaml string into a TcrEventYaml struct
func unmarshal(yamlString string) (out TcrEventYaml) {
	if err := yaml.Unmarshal([]byte(yamlString), &out); err != nil {
		report.PostWarning(err)
	}
	return
}

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

type (
	// ChangedLinesYaml provides the YAML structure containing info related to the lines changes in src and test
	ChangedLinesYaml struct {
		Src  int `yaml:"src"`
		Test int `yaml:"test"`
	}

	// TestStatsYaml provides the YAML structure containing info related to the tests execution
	TestStatsYaml struct {
		Run      int           `yaml:"run"`
		Passed   int           `yaml:"passed"`
		Failed   int           `yaml:"failed"`
		Skipped  int           `yaml:"skipped"`
		Error    int           `yaml:"error"`
		Duration time.Duration `yaml:"duration"`
	}

	// TcrEventYaml provides the YAML structure containing information related to a TCR event
	TcrEventYaml struct {
		Changes ChangedLinesYaml `yaml:"changed-lines"`
		Tests   TestStatsYaml    `yaml:"test-stats"`
	}
)

func tcrEventToYaml(event TcrEvent) string {
	return newTcrEventYaml(event).marshal()
}

func yamlToTcrEvent(yaml string) TcrEvent {
	return unmarshal(yaml).toTcrEvent()
}

func newTcrEventYaml(event TcrEvent) TcrEventYaml {
	return TcrEventYaml{
		Changes: ChangedLinesYaml(event.Changes),
		Tests:   TestStatsYaml(event.Tests),
	}
}

func (event TcrEventYaml) toTcrEvent() TcrEvent {
	return TcrEvent{
		Changes: ChangedLines(event.Changes),
		Tests:   TestStats(event.Tests),
	}
}

func (event TcrEventYaml) marshal() string {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(0)
	err := yamlEncoder.Encode(&event)
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

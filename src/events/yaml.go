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
	"time"

	"github.com/murex/tcr/report"
	"gopkg.in/yaml.v3"
)

type (
	// ChangedLinesYAML provides the YAML structure containing info related to the lines changes in src and test
	ChangedLinesYAML struct {
		Src  int `yaml:"src"`
		Test int `yaml:"test"`
	}

	// TestStatsYAML provides the YAML structure containing info related to the tests execution
	TestStatsYAML struct {
		Run      int           `yaml:"run"`
		Passed   int           `yaml:"passed"`
		Failed   int           `yaml:"failed"`
		Skipped  int           `yaml:"skipped"`
		Error    int           `yaml:"error"`
		Duration time.Duration `yaml:"duration"`
	}

	// TCREventYAML provides the YAML structure containing information related to a TCR event
	TCREventYAML struct {
		Changes ChangedLinesYAML `yaml:"changed-lines"`
		Tests   TestStatsYAML    `yaml:"test-stats"`
	}
)

func tcrEventToYAML(event TCREvent) string {
	return newTCREventYAML(event).marshal()
}

func yamlToTCREvent(yaml string) TCREvent {
	return unmarshal(yaml).toTCREvent()
}

func newTCREventYAML(event TCREvent) TCREventYAML {
	return TCREventYAML{
		Changes: ChangedLinesYAML(event.Changes),
		Tests:   TestStatsYAML(event.Tests),
	}
}

func (event TCREventYAML) toTCREvent() TCREvent {
	return NewTCREvent(StatusUnknown, ChangedLines(event.Changes), TestStats(event.Tests))
}

func (event TCREventYAML) marshal() string {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(0)
	err := yamlEncoder.Encode(&event)
	if err != nil {
		report.PostWarning(err)
	}
	return b.String()
}

// unmarshal un-marshals a yaml string into a TCREventYAML struct
func unmarshal(yamlString string) (out TCREventYAML) {
	if err := yaml.Unmarshal([]byte(yamlString), &out); err != nil {
		report.PostWarning(err)
	}
	return out
}

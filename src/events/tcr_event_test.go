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
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_yaml_conversion_on_empty_tcr_event(t *testing.T) {
	event := *ATcrEvent()
	expected := buildYAMLString("0", "0", "0", "0", "0", "0", "0", "0s")

	yaml := event.ToYAML()
	assert.Equal(t, expected, yaml)
	assert.Equal(t, event, FromYAML(yaml))
}

func Test_yaml_conversion_on_sample_tcr_event(t *testing.T) {
	event := ATcrEvent(
		WithModifiedSrcLines(1),
		WithModifiedTestLines(2),
		WithTotalTestsRun(12),
		WithTestsPassed(3),
		WithTestsFailed(4),
		WithTestsSkipped(5),
		WithTestsWithErrors(6),
		WithTestsDuration(10*time.Second),
	)
	expected := buildYAMLString("1", "2", "12", "3", "4", "5", "6", "10s")

	yaml := event.ToYAML()
	assert.Equal(t, expected, yaml)
	assert.Equal(t, *event, FromYAML(yaml))
}

func Test_ChangedLines_can_sum_its_total_line_changes(t *testing.T) {
	changedLines := ChangedLines{1, 2}

	assert.Equal(t, 1+2, changedLines.All())
}

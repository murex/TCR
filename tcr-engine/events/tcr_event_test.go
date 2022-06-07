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

func Test_empty_tcr_event_conversion_to_yaml(t *testing.T) {
	event := *ATcrEvent()
	expected := buildYamlString("0", "0", "0", "0", "0", "0", "0", "0s")
	assert.Equal(t, expected, event.ToYaml())
}

func Test_sample_tcr_event_conversion_to_yaml(t *testing.T) {
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
	expected := buildYamlString("1", "2", "12", "3", "4", "5", "6", "10s")
	assert.Equal(t, expected, event.ToYaml())
}

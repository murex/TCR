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

package stats

import (
	"github.com/murex/tcr/tcr-engine/events"
	"github.com/murex/tcr/tcr-engine/report"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_print_stat(t *testing.T) {
	testFlags := []struct {
		desc     string
		name     string
		value    interface{}
		expected string
	}{
		{
			"int value",
			"some stat",
			5,
			"- some stat:           5",
		},
		{
			"string value",
			"some stat",
			"some value",
			"- some stat:           some value",
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			sniffer := report.NewSniffer()
			printStat(tt.name, tt.value)
			sniffer.Stop()
			assert.Equal(t, 1, sniffer.GetMatchCount())
			assert.Equal(t, tt.expected, sniffer.GetAllMatches()[0].Text)
		})
	}
}

func Test_print_stat_with_percentage(t *testing.T) {
	testFlags := []struct {
		desc       string
		name       string
		value      interface{}
		percentage int
		expected   string
	}{
		{
			"int value at 0%",
			"some stat",
			5,
			0,
			"- some stat:           5 (0%)",
		},
		{
			"string value at 12%",
			"some stat",
			"some value",
			12,
			"- some stat:           some value (12%)",
		},
		{
			"boolean value at 100%",
			"some stat",
			false,
			100,
			"- some stat:           false (100%)",
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			sniffer := report.NewSniffer()
			printStatWithPercentage(tt.name, tt.value, tt.percentage)
			sniffer.Stop()
			assert.Equal(t, 1, sniffer.GetMatchCount())
			assert.Equal(t, tt.expected, sniffer.GetAllMatches()[0].Text)
		})
	}
}

func Test_print(t *testing.T) {
	branch := "some-branch"
	inputEvents := events.TcrEvents{
		*events.ADatedTcrEvent(
			events.WithTimestamp(time.Date(2022, 9, 24, 13, 24, 35, 0, time.UTC)),
			events.WithTcrEvent(*events.ATcrEvent(events.WithCommandStatus(events.StatusFail))),
		),
		*events.ADatedTcrEvent(
			events.WithTimestamp(time.Date(2022, 9, 24, 13, 51, 12, 0, time.UTC)),
			events.WithTcrEvent(*events.ATcrEvent(events.WithCommandStatus(events.StatusPass))),
		),
		*events.ADatedTcrEvent(
			events.WithTimestamp(time.Date(2022, 9, 24, 14, 42, 33, 0, time.UTC)),
			events.WithTcrEvent(*events.ATcrEvent(events.WithCommandStatus(events.StatusFail))),
		),
	}
	expected := []string{
		"- Branch:              some-branch",
		"- First commit:        2022-09-24 13:24:35 +0000 UTC",
		"- Last commit:         2022-09-24 14:42:33 +0000 UTC",
		"- Number of commits:   3",
		"- Passing commits:     1 (33%)",
		"- Failing commits:     2 (67%)",
		"- Time span:           1h17m58s",
		"- Time in green:       51m21s (66%)",
		"- Time in red:         26m37s (34%)",
	}
	sniffer := report.NewSniffer()
	Print(branch, inputEvents)
	time.Sleep(1 * time.Second)
	sniffer.Stop()

	assert.Equal(t, len(expected), sniffer.GetMatchCount())
	var result []string
	for _, line := range sniffer.GetAllMatches() {
		result = append(result, line.Text)
	}
	assert.Equal(t, expected, result)
}

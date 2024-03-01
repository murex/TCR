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
	"github.com/murex/tcr/events"
	"github.com/murex/tcr/report"
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
			"- some stat:                 5",
		},
		{
			"string value",
			"some stat",
			"some value",
			"- some stat:                 some value",
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			sniffer := report.NewSniffer()
			printStat(tt.name, tt.value)
			sniffer.Stop()
			assert.Equal(t, 1, sniffer.GetMatchCount())
			assert.Equal(t, tt.expected, sniffer.GetAllMatches()[0].Payload.ToString())
		})
	}
}

func Test_print_stat_value_and_ratio(t *testing.T) {
	testFlags := []struct {
		desc     string
		name     string
		stat     events.ValueAndRatio
		expected string
	}{
		{
			"int value",
			"some stat",
			events.IntValueAndRatio{},
			"- some stat:                 0 (0%)",
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			sniffer := report.NewSniffer()
			printStatValueAndRatio(tt.name, tt.stat)
			sniffer.Stop()
			assert.Equal(t, 1, sniffer.GetMatchCount())
			assert.Equal(t, tt.expected, sniffer.GetAllMatches()[0].Payload.ToString())
		})
	}
}

func Test_print_stat_min_max_avg(t *testing.T) {
	testFlags := []struct {
		desc     string
		name     string
		stat     events.Aggregates
		expected string
	}{
		{
			"duration aggregates",
			"some stat",
			events.DurationAggregates{},
			"- some stat:                 0s (min) / 0s (avg) / 0s (max)",
		},
		{
			"int aggregates",
			"some stat",
			events.IntAggregates{},
			"- some stat:                 0 (min) / 0 (avg) / 0 (max)",
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			sniffer := report.NewSniffer()
			printStatMinMaxAvg(tt.name, tt.stat)
			time.Sleep(1 * time.Millisecond)
			sniffer.Stop()
			assert.Equal(t, 1, sniffer.GetMatchCount())
			assert.Equal(t, tt.expected, sniffer.GetAllMatches()[0].Payload.ToString())
		})
	}
}

func Test_print_stat_evolution(t *testing.T) {
	testFlags := []struct {
		desc     string
		name     string
		value    events.ValueEvolution
		expected string
	}{
		{
			desc:     "int value evolution",
			name:     "some stat",
			value:    events.IntValueEvolution{},
			expected: "- some stat:                 0 --> 0",
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			sniffer := report.NewSniffer()
			printStatEvolution(tt.name, tt.value)
			time.Sleep(1 * time.Millisecond)
			sniffer.Stop()
			assert.Equal(t, 1, sniffer.GetMatchCount())
			assert.Equal(t, tt.expected, sniffer.GetAllMatches()[0].Payload.ToString())
		})
	}
}

func Test_print_all_stats(t *testing.T) {
	branch := "some-branch"
	inputEvents := events.TcrEvents{
		*events.ADatedTcrEvent(
			events.WithTimestamp(time.Date(2022, 9, 22, 13, 24, 35, 0, time.UTC)),
			events.WithTcrEvent(*events.ATcrEvent(
				events.WithCommandStatus(events.StatusFail),
				events.WithModifiedSrcLines(1),
				events.WithModifiedTestLines(0),
				events.WithTestsPassed(2),
				events.WithTestsFailed(1),
				events.WithTestsSkipped(5),
				events.WithTestsDuration(500*time.Millisecond),
			)),
		),
		*events.ADatedTcrEvent(
			events.WithTimestamp(time.Date(2022, 9, 22, 13, 51, 12, 0, time.UTC)),
			events.WithTcrEvent(*events.ATcrEvent(
				events.WithCommandStatus(events.StatusPass),
				events.WithModifiedSrcLines(10),
				events.WithModifiedTestLines(3),
				events.WithTestsPassed(3),
				events.WithTestsFailed(1),
				events.WithTestsSkipped(3),
				events.WithTestsDuration(1*time.Second),
			)),
		),
		*events.ADatedTcrEvent(
			events.WithTimestamp(time.Date(2022, 9, 22, 14, 42, 33, 0, time.UTC)),
			events.WithTcrEvent(*events.ATcrEvent(
				events.WithCommandStatus(events.StatusFail),
				events.WithModifiedSrcLines(4),
				events.WithModifiedTestLines(1),
				events.WithTestsPassed(8),
				events.WithTestsFailed(2),
				events.WithTestsSkipped(1),
				events.WithTestsDuration(2*time.Second),
			)),
		),
	}
	expected := []string{
		"- Branch:                    some-branch",
		"- First commit:              Thursday 22 Sep 2022 at 13:24:35",
		"- Last commit:               Thursday 22 Sep 2022 at 14:42:33",
		"- Number of commits:         3",
		"- Passing commits:           1 (33%)",
		"- Failing commits:           2 (67%)",
		"- Time span:                 1h17m58s",
		"- Time in green:             51m21s (66%)",
		"- Time in red:               26m37s (34%)",
		"- Time between commits:      26m37s (min) / 38m59s (avg) / 51m21s (max)",
		"- Changes per commit (src):  1 (min) / 5 (avg) / 10 (max)",
		"- Changes per commit (test): 0 (min) / 1.3 (avg) / 3 (max)",
		"- Passing tests count:       2 --> 8",
		"- Failing tests count:       1 --> 2",
		"- Skipped tests count:       5 --> 1",
		"- Test execution duration:   500ms --> 2s",
	}
	sniffer := report.NewSniffer()
	Print(branch, inputEvents)
	time.Sleep(1 * time.Millisecond)
	sniffer.Stop()

	assert.Equal(t, len(expected), sniffer.GetMatchCount())
	var result []string
	for _, line := range sniffer.GetAllMatches() {
		result = append(result, line.Payload.ToString())
	}
	assert.Equal(t, expected, result)
}

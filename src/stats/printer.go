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
	"fmt"
	"github.com/murex/tcr/events"
	"github.com/murex/tcr/report"
	"time"
)

// Print prints all TCR stats for the provided list of TCR events.
func Print(branch string, tcrEvents events.TcrEvents) {
	printStat("Branch", branch)
	printHumanDate("First commit", tcrEvents.StartingTime())
	printHumanDate("Last commit", tcrEvents.EndingTime())
	printStat("Number of commits", tcrEvents.NbRecords())
	printStatValueAndRatio("Passing commits", tcrEvents.PassingRecords())
	printStatValueAndRatio("Failing commits", tcrEvents.FailingRecords())
	printStat("Time span", tcrEvents.TimeSpan())
	printStatValueAndRatio("Time in green", tcrEvents.DurationInGreen())
	printStatValueAndRatio("Time in red", tcrEvents.DurationInRed())
	printStatMinMaxAvg("Time between commits", tcrEvents.TimeBetweenCommits())
	printStatMinMaxAvg("Changes per commit (src)", tcrEvents.SrcLineChangesPerCommit())
	printStatMinMaxAvg("Changes per commit (test)", tcrEvents.TestLineChangesPerCommit())
	printStatEvolution("Passing tests count", tcrEvents.PassingTestsEvolution())
	printStatEvolution("Failing tests count", tcrEvents.FailingTestsEvolution())
	printStatEvolution("Skipped tests count", tcrEvents.SkippedTestsEvolution())
	printStatEvolution("Test execution duration", tcrEvents.TestDurationEvolution())
}

func printStatEvolution(name string, stat events.ValueEvolution) {
	// printStat(name, "from ", stat.From(), " to ", stat.To())
	printStat(name, stat.From(), " --> ", stat.To())
}

func printStatMinMaxAvg(name string, stat events.Aggregates) {
	printStat(name, stat.Min(), " (min) / ", stat.Avg(), " (avg) / ", stat.Max(), " (max)")
}

func printStatValueAndRatio(name string, stat events.ValueAndRatio) {
	printStat(name, stat.Value(), " (", stat.Percentage(), "%)")
}

func printStat(name string, stat ...any) {
	report.PostInfo(
		fmt.Sprintf("- %-26s ", name+":"),
		fmt.Sprint(stat...),
	)
}

func printHumanDate(name string, t time.Time) {
	printStat(name, humanDate(t))
}

// humanDate function returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("Monday 02 Jan 2006 at 15:04:05")
}

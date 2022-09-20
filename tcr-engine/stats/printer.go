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
	"github.com/murex/tcr/tcr-engine/events"
	"github.com/murex/tcr/tcr-engine/report"
)

// Print prints all TCR stats for the provided list of TCR events.
func Print(branch string, tcrEvents events.TcrEvents) {
	printStat("Branch", branch)
	printStat("First commit", tcrEvents.StartingTime())
	printStat("Last commit", tcrEvents.EndingTime())
	printStat("Number of commits", tcrEvents.NbRecords())
	printStatWithPercentage("Passing commits", tcrEvents.NbPassingRecords(), tcrEvents.PercentPassing())
	printStatWithPercentage("Failing commits", tcrEvents.NbFailingRecords(), tcrEvents.PercentFailing())
	printStat("Time span", tcrEvents.TimeSpan())
	printStatWithPercentage("Time in green", tcrEvents.DurationInGreen(), tcrEvents.PercentDurationInGreen())
	printStatWithPercentage("Time in red", tcrEvents.DurationInRed(), tcrEvents.PercentDurationInRed())
	printStatMinMaxAvg("Time between commits", tcrEvents.TimeBetweenCommits())
	printStatMinMaxAvg("Changes per commit (src)", tcrEvents.SrcLineChangesPerCommit())
	printStatMinMaxAvg("Changes per commit (test)", tcrEvents.TestLineChangesPerCommit())
	// TODO add stats on test cases (executed, failed, passed, percentage)
}

func printStatMinMaxAvg(name string, values events.Aggregates) {
	printStat(name, values.Min(), " (min) / ", values.Avg(), " (avg) / ", values.Max(), " (max)")
}

func printStatWithPercentage(name string, value interface{}, percentage int) {
	printStat(name, value, " (", percentage, "%)")
}

func printStat(name string, value ...interface{}) {
	report.PostInfo(
		fmt.Sprintf("- %-26s ", name+":"),
		fmt.Sprint(value...),
	)
}

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

package cmd

import (
	"github.com/murex/tcr/cli"
	"github.com/murex/tcr/engine"
	"github.com/murex/tcr/runmode"
	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Print TCR stats",
	Long: `
TCR stats subcommand prints out TCR usage stats based on commit history.

The commit history is retrieved for the repository containing
TCR base directory (cf. -b option). The branch is the current working
branch set for this repository.

The following stats are reported:

- First commit date and time
- Last commit date and time
- Number of commits
- Number of passing commits (absolute value and percentage) (*)
- Number of failing commits, (absolute value and percentage) (*)
- Time span between the first and last commit
- Time in green: total time where all tests passed (absolute value and percentage) (*)
- Time in red: total time where one or more tests failed (absolute value and percentage) (*)
- Time between commits (minimum, average and maximum values)
- Changes per commit (src): number of lines of source code changed per commit (minimum, average and maximum values)
- Changes per commit (test): number of lines of test code changed per commit (minimum, average and maximum values)
- Passing tests count evolution (values for first and last commit)
- Failing tests count evolution (values for first and last commit) (*)
- Skipped tests count evolution (values for first and last commit)
- Test execution duration cumulated for all tests (values for first and last commit)

> (*) These metrics are relevant only if TCR commit history was created while running TCR with "commit-failures" option.
> Without this option there is no record of test failures in TCR commit history, thus:
> - "Number of passing commits" and "time in green" will always be at 100%
> - "Number of failing commits" and "time in red" will always be at 0%
> - "Failing tests" will always be at 0

This subcommand does not start TCR engine.`,
	Run: func(cmd *cobra.Command, args []string) {
		parameters.Mode = runmode.Stats{}
		u := cli.New(parameters, engine.NewTcrEngine())
		u.Start()
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

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
	"github.com/murex/tcr/tcr-cli/cli"
	"github.com/murex/tcr/tcr-engine/engine"
	"github.com/murex/tcr/tcr-engine/runmode"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var oneShotCmd = &cobra.Command{
	Use:   "one-shot",
	Short: "Run one TCR cycle and exit",
	Long: `
When used in "one-shot" mode, TCR executes instantaneously one TCR cycle without waiting
for any file system change, then exits.

The return code of TCR "one-shot" is one of the following:

| RC  | Meaning                                                                        |
|-----|--------------------------------------------------------------------------------|
| 0   | Build and Test Passed and changes were successfully committed                  |
| 1   | Build failed                                                                   |
| 2   | Build passed, one or more tests failed, and changes were successfully reverted |
| 3   | Error in configuration or parameter values                                     |
| 4   | Error while interacting with git                                               |
| 5   | Any other error                                                                |

This subcommand runs directly in the terminal (no GUI).`,
	Run: func(cmd *cobra.Command, args []string) {
		parameters.Mode = runmode.OneShot{}
		parameters.AutoPush = parameters.Mode.AutoPushDefault()
		u := cli.New(parameters, engine.NewTcrEngine())
		u.Start()
	},
}

func init() {
	rootCmd.AddCommand(oneShotCmd)
}

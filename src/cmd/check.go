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

// runCmd represents the run command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check TCR configuration and parameters and exit",
	Long: `
When used in "check" mode, TCR performs a series of verifications on provided parameters,
configuration and local environment, then exits.

Its purpose is to verify that TCR is ready to run. It does not trigger any TCR cycle execution.

The main checkpoints are organized in sections as follows:

- Configuration directory
- Base directory
- Work Directory
- Language settings
- Toolchain settings
- Git environment
- Auto-push settings
- Mob timer settings (for driver role)
- Polling period settings (for navigator role)

The return code of TCR "check" is one of the following:

| RC  | Meaning                                                                    |
|-----|----------------------------------------------------------------------------|
| 0   | All checks passed without any warning or error                             |
| 1   | One or more warnings were raised. This should not prevent TCR from running |
| 2   | One or more errors were raised. TCR will not be able to run properly       |
`,
	Run: func(cmd *cobra.Command, args []string) {
		parameters.Mode = runmode.Check{}
		u := cli.New(parameters, engine.NewTcrEngine())
		u.Start()
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

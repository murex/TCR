/*
Copyright (c) 2021 Murex

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
	"github.com/murex/tcr/tcr-engine/runmode"
	"github.com/murex/tcr/tcr-engine/settings"

	"github.com/spf13/cobra"
)

// soloCmd represents the solo command
var soloCmd = &cobra.Command{
	Use:   "solo",
	Short: "Run TCR in solo mode",
	Long: `
When used in "solo" mode, TCR only commits changes locally.
It never pushes or pulls to a remote repository.

This subcommand runs directly in the terminal (no GUI).`,
	Run: func(cmd *cobra.Command, args []string) {
		params.Mode = runmode.Solo{}
		params.AutoPush = params.Mode.AutoPushDefault()
		params.PollingPeriod = settings.DefaultPollingPeriod
		u := cli.New(params)
		u.Start()
	},
}

func init() {
	rootCmd.AddCommand(soloCmd)
}

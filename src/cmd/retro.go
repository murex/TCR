/*
Copyright (c) 2023 Murex

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

// retroCmd represents the retro command
var retroCmd = &cobra.Command{
	Use:   "retro",
	Short: "Generate retrospective template with stats",
	Long: `
TCR retro subcommand generates a retrospective template in markdown format prefilled with TCR execution info.
The markdown file is saved into the TCR base directory with the name 'tcr-retro.md'. 

The following information is included in the markdown:

- Average size of changes per passing commit 
- Average size of changes per failing commit

These stats are extracted for the repository containing TCR base directory (cf. -b option). 
The branch is the current working branch set for this repository.

This subcommand does not start TCR engine.`,
	Run: func(_ *cobra.Command, _ []string) {
		parameters.Mode = runmode.Retro{}
		u := cli.New(parameters, engine.NewTCREngine())
		u.Start()
	},
}

func init() {
	rootCmd.AddCommand(retroCmd)
}

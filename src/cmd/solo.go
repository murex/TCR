/*
Copyright (c) 2024 Murex

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
	"github.com/murex/tcr/http"
	"github.com/murex/tcr/runmode"
	"github.com/spf13/cobra"
)

var soloCmd = &cobra.Command{
	Use:   "solo",
	Short: "Run TCR in solo mode",
	Long: `
When used in "solo" mode, TCR only commits changes locally.
It never pushes or pulls to a remote repository.
`,
	Run: func(_ *cobra.Command, _ []string) {
		parameters.Mode = runmode.Solo{}
		parameters.AutoPush = parameters.Mode.AutoPushDefault()

		// Create TCR engine and UI instances
		tcr := engine.NewTCREngine()
		h := http.New(8080, tcr)
		u := cli.New(parameters, tcr)

		// Initialize TCR engine and start UIs
		tcr.Init(parameters)
		h.Start()
		u.Start()
	},
}

func init() {
	rootCmd.AddCommand(soloCmd)
}

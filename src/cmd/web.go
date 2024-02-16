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

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Run TCR with web user interface (experimental)",
	Long: `
When used in "web" mode, TCR starts an HTTP server
allowing to interact with it through a web browser in
addition to the command line interface.

TCR is listening by default to port 8483. The port number
can be changed through the port-number option.

IMPORTANT: This feature is still at an experimental stage!
`,
	Run: func(_ *cobra.Command, _ []string) {
		// Default running mode is mob
		parameters.Mode = runmode.Mob{}
		parameters.AutoPush = parameters.Mode.AutoPushDefault()

		// Forcing port number to 8483 if not set for now (so that HTTP server and webapp options
		// are available only when port number set to a non-0 value
		// Later on if we decide to run HTTP server by default, this value
		// can be set directly in config/param_port_number.go
		// Why 8483? TCR = T&C|R = "84" + "67 | 82" = "84" + "83"
		if parameters.PortNumber == 0 {
			parameters.PortNumber = 8483
		}

		// Create TCR engine and UI instances
		tcr := engine.NewTCREngine()
		h := http.New(parameters, tcr)
		u := cli.New(parameters, tcr)

		// Initialize TCR engine and start UIs
		tcr.Init(parameters)
		h.Start()
		u.Start()
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
}

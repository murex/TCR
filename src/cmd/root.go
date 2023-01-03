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
	"github.com/murex/tcr/cli"
	"github.com/murex/tcr/config"
	"github.com/murex/tcr/engine"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/runmode"
	"github.com/murex/tcr/settings"
	"github.com/spf13/cobra"
)

var parameters params.Params

var (
	rootCmd = &cobra.Command{
		Use:     "tcr",
		Version: settings.BuildVersion,
		Short:   settings.ApplicationShortDescription,
		Long: `
This application is a tool for practicing TCR (Test && Commit || Revert).
It can be used either in solo, or as a group within a mob or pair session.
`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.UpdateEngineParams(&parameters)
		},
		Run: func(cmd *cobra.Command, args []string) {
			parameters.Mode = runmode.Mob{}
			parameters.AutoPush = parameters.Mode.AutoPushDefault()
			u := cli.New(parameters, engine.NewTCREngine())
			u.Start()
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(config.StandardInit)
	// When running TCR-CLI, default base and work dir is the current directory
	config.AddParameters(rootCmd, ".")
}

// GetRootCmd returns the root command. This function is used by the doc package to generate
// the application help markdown files
func GetRootCmd() *cobra.Command {
	return rootCmd
}

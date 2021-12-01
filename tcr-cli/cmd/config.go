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
	"github.com/murex/tcr/tcr-engine/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage TCR configuration",
	Long: `
TCR config subcommand provides management of TCR configuration.

This subcommand does not start TCR engine.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Usage()
	},
}

// showCmd represents the config add command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show TCR configuration",
	Long: `
config show subcommand displays TCR configuration.

This subcommand does not start TCR engine.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.Show()
	},
}

// saveCmd represents the config save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save TCR configuration",
	Long: `
config save subcommand saves TCR configuration into a file.

This subcommand does not start TCR engine.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.Save()
	},
}

// resetCmd represents the config reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset TCR configuration",
	Long: `
config reset subcommand resets TCR configuration to default values.

This subcommand does not start TCR engine.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.Reset()
	},
}

func init() {
	configCmd.AddCommand(showCmd)
	configCmd.AddCommand(resetCmd)
	configCmd.AddCommand(saveCmd)

	rootCmd.AddCommand(configCmd)
}

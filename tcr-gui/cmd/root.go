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
	"github.com/mitchellh/go-homedir"
	"github.com/murex/tcr-engine/engine"
	"github.com/murex/tcr-engine/runmode"
	"github.com/murex/tcr-engine/settings"
	"github.com/murex/tcr-gui/gui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fmt"
	"os"
)

// Command Line Options placeholders

var params engine.Params
var infoFlag bool

var rootCmd = &cobra.Command{
	Use:     "tcr-gui",
	Version: settings.BuildVersion,
	Short:   settings.ApplicationShortDescription,
	Long: `
This application is a tool for practicing TCR (Test && Commit || Revert).
It can be used either in solo, or as a group within a mob or pair session.

This application runs within a GUI.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		printBuildInfo()
	},
	Run: func(cmd *cobra.Command, args []string) {
		// When run without a subcommand, we open the gui by default
		// so that the user can decide what they want to do
		params.Mode = runmode.Mob{}
		params.AutoPush = params.Mode.AutoPushDefault()
		params.PollingPeriod = settings.DefaultPollingPeriod
		u := gui.New(params)
		u.Start()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is $HOME/.tcr.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&params.Toolchain,
		"toolchain",
		"t",
		"",
		"Indicate the toolchain to be used by TCR")

	rootCmd.PersistentFlags().StringVarP(&params.BaseDir,
		"base-dir",
		"b",
		"",
		"Indicate the base directory from which TCR is running")

	rootCmd.Flags().BoolVarP(&params.AutoPush,
		"auto-push",
		"p",
		false,
		"Enable git push after every commit")

	rootCmd.PersistentFlags().DurationVarP(&params.MobTurnDuration,
		"duration",
		"d",
		settings.DefaultMobTurnDuration,
		"Set the duration for role rotation countdown timer")

	rootCmd.PersistentFlags().BoolVarP(&infoFlag,
		"info",
		"i",
		false,
		"show build information about TCR application")
}

func printBuildInfo() {
	if infoFlag {
		settings.PrintBuildInfo()
		os.Exit(0)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if params.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(params.CfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".tcr" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tcr")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// GetRootCmd returns the root command. This function is used by the doc package to generate
// the application help markdown files
func GetRootCmd() *cobra.Command {
	return rootCmd
}

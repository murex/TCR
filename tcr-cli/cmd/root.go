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
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/murex/tcr/tcr-cli/cli"
	"github.com/murex/tcr/tcr-engine/engine"
	"github.com/murex/tcr/tcr-engine/runmode"
	"github.com/murex/tcr/tcr-engine/settings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// Command Line Options placeholders

var (
	configFileParam       *param.StringParam
	mobTimerDurationParam *param.DurationParam
	toolchainParam        *param.StringParam
	pollingPeriod         *param.DurationParam
)

var params engine.Params
var infoFlag bool
var saveCfgFlag bool

var (
	rootCmd = &cobra.Command{
		Use:     "tcr",
		Version: settings.BuildVersion,
		Short:   settings.ApplicationShortDescription,
		Long: `
This application is a tool for practicing TCR (Test && Commit || Revert).
It can be used either in solo, or as a group within a mob or pair session.

This application runs within a terminal.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			printBuildInfo()
			retrieveParams()
		},
		Run: func(cmd *cobra.Command, args []string) {
			//retrieveParams()
			params.Mode = runmode.Mob{}
			params.AutoPush = params.Mode.AutoPushDefault()
			saveConfig()
			u := cli.New(params)
			u.Start()
		},
	}
)

func retrieveParams() {
	params.ConfigFile = configFileParam.GetValue()
	params.MobTurnDuration = mobTimerDurationParam.GetValue()
	params.Toolchain = toolchainParam.GetValue()
	params.PollingPeriod = pollingPeriod.GetValue()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	configFileParam = param.NewConfigFileParam(rootCmd)
	toolchainParam = param.NewToolchainParam(rootCmd)
	pollingPeriod = param.NewPollingPeriodParam(rootCmd)
	mobTimerDurationParam = param.NewMobTimerDurationParam(rootCmd)
	initBaseDirParam()
	initAutoPushParam()
	initInfoParam()
	initSaveCfgParam()
}

func initSaveCfgParam() {
	rootCmd.PersistentFlags().BoolVarP(&saveCfgFlag,
		"save",
		"s",
		false,
		"save configuration file on exit")
}

func initInfoParam() {
	rootCmd.PersistentFlags().BoolVarP(&infoFlag,
		"info",
		"i",
		false,
		"show build information about TCR application")
}

func initAutoPushParam() {
	rootCmd.Flags().BoolVarP(&params.AutoPush,
		"auto-push",
		"p",
		false,
		"enable git push after every commit")
	_ = viper.BindPFlag("params.auto-push", rootCmd.PersistentFlags().Lookup("auto-push"))
}

func initBaseDirParam() {
	rootCmd.PersistentFlags().StringVarP(&params.BaseDir,
		"base-dir",
		"b",
		"",
		"indicate the base directory from which TCR is running")
}

// printBuildInfo prints out application's build information and exits
func printBuildInfo() {
	if infoFlag {
		settings.PrintBuildInfo()
		os.Exit(0)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if params.ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(params.ConfigFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name "tcr.yaml".
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("tcr.yaml")
		viper.SetConfigFile(filepath.Join(home, "tcr.yaml"))
	}

	viper.SetEnvPrefix("TCR")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "["+settings.ApplicationName+"]", "Loading configuration:", viper.ConfigFileUsed())
	}
}

func saveConfig() {
	if saveCfgFlag {
		_, _ = fmt.Fprintln(os.Stderr, "["+settings.ApplicationName+"]", "Saving configuration:", viper.ConfigFileUsed())
		if err := viper.WriteConfig(); err != nil {
			if os.IsNotExist(err) {
				err = viper.WriteConfigAs(viper.ConfigFileUsed())
			} else {
				_, _ = fmt.Fprintln(os.Stderr, "["+settings.ApplicationName+"]", "Error while saving config file:", err)
			}
		}
	}
}

// GetRootCmd returns the root command. This function is used by the doc package to generate
// the application help markdown files
func GetRootCmd() *cobra.Command {
	return rootCmd
}

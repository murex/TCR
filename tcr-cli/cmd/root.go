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
	"github.com/murex/tcr/tcr-engine/config"
	"github.com/murex/tcr/tcr-engine/engine"
	"github.com/murex/tcr/tcr-engine/runmode"
	"github.com/spf13/cobra"
)

// Command Line Options placeholders

var (
	baseDirParam          *config.StringParam
	saveConfigParam       *config.BoolParam
	configFileParam       *config.StringParam
	languageParam         *config.StringParam
	toolchainParam        *config.StringParam
	pollingPeriod         *config.DurationParam
	mobTimerDurationParam *config.DurationParam
	autoPushParam         *config.BoolParam
	buildInfoParam        *config.BoolParam
)

var params engine.Params

var (
	rootCmd = &cobra.Command{
		Use:     "tcr",
		Version: config.BuildVersion,
		Short:   config.ApplicationShortDescription,
		Long: `
This application is a tool for practicing TCR (Test && Commit || Revert).
It can be used either in solo, or as a group within a mob or pair session.

This application runs within a terminal.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			retrieveConfig()
			config.PrintBuildInfo()
		},
		Run: func(cmd *cobra.Command, args []string) {
			params.Mode = runmode.Mob{}
			params.AutoPush = params.Mode.AutoPushDefault()
			u := cli.New(params)
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
	cobra.OnInitialize(initConfig)
	addParameters()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.Init(params.ConfigFile)
}

func addParameters() {
	// When running TCR-CLI, default base dir is the current directory
	baseDirParam = config.AddBaseDirParamWithDefault(rootCmd, ".")
	configFileParam = config.AddConfigFileParam(rootCmd)
	languageParam = config.AddLanguageParam(rootCmd)
	toolchainParam = config.AddToolchainParam(rootCmd)
	pollingPeriod = config.AddPollingPeriodParam(rootCmd)
	mobTimerDurationParam = config.AddMobTimerDurationParam(rootCmd)
	autoPushParam = config.AddAutoPushParam(rootCmd)
	buildInfoParam = config.AddBuildInfoParam(rootCmd)
	saveConfigParam = config.AddSaveConfigParam(rootCmd)
}

func retrieveConfig() {
	params.BaseDir = baseDirParam.GetValue()
	params.ConfigFile = configFileParam.GetValue()
	params.MobTurnDuration = mobTimerDurationParam.GetValue()
	params.Language = languageParam.GetValue()
	params.Toolchain = toolchainParam.GetValue()
	params.PollingPeriod = pollingPeriod.GetValue()
	params.AutoPush = autoPushParam.GetValue()

	config.BuildInfoFlag = buildInfoParam.GetValue()
	config.SaveConfigFlag = saveConfigParam.GetValue()
}

// GetRootCmd returns the root command. This function is used by the doc package to generate
// the application help markdown files
func GetRootCmd() *cobra.Command {
	return rootCmd
}

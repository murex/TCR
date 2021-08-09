package cmd

import (
	"github.com/mengdaming/tcr/tcr"
	"github.com/mengdaming/tcr/tcr/engine"
	"github.com/mengdaming/tcr/tcr/runmode"
	"github.com/mengdaming/tcr/tcr/ui/cli"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fmt"
	"os"
)

// Command Line Options placeholders

var params tcr.Params

var rootCmd = &cobra.Command{
	Use:     "tcr",
	Version: "0.2.3",
	Short:   "TCR (Test && Commit || Revert)",
	Long: `
This application is a tool for practicing TCR (Test && Commit || Revert).
It can be used either in solo, or as a group within a mob or pair session.`,
	Run: func(cmd *cobra.Command, args []string) {
		// When run without a subcommand, we run in mob mode by default
		// so that the user can have access to the menu
		u := cli.New()
		params.Mode = runmode.Mob{}
		params.PollingPeriod = tcr.DefaultPollingPeriod
		engine.Init(u, params)
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
		".",
		"Indicate the base directory from which TCR is running")

	rootCmd.Flags().BoolVarP(&params.AutoPush,
		"auto-push",
		"p",
		false,
		"Enable git push after every commit")
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
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

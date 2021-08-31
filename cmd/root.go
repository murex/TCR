package cmd

import (
	"github.com/mengdaming/tcr/engine"
	"github.com/mengdaming/tcr/runmode"
	"github.com/mengdaming/tcr/ui/gui"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fmt"
	"os"
)

// Command Line Options placeholders

var params engine.Params

var rootCmd = &cobra.Command{
	Use:     "tcr",
	Version: engine.Version,
	Short:   "TCR (Test && Commit || Revert)",
	Long: `
This application is a tool for practicing TCR (Test && Commit || Revert).
It can be used either in solo, or as a group within a mob or pair session.

When called directly without any subcommand, the application opens a GUI.`,
	Run: func(cmd *cobra.Command, args []string) {
		// When run without a subcommand, we open the gui by default
		// so that the user can decide what they want to do
		params.Mode = runmode.Mob{}
		params.AutoPush = params.Mode.AutoPushDefault()
		params.PollingPeriod = engine.DefaultPollingPeriod
		u := gui.New(params)
		u.Start(params.Mode)
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

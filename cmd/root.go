package cmd

import (
	"github.com/mengdaming/tcr/sandbox"
	"github.com/mengdaming/tcr/trace"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fmt"
	"os"
)

var cfgFile string
var toolchain string
var autoPush bool

var rootCmd = &cobra.Command{
	Use:   "tcr",
	Short: "TCR (Test && Commit || Revert)",
	Long: `
This application is a tool to practice TCR.
It can be used either in solo, or as a group within a mob or pair session.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO hook up real application here
		trace.HorizontalLine()
		trace.Info("Running in default mode (sandbox)")

		trace.HorizontalLine()
		trace.Info("Command Line Options")

		var toolchainTrace = fmt.Sprintf("Toolchain = %v", toolchain)
		trace.Info(toolchainTrace)
		var autoPushTrace = fmt.Sprintf("Auto-Push = %v", autoPush)
		trace.Info(autoPushTrace)

		sandbox.KeyboardSandbox()
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tcr.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&toolchain, "toolchain", "t",
		"maven", "indicate the toolchain to be used by TCR")
	rootCmd.Flags().BoolVarP(&autoPush, "auto-push", "p", false, "Enable git push after every commit")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
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

package cmd

import (
	"github.com/mengdaming/tcr/sandbox"
	"github.com/mengdaming/tcr/trace"

	"github.com/spf13/cobra"
)

// sandboxCmd represents the sandbox command
var sandboxCmd = &cobra.Command{
	Use:   "sandbox",
	Short: "Sandbox mode",
	Long: `
Used to experiment with new functionalities or modules.
Don't use unless you know what you're doing!`,
	Run: func(cmd *cobra.Command, args []string) {
		trace.Info("Running in sandbox mode")
		sandbox.ConsoleSandbox()
	},
}

func init() {
	rootCmd.AddCommand(sandboxCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sandboxCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sandboxCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

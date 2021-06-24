
package cmd

import (
	"github.com/mengdaming/tcr/sandbox"
	"github.com/mengdaming/tcr/trace"

	"github.com/spf13/cobra"
)

// keyboardCmd represents the keyboard command
var keyboardCmd = &cobra.Command{
	Use:   "keyboard",
	Short: "Sandbox Keyboard Mode",
	Long: `
Used to experiment with keyboard-related functionalities or modules.
Don't use unless you know what you're doing!`,
	Run: func(cmd *cobra.Command, args []string) {
		trace.Info("Running in sandbox/keyboard mode")
		sandbox.KeyboardSandbox()
	},
}

func init() {
	sandboxCmd.AddCommand(keyboardCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// keyboardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// keyboardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

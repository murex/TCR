
package cmd

import (
	"github.com/mengdaming/tcr/sandbox"
	"github.com/mengdaming/tcr/trace"
	"github.com/spf13/cobra"
)

// cmdCmd represents the cmd command
var cmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "A brief description of your command",
	Long: `
Used to experiment with command-related functionalities or modules.
Don't use unless you know what you're doing!`,
	Run: func(cmd *cobra.Command, args []string) {
		trace.Info("Running in sandbox/cmd mode")
		sandbox.CmdSandbox()
	},
}

func init() {
	sandboxCmd.AddCommand(cmdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cmdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cmdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
